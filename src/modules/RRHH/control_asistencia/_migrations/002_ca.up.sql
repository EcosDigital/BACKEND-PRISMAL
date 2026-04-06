-- ============================================================
-- MIGRACIÓN: 002_CA_UP.sql
-- MÓDULO: Control de Asistencia (CA)
-- SISTEMA: Prismar ERP
-- DESCRIPCIÓN: Creación de tabla de configuración de horarios
--              laborales (cabecera + detalle por bloques) y
--              tabla de asignación de horarios a empleados.
--              Soporta jornadas continuas y jornadas partidas.
--              La asociación se origina desde el módulo de
--              Recursos Humanos al registrar contratos.
-- DEPENDE DE: 001_CA_UP.sql
-- AUTOR: Prismar ERP - Ingeniería de Datos
-- FECHA: 2026-03-16
-- ============================================================


-- ============================================================
-- ESQUEMA
-- ============================================================

CREATE SCHEMA IF NOT EXISTS control_asistencia;


-- ============================================================
-- TABLAS DE CONFIGURACIÓN
-- ============================================================

-- ------------------------------------------------------------
-- CFG: Horarios de asistencia (cabecera)
-- Define los distintos horarios laborales de la organización.
-- Actúa como cabecera del horario. Los bloques de entrada y
-- salida se registran en cfg_horarios_detalle, permitiendo
-- soportar tanto jornadas continuas como jornadas partidas.
-- La asociación con empleados se realiza desde el módulo de
-- Recursos Humanos al registrar contratos.
-- ------------------------------------------------------------

CREATE TABLE IF NOT EXISTS control_asistencia.cfg_horarios_asistencia (
    id                  SERIAL          PRIMARY KEY,
    codigo              VARCHAR(20)     NOT NULL,
    nombre              VARCHAR(150)    NOT NULL,
    descripcion         TEXT,
    minutos_tolerancia  INTEGER         NOT NULL DEFAULT 0,
    activo              BOOLEAN         NOT NULL DEFAULT TRUE,
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_cfg_horarios_asistencia_codigo
        UNIQUE (codigo)
);

COMMENT ON TABLE  control_asistencia.cfg_horarios_asistencia                       IS 'Cabecera de horarios laborales de la organización. Los bloques horarios se definen en cfg_horarios_detalle. Soporta jornadas continuas y partidas.';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.id                    IS 'Identificador único del horario laboral';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.codigo                IS 'Código único interno del horario (ej: ADM, MED_AM)';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.nombre                IS 'Nombre descriptivo del horario laboral';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.descripcion           IS 'Descripción adicional sobre el horario o su aplicación';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.minutos_tolerancia    IS 'Minutos de gracia permitidos después de la hora de entrada antes de marcar tardanza';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.activo                IS 'Indica si el horario está disponible para ser asignado a empleados';
COMMENT ON COLUMN control_asistencia.cfg_horarios_asistencia.fecha_creacion        IS 'Fecha y hora de registro del horario en el sistema';


-- ------------------------------------------------------------
-- CFG: Detalle de bloques horarios
-- Almacena los bloques de entrada/salida de cada horario.
-- Un horario continuo tendrá un único registro en esta tabla.
-- Un horario partido tendrá dos o más registros ordenados.
-- El campo orden define la secuencia de los bloques del día.
-- ------------------------------------------------------------

CREATE TABLE IF NOT EXISTS control_asistencia.cfg_horarios_detalle (
    id              SERIAL      PRIMARY KEY,
    id_horario      INTEGER     NOT NULL,
    hora_entrada    TIME        NOT NULL,
    hora_salida     TIME        NOT NULL,
    orden           INTEGER     NOT NULL,
    activo          BOOLEAN     NOT NULL DEFAULT TRUE,
    fecha_creacion  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_cfg_horarios_detalle_horario
        FOREIGN KEY (id_horario)
        REFERENCES control_asistencia.cfg_horarios_asistencia (id),

    CONSTRAINT chk_cfg_horarios_detalle_hora_salida
        CHECK (hora_salida > hora_entrada),

    CONSTRAINT uq_cfg_horarios_detalle_horario_orden
        UNIQUE (id_horario, orden)
);

COMMENT ON TABLE  control_asistencia.cfg_horarios_detalle                      IS 'Bloques de entrada/salida que componen un horario laboral. Permite definir jornadas continuas (1 bloque) y jornadas partidas (2 o más bloques).';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.id                   IS 'Identificador único del bloque horario';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.id_horario           IS 'Horario al que pertenece este bloque';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.hora_entrada         IS 'Hora de inicio del bloque horario';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.hora_salida          IS 'Hora de fin del bloque horario. Debe ser posterior a hora_entrada.';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.orden                IS 'Número de orden del bloque dentro del horario. Define la secuencia del día.';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.activo               IS 'Indica si este bloque horario está vigente';
COMMENT ON COLUMN control_asistencia.cfg_horarios_detalle.fecha_creacion       IS 'Fecha y hora de registro del bloque en el sistema';


