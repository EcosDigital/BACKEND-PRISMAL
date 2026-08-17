-- ============================================================================
-- MIGRACIÓN: Módulo Períodos Contables (Meses Contables)
-- Proyecto : PRISMAR ERP Contable
-- Esquema  : contabilidad
-- Tablas   :
--   contabilidad.ref_estado_periodo       → catálogo de estados de período
--   contabilidad.ref_modulo_contable      → catálogo de módulos del sistema
--   contabilidad.ref_estado_modulo        → catálogo de estados por módulo
--   contabilidad.cfg_periodos_contables   → períodos mensuales por empresa/vigencia
--   contabilidad.cfg_control_modulos_periodo → control individual por módulo
-- ----------------------------------------------------------------------------
-- Patrón: estados mediante FK a tablas ref_*, nunca VARCHAR directo.
-- Idempotente: usa IF NOT EXISTS y ON CONFLICT DO NOTHING.
-- Ejecutar en una sola transacción.
-- ============================================================================

BEGIN;
 
-- ── 1. Catálogo de estados posibles para un período contable ─────────────────
--
-- Patrón ref_* del sistema PRISMAR.
-- El estado del período NUNCA se almacena como texto en la tabla principal.

CREATE TABLE IF NOT EXISTS contabilidad.ref_estado_periodo (
    id          SERIAL      NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    descripcion TEXT,
    color       VARCHAR(20),        -- sugerencia de color para la UI (ej: 'emerald', 'amber')
    CONSTRAINT pk_ref_estado_periodo
        PRIMARY KEY (id),
    CONSTRAINT uq_ref_estado_periodo_nombre
        UNIQUE (nombre)
);

COMMENT ON TABLE contabilidad.ref_estado_periodo IS
    'Catálogo de estados posibles para un período contable mensual. Patrón ref_* del sistema PRISMAR.';
COMMENT ON COLUMN contabilidad.ref_estado_periodo.color IS
    'Sugerencia semántica de color para el frontend (ej: emerald, red, amber, blue).';

-- Semilla: estados iniciales con IDs fijos para referencia segura desde código
INSERT INTO contabilidad.ref_estado_periodo (id, nombre, descripcion, color)
VALUES
    (1, 'Abierto',   'Período activo. Permite registrar movimientos contables en todos los módulos habilitados.', 'emerald'),
    (2, 'Cerrado',   'Período cerrado. Solo permite consultas, sin movimientos.',                                 'gray'),
    (3, 'Ajuste',    'Período en fase de ajuste contable. Solo usuarios autorizados pueden registrar.',          'amber'),
    (4, 'Bloqueado', 'Período completamente bloqueado. No permite ningún movimiento ni ajuste.',                 'red')
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('contabilidad.ref_estado_periodo', 'id'),
    (SELECT MAX(id) FROM contabilidad.ref_estado_periodo)
);

-- ── 2. Catálogo de módulos del sistema que pueden controlarse por período ─────

CREATE TABLE IF NOT EXISTS contabilidad.ref_modulo_contable (
    id          SERIAL      NOT NULL,
    codigo      VARCHAR(30) NOT NULL,
    nombre      VARCHAR(100) NOT NULL,
    descripcion TEXT,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    CONSTRAINT pk_ref_modulo_contable
        PRIMARY KEY (id),
    CONSTRAINT uq_ref_modulo_contable_codigo
        UNIQUE (codigo)
);

COMMENT ON TABLE contabilidad.ref_modulo_contable IS
    'Catálogo de módulos del ERP PRISMAR que pueden ser controlados individualmente por período contable.';

INSERT INTO contabilidad.ref_modulo_contable (id, codigo, nombre, descripcion, is_active)
VALUES
    (1, '001',   'Facturación',           'Control de facturas de venta, notas débito y crédito.',         TRUE),
    (2, '002',  'Contabilidad',           'Control de asientos contables y comprobantes.',                 TRUE),
    (3, '003',    'Inventario',             'Control de entradas, salidas y ajustes de inventario.',         TRUE),
    (4, '004',       'Compras',                'Control de órdenes de compra y facturas de proveedores.',       TRUE),
    (5, '005',     'Tesorería',              'Control de pagos, cobros y movimientos de caja/bancos.',        TRUE),
    (6, '006',        'Nómina',                 'Control de liquidaciones y pagos de nómina.',                   TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('contabilidad.ref_modulo_contable', 'id'),
    (SELECT MAX(id) FROM contabilidad.ref_modulo_contable)
);

-- ── 3. Catálogo de estados operacionales de un módulo dentro de un período ────

