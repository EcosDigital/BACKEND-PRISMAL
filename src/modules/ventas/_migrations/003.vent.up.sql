-- ============================================================
-- 003_fiados_up.sql
-- Funciones de consulta y vista principal
-- ============================================================

-- ─────────────────────────────────────────────────────────────
-- fn_historial_cuenta
-- Kardex de movimientos de una cuenta.
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION fiados.fn_historial_cuenta(
    p_id_cuenta    INT,
    p_id_empresa   INT,
    p_fecha_inicio DATE DEFAULT NULL,
    p_fecha_fin    DATE DEFAULT NULL
)
RETURNS TABLE (
    id_movimiento       INT,
    fecha               DATE,
    tipo_codigo         VARCHAR(20),
    tipo_nombre         TEXT,
    valor               NUMERIC(14,2),
    observaciones       TEXT,
    saldo_anterior      NUMERIC(14,2),
    saldo_movimiento    NUMERIC(14,2),
    saldo_despues       NUMERIC(14,2),
    documento           TEXT,
    estado_comprobante  TEXT,
    es_anulado          BOOLEAN
)
LANGUAGE plpgsql AS
$$
BEGIN
    RETURN QUERY
    SELECT
        mf.id,
        mf.fecha_movimiento,
        rtm.codigo                                                     AS tipo_codigo,
        rtm.nombre::TEXT                                               AS tipo_nombre,
        mf.valor,
        COALESCE(mf.observaciones, '')::TEXT,
        mf.saldo_anterior,
        mf.saldo_movimiento,
        mf.saldo_despues,
        (mg.prefijo_comprobante || '-' ||
            LPAD(mg.consecutivo_comprobante::TEXT, 6, '0'))::TEXT      AS documento,
        rec.nombre::TEXT                                               AS estado_comprobante,
        (rec.nombre = 'ANULADO')                                       AS es_anulado
    FROM fiados.mov_fiados mf
    INNER JOIN fiados.ref_tipo_movimiento rtm
        ON rtm.id = mf.id_tipo
    INNER JOIN comprobantes.mov_gestion_comprobantes mg
        ON mg.id = mf.id_mov_comprobante
    INNER JOIN comprobantes.ref_estado_comprobante rec
        ON rec.id = mg.id_estado_comprobante
    WHERE mf.id_cuenta   = p_id_cuenta
      AND mf.id_empresa  = p_id_empresa
      AND (p_fecha_inicio IS NULL OR mf.fecha_movimiento >= p_fecha_inicio)
      AND (p_fecha_fin    IS NULL OR mf.fecha_movimiento <= p_fecha_fin)
    ORDER BY mf.fecha_movimiento DESC, mf.id DESC;
END;
$$;

COMMENT ON FUNCTION fiados.fn_historial_cuenta IS
'Kardex de una cuenta de fiado.
Retorna saldo_anterior, saldo_movimiento y saldo_despues por fila.
Los movimientos anulados aparecen con es_anulado=TRUE.
Documento en formato PREFIJO-000001.';


-- ─────────────────────────────────────────────────────────────
-- fn_resumen_cartera
-- Stats del dashboard: una fila por empresa/sede.
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION fiados.fn_resumen_cartera(
    p_id_empresa INT,
    p_id_sede    INT DEFAULT NULL
)
RETURNS TABLE (
    total_cuentas        INT,
    cuentas_al_dia       INT,
    cuentas_pendientes   INT,
    cuentas_atrasadas    INT,
    total_cartera        NUMERIC(14,2),
    total_fiados_mes     NUMERIC(14,2),
    total_abonos_mes     NUMERIC(14,2)
)
LANGUAGE plpgsql AS
$$
DECLARE
    v_inicio_mes DATE := DATE_TRUNC('month', CURRENT_DATE)::DATE;
    v_cod_fiado  VARCHAR(20) := '01';
    v_cod_abono  VARCHAR(20) := '02';