-- ------------------------------------------------------------
-- CFG: Asignación de horarios a empleados
-- Relaciona cada empleado con su horario laboral vigente.
-- Un empleado puede tener solo un horario activo a la vez,
-- garantizado por el índice único parcial definido abajo.
-- La referencia a id_empleado es externa al módulo (RR.HH.).
-- ------------------------------------------------------------

CREATE TABLE IF NOT EXISTS control_asistencia.cfg_empleado_horario (
    id              SERIAL      PRIMARY KEY,
    id_empleado     INTEGER     NOT NULL REFERENCES configuración.cfg_terceros(id),
    id_horario      INTEGER     NOT NULL REFERENCES control_asistencia.cfg_horarios_asistencia(id),
    fecha_inicio    DATE        NOT NULL,
    fecha_fin       DATE,
    activo          BOOLEAN     NOT NULL DEFAULT TRUE,
    observaciones   TEXT,
    registrado_por  INTEGER,
    fecha_creacion  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_cfg_empleado_horario_horario
        FOREIGN KEY (id_horario)
        REFERENCES control_asistencia.cfg_horarios_asistencia (id)
);

COMMENT ON TABLE  control_asistencia.cfg_empleado_horario                      IS 'Asignación de horarios laborales a empleados. Referencia externa a la tabla de empleados del módulo de Recursos Humanos.';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.id                   IS 'Identificador único de la asignación';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.id_empleado          IS 'Referencia externa al empleado en el módulo de Recursos Humanos';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.id_horario           IS 'Horario laboral asignado al empleado';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.activo               IS 'Indica si esta asignación está vigente. Solo debe existir una asignación activa por empleado.';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.observaciones        IS 'Notas u observaciones sobre la asignación del horario';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.registrado_por       IS 'Referencia externa al usuario del sistema que registró la asignación';
COMMENT ON COLUMN control_asistencia.cfg_empleado_horario.fecha_creacion       IS 'Fecha y hora en que se registró la asignación';


-- ============================================================
-- ÍNDICES
-- ============================================================

-- Índice por estado activo en horarios: optimiza la carga de horarios disponibles
CREATE INDEX IF NOT EXISTS idx_cfg_horarios_asistencia_activo
    ON control_asistencia.cfg_horarios_asistencia (activo);

-- Índice por horario en detalle: optimiza la carga de bloques de un horario
CREATE INDEX IF NOT EXISTS idx_cfg_horarios_detalle_id_horario
    ON control_asistencia.cfg_horarios_detalle (id_horario);

-- Índice por horario y orden: optimiza la lectura ordenada de bloques del día
CREATE INDEX IF NOT EXISTS idx_cfg_horarios_detalle_horario_orden
    ON control_asistencia.cfg_horarios_detalle (id_horario, orden);

-- Índice por empleado: optimiza la consulta del historial de horarios de un empleado
CREATE INDEX IF NOT EXISTS idx_cfg_empleado_horario_id_empleado
    ON control_asistencia.cfg_empleado_horario (id_empleado);

-- Índice por horario: optimiza la consulta de empleados asignados a un horario
CREATE INDEX IF NOT EXISTS idx_cfg_empleado_horario_id_horario
    ON control_asistencia.cfg_empleado_horario (id_horario);

-- Índice único parcial: garantiza que un empleado tenga como máximo
-- una asignación de horario activa en cualquier momento
CREATE UNIQUE INDEX IF NOT EXISTS uix_cfg_empleado_horario_activo
    ON control_asistencia.cfg_empleado_horario (id_empleado)
    WHERE activo = TRUE;

-- ============================================================
-- DATOS INICIALES
-- ============================================================

-- ------------------------------------------------------------
-- Cabecera de horario único
-- ------------------------------------------------------------

INSERT INTO control_asistencia.cfg_horarios_asistencia
    (codigo, nombre, descripcion, minutos_tolerancia)
VALUES
    (
        '001',
        'Horario General',
        'Horario estándar con jornada partida: mañana y tarde. Configuración base para operación inicial.',
        10
    )
ON CONFLICT (codigo) DO NOTHING;

-- ------------------------------------------------------------
-- Bloques de detalle (jornada partida)
-- ------------------------------------------------------------

-- Bloque mañana: 08:00–12:00
INSERT INTO control_asistencia.cfg_horarios_detalle
    (id_horario, hora_entrada, hora_salida, orden)
SELECT h.id, '08:00', '12:00', 1
FROM control_asistencia.cfg_horarios_asistencia h
WHERE h.codigo = '001'
ON CONFLICT (id_horario, orden) DO NOTHING;

-- Bloque tarde: 14:00–18:00
INSERT INTO control_asistencia.cfg_horarios_detalle
    (id_horario, hora_entrada, hora_salida, orden)
SELECT h.id, '14:00', '18:00', 2
FROM control_asistencia.cfg_horarios_asistencia h
WHERE h.codigo = '001'
ON CONFLICT (id_horario, orden) DO NOTHING;


-- ============================================================
-- FIN DE MIGRACIÓN: 002_CA_UP.sql
-- ============================================================