CREATE TABLE IF NOT EXISTS contabilidad.ref_estado_modulo (
    id          SERIAL      NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    descripcion TEXT,
    permite_op  BOOLEAN     NOT NULL DEFAULT FALSE,  -- ¿permite operaciones?
    color       VARCHAR(20),
    CONSTRAINT pk_ref_estado_modulo
        PRIMARY KEY (id),
    CONSTRAINT uq_ref_estado_modulo_nombre
        UNIQUE (nombre)
);

COMMENT ON TABLE contabilidad.ref_estado_modulo IS
    'Catálogo de estados operacionales de un módulo dentro de un período contable. Patrón ref_* del sistema PRISMAR.';
COMMENT ON COLUMN contabilidad.ref_estado_modulo.permite_op IS
    'Indica si en este estado se permite ejecutar operaciones sobre el módulo.';


INSERT INTO contabilidad.ref_estado_modulo (id, nombre, descripcion, permite_op, color)
VALUES
    (1, 'Abierto',    'El módulo opera con normalidad en este período.',            TRUE,  'emerald'),
    (2, 'Cerrado',    'El módulo no permite operaciones en este período.',          FALSE, 'gray'),
    (3, 'Solo Lectura','El módulo permite consultas pero no registros nuevos.',     FALSE, 'blue'),
    (4, 'Restringido','El módulo solo permite operaciones a usuarios autorizados.', TRUE,  'amber')
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('contabilidad.ref_estado_modulo', 'id'),
    (SELECT MAX(id) FROM contabilidad.ref_estado_modulo)
);

-- ── 4. Tabla principal: períodos contables (meses) ───────────────────────────
--
-- Cada fila representa un mes calendario dentro de una vigencia fiscal.
-- Generados automáticamente (12 por año fiscal) por la capa de servicios.


