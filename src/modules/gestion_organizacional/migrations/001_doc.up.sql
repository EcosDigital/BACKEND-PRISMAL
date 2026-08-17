-- ============================================================================
-- MIGRACIÓN: Módulo SOP (Standard Operating Procedures)
-- Proyecto : PRISMAR ERP
-- Esquema  : documentacion
-- Tablas   :
--   documentacion.ref_categoria_sop   → catálogo de categorías de SOP
--   documentacion.ref_estado_sop      → catálogo de estados de SOP
--   documentacion.cfg_sop             → tabla principal de procedimientos
-- ----------------------------------------------------------------------------
-- Patrón: estados y categorías mediante FK a tablas ref_*, nunca VARCHAR directo.
-- Soft delete mediante columna is_active BOOLEAN.
-- Idempotente: usa IF NOT EXISTS y ON CONFLICT DO NOTHING.
-- Ejecutar en una sola transacción.
-- ============================================================================
 
 BEGIN;

 -- ── 1. Crear esquema si no existe ─────────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS documentacion;

COMMENT ON SCHEMA documentacion IS
    'Esquema del módulo de Documentación de PRISMAR ERP. Centraliza SOPs, manuales y políticas.';

-- ── 2. Catálogo de categorías para un SOP ────────────────────────────────────
--
-- Patrón ref_* del sistema PRISMAR.
-- La categoría NUNCA se almacena como texto en la tabla principal.

CREATE TABLE IF NOT EXISTS documentacion.ref_categoria_sop (
    id          SERIAL      NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    descripcion TEXT,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
 
    CONSTRAINT pk_ref_categoria_sop
        PRIMARY KEY (id),
 
    CONSTRAINT uq_ref_categoria_sop_nombre
        UNIQUE (nombre)
);

COMMENT ON TABLE documentacion.ref_categoria_sop IS
    'Catálogo de categorías posibles para un SOP. Patrón ref_* del sistema PRISMAR.';

INSERT INTO documentacion.ref_categoria_sop (id, nombre, descripcion, is_active)
VALUES
    (1, 'Operación',       'Procedimientos relacionados con operaciones del negocio.',      TRUE),
    (2, 'Soporte',         'Procedimientos de atención y soporte a usuarios/clientes.',     TRUE),
    (3, 'Tecnología',      'Procedimientos técnicos de sistemas e infraestructura.',        TRUE),
    (4, 'Pagos',           'Procedimientos relacionados con gestión de pagos y cobros.',    TRUE),
    (5, 'Comercial',       'Procedimientos del área comercial y ventas.',                   TRUE),
    (6, 'Marketing',       'Procedimientos del área de marketing y comunicaciones.',        TRUE),
    (7, 'Administración',  'Procedimientos administrativos y de gestión interna.',          TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('documentacion.ref_categoria_sop', 'id'),
    (SELECT MAX(id) FROM documentacion.ref_categoria_sop)
);

-- ── 3. Catálogo de estados posibles para un SOP ───────────────────────────────
--
-- Patrón ref_* del sistema PRISMAR.
-- El estado NUNCA se almacena como texto en la tabla principal.

CREATE TABLE IF NOT EXISTS documentacion.ref_estado_sop (
    id          SERIAL      NOT NULL,
    nombre      VARCHAR(50) NOT NULL,
    descripcion TEXT,
    color       VARCHAR(20),    -- sugerencia de color semántico para el frontend
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
 
    CONSTRAINT pk_ref_estado_sop
        PRIMARY KEY (id),
 
    CONSTRAINT uq_ref_estado_sop_nombre
        UNIQUE (nombre)
);

COMMENT ON TABLE documentacion.ref_estado_sop IS
    'Catálogo de estados del ciclo de vida de un SOP. Patrón ref_* del sistema PRISMAR.';

COMMENT ON COLUMN documentacion.ref_estado_sop.color IS
    'Sugerencia semántica de color para el frontend (ej: amber, emerald, blue, gray).';

INSERT INTO documentacion.ref_estado_sop (id, nombre, descripcion, color, is_active)
VALUES
    (1, 'Draft',     'Borrador. El SOP está siendo redactado y no es oficial.',         'amber',   TRUE),
    (2, 'Review',    'En revisión. El SOP está siendo evaluado por los responsables.',  'blue',    TRUE),
    (3, 'Approved',  'Aprobado. El SOP es oficial y está vigente.',                     'emerald', TRUE),
    (4, 'Obsolete',  'Obsoleto. El SOP ya no aplica y fue reemplazado o descontinuado.','gray',    TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('documentacion.ref_estado_sop', 'id'),
    (SELECT MAX(id) FROM documentacion.ref_estado_sop)
);

-- ── 4. Tabla principal: SOPs por empresa ──────────────────────────────────────
--
-- Soft delete mediante is_active = FALSE.
-- El código del SOP es único por empresa (ej: SOP-OPE-001).
-- La versión es un string libre para soportar esquemas como "1.0", "2.1.3".