BEGIN
    RETURN QUERY
    SELECT
        COUNT(cf.id)::INT,
        COUNT(cf.id) FILTER (WHERE re.codigo = '03')::INT,  -- Al día
        COUNT(cf.id) FILTER (WHERE re.codigo = '02')::INT,  -- Pendiente
        COUNT(cf.id) FILTER (WHERE re.codigo = '01')::INT,  -- Atrasado
        COALESCE(SUM(cf.saldo_actual) FILTER (WHERE cf.saldo_actual > 0), 0),

        -- Fiados del mes (excluyendo comprobantes anulados)
        COALESCE((
            SELECT SUM(mf.valor)
            FROM fiados.mov_fiados mf
            INNER JOIN fiados.ref_tipo_movimiento rtm ON rtm.id = mf.id_tipo
            INNER JOIN comprobantes.mov_gestion_comprobantes mg
                ON mg.id = mf.id_mov_comprobante
            INNER JOIN comprobantes.ref_estado_comprobante rec
                ON rec.id = mg.id_estado_comprobante
            WHERE mf.id_empresa        = p_id_empresa
              AND (p_id_sede IS NULL OR mf.id_sede = p_id_sede)
              AND rtm.codigo           = v_cod_fiado
              AND rec.nombre          != 'ANULADO'
              AND mf.fecha_movimiento >= v_inicio_mes
        ), 0),

        -- Abonos del mes
        COALESCE((
            SELECT SUM(mf.valor)
            FROM fiados.mov_fiados mf
            INNER JOIN fiados.ref_tipo_movimiento rtm ON rtm.id = mf.id_tipo
            INNER JOIN comprobantes.mov_gestion_comprobantes mg
                ON mg.id = mf.id_mov_comprobante
            INNER JOIN comprobantes.ref_estado_comprobante rec
                ON rec.id = mg.id_estado_comprobante
            WHERE mf.id_empresa        = p_id_empresa
              AND (p_id_sede IS NULL OR mf.id_sede = p_id_sede)
              AND rtm.codigo           = v_cod_abono
              AND rec.nombre          != 'ANULADO'
              AND mf.fecha_movimiento >= v_inicio_mes
        ), 0)

    FROM fiados.cfg_cuentas_fiado cf
    INNER JOIN fiados.ref_estado_credito re ON re.id = cf.id_estado
    WHERE cf.id_empresa = p_id_empresa
      AND (p_id_sede IS NULL OR cf.id_sede = p_id_sede);
END;
$$;

COMMENT ON FUNCTION fiados.fn_resumen_cartera IS
'Dashboard de fiados: conteos por estado (Al día/Pendiente/Atrasado),
cartera total y movimientos del mes. Una sola fila de respuesta.
Excluye comprobantes anulados de los totales del mes.';


-- ─────────────────────────────────────────────────────────────
-- VISTA PRINCIPAL — cuentas con datos del tercero (cliente)
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE VIEW fiados.v_cuentas AS
SELECT
    cf.id                                              AS id_cuenta,
    cf.id_tercero,
    -- Nombre completo del cliente desde cfg_terceros
    t.primer_nombre || ' ' ||
        COALESCE(t.primer_apellido, '')                      AS nombre_cliente,
    t.numero_documento                         AS cedula,
    t.telefono,
    t.direccion,
    cf.saldo_actual,
    re.codigo                                          AS estado_codigo,
    re.nombre                                          AS estado_nombre,
    cf.limite_credito,
    cf.fecha_apertura,
    cf.fecha_cierre,
    cf.observaciones,
    -- Último movimiento activo (no anulado)
    (
        SELECT mf.fecha_movimiento
        FROM fiados.mov_fiados mf
        INNER JOIN comprobantes.mov_gestion_comprobantes mg
            ON mg.id = mf.id_mov_comprobante
        INNER JOIN comprobantes.ref_estado_comprobante rec
            ON rec.id = mg.id_estado_comprobante
        WHERE mf.id_cuenta  = cf.id
          AND rec.nombre    != 'ANULADO'
        ORDER BY mf.fecha_movimiento DESC, mf.id DESC
        LIMIT 1
    )                                                  AS fecha_ultimo_movimiento,
    -- Tipo del último movimiento
    (
        SELECT rtm.nombre
        FROM fiados.mov_fiados mf
        INNER JOIN fiados.ref_tipo_movimiento rtm ON rtm.id = mf.id_tipo
        INNER JOIN comprobantes.mov_gestion_comprobantes mg
            ON mg.id = mf.id_mov_comprobante
        INNER JOIN comprobantes.ref_estado_comprobante rec
            ON rec.id = mg.id_estado_comprobante
        WHERE mf.id_cuenta  = cf.id
          AND rec.nombre    != 'ANULADO'
        ORDER BY mf.fecha_movimiento DESC, mf.id DESC
        LIMIT 1
    )                                                  AS ultimo_tipo_movimiento,
    cf.id_empresa,
    cf.id_sede,
    cf.created_at
FROM fiados.cfg_cuentas_fiado cf
INNER JOIN configuracion.cfg_terceros    t  ON t.id  = cf.id_tercero
INNER JOIN fiados.ref_estado_credito     re ON re.id = cf.id_estado;

COMMENT ON VIEW fiados.v_cuentas IS
'Vista principal del módulo de fiados.
Une cfg_cuentas_fiado + cfg_terceros + ref_estado_credito.
Expone nombre, cédula, teléfono, dirección del cliente y estado de la cuenta.
Usada por el endpoint GET /fiados/cuentas y el buscador del frontend.';