-- ============================================================================
-- MIGRACIÓN: Módulo Gestión Organizacional — Proyectos
-- Proyecto : PRISMAR ERP
-- Esquema  : organizacion
-- Tablas   :
--   organizacion.ref_estado_proyecto     → catálogo de estados
--   organizacion.ref_prioridad_proyecto  → catálogo de prioridades
--   organizacion.cfg_proyectos           → tabla principal de proyectos
--   organizacion.cfg_partes_interesadas  → stakeholders por proyecto
--   organizacion.bit_comentarios         → bitácora cronológica
--   organizacion.bit_evidencias          → evidencias por comentario
-- ----------------------------------------------------------------------------
-- Convención: todos los nombres de columna en español.
-- Soft delete: columna "es_activo" BOOLEAN en proyectos y partes interesadas.
-- Bitácora  : los comentarios NO tienen soft delete (historial inmutable).
-- Idempotente: IF NOT EXISTS y ON CONFLICT DO NOTHING en todos los objetos.
-- Ejecutar en una sola transacción.
-- ============================================================================
 

 BEGIN;
 
-- ── 1. Esquema ────────────────────────────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS organizacion;

COMMENT ON SCHEMA organizacion IS
    'Esquema del módulo de Gestión Organizacional de PRISMAR ERP.';

-- ── 2. Catálogo: estados de proyecto ──────────────────────────────────────────
CREATE TABLE IF NOT EXISTS organizacion.ref_estado_proyecto (
    id          SERIAL       NOT NULL,
    nombre      VARCHAR(50)  NOT NULL,
    descripcion TEXT,
    color       VARCHAR(20),
    es_activo   BOOLEAN      NOT NULL DEFAULT TRUE,
 
    CONSTRAINT pk_ref_estado_proyecto        PRIMARY KEY (id),
    CONSTRAINT uq_ref_estado_proyecto_nombre UNIQUE (nombre)
);

COMMENT ON TABLE  organizacion.ref_estado_proyecto          IS 'Catálogo de estados posibles para un proyecto. Patrón ref_* de PRISMAR.';
COMMENT ON COLUMN organizacion.ref_estado_proyecto.color    IS 'Sugerencia semántica de color para el frontend (ej: blue, emerald, amber, red).';
COMMENT ON COLUMN organizacion.ref_estado_proyecto.es_activo IS 'Indica si el estado está disponible para ser asignado.';

INSERT INTO organizacion.ref_estado_proyecto (id, nombre, descripcion, color, es_activo)
VALUES
    (1, 'Planeación',  'El proyecto está siendo planificado, aún no iniciado.',  'blue',    TRUE),
    (2, 'En Progreso', 'El proyecto está activo y en ejecución.',                'emerald', TRUE),
    (3, 'En Pausa',    'El proyecto fue pausado temporalmente.',                 'amber',   TRUE),
    (4, 'Finalizado',  'El proyecto concluyó exitosamente.',                     'indigo',  TRUE),
    (5, 'Cancelado',   'El proyecto fue cancelado y no continuará.',             'red',     TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('organizacion.ref_estado_proyecto', 'id'),
    (SELECT MAX(id) FROM organizacion.ref_estado_proyecto)
);

-- ── 3. Catálogo: prioridades de proyecto ──────────────────────────────────────
CREATE TABLE IF NOT EXISTS organizacion.ref_prioridad_proyecto (
    id          SERIAL       NOT NULL,
    nombre      VARCHAR(50)  NOT NULL,
    descripcion TEXT,
    color       VARCHAR(20),
    es_activo   BOOLEAN      NOT NULL DEFAULT TRUE,
 
    CONSTRAINT pk_ref_prioridad_proyecto        PRIMARY KEY (id),
    CONSTRAINT uq_ref_prioridad_proyecto_nombre UNIQUE (nombre)
);