CREATE TABLE IF NOT EXISTS documentacion.cfg_sop (
    id              BIGSERIAL       NOT NULL,
    code            VARCHAR(50)     NOT NULL,
    title           VARCHAR(500)    NOT NULL,
    id_categoria    INTEGER         NOT NULL,
    area            VARCHAR(200)    NOT NULL,
    id_estado       INTEGER         NOT NULL,
    content         TEXT            NOT NULL,
    version         VARCHAR(20)     NOT NULL DEFAULT '1.0',
    effective_date  DATE            NOT NULL,
    review_date     DATE,
    is_active       BOOLEAN         NOT NULL DEFAULT TRUE,
    empresa_id      BIGINT          NOT NULL,
    created_by      BIGINT          NOT NULL,
    updated_by      BIGINT          NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_cfg_sop
        PRIMARY KEY (id),
 
    -- El código del SOP es único por empresa (soft delete excluido)
    CONSTRAINT uq_cfg_sop_code_empresa
        UNIQUE (empresa_id, code),
 
    -- FK al catálogo de categorías
    CONSTRAINT fk_cfg_sop_categoria
        FOREIGN KEY (id_categoria)
        REFERENCES documentacion.ref_categoria_sop (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
 
    -- FK al catálogo de estados
    CONSTRAINT fk_cfg_sop_estado
        FOREIGN KEY (id_estado)
        REFERENCES documentacion.ref_estado_sop (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
 
    -- Coherencia de fechas: la fecha de revisión debe ser posterior a la efectiva
    CONSTRAINT chk_cfg_sop_fechas
        CHECK (review_date IS NULL OR review_date >= effective_date)
);

COMMENT ON TABLE documentacion.cfg_sop IS
    'Procedimientos Operativos Estándar (SOP) de la organización. Soft delete mediante is_active.';

COMMENT ON COLUMN documentacion.cfg_sop.code IS
    'Código único del SOP por empresa (ej: SOP-OPE-001). Definido por el usuario.';

COMMENT ON COLUMN documentacion.cfg_sop.version IS
    'Versión del SOP en formato libre (ej: 1.0, 2.1.3). Se incrementa manualmente al actualizar.';

COMMENT ON COLUMN documentacion.cfg_sop.id_categoria IS
    'FK a documentacion.ref_categoria_sop. No almacenar texto de categoría directamente.';

COMMENT ON COLUMN documentacion.cfg_sop.id_estado IS
    'FK a documentacion.ref_estado_sop. No almacenar texto de estado directamente.';

COMMENT ON COLUMN documentacion.cfg_sop.is_active IS
    'Soft delete: FALSE indica que el SOP fue eliminado lógicamente.';

-- ── 5. Índices de rendimiento ─────────────────────────────────────────────────

-- Filtro por empresa (el más frecuente)
CREATE INDEX IF NOT EXISTS idx_cfg_sop_empresa
    ON documentacion.cfg_sop (empresa_id)
    WHERE is_active = TRUE;

-- Filtro por empresa + estado (listar por estado)
CREATE INDEX IF NOT EXISTS idx_cfg_sop_empresa_estado
    ON documentacion.cfg_sop (empresa_id, id_estado)
    WHERE is_active = TRUE;

-- Filtro por empresa + categoría
CREATE INDEX IF NOT EXISTS idx_cfg_sop_empresa_categoria
    ON documentacion.cfg_sop (empresa_id, id_categoria)
    WHERE is_active = TRUE;

-- Búsqueda de texto en título (ILIKE queries del listado)
CREATE INDEX IF NOT EXISTS idx_cfg_sop_title_gin
    ON documentacion.cfg_sop USING gin (to_tsvector('spanish', title))
    WHERE is_active = TRUE;

-- ── 6. Vista de consulta: SOPs con catálogos resueltos ───────────────────────
--
-- Facilita queries del frontend y reportes sin JOINs repetidos.

 
CREATE OR REPLACE VIEW documentacion.v_sop AS
SELECT
    s.id,
    s.code,
    s.title,
    s.id_categoria,
    c.nombre        AS categoria,
    s.area,
    s.id_estado,
    e.nombre        AS estado,
    e.color         AS estado_color,
    s.content,
    s.version,
    s.effective_date,
    s.review_date,
    s.is_active,
    s.empresa_id,
    s.created_by,
    s.updated_by,
    s.created_at,
    s.updated_at
FROM documentacion.cfg_sop s
INNER JOIN documentacion.ref_categoria_sop c ON c.id = s.id_categoria
INNER JOIN documentacion.ref_estado_sop    e ON e.id = s.id_estado;
 
COMMENT ON VIEW documentacion.v_sop IS
    'Vista de SOPs con categoría y estado resueltos por JOIN. Facilita consultas del frontend.';
 
COMMIT;