CREATE TABLE IF NOT EXISTS contabilidad.cfg_periodos_contables (
    id              BIGSERIAL       NOT NULL,
    id_año_fiscal   BIGINT          NOT NULL,
    empresa_id      BIGINT          NOT NULL,
    numero_mes      SMALLINT        NOT NULL,       -- 1..12
    nombre_mes      VARCHAR(20)     NOT NULL,       -- 'Enero', 'Febrero'...
    fecha_inicio    DATE            NOT NULL,
    fecha_fin       DATE            NOT NULL,
    id_estado       INTEGER         NOT NULL,
    notas           TEXT,
    created_by      BIGINT          NOT NULL,
    updated_by      BIGINT          NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_cfg_periodos_contables
        PRIMARY KEY (id),
 
    -- Unicidad: un solo período por mes/año/empresa
    CONSTRAINT uq_cfg_periodos_contables_mes_año_empresa
        UNIQUE (empresa_id, id_año_fiscal, numero_mes),
 
    -- FK al año fiscal padre
    CONSTRAINT fk_cfg_periodos_contables_año_fiscal
        FOREIGN KEY (id_año_fiscal)
        REFERENCES contabilidad.cfg_años_fiscales (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
 
    -- FK al catálogo de estados de período
    CONSTRAINT fk_cfg_periodos_contables_estado
        FOREIGN KEY (id_estado)
        REFERENCES contabilidad.ref_estado_periodo (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
 
    -- El mes debe ser un valor válido
    CONSTRAINT chk_cfg_periodos_contables_mes
        CHECK (numero_mes BETWEEN 1 AND 12),
 
    -- Coherencia de fechas
    CONSTRAINT chk_cfg_periodos_contables_fechas
        CHECK (fecha_fin >= fecha_inicio)
);

COMMENT ON TABLE contabilidad.cfg_periodos_contables IS
    'Períodos contables mensuales asociados a un año fiscal y empresa. Motor de control transaccional de PRISMAR.';
COMMENT ON COLUMN contabilidad.cfg_periodos_contables.id_año_fiscal IS
    'FK a contabilidad.cfg_años_fiscales. Vigencia fiscal a la que pertenece este período.';
COMMENT ON COLUMN contabilidad.cfg_periodos_contables.numero_mes IS
    'Número de mes: 1 = Enero, 12 = Diciembre.';
COMMENT ON COLUMN contabilidad.cfg_periodos_contables.id_estado IS
    'FK a contabilidad.ref_estado_periodo. No almacenar estado como texto directo.';

-- Índice por empresa + año fiscal (consulta más frecuente: "dame los 12 meses de X año")
CREATE INDEX IF NOT EXISTS idx_cfg_periodos_contables_empresa_año
    ON contabilidad.cfg_periodos_contables (empresa_id, id_año_fiscal);

-- Índice para buscar el período activo por empresa y rango de fechas
CREATE INDEX IF NOT EXISTS idx_cfg_periodos_contables_empresa_estado
    ON contabilidad.cfg_periodos_contables (empresa_id, id_estado);

-- Índice para validación por fecha (CanCreate* queries)
CREATE INDEX IF NOT EXISTS idx_cfg_periodos_contables_fechas
    ON contabilidad.cfg_periodos_contables (empresa_id, fecha_inicio, fecha_fin);

-- Constraint avanzado: solo UN período ABIERTO por empresa/año fiscal
-- Segunda línea de defensa tras la validación del servicio.
-- id_estado = 1 → 'Abierto' (semilla fija)
CREATE UNIQUE INDEX IF NOT EXISTS uq_cfg_periodos_contables_un_abierto_por_año
    ON contabilidad.cfg_periodos_contables (empresa_id, id_año_fiscal)
    WHERE id_estado = 1;
 
COMMENT ON INDEX contabilidad.uq_cfg_periodos_contables_un_abierto_por_año IS
    'Garantiza a nivel de BD que solo exista UN período Abierto por empresa y año fiscal.';

-- ── 5. Tabla de control individual de módulos por período ────────────────────
--
-- Permite que cada módulo tenga su propio estado dentro de un período.
-- Ejemplo: Facturación=Cerrado, Contabilidad=Abierto en el mismo período.

CREATE TABLE IF NOT EXISTS contabilidad.cfg_control_modulos_periodo (
    id              BIGSERIAL   NOT NULL,
    id_periodo      BIGINT      NOT NULL,
    id_modulo       INTEGER     NOT NULL,
    id_estado       INTEGER     NOT NULL,
    empresa_id      BIGINT      NOT NULL,
    notas           TEXT,
    updated_by      BIGINT      NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_cfg_control_modulos_periodo
        PRIMARY KEY (id),
 
    -- Un módulo tiene un único estado por período
    CONSTRAINT uq_cfg_control_modulos_periodo_modulo
        UNIQUE (id_periodo, id_modulo),
 
    CONSTRAINT fk_cfg_control_modulos_periodo_periodo
        FOREIGN KEY (id_periodo)
        REFERENCES contabilidad.cfg_periodos_contables (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
 
    CONSTRAINT fk_cfg_control_modulos_periodo_modulo
        FOREIGN KEY (id_modulo)
        REFERENCES contabilidad.ref_modulo_contable (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
 
    CONSTRAINT fk_cfg_control_modulos_periodo_estado
        FOREIGN KEY (id_estado)
        REFERENCES contabilidad.ref_estado_modulo (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

COMMENT ON TABLE contabilidad.cfg_control_modulos_periodo IS
    'Control granular del estado operacional de cada módulo dentro de un período contable. Permite configuración independiente por módulo.';
COMMENT ON COLUMN contabilidad.cfg_control_modulos_periodo.id_estado IS
    'FK a contabilidad.ref_estado_modulo. Estado operacional del módulo en este período.';
 
CREATE INDEX IF NOT EXISTS idx_cfg_control_modulos_periodo_periodo
    ON contabilidad.cfg_control_modulos_periodo (id_periodo);
 
CREATE INDEX IF NOT EXISTS idx_cfg_control_modulos_periodo_empresa_modulo
    ON contabilidad.cfg_control_modulos_periodo (empresa_id, id_modulo);

-- ── 6. Vista de consulta: períodos con sus estados resueltos ─────────────────
--
-- Facilita queries del frontend y reportes sin JOINs repetidos.

CREATE OR REPLACE VIEW contabilidad.v_periodos_contables AS
SELECT
    p.id,
    p.id_año_fiscal,
    p.empresa_id,
    p.numero_mes,
    p.nombre_mes,
    p.fecha_inicio,
    p.fecha_fin,
    p.id_estado,
    ep.nombre        AS estado,
    ep.color         AS estado_color,
    af.year          AS año_fiscal,
    p.notas,
    p.created_by,
    p.updated_by,
    p.created_at,
    p.updated_at
FROM contabilidad.cfg_periodos_contables p
INNER JOIN contabilidad.ref_estado_periodo ep ON ep.id = p.id_estado
INNER JOIN contabilidad.cfg_años_fiscales af  ON af.id = p.id_año_fiscal;
 
COMMENT ON VIEW contabilidad.v_periodos_contables IS
    'Vista de períodos contables con estados y año fiscal resueltos por JOIN.';
 
COMMIT;