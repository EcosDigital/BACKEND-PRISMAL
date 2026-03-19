-- ============================================================
-- 008_inv_up.sql
-- Reemplaza fn_kardex_articulo para que el parámetro p_lote
-- sea OPCIONAL: NULL o '' retorna todos los lotes.
-- No depende de ninguna vista — usa CTEs directamente sobre
-- las tablas base (igual que 007_inv_up.sql).
-- ============================================================

CREATE OR REPLACE FUNCTION inventario.fn_kardex_articulo(
    p_id_articulo    INT,
    p_id_bodega      INT,
    p_lote           VARCHAR(100),   -- NULL o '' = todos los lotes
    p_fecha_inicio   DATE,
    p_fecha_fin      DATE,
    p_id_empresa     INT
)
RETURNS TABLE (
    fila_orden         INT,
    fecha              TEXT,
    tipo_movimiento    TEXT,
    documento          TEXT,
    referencia_externa TEXT,
    observaciones      TEXT,
    cantidad_entrada   NUMERIC(14,4),
    cantidad_salida    NUMERIC(14,4),
    costo_unitario     NUMERIC(14,4),
    saldo_cantidad     NUMERIC(14,4),
    es_anulado         BOOLEAN,
    motivo_anulacion   TEXT,
    id_movimiento      INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_filtrar_lote BOOLEAN;
BEGIN

    -- Si p_lote es NULL o vacío → no filtrar por lote (todos)
    v_filtrar_lote := (p_lote IS NOT NULL AND p_lote <> '');

    RETURN QUERY
    WITH

    -- ── 1. Saldo acumulado ANTES del período ──────────────────
    -- Solo movimientos ejecutados (no anulados) construyen el saldo inicial.
    -- Si v_filtrar_lote=FALSE se suman todos los lotes de esa bodega.
    saldo_inicial AS (
        SELECT
            COALESCE(SUM(
                CASE
                    WHEN tm.nombre IN ('Entradas', 'Saldo Inicial')
                         AND m.id_bodega_destino = p_id_bodega
                    THEN d.cantidad

                    WHEN tm.nombre = 'Traslados'
                         AND m.id_bodega_destino = p_id_bodega
                    THEN d.cantidad

                    WHEN tm.nombre = 'Traslados'
                         AND m.id_bodega_origen = p_id_bodega
                    THEN -d.cantidad

                    WHEN tm.nombre IN ('Bajas', 'Despachos')
                         AND m.id_bodega_origen = p_id_bodega
                    THEN -d.cantidad

                    ELSE 0
                END
            ), 0)::NUMERIC(14,4) AS cantidad_acumulada,

            COALESCE((
                SELECT e.costo_promedio
                FROM inventario.inv_existencias_articulos e
                WHERE e.id_articulo = p_id_articulo
                  AND e.id_bodega   = p_id_bodega
                  AND e.id_empresa  = p_id_empresa
                  AND (NOT v_filtrar_lote OR COALESCE(e.lote, '') = COALESCE(p_lote, ''))
                LIMIT 1
            ), 0)::NUMERIC(14,4) AS costo_promedio

        FROM inventario.mov_movimientos m
        INNER JOIN inventario.mov_movimientos_detalle d
            ON d.id_movimiento = m.id
        INNER JOIN inventario.ref_tipo_movimiento tm
            ON tm.id = m.id_tipo_movimiento
        INNER JOIN inventario.ref_estado_movimiento em
            ON em.id = m.id_estado
        WHERE d.id_articulo        = p_id_articulo
          AND m.id_empresa         = p_id_empresa
          AND m.fecha_movimiento   < p_fecha_inicio
          AND (m.id_bodega_destino = p_id_bodega OR m.id_bodega_origen = p_id_bodega)
          AND em.codigo           != '05'   -- excluir anulados del saldo inicial
          -- Filtro de lote condicional
          AND (NOT v_filtrar_lote OR COALESCE(d.lote, '') = COALESCE(p_lote, ''))
    ),

    -- ── 2. Movimientos dentro del período ─────────────────────
    movimientos_periodo AS (
        SELECT
            m.id                                                            AS id_movimiento,
            m.fecha_movimiento,
            tm.nombre                                                       AS tipo_movimiento,
            mg.prefijo_comprobante || '-' ||
                LPAD(mg.consecutivo_comprobante::TEXT, 6, '0')             AS documento,
            COALESCE(mg.documento_soporte, '')                             AS referencia_externa,
            COALESCE(mg.observaciones, '')                                 AS observaciones,

            -- Entrada a la bodega (0 si anulado)
            CASE
                WHEN em.codigo = '05' THEN 0::NUMERIC(14,4)
                WHEN tm.nombre IN ('Entradas', 'Saldo Inicial')
                     AND m.id_bodega_destino = p_id_bodega
                THEN d.cantidad
                WHEN tm.nombre = 'Traslados'
                     AND m.id_bodega_destino = p_id_bodega
                THEN d.cantidad
                ELSE 0::NUMERIC(14,4)
            END AS cantidad_entrada,

            -- Salida de la bodega (0 si anulado)
            CASE
                WHEN em.codigo = '05' THEN 0::NUMERIC(14,4)
                WHEN tm.nombre = 'Traslados'
                     AND m.id_bodega_origen = p_id_bodega
                THEN d.cantidad
                WHEN tm.nombre IN ('Bajas', 'Despachos')
                     AND m.id_bodega_origen = p_id_bodega
                THEN d.cantidad
                ELSE 0::NUMERIC(14,4)
            END AS cantidad_salida,

            COALESCE(d.valor_unitario, 0)::NUMERIC(14,4)                  AS costo_unitario,
            (em.codigo = '05')                                             AS es_anulado,
            COALESCE(m.motivo_anulacion, '')                               AS motivo_anulacion

        FROM inventario.mov_movimientos m
        INNER JOIN inventario.mov_movimientos_detalle d
            ON d.id_movimiento = m.id
        INNER JOIN inventario.ref_tipo_movimiento tm
            ON tm.id = m.id_tipo_movimiento
        INNER JOIN inventario.ref_estado_movimiento em
            ON em.id = m.id_estado
        INNER JOIN comprobantes.mov_gestion_comprobantes mg
            ON mg.id = m.id_mov_comprobante
        WHERE d.id_articulo        = p_id_articulo
          AND m.id_empresa         = p_id_empresa
          AND m.fecha_movimiento BETWEEN p_fecha_inicio AND p_fecha_fin
          AND (m.id_bodega_destino = p_id_bodega OR m.id_bodega_origen = p_id_bodega)
          -- Filtro de lote condicional
          AND (NOT v_filtrar_lote OR COALESCE(d.lote, '') = COALESCE(p_lote, ''))
    ),

    -- ── 3. Unión: saldo inicial + movimientos del período ──────
    union_filas AS (
        -- Fila 0: saldo al inicio del período
        SELECT
            0                                        AS orden,
            p_fecha_inicio::TEXT                     AS fecha,
            'Saldo Inicial'::TEXT                    AS tipo_movimiento,
            'SALDO INICIAL'::TEXT                    AS documento,
            ''::TEXT                                 AS referencia_externa,
            'Saldo acumulado al inicio del período'::TEXT AS observaciones,
            0::NUMERIC(14,4)                         AS cantidad_entrada,
            0::NUMERIC(14,4)                         AS cantidad_salida,
            si.costo_promedio                        AS costo_unitario,
            si.cantidad_acumulada                    AS delta,
            FALSE                                    AS es_anulado,
            ''::TEXT                                 AS motivo_anulacion,
            0                                        AS id_movimiento
        FROM saldo_inicial si

        UNION ALL

        SELECT
            ROW_NUMBER() OVER (
                ORDER BY mp.fecha_movimiento, mp.id_movimiento
            )::INT                                   AS orden,
            mp.fecha_movimiento::TEXT                AS fecha,
            mp.tipo_movimiento,
            mp.documento,
            mp.referencia_externa,
            mp.observaciones,
            mp.cantidad_entrada,
            mp.cantidad_salida,
            mp.costo_unitario,
            (mp.cantidad_entrada - mp.cantidad_salida) AS delta,
            mp.es_anulado,
            mp.motivo_anulacion,
            mp.id_movimiento
        FROM movimientos_periodo mp
    ),

    -- ── 4. Saldo acumulado progresivo ──────────────────────────
    kardex_saldos AS (
        SELECT
            uf.orden,
            uf.fecha,
            uf.tipo_movimiento,
            uf.documento,
            uf.referencia_externa,
            uf.observaciones,
            uf.cantidad_entrada,
            uf.cantidad_salida,
            uf.costo_unitario,
            uf.es_anulado,
            uf.motivo_anulacion,
            uf.id_movimiento,
            SUM(uf.delta) OVER (
                ORDER BY uf.orden
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS saldo_cantidad
        FROM union_filas uf
    )

    SELECT
        ks.orden::INT,
        ks.fecha,
        ks.tipo_movimiento,
        ks.documento,
        ks.referencia_externa,
        ks.observaciones,
        ks.cantidad_entrada,
        ks.cantidad_salida,
        ks.costo_unitario,
        ks.saldo_cantidad,
        ks.es_anulado,
        ks.motivo_anulacion,
        ks.id_movimiento
    FROM kardex_saldos ks
    ORDER BY ks.orden ASC;

END;
$$;

COMMENT ON FUNCTION inventario.fn_kardex_articulo IS
'Kardex de inventario por artículo/bodega con lote OPCIONAL.
p_lote = NULL o ""  → agrega todos los lotes de la bodega.
p_lote = "LOT01"    → filtra exclusivamente ese lote.
Fila 0 = saldo inicial acumulado ANTES de p_fecha_inicio (sin anulados).
Filas 1..N = movimientos del período en orden cronológico.
Movimientos anulados: visibles con es_anulado=TRUE, delta=0 (no afectan saldo).
No depende de vistas — usa CTEs directamente sobre tablas base.';