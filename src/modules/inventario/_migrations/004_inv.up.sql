-- ============================================================
-- 004_inv_up.sql
-- Función para anular una entrada al inventario
-- Revierte stock, actualiza estados y genera nota en contabilidad
-- ============================================================

CREATE OR REPLACE FUNCTION inventario.fn_anular_entrada(
    p_id_movimiento  INT,
    p_motivo         TEXT,
    p_anulado_por    INT,
    p_id_empresa     INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    v_movimiento      inventario.mov_movimientos%ROWTYPE;
    v_comprobante     comprobantes.cfg_comprobante%ROWTYPE;
    v_id_mov_comp     INT;
    v_id_estado_anulado_inv  INT;
    v_id_estado_anulado_comp INT;

    -- Para iterar el detalle
    v_linea           inventario.mov_movimientos_detalle%ROWTYPE;
    v_stock_actual    NUMERIC(14,4);
    v_costo_actual    NUMERIC(14,4);
    v_nuevo_stock     NUMERIC(14,4);
    v_nuevo_costo     NUMERIC(14,4);
BEGIN

    -- ── 1. Leer y bloquear el movimiento ─────────────────────────
    SELECT * INTO v_movimiento
    FROM inventario.mov_movimientos
    WHERE id = p_id_movimiento
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'No se encontró el movimiento con id=%', p_id_movimiento;
    END IF;

    -- ── 2. Verificar que es una Entrada ──────────────────────────
    IF (SELECT nombre FROM inventario.ref_tipo_movimiento WHERE id = v_movimiento.id_tipo_movimiento) != 'Entradas' THEN
        RAISE EXCEPTION 'Solo se pueden anular movimientos de tipo Entrada';
    END IF;

    -- ── 3. Verificar que no esté ya anulado ──────────────────────
    IF (SELECT codigo FROM inventario.ref_estado_movimiento WHERE id = v_movimiento.id_estado) = '05' THEN
        RAISE EXCEPTION 'Este movimiento ya fue anulado';
    END IF;

    -- ── 4. Leer el comprobante y verificar que permite anulación ─
    v_id_mov_comp := v_movimiento.id_mov_comprobante;

    SELECT c.* INTO v_comprobante
    FROM comprobantes.cfg_comprobante c
    INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id_comprobante = c.id
    WHERE mg.id = v_id_mov_comp
    FOR UPDATE;

    IF v_comprobante.permite_anulacion = FALSE THEN
        RAISE EXCEPTION 'El comprobante "%" no permite anulación', v_comprobante.nombre_comprobante;
    END IF;

    -- ── 5. Obtener IDs de estado Anulado ─────────────────────────
    SELECT id INTO v_id_estado_anulado_inv
    FROM inventario.ref_estado_movimiento
    WHERE codigo = '05'; -- Anulado

    SELECT id INTO v_id_estado_anulado_comp
    FROM comprobantes.ref_estado_comprobante
    WHERE nombre = 'ANULADO';

    -- ── 6. Revertir existencias por cada línea del detalle ────────
    FOR v_linea IN
        SELECT * FROM inventario.mov_movimientos_detalle
        WHERE id_movimiento = p_id_movimiento
    LOOP
        -- Leer stock y costo actual
        SELECT stock_actual, COALESCE(costo_promedio, 0)
        INTO v_stock_actual, v_costo_actual
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega    = v_movimiento.id_bodega_destino
          AND id_empresa   = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        IF NOT FOUND THEN
            -- No hay registro de existencia — no hay nada que revertir
            CONTINUE;
        END IF;

        -- Calcular nuevo stock (puede quedar negativo — opción B)
        v_nuevo_stock := v_stock_actual - v_linea.cantidad;

        -- Recalcular costo promedio inverso
        -- Si el stock resultante es 0 o negativo dejamos el costo en 0
        IF v_nuevo_stock > 0 THEN
            v_nuevo_costo := (
                (v_stock_actual * v_costo_actual) - (v_linea.cantidad * v_costo_actual)
            ) / v_nuevo_stock;
        ELSE
            v_nuevo_costo := 0;
        END IF;

        UPDATE inventario.inv_existencias_articulos
        SET
            stock_actual   = v_nuevo_stock,
            costo_promedio = v_nuevo_costo,
            updated_at     = NOW()
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega    = v_movimiento.id_bodega_destino
          AND id_empresa   = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

    END LOOP;

    -- ── 7. Actualizar estado del movimiento ───────────────────────
    UPDATE inventario.mov_movimientos
    SET
        id_estado        = v_id_estado_anulado_inv,
        anulado_por      = p_anulado_por,
        anulado_at       = NOW(),
        motivo_anulacion = p_motivo,
        updated_at       = NOW(),
        updated_by       = p_anulado_por
    WHERE id = p_id_movimiento;

    -- ── 8. Actualizar estado del comprobante ──────────────────────
    UPDATE comprobantes.mov_gestion_comprobantes
    SET
        id_estado_comprobante = v_id_estado_anulado_comp,
        id_user_anulo         = p_anulado_por,
        fecha_anulacion       = NOW(),
        motivo_anulacion      = p_motivo,
        updated_at            = NOW(),
        updated_by            = p_anulado_por
    WHERE id = v_id_mov_comp;

    -- ── 9. Revertir movimientos contables si aplica ───────────────
    -- Simplemente eliminamos los asientos de esta entrada
    -- ya que se anulan completamente (débito y crédito desaparecen)
    IF v_comprobante.aplica_mov_contable = TRUE THEN
        DELETE FROM contabilidad.mov_movimientos_contables
        WHERE id_mov_comprobante = v_id_mov_comp;
    END IF;

END;
$$;

COMMENT ON FUNCTION inventario.fn_anular_entrada IS
'Anula una entrada al inventario dentro de una sola transacción.
Pasos: verifica estado → revierte existencias (stock negativo permitido) →
       actualiza mov_movimientos → actualiza mov_gestion_comprobantes →
       elimina asientos contables si aplica.
En caso de error en cualquier paso se hace ROLLBACK automático.';