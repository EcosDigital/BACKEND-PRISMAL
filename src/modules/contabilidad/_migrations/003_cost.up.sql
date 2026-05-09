-- ============================================================================
-- MIGRACIÓN: Módulo Año Fiscal
-- Proyecto : PRISMAR ERP Contable
-- Tablas   : contabilidad.ref_estado_año_fiscal
--            contabilidad.cfg_años_fiscales
-- ----------------------------------------------------------------------------
-- Idempotente: usa IF NOT EXISTS en todos los objetos.
-- Ejecutar en orden secuencial dentro de una sola transacción.
-- ============================================================================

BEGIN;

-- ── 1. Tabla de referencia: estados posibles de un año fiscal ─────────────────
--
-- Patrón relacional: el estado NO se guarda como texto en cfg_años_fiscales.
-- Se guarda id_estado INTEGER → FK a esta tabla catálogo.
-- Esto permite agregar/renombrar estados sin tocar la tabla principal.

CREATE TABLE IF NOT EXISTS contabilidad.ref_estado_año_fiscal (
    id          SERIAL      NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    descripcion TEXT,
    CONSTRAINT pk_ref_estado_año_fiscal
        PRIMARY KEY (id),
 
    CONSTRAINT uq_ref_estado_año_fiscal_nombre
        UNIQUE (nombre)
);

COMMENT ON TABLE contabilidad.ref_estado_año_fiscal IS
    'Catálogo de estados posibles para un año fiscal. Patrón ref_* del sistema PRISMAR.';
COMMENT ON COLUMN contabilidad.ref_estado_año_fiscal.nombre IS
    'Nombre legible del estado (ej: Abierto, Cerrado).';

-- ── 2. Datos semilla: estados iniciales ───────────────────────────────────────
--
-- INSERT ... ON CONFLICT DO NOTHING para idempotencia total.
-- Los IDs fijos (1, 2) permiten referenciarlos con seguridad desde el código.

INSERT INTO contabilidad.ref_estado_año_fiscal (id, nombre, descripcion)
VALUES
    (1, 'Abierto',  'Año fiscal activo. Permite registrar movimientos contables.'),
    (2, 'Cerrado',  'Año fiscal cerrado. Solo permite consultas, sin movimientos.')
ON CONFLICT (id) DO NOTHING;

-- Reiniciar la secuencia para que el próximo INSERT automático no colisione
-- con los IDs que acabamos de insertar manualmente.
SELECT setval(
    pg_get_serial_sequence('contabilidad.ref_estado_año_fiscal', 'id'),
    (SELECT MAX(id) FROM contabilidad.ref_estado_año_fiscal)
);

-- ── 3. Tabla principal: años fiscales por empresa ─────────────────────────────
CREATE TABLE IF NOT EXISTS contabilidad.cfg_años_fiscales (
    id          BIGSERIAL   NOT NULL,
    year        INTEGER     NOT NULL,
    id_estado   INTEGER     NOT NULL,
    empresa_id  BIGINT      NOT NULL,
    created_by  BIGINT      NOT NULL,
    updated_by  BIGINT      NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_cfg_años_fiscales
        PRIMARY KEY (id),
 
    -- Un año es único por empresa
    CONSTRAINT uq_cfg_años_fiscales_year_empresa
        UNIQUE (year, empresa_id),
 
    -- El estado debe existir en el catálogo
    CONSTRAINT fk_cfg_años_fiscales_estado
        FOREIGN KEY (id_estado)
        REFERENCES contabilidad.ref_estado_año_fiscal (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
 
    -- El año debe ser un valor razonable
    CONSTRAINT chk_cfg_años_fiscales_year
        CHECK (year BETWEEN 1900 AND 2200)
);

COMMENT ON TABLE contabilidad.cfg_años_fiscales IS
    'Años fiscales (vigencias contables) por empresa. El estado se gestiona mediante FK a ref_estado_año_fiscal.';
COMMENT ON COLUMN contabilidad.cfg_años_fiscales.year IS
    'Año calendario de la vigencia (ej: 2025).';
COMMENT ON COLUMN contabilidad.cfg_años_fiscales.id_estado IS
    'FK a contabilidad.ref_estado_año_fiscal. No almacenar texto de estado directamente.';
COMMENT ON COLUMN contabilidad.cfg_años_fiscales.empresa_id IS
    'Empresa propietaria de la vigencia.';

-- ── 4. Índices de rendimiento ─────────────────────────────────────────────────

-- Consultas por empresa (filtro más frecuente)
CREATE INDEX IF NOT EXISTS idx_cfg_años_fiscales_empresa
    ON contabilidad.cfg_años_fiscales (empresa_id);
 
-- Consultas por empresa + estado (buscar el año abierto de una empresa)
CREATE INDEX IF NOT EXISTS idx_cfg_años_fiscales_empresa_estado
    ON contabilidad.cfg_años_fiscales (empresa_id, id_estado);
 
-- ── 5. Constraint avanzado: máximo UN año fiscal ABIERTO por empresa ──────────
--
-- Índice único parcial de PostgreSQL: segunda línea de defensa a nivel de motor
-- además de la validación que implementa la capa de servicios.
-- id_estado = 1 → 'Abierto' (semilla fija)
 
CREATE UNIQUE INDEX IF NOT EXISTS uq_cfg_años_fiscales_un_abierto_por_empresa
    ON contabilidad.cfg_años_fiscales (empresa_id)
    WHERE id_estado = 1;
 
COMMENT ON INDEX contabilidad.uq_cfg_años_fiscales_un_abierto_por_empresa IS
    'Garantiza a nivel de BD que solo exista UN año fiscal Abierto por empresa. Segunda línea de defensa tras la validación del servicio.';
 
COMMIT;