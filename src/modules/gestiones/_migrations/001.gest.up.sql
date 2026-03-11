-- ============================================================
-- PRISMAR ERP - MÓDULO GESTIONES
-- Esquema: gestiones
-- Base de datos: PostgreSQL
-- Versión: 1.0
-- ============================================================

CREATE SCHEMA IF NOT EXISTS gestiones;

-- ----------------------------------------------------------
-- REFERENCIAL: Área funcional del grupo
-- ----------------------------------------------------------

CREATE TABLE gestiones.cfg_tipos_grupo (
    id     SMALLINT    PRIMARY KEY,
    codigo VARCHAR(30) NOT NULL UNIQUE,
    nombre VARCHAR(80) NOT NULL
);

INSERT INTO gestiones.cfg_tipos_grupo (id, codigo, nombre)
VALUES
    (1, '001',        'Soporte'),
    (2, '002',     'Desarrollo'),
    (3, '003',   'Contabilidad'),
    (4, '004', 'Administrativo'),
    (5, '005',    'Operaciones')
ON CONFLICT (id) DO NOTHING;

-- ----------------------------------------------------------
-- CONFIGURACIÓN: Grupos de trabajo
-- ----------------------------------------------------------

CREATE TABLE gestiones.cfg_grupos_trabajo (
    id        SERIAL       PRIMARY KEY,
    codigo    VARCHAR(30)  NOT NULL,
    nombre    VARCHAR(100) NOT NULL,
    id_tercero_jefe INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    id_tipo   SMALLINT  NOT NULL REFERENCES gestiones.cfg_tipos_grupo(id),
    is_active BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by INT         NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT         NULL,

    CONSTRAINT uq_cfg_grupos_trabajo_codigo UNIQUE (codigo),
    CONSTRAINT uq_cfg_grupos_trabajo_nombre UNIQUE (nombre)
);

COMMENT ON TABLE  gestiones.cfg_grupos_trabajo          IS 'Grupos o equipos de trabajo que atienden gestiones';
COMMENT ON COLUMN gestiones.cfg_grupos_trabajo.id_tipo  IS 'FK → ref_tipos_grupo (área funcional del grupo)';
COMMENT ON COLUMN gestiones.cfg_grupos_trabajo.is_active IS 'TRUE = activo, FALSE = inactivo';
COMMENT ON COLUMN gestiones.cfg_grupos_trabajo.id_tercero_jefe IS 'FK → cfg_colaboradores. Colaborador designado como jefe del grupo';

-- ==========================================================
-- MÓDULO GESTIONES — Colaboradores por Grupo de Trabajo
-- ==========================================================

-- Un mismo tercero puede pertenecer a varios grupos.
-- La clave primaria compuesta garantiza que no se duplique
-- la misma combinación tercero + grupo.

CREATE TABLE gestiones.cfg_colaboradores (
    id SERIAL PRIMARY KEY,
    id_tercero INT NOT NULL REFERENCES configuracion.cfg_terceros(id),  
    id_grupo   INT NOT NULL REFERENCES gestiones.cfg_grupos_trabajo(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by INT         NOT NULL
);

COMMENT ON TABLE  gestiones.cfg_colaboradores            IS 'Relación entre terceros y grupos de trabajo. Un tercero puede pertenecer a varios grupos.';
COMMENT ON COLUMN gestiones.cfg_colaboradores.id_tercero IS 'FK → configuracion.cfg_terceros';
COMMENT ON COLUMN gestiones.cfg_colaboradores.id_grupo   IS 'FK → gestiones.cfg_grupos_trabajo';
COMMENT ON COLUMN gestiones.cfg_colaboradores.created_at IS 'Fecha y hora en que se asignó el tercero al grupo';
COMMENT ON COLUMN gestiones.cfg_colaboradores.created_by IS 'ID del usuario que registró la asignación';


-- ----------------------------------------------------------
-- CONFIGURACIÓN: Estados del ticket
-- ----------------------------------------------------------

CREATE TABLE gestiones.cfg_estados_ticket (
    id          SERIAL      PRIMARY KEY,
    codigo      VARCHAR(30) NOT NULL UNIQUE,
    nombre      VARCHAR(80) NOT NULL,
    descripcion TEXT   NULL,
    -- Apariencia en UI
    color_hex   CHAR(10)  NOT NULL DEFAULT '#6B7280', 
    -- Comportamiento
    is_gestionable BOOLEAN     NOT NULL DEFAULT TRUE,
    is_encuesta BOOLEAN     NOT NULL DEFAULT FALSE,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    orden       SMALLINT    NOT NULL DEFAULT 0, 
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  INT         NOT NULL,
    updated_at  TIMESTAMPTZ NULL,
    updated_by  INT         NULL
);

COMMENT ON TABLE  gestiones.cfg_estados_ticket                IS 'Estados parametrizables del ciclo de vida de un ticket';
COMMENT ON COLUMN gestiones.cfg_estados_ticket.color_hex      IS 'Color hexadecimal para badges y vistas kanban. Ej: #22C55E';
COMMENT ON COLUMN gestiones.cfg_estados_ticket.is_gestionable IS 'TRUE = el cliente, colaborador o jefe de grupo puede anexar comentarios al caso';
COMMENT ON COLUMN gestiones.cfg_estados_ticket.is_encuesta    IS 'TRUE = habilita la encuesta de satisfacción para el cliente (típicamente en estado RESUELTO)';
COMMENT ON COLUMN gestiones.cfg_estados_ticket.is_active      IS 'FALSE = estado deshabilitado, no aparece en el flujo';
COMMENT ON COLUMN gestiones.cfg_estados_ticket.orden          IS 'Posición visual en el flujo o tablero kanban';


-- ----------------------------------------------------------
-- DATOS INICIALES
-- ----------------------------------------------------------

INSERT INTO gestiones.cfg_estados_ticket
    (codigo, nombre, descripcion, color_hex, is_gestionable, is_encuesta, is_active, orden, created_by)
VALUES
    ('001',    'Propuesto',       'Ticket recién creado, sin asignar',                                                   '#3B82F6', TRUE,  FALSE, TRUE, 1, 1),
    ('002',    'Asignado',        'Asignado a un colaborador, pendiente de inicio',                                      '#8B5CF6', TRUE,  FALSE, TRUE, 2, 1),
    ('003',    'En Progreso',     'Colaborador trabajando activamente en el caso',                                       '#F59E0B', TRUE,  FALSE, TRUE, 3, 1),
    ('004',    'Resuelto',        'Solución entregada — el cliente puede llenar la encuesta y confirmar o devolver',     '#22C55E', TRUE,  TRUE,  TRUE, 4, 1),
    ('005',    'En Espera',       'Pausado, esperando información o respuesta del cliente',                              '#6B7280', TRUE,  FALSE, TRUE, 5, 1),
    ('006',    'Cerrado',     'Confirmado y cerrado definitivamente',                                               '#15803D', FALSE, FALSE, TRUE, 6, 1),
    ('007',    'Devuelto',    'Anulado sin resolución',                                                             '#EF4444', FALSE, FALSE, TRUE, 7, 1)
ON CONFLICT (codigo) DO NOTHING;