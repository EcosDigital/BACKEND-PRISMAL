

	-- ============================================================
-- FUNCIÓN: fn_historial_asistencia
-- MÓDULO:  Control de Asistencia
-- DESCRIPCIÓN: Retorna el historial de marcaciones de un
--              empleado en un rango de fechas, calculando el
--              estado (PUNTUAL / TOLERANCIA / TARDE / SIN_HORARIO)
--              mediante cruce con el horario asignado vigente.
--
-- El estado se deriva comparando la hora de marcación contra
-- la hora_entrada del bloque orden=1 del horario asignado,
-- aplicando minutos_tolerancia de cfg_horarios_asistencia.
--
-- PARÁMETROS:
--   p_id_tercero   BIGINT   — id del empleado (cfg_terceros)
--   p_fecha_inicio DATE     — inicio del período
--   p_fecha_fin    DATE     — fin del período
--
-- RETORNA TABLE con las columnas que consume el backend Go.
-- ============================================================

CREATE OR REPLACE FUNCTION control_asistencia.fn_historial_asistencia(
    p_id_tercero   BIGINT,
    p_fecha_inicio DATE,
    p_fecha_fin    DATE
)
RETURNS TABLE (
    fila_orden          INT,
    fecha               TEXT,
    hora                TEXT,
    dia_semana          TEXT,
    tipo_marcacion      TEXT,
    punto_marcacion     TEXT,
    observaciones       TEXT,
    horario_asignado    TEXT,
    hora_programada     TEXT,
    minutos_diferencia  INT,
    estado              TEXT
)
LANGUAGE sql
STABLE
AS $$
    SELECT
        ROW_NUMBER() OVER (ORDER BY m.fecha_hora ASC)::INT  AS fila_orden,
        m.fecha_hora::date::text                            AS fecha,
        TO_CHAR(m.fecha_hora, 'HH24:MI:SS')                 AS hora,
        TRIM(TO_CHAR(m.fecha_hora, 'Day'))                   AS dia_semana,
        COALESCE(m.tipo_evento, 'ASISTENCIA')                AS tipo_marcacion,
        COALESCE(p.nombre,      '')                          AS punto_marcacion,
        COALESCE(m.observacion, '')                          AS observaciones,

        -- Nombre del horario asignado vigente en la fecha de marcación
        COALESCE(h.nombre, '')                               AS horario_asignado,

        -- Hora de entrada del primer bloque del horario (orden = 1)
        COALESCE(TO_CHAR(hd.hora_entrada, 'HH24:MI'), '')    AS hora_programada,

        -- Diferencia en minutos: positivo = tarde, negativo = temprano
        CASE
            WHEN hd.hora_entrada IS NULL THEN 0
            ELSE EXTRACT(
                EPOCH FROM (m.fecha_hora::time - hd.hora_entrada)
            )::int / 60
        END                                                  AS minutos_diferencia,

        -- Estado calculado con tolerancia definida en el horario
        CASE
            WHEN eh.id IS NULL
              OR h.id  IS NULL
              OR hd.id IS NULL
                THEN 'SIN_HORARIO'
            WHEN m.fecha_hora::time
                 <= hd.hora_entrada
                  + (h.minutos_tolerancia || ' minutes')::interval
                THEN 'PUNTUAL'
            WHEN m.fecha_hora::time
                 <= hd.hora_entrada
                  + (h.minutos_tolerancia * 2 || ' minutes')::interval
                THEN 'TOLERANCIA'
            ELSE 'TARDE'
        END                                                  AS estado

    FROM control_asistencia.mov_marcaciones m

    -- Punto de marcación QR
    LEFT JOIN control_asistencia.cfg_puntos_marcacion p
        ON p.id = m.id_punto_marcacion

    -- Horario vigente en la fecha exacta de cada marcación
    LEFT JOIN control_asistencia.cfg_empleado_horario eh
        ON  eh.id_empleado   = m.id_empleado
        AND eh.activo        = TRUE
        AND eh.fecha_inicio <= m.fecha_hora::date
        AND (eh.fecha_fin IS NULL OR eh.fecha_fin >= m.fecha_hora::date)

    -- Cabecera del horario (nombre + tolerancia)
    LEFT JOIN control_asistencia.cfg_horarios_asistencia h
        ON h.id = eh.id_horario

    -- Primer bloque del horario (orden=1) — define la hora de entrada del día
    -- cfg_horarios_detalle NO tiene dia_semana; orden=1 es el bloque de apertura
    LEFT JOIN control_asistencia.cfg_horarios_detalle hd
        ON  hd.id_horario = h.id
        AND hd.orden       = 1
        AND hd.activo      = TRUE

    WHERE m.id_empleado        = p_id_tercero
      AND m.fecha_hora::date BETWEEN p_fecha_inicio AND p_fecha_fin

    ORDER BY m.fecha_hora ASC;
$$;

COMMENT ON FUNCTION control_asistencia.fn_historial_asistencia(BIGINT, DATE, DATE)
    IS 'Historial de marcaciones de un empleado con estado calculado (PUNTUAL/TOLERANCIA/TARDE/SIN_HORARIO). El estado se deriva del horario asignado vigente y su tolerancia configurada.';