COMMENT ON TABLE  organizacion.ref_prioridad_proyecto           IS 'Catálogo de prioridades para un proyecto. Patrón ref_* de PRISMAR.';
COMMENT ON COLUMN organizacion.ref_prioridad_proyecto.es_activo IS 'Indica si la prioridad está disponible para ser asignada.';

INSERT INTO organizacion.ref_prioridad_proyecto (id, nombre, descripcion, color, es_activo)
VALUES
    (1, 'Baja',    'El proyecto puede postergarse sin impacto significativo.',   'gray',  TRUE),
    (2, 'Media',   'El proyecto tiene importancia moderada.',                    'blue',  TRUE),
    (3, 'Alta',    'El proyecto es prioritario y requiere atención pronta.',     'amber', TRUE),
    (4, 'Crítica', 'El proyecto es urgente y estratégico para la organización.', 'red',   TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('organizacion.ref_prioridad_proyecto', 'id'),
    (SELECT MAX(id) FROM organizacion.ref_prioridad_proyecto)
);

-- ── 4. Tabla principal: proyectos ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS organizacion.cfg_proyectos (
    id                 BIGSERIAL    NOT NULL,
    codigo             VARCHAR(50)  NOT NULL,
    nombre             VARCHAR(500) NOT NULL,
    descripcion        TEXT,
    fecha_inicio       DATE         NOT NULL,
    fecha_fin          DATE,
    id_estado          INTEGER      NOT NULL,
    id_prioridad       INTEGER      NOT NULL,
    id_responsable     BIGINT       NOT NULL,
    es_activo          BOOLEAN      NOT NULL DEFAULT TRUE,
    id_empresa         BIGINT       NOT NULL,
    created_by          BIGINT       NOT NULL,
    updated_by         BIGINT       NOT NULL,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_cfg_proyectos
        PRIMARY KEY (id),
 
    -- El código del proyecto es único por empresa
    CONSTRAINT uq_cfg_proyectos_codigo_empresa
        UNIQUE (id_empresa, codigo),
 
    CONSTRAINT fk_cfg_proyectos_estado
        FOREIGN KEY (id_estado)
        REFERENCES organizacion.ref_estado_proyecto (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
 
    CONSTRAINT fk_cfg_proyectos_prioridad
        FOREIGN KEY (id_prioridad)
        REFERENCES organizacion.ref_prioridad_proyecto (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
 
    -- La fecha de fin, si existe, no puede ser anterior a la fecha de inicio
    CONSTRAINT chk_cfg_proyectos_fechas
        CHECK (fecha_fin IS NULL OR fecha_fin >= fecha_inicio)
);

COMMENT ON TABLE  organizacion.cfg_proyectos                IS 'Proyectos internos de la organización. Soft delete mediante es_activo.';
COMMENT ON COLUMN organizacion.cfg_proyectos.codigo         IS 'Código único del proyecto por empresa (ej: PRY-2026-001).';
COMMENT ON COLUMN organizacion.cfg_proyectos.id_responsable IS 'ID del usuario responsable del proyecto.';
COMMENT ON COLUMN organizacion.cfg_proyectos.es_activo      IS 'Soft delete: FALSE indica que el proyecto fue eliminado lógicamente.';
 
 -- Índices de rendimiento
CREATE INDEX IF NOT EXISTS idx_cfg_proyectos_empresa
    ON organizacion.cfg_proyectos (id_empresa) WHERE es_activo = TRUE;
 
CREATE INDEX IF NOT EXISTS idx_cfg_proyectos_empresa_estado
    ON organizacion.cfg_proyectos (id_empresa, id_estado) WHERE es_activo = TRUE;
 
CREATE INDEX IF NOT EXISTS idx_cfg_proyectos_empresa_prioridad
    ON organizacion.cfg_proyectos (id_empresa, id_prioridad) WHERE es_activo = TRUE;
 
CREATE INDEX IF NOT EXISTS idx_cfg_proyectos_responsable
    ON organizacion.cfg_proyectos (id_empresa, id_responsable) WHERE es_activo = TRUE;


-- ── 5. Tabla: partes interesadas por proyecto ─────────────────────────────────
CREATE TABLE IF NOT EXISTS organizacion.cfg_partes_interesadas (
    id              BIGSERIAL    NOT NULL,
    id_proyecto     BIGINT       NOT NULL,
    nombre          VARCHAR(300) NOT NULL,
    rol             VARCHAR(200),
    organizacion    VARCHAR(300),
    correo          VARCHAR(200),
    telefono        VARCHAR(50),
    observaciones   TEXT,
    es_activo       BOOLEAN      NOT NULL DEFAULT TRUE,
    id_empresa      BIGINT       NOT NULL,
    creado_por      BIGINT       NOT NULL,
    actualizado_por BIGINT       NOT NULL,
    creado_en       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    actualizado_en  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_cfg_partes_interesadas
        PRIMARY KEY (id),
 
    CONSTRAINT fk_cfg_partes_interesadas_proyecto
        FOREIGN KEY (id_proyecto)
        REFERENCES organizacion.cfg_proyectos (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

 
COMMENT ON TABLE  organizacion.cfg_partes_interesadas          IS 'Partes interesadas asociadas a un proyecto. Soft delete mediante es_activo.';
COMMENT ON COLUMN organizacion.cfg_partes_interesadas.es_activo IS 'Soft delete: FALSE indica que la parte interesada fue eliminada lógicamente.';
 
CREATE INDEX IF NOT EXISTS idx_cfg_partes_interesadas_proyecto
    ON organizacion.cfg_partes_interesadas (id_proyecto) WHERE es_activo = TRUE;
 
CREATE INDEX IF NOT EXISTS idx_cfg_partes_interesadas_empresa
    ON organizacion.cfg_partes_interesadas (id_empresa) WHERE es_activo = TRUE;
 
 -- ── 6. Tabla: bitácora de comentarios ─────────────────────────────────────────
--
-- Registro histórico inmutable: NO tiene soft delete.
-- Una vez creado, un comentario no puede eliminarse para preservar la trazabilidad.

CREATE TABLE IF NOT EXISTS organizacion.bit_comentarios (
    id              BIGSERIAL   NOT NULL,
    id_proyecto     BIGINT      NOT NULL,
    fecha_actividad DATE        NOT NULL,
    comentario      TEXT        NOT NULL,
    id_empresa      BIGINT      NOT NULL,
    created_by          BIGINT       NOT NULL,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_bit_comentarios
        PRIMARY KEY (id),
 
    CONSTRAINT fk_bit_comentarios_proyecto
        FOREIGN KEY (id_proyecto)
        REFERENCES organizacion.cfg_proyectos (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

COMMENT ON TABLE  organizacion.bit_comentarios                  IS 'Bitácora cronológica de eventos por proyecto. Sin soft delete: historial inmutable.';
COMMENT ON COLUMN organizacion.bit_comentarios.fecha_actividad  IS 'Fecha real del evento (puede diferir de creado_en si se registra con retraso).';
COMMENT ON COLUMN organizacion.bit_comentarios.comentario       IS 'Descripción del evento o actividad registrada.';

CREATE INDEX IF NOT EXISTS idx_bit_comentarios_proyecto
    ON organizacion.bit_comentarios (id_proyecto);
 
CREATE INDEX IF NOT EXISTS idx_bit_comentarios_fecha
    ON organizacion.bit_comentarios (id_proyecto, fecha_actividad DESC);

-- ── 7. Catálogo: tipos de evidencia ───────────────────────────────────────────
CREATE TABLE IF NOT EXISTS organizacion.ref_tipo_evidencia (
    id        SERIAL      NOT NULL,
    nombre    VARCHAR(50) NOT NULL,
    es_activo BOOLEAN     NOT NULL DEFAULT TRUE,
 
    CONSTRAINT pk_ref_tipo_evidencia        PRIMARY KEY (id),
    CONSTRAINT uq_ref_tipo_evidencia_nombre UNIQUE (nombre)
);

COMMENT ON TABLE organizacion.ref_tipo_evidencia IS 'Catálogo de tipos de archivo permitidos como evidencia. Patrón ref_* de PRISMAR.';

INSERT INTO organizacion.ref_tipo_evidencia (id, nombre, es_activo)
VALUES
    (1, 'Imagen',          TRUE),
    (2, 'Audio',           TRUE),
    (3, 'Video',           TRUE),
    (4, 'Documento',       TRUE),
    (5, 'PDF',             TRUE),
    (6, 'Archivo General', TRUE)
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('organizacion.ref_tipo_evidencia', 'id'),
    (SELECT MAX(id) FROM organizacion.ref_tipo_evidencia)
);

-- ── 8. Tabla: evidencias por comentario ───────────────────────────────────────
--
-- Las evidencias siguen el ciclo de vida del comentario padre (CASCADE).
-- No tienen soft delete: si se elimina el comentario, sus evidencias desaparecen.

CREATE TABLE IF NOT EXISTS organizacion.bit_evidencias (
    id               BIGSERIAL    NOT NULL,
    id_comentario    BIGINT       NOT NULL,
    nombre_archivo   VARCHAR(500) NOT NULL,
    ruta_archivo     TEXT         NOT NULL,
    id_tipo          INTEGER      NOT NULL,
    id_empresa       BIGINT       NOT NULL,
    created_by       BIGINT       NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
 
    CONSTRAINT pk_bit_evidencias
        PRIMARY KEY (id),
 
    -- Al eliminar un comentario, sus evidencias se eliminan automáticamente
    CONSTRAINT fk_bit_evidencias_comentario
        FOREIGN KEY (id_comentario)
        REFERENCES organizacion.bit_comentarios (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
 
    CONSTRAINT fk_bit_evidencias_tipo
        FOREIGN KEY (id_tipo)
        REFERENCES organizacion.ref_tipo_evidencia (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

COMMENT ON TABLE  organizacion.bit_evidencias               IS 'Archivos adjuntos a comentarios de la bitácora. ON DELETE CASCADE hereda el ciclo del comentario padre.';
COMMENT ON COLUMN organizacion.bit_evidencias.ruta_archivo  IS 'Ruta relativa o URL del archivo en el sistema de almacenamiento.';
COMMENT ON COLUMN organizacion.bit_evidencias.id_tipo       IS 'FK a organizacion.ref_tipo_evidencia.';
 
CREATE INDEX IF NOT EXISTS idx_bit_evidencias_comentario
    ON organizacion.bit_evidencias (id_comentario);

-- ── 9. Vista: proyectos con catálogos resueltos ───────────────────────────────
--
-- Facilita consultas del frontend sin repetir JOINs en cada query.
CREATE OR REPLACE VIEW organizacion.v_proyectos AS
SELECT
    p.id,
    p.codigo,
    p.nombre,
    p.descripcion,
    p.fecha_inicio,
    p.fecha_fin,
    p.id_estado,
    ep.nombre       AS estado,
    ep.color        AS color_estado,
    p.id_prioridad,
    pr.nombre       AS prioridad,
    pr.color        AS color_prioridad,
    p.id_responsable,
    p.es_activo,
    p.id_empresa,
    p.created_by,
    p.updated_by,
    p.created_at,
    p.updated_at
FROM organizacion.cfg_proyectos p
INNER JOIN organizacion.ref_estado_proyecto    ep ON ep.id = p.id_estado
INNER JOIN organizacion.ref_prioridad_proyecto pr ON pr.id = p.id_prioridad;
 
COMMENT ON VIEW organizacion.v_proyectos IS 'Vista de proyectos con estado y prioridad resueltos por JOIN. Facilita consultas del frontend.';
 
COMMIT;