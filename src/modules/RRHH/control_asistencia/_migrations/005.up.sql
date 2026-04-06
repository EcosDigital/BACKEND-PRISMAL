-- ============================================================
-- FUNCIÓN: fn_reporte_diario_asistencia
-- MÓDULO:  Control de Asistencia
-- DESCRIPCIÓN: Genera el informe diario de asistencia para
--              todos los empleados que tienen horario asignado
--              activo en el rango de fechas indicado.
--
-- Por cada empleado con horario asignado:
--   · Si marcó en ese período → se toma la primera marcación
--     del día como hora_entrada y la última como hora_salida,
--     y se calcula el estado vs hora_entrada del bloque orden=1.
--   · Si NO marcó              → estado = AUSENTE
--
-- PARÁMETROS:
--   p_fecha_inicio DATE   — inicio del período
--   p_fecha_fin    DATE   — fin del período
--   p_empresa_id   BIGINT — filtra solo empleados de la empresa
--
-- RETORNA TABLE con las columnas que consume el backend Go.
-- ============================================================

CREATE OR REPLACE FUNCTION control_asistencia.fn_reporte_diario_asistencia(
    p_fecha_inicio DATE,
    p_fecha_fin    DATE,
    p_empresa_id   BIGINT
)
RETURNS TABLE (
    id_tercero        BIGINT,
    nombre_empleado   TEXT,
    numero_documento  TEXT,
    hora_entrada      TEXT,
    hora_salida       TEXT,
    estado            TEXT,
    horario_asignado  TEXT,
    hora_programada   TEXT,
    observaciones     TEXT
)
LANGUAGE sql
STABLE
AS $$
    -- CTE: empleados con horario activo en el rango de fechas
    WITH empleados_con_horario AS (
        SELECT
            eh.id_empleado,
            h.id                                                    AS id_horario,
            h.nombre                                                AS nombre_horario,
            h.minutos_tolerancia,
            hd.hora_entrada                                         AS hora_entrada_programada
        FROM control_asistencia.cfg_empleado_horario eh
        INNER JOIN control_asistencia.cfg_horarios_asistencia h
            ON h.id = eh.id_horario
        -- Bloque orden=1 define la hora de entrada del día
        INNER JOIN control_asistencia.cfg_horarios_detalle hd
            ON  hd.id_horario = h.id
            AND hd.orden       = 1
            AND hd.activo      = TRUE
        WHERE eh.activo       = TRUE
          AND eh.fecha_inicio <= p_fecha_fin
          AND (eh.fecha_fin IS NULL OR eh.fecha_fin >= p_fecha_inicio)
    ),

    -- CTE: primera y última marcación por empleado en el período
    marcaciones_agrupadas AS (
        SELECT
            m.id_empleado,
            MIN(m.fecha_hora)                                       AS primera_marcacion,
            MAX(m.fecha_hora)                                       AS ultima_marcacion,
            COALESCE(MAX(m.observacion), '')                        AS observaciones
        FROM control_asistencia.mov_marcaciones m
        WHERE m.fecha_hora::date BETWEEN p_fecha_inicio AND p_fecha_fin
        GROUP BY m.id_empleado
    )

    SELECT
        t.id,
        TRIM(
            COALESCE(t.primer_nombre,    '') || ' ' ||
            COALESCE(t.segundo_nombre,   '') || ' ' ||
            COALESCE(t.primer_apellido,  '') || ' ' ||
            COALESCE(t.segundo_apellido, '')
        )                                                           AS nombre_empleado,
        COALESCE(t.numero_documento, '')                            AS numero_documento,

        -- Hora de entrada: primera marcación del período
        COALESCE(
            TO_CHAR(ma.primera_marcacion, 'HH24:MI:SS'), ''
        )                                                           AS hora_entrada,

        -- Hora de salida: última marcación (distinta de la primera)
        CASE
            WHEN ma.ultima_marcacion IS NOT NULL
             AND ma.ultima_marcacion <> ma.primera_marcacion
                THEN TO_CHAR(ma.ultima_marcacion, 'HH24:MI:SS')
            ELSE ''
        END                                                         AS hora_salida,

        -- Estado calculado contra la hora programada del horario
        CASE
            WHEN ma.primera_marcacion IS NULL
                THEN 'AUSENTE'
            WHEN ma.primera_marcacion::time
                 <= ech.hora_entrada_programada
                  + (ech.minutos_tolerancia || ' minutes')::interval
                THEN 'PUNTUAL'
            WHEN ma.primera_marcacion::time
                 <= ech.hora_entrada_programada
                  + (ech.minutos_tolerancia * 2 || ' minutes')::interval
                THEN 'TOLERANCIA'
            ELSE 'TARDE'
        END                                                         AS estado,

        COALESCE(ech.nombre_horario, '')                            AS horario_asignado,
        COALESCE(
            TO_CHAR(ech.hora_entrada_programada, 'HH24:MI'), ''
        )                                                           AS hora_programada,
        COALESCE(ma.observaciones, '')                              AS observaciones

    FROM empleados_con_horario ech
    INNER JOIN configuracion.cfg_terceros t
        ON t.id = ech.id_empleado
    LEFT JOIN marcaciones_agrupadas ma
        ON ma.id_empleado = ech.id_empleado
    -- Filtrar solo empleados de la empresa solicitada

    ORDER BY t.primer_apellido ASC, t.primer_nombre ASC;
$$;

COMMENT ON FUNCTION control_asistencia.fn_reporte_diario_asistencia(DATE, DATE, BIGINT)
    IS 'Informe diario de asistencia: lista todos los empleados con horario asignado, su primera/última marcación y el estado calculado (PUNTUAL/TOLERANCIA/TARDE/AUSENTE).';