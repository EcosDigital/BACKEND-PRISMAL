-- ============================================================
-- MIGRACIÓN: 001_CA_UP.sql
-- MÓDULO: Control de Asistencia (CA)
-- SISTEMA: Prismar ERP
-- DESCRIPCIÓN: Creación de tablas CFG y REF para gestión de
--              puntos de marcación mediante códigos QR
-- AUTOR: Andres Pastrana - Ingeniero de Software
-- FECHA: 2026-03-16
-- ============================================================

-- ============================================================
-- ESQUEMA
-- ============================================================

CREATE SCHEMA IF NOT EXISTS control_asistencia;

-- ============================================================
-- TABLAS REFERENCIALES
-- ============================================================

-- ------------------------------------------------------------
-- REF: Tipos de punto de marcación
-- Define el propósito o uso de cada punto de marcación QR
-- ------------------------------------------------------------

CREATE TABLE IF NOT EXISTS control_asistencia.ref_tipo_punto_marcacion (
    id          SERIAL          PRIMARY KEY,
    codigo      VARCHAR(60)     NOT NULL,
    nombre      VARCHAR(150)    NOT NULL,
    descripcion TEXT,
    activo      BOOLEAN         NOT NULL DEFAULT TRUE,
    CONSTRAINT uq_ref_tipo_punto_marcacion_codigo UNIQUE (codigo)
);

COMMENT ON TABLE  control_asistencia.ref_tipo_punto_marcacion              IS 'Tipos de punto de marcación QR disponibles en el sistema';
COMMENT ON COLUMN control_asistencia.ref_tipo_punto_marcacion.id           IS 'Identificador único del tipo de punto';
COMMENT ON COLUMN control_asistencia.ref_tipo_punto_marcacion.codigo       IS 'Código único interno del tipo de punto';
COMMENT ON COLUMN control_asistencia.ref_tipo_punto_marcacion.nombre       IS 'Nombre descriptivo del tipo de punto';
COMMENT ON COLUMN control_asistencia.ref_tipo_punto_marcacion.descripcion  IS 'Descripción detallada del uso del tipo de punto';
COMMENT ON COLUMN control_asistencia.ref_tipo_punto_marcacion.activo       IS 'Indica si el tipo de punto está habilitado en el sistema';
 
 -- ------------------------------------------------------------
-- Datos iniciales: tipos de punto de marcación
-- ------------------------------------------------------------

INSERT INTO control_asistencia.ref_tipo_punto_marcacion (codigo, nombre, descripcion)
VALUES
    (
        '001',
        'Asistencia General',
        'Punto de marcación de uso general para el registro de entrada y salida del personal en las instalaciones de la institución.'
    ),
    (
        '002',
        'Ingreso a Consultorio',
        'Punto de marcación ubicado en la entrada de consultorios médicos o unidades de atención para controlar el acceso del personal asistencial.'
    ),
    (
        '003',
        'Área Restringida',
        'Punto de marcación para el control de acceso a zonas con restricción de ingreso, como laboratorios, UCI, quirófanos o áreas de bioseguridad.'
    ),
    (
        'C004',
        'Control Operativo',
        'Punto de marcación destinado al seguimiento de rondas, tareas operativas o verificación de presencia del personal en puntos específicos de la institución.'
    )
ON CONFLICT (codigo) DO NOTHING;

-- ============================================================
-- TABLAS DE CONFIGURACIÓN
-- ============================================================

-- ------------------------------------------------------------
-- CFG: Puntos de marcación
-- Registra los puntos físicos donde se realiza el marcado QR.
-- El campo id_sede referencia la tabla de sedes del módulo
-- correspondiente del sistema Prismar (relación externa).
-- ------------------------------------------------------------

CREATE TABLE IF NOT EXISTS control_asistencia.cfg_puntos_marcacion (
    id              SERIAL          PRIMARY KEY,
    id_empresa      INTEGER         NOT NULL,
    id_sede         INTEGER         NOT NULL,
    id_tipo_punto   INTEGER         NOT NULL,
    nombre          VARCHAR(150)    NOT NULL,
    descripcion     TEXT,
    token_qr VARCHAR(100) UNIQUE NOT NULL,
    latitud         NUMERIC(10, 7)  NOT NULL,
    longitud        NUMERIC(10, 7)  NOT NULL,
    radio_metros    INTEGER         NOT NULL DEFAULT 80,
    activo          BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
 
    CONSTRAINT fk_cfg_puntos_marcacion_tipo_punto
        FOREIGN KEY (id_tipo_punto)
        REFERENCES control_asistencia.ref_tipo_punto_marcacion (id)
);

COMMENT ON TABLE  control_asistencia.cfg_puntos_marcacion                  IS 'Puntos físicos de marcación QR configurados por sede en la institución';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.id               IS 'Identificador único del punto de marcación';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.id_sede          IS 'Referencia a la sede del sistema Prismar (tabla externa al módulo)';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.id_tipo_punto    IS 'Tipo de punto de marcación según ref_tipo_punto_marcacion';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.nombre           IS 'Nombre descriptivo del punto de marcación';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.descripcion      IS 'Descripción adicional sobre la ubicación o propósito del punto';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.latitud          IS 'Latitud geográfica del punto de marcación (precisión 7 decimales)';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.longitud         IS 'Longitud geográfica del punto de marcación (precisión 7 decimales)';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.radio_metros     IS 'Radio en metros permitido para validar la proximidad al marcar (geofencing)';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.activo           IS 'Indica si el punto de marcación está habilitado para recibir marcaciones';
COMMENT ON COLUMN control_asistencia.cfg_puntos_marcacion.fecha_creacion
IS 'Fecha de creación del registro del punto de marcación';

-- ============================================================
-- ÍNDICES
-- ============================================================
 
-- Índice por sede: optimiza consultas de puntos asociados a una sede específica
CREATE INDEX IF NOT EXISTS idx_cfg_puntos_marcacion_id_sede
    ON control_asistencia.cfg_puntos_marcacion (id_sede);
 
-- Índice por tipo de punto: optimiza filtros y reportes por categoría de punto
CREATE INDEX IF NOT EXISTS idx_cfg_puntos_marcacion_id_tipo_punto
    ON control_asistencia.cfg_puntos_marcacion (id_tipo_punto);
 
-- Índice por sede y estado activo: optimiza la carga de puntos activos por sede
CREATE INDEX IF NOT EXISTS idx_cfg_puntos_marcacion_sede_activo
    ON control_asistencia.cfg_puntos_marcacion (id_sede, activo);
 
 
-- ============================================================
-- FIN DE MIGRACIÓN: 001_CA_UP.sql
-- ============================================================