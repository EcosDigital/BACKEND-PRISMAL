-- ============================================================
-- Migration: report_server_connections
-- Schema:    administracion
-- Proyecto:  PRISMAR
-- Descripción: Almacena las conexiones a servidores de informes
--              configuradas por los usuarios de la plataforma.
-- Idempotente: puede ejecutarse múltiples veces sin error.
-- ============================================================

-- Crear schema si no existe
CREATE SCHEMA IF NOT EXISTS administracion;

-- Tabla de conexiones a servidores de informes
CREATE TABLE IF NOT EXISTS administracion.report_server_connections (
    id          SERIAL          PRIMARY KEY,
    name        TEXT            NOT NULL,
    base_url    TEXT            NOT NULL,
    api_key     TEXT            NOT NULL,
    is_active   BOOLEAN         NOT NULL DEFAULT FALSE,
    created_by  TEXT            NOT NULL,           -- usuario que creó la conexión
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_rsc_name_user UNIQUE (name, created_by)
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_rsc_created_by
    ON administracion.report_server_connections(created_by);

CREATE INDEX IF NOT EXISTS idx_rsc_is_active
    ON administracion.report_server_connections(is_active, created_by);

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION administracion.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_rsc_updated_at
    ON administracion.report_server_connections;

CREATE TRIGGER trg_rsc_updated_at
    BEFORE UPDATE ON administracion.report_server_connections
    FOR EACH ROW EXECUTE FUNCTION administracion.set_updated_at();

-- Garantiza que solo una conexión esté activa por usuario a la vez
CREATE OR REPLACE FUNCTION administracion.enforce_single_active_connection()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_active = TRUE THEN
        UPDATE administracion.report_server_connections
        SET    is_active = FALSE
        WHERE  created_by = NEW.created_by
          AND  id <> NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_rsc_single_active
    ON administracion.report_server_connections;

CREATE TRIGGER trg_rsc_single_active
    AFTER INSERT OR UPDATE OF is_active
    ON administracion.report_server_connections
    FOR EACH ROW
    WHEN (NEW.is_active = TRUE)
    EXECUTE FUNCTION administracion.enforce_single_active_connection();

-- ============================================================
-- Comentarios de documentación
-- ============================================================

COMMENT ON SCHEMA administracion
    IS 'Configuración y administración general de la plataforma PRISMAR';

COMMENT ON TABLE administracion.report_server_connections
    IS 'Conexiones a servidores de informes (Report Server) por usuario';

COMMENT ON COLUMN administracion.report_server_connections.name
    IS 'Nombre descriptivo de la conexión, único por usuario';

COMMENT ON COLUMN administracion.report_server_connections.base_url
    IS 'URL base del servidor de informes (ej: http://192.168.1.10:8080)';

COMMENT ON COLUMN administracion.report_server_connections.api_key
    IS 'API Key de acceso al servidor de informes';

COMMENT ON COLUMN administracion.report_server_connections.is_active
    IS 'Indica si esta es la conexión activa del usuario. Solo una puede estar activa por usuario.';

COMMENT ON COLUMN administracion.report_server_connections.created_by
    IS 'Identificador del usuario propietario de la conexión';