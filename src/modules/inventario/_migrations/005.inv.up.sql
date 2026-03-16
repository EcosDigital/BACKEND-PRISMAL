-- ============================================================
-- 005_inv_up.sql
-- Funciones para traslados de inventario entre bodegas
-- ============================================================

-- ─── Tipo para las líneas del traslado ───────────────────────
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 't_traslado_detalle'
        AND n.nspname = 'inventario'
    ) THEN
        CREATE TYPE inventario.t_traslado_detalle AS (
            id_articulo  INT,
            lote         VARCHAR(100),
            cantidad     NUMERIC(14,4),
            observacion  TEXT
        );
    END IF;
END
$$;

-- ─── Función principal de traslado ───────────────────────────
CREATE OR REPLACE FUNCTION inventario.fn_registrar_traslado(
    p_id_comprobante      INT,
    p_id_bodega_origen    INT,
    p_id_bodega_destino   INT,
    p_referencia_externa  VARCHAR(100),
    p_observaciones       TEXT,
    p_fecha_movimiento    DATE,
    p_created_by          INT,
    p_id_empresa          INT,
    p_id_sede             INT,
    p_detalle             inventario.t_traslado_detalle[]
)
RETURNS TABLE (
    id_movimiento        INT,
    id_mov_comprobante   INT,
    consecutivo_generado INT,
    prefijo_comprobante  VARCHAR(10)
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_comprobante        comprobantes.cfg_comprobante%ROWTYPE;
    v_consecutivo        INT;
    v_id_mov_comprobante INT;
    v_id_movimiento      INT;
    v_id_estado_aplicado INT;
    v_id_tipo_traslado   INT;

    v_linea              inventario.t_traslado_detalle;

    -- Existencias origen
    v_stock_origen       NUMERIC(14,4);
    v_costo_origen       NUMERIC(14,4);

    -- Existencias destino
    v_stock_destino      NUMERIC(14,4);
    v_costo_destino      NUMERIC(14,4);
    v_nuevo_costo_destino NUMERIC(14,4);

    -- Contabilidad
    v_id_cuenta_inventario_origen  INT;
    v_id_cuenta_inventario_destino INT;
    v_valor_linea                  NUMERIC(14,2);
BEGIN

    -- ── 1. Leer y bloquear comprobante ────────────────────────────
    SELECT * INTO v_comprobante
    FROM comprobantes.cfg_comprobante
    WHERE id = p_id_comprobante
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Comprobante no encontrado (id=%)', p_id_comprobante;
    END IF;

    IF v_comprobante.is_active = FALSE THEN
        RAISE EXCEPTION 'El comprobante "%" está inactivo', v_comprobante.nombre_comprobante;
    END IF;

    -- ── 2. Validar que las bodegas son distintas ──────────────────
    IF p_id_bodega_origen = p_id_bodega_destino THEN
        RAISE EXCEPTION 'La bodega origen y destino no pueden ser la misma';
    END IF;

    -- ── 3. Validar stock suficiente en origen ANTES de ejecutar ──
    FOREACH v_linea IN ARRAY p_detalle LOOP

        SELECT COALESCE(stock_actual, 0)
        INTO v_stock_origen
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = p_id_bodega_origen
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        IF NOT FOUND THEN
            v_stock_origen := 0;
        END IF;

        IF v_stock_origen < v_linea.cantidad THEN
            RAISE EXCEPTION
                'Stock insuficiente para el artículo id=% lote=%. Disponible: %, solicitado: %',
                v_linea.id_articulo,
                COALESCE(v_linea.lote, 'Sin lote'),
                v_stock_origen,
                v_linea.cantidad;
        END IF;

    END LOOP;

    -- ── 4. Calcular consecutivo ───────────────────────────────────
    IF v_comprobante.consecutivo_actual = 0 THEN
        v_consecutivo := v_comprobante.consecutivo_inicial;
    ELSE
        v_consecutivo := v_comprobante.consecutivo_actual + 1;
    END IF;

    IF v_comprobante.consecutivo_fin IS NOT NULL
       AND v_consecutivo > v_comprobante.consecutivo_fin THEN
        RAISE EXCEPTION 'El comprobante "%" ha agotado su rango de consecutivos (máx: %)',
            v_comprobante.nombre_comprobante,
            v_comprobante.consecutivo_fin;
    END IF;

    -- ── 5. Obtener estado APLICADO ────────────────────────────────
    SELECT id INTO v_id_estado_aplicado
    FROM comprobantes.ref_estado_comprobante
    WHERE nombre = 'APLICADO';

    -- ── 6. Registrar en mov_gestion_comprobantes ──────────────────
    INSERT INTO comprobantes.mov_gestion_comprobantes (
        id_comprobante,
        id_estado_comprobante,
        autorizado,
        documento_soporte,
        fecha_movimiento,
        fecha_creacion,
        valor_impuesto,
        valor_total_comprobante,
        consecutivo_comprobante,
        prefijo_comprobante,
        observaciones,
        created_at,
        created_by,
        id_empresa,
        id_sede
    ) VALUES (
        p_id_comprobante,
        v_id_estado_aplicado,
        TRUE,
        p_referencia_externa,
        p_fecha_movimiento,
        NOW(),
        0,  -- traslados no tienen impuesto
        0,  -- valor calculado por costo promedio, no precio de venta
        v_consecutivo,
        v_comprobante.prefijo_comprobante,
        p_observaciones,
        NOW(),
        p_created_by,
        p_id_empresa,
        p_id_sede
    )
    RETURNING id INTO v_id_mov_comprobante;

    -- ── 7. Actualizar consecutivo ─────────────────────────────────
    UPDATE comprobantes.cfg_comprobante
    SET consecutivo_actual = v_consecutivo,
        updated_at         = NOW(),
        updated_by         = p_created_by
    WHERE id = p_id_comprobante;

    -- ── 8. Obtener id tipo Traslados ──────────────────────────────
    SELECT id INTO v_id_tipo_traslado
    FROM inventario.ref_tipo_movimiento
    WHERE nombre = 'Traslados';

    -- ── 9. Registrar mov_movimientos ──────────────────────────────
    INSERT INTO inventario.mov_movimientos (
        id_tipo_movimiento,
        id_estado,
        fecha_movimiento,
        id_bodega_origen,
        id_bodega_destino,
        id_mov_comprobante,
        referencia_externa,
        observaciones,
        id_empresa,
        id_sede,
        created_at,
        created_by
    ) VALUES (
        v_id_tipo_traslado,
        (SELECT id FROM inventario.ref_estado_movimiento WHERE codigo = '04'),
        p_fecha_movimiento,
        p_id_bodega_origen,
        p_id_bodega_destino,
        v_id_mov_comprobante,
        p_referencia_externa,
        p_observaciones,
        p_id_empresa,
        p_id_sede,
        NOW(),
        p_created_by
    )
    RETURNING id INTO v_id_movimiento;

    -- ── 10. Procesar cada línea ───────────────────────────────────
    FOREACH v_linea IN ARRAY p_detalle LOOP

        -- 10.1 Leer costo y stock en origen
        SELECT stock_actual, COALESCE(costo_promedio, 0)
        INTO v_stock_origen, v_costo_origen
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = p_id_bodega_origen
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        -- 10.2 Restar en bodega origen
        UPDATE inventario.inv_existencias_articulos
        SET stock_actual = stock_actual - v_linea.cantidad,
            updated_at   = NOW()
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = p_id_bodega_origen
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        -- 10.3 Leer stock y costo actuales en destino
        SELECT COALESCE(stock_actual, 0), COALESCE(costo_promedio, 0)
        INTO v_stock_destino, v_costo_destino
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = p_id_bodega_destino
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        IF NOT FOUND THEN
            v_stock_destino := 0;
            v_costo_destino := 0;
        END IF;

        -- 10.4 Calcular nuevo costo promedio ponderado en destino
        -- usando el costo de origen como valor de entrada
        IF (v_stock_destino + v_linea.cantidad) > 0 THEN
            v_nuevo_costo_destino := (
                (v_stock_destino * v_costo_destino)
                + (v_linea.cantidad * v_costo_origen)
            ) / (v_stock_destino + v_linea.cantidad);
        ELSE
            v_nuevo_costo_destino := v_costo_origen;
        END IF;

        -- 10.5 Upsert en bodega destino
        INSERT INTO inventario.inv_existencias_articulos (
            id_articulo,
            id_bodega,
            lote,
            stock_actual,
            costo_promedio,
            id_empresa,
            updated_at
        ) VALUES (
            v_linea.id_articulo,
            p_id_bodega_destino,
            v_linea.lote,
            v_linea.cantidad,
            v_nuevo_costo_destino,
            p_id_empresa,
            NOW()
        )
        ON CONFLICT ON CONSTRAINT uix_existencias_articulo
        DO UPDATE SET
            stock_actual   = inventario.inv_existencias_articulos.stock_actual + v_linea.cantidad,
            costo_promedio = v_nuevo_costo_destino,
            updated_at     = NOW();

        -- 10.6 Insertar detalle del movimiento
        -- valor_unitario = costo_origen (el costo al que sale de la bodega origen)
        -- valor_impuesto = 0 (traslados no generan impuesto)
        INSERT INTO inventario.mov_movimientos_detalle (
            id_movimiento,
            id_articulo,
            lote,
            cantidad,
            valor_unitario,
            valor_impuesto,
            observacion,
            created_at,
            created_by
        ) VALUES (
            v_id_movimiento,
            v_linea.id_articulo,
            v_linea.lote,
            v_linea.cantidad,
            v_costo_origen,   -- costo promedio de bodega origen al momento del traslado
            0,                -- traslados no tienen impuesto
            v_linea.observacion,
            NOW(),
            p_created_by
        );

        -- 10.7 Movimientos contables si aplica
        IF v_comprobante.aplica_mov_contable = TRUE THEN

            v_valor_linea := ROUND(v_linea.cantidad * v_costo_origen, 2);

            -- Cuenta inventario bodega ORIGEN
            SELECT cc.id_cuenta_contable INTO v_id_cuenta_inventario_origen
            FROM contabilidad.cfg_cuentas_grupo_articulos cc
            INNER JOIN inventario.cfg_articulos a ON a.id_tipo_articulo = cc.id_tipo_articulo
            WHERE a.id            = v_linea.id_articulo
              AND cc.id_bodega    = p_id_bodega_origen
              AND cc.id_empresa   = p_id_empresa
              AND cc.id_concepto_articulo = (
                  SELECT id FROM contabilidad.ref_conceptos_articulos
                  WHERE nombre = 'Cuenta Inventario'
              )
            LIMIT 1;

            -- Cuenta inventario bodega DESTINO
            SELECT cc.id_cuenta_contable INTO v_id_cuenta_inventario_destino
            FROM contabilidad.cfg_cuentas_grupo_articulos cc
            INNER JOIN inventario.cfg_articulos a ON a.id_tipo_articulo = cc.id_tipo_articulo
            WHERE a.id            = v_linea.id_articulo
              AND cc.id_bodega    = p_id_bodega_destino
              AND cc.id_empresa   = p_id_empresa
              AND cc.id_concepto_articulo = (
                  SELECT id FROM contabilidad.ref_conceptos_articulos
                  WHERE nombre = 'Cuenta Inventario'
              )
            LIMIT 1;

            IF v_id_cuenta_inventario_origen IS NOT NULL THEN
                -- CRÉDITO en origen (sale mercancía)
                INSERT INTO contabilidad.mov_movimientos_contables (
                    id_mov_comprobante, id_cuenta_contable, naturaleza, valor
                ) VALUES (
                    v_id_mov_comprobante, v_id_cuenta_inventario_origen, 'C', v_valor_linea
                );
            END IF;

            IF v_id_cuenta_inventario_destino IS NOT NULL THEN
                -- DÉBITO en destino (entra mercancía)
                INSERT INTO contabilidad.mov_movimientos_contables (
                    id_mov_comprobante, id_cuenta_contable, naturaleza, valor
                ) VALUES (
                    v_id_mov_comprobante, v_id_cuenta_inventario_destino, 'D', v_valor_linea
                );
            END IF;

        END IF;

    END LOOP;

    -- ── 11. Retornar resultado ────────────────────────────────────
    RETURN QUERY
    SELECT
        v_id_movimiento,
        v_id_mov_comprobante,
        v_consecutivo,
        v_comprobante.prefijo_comprobante;

END;
$$;

COMMENT ON FUNCTION inventario.fn_registrar_traslado IS
'Registra un traslado entre bodegas en una sola transacción.
Valida stock suficiente en origen ANTES de ejecutar.
Resta en bodega origen, suma en destino con costo promedio ponderado.
Genera asientos contables si aplica_mov_contable = true.
ROLLBACK automático si cualquier paso falla.';

-- ─── Función de anulación de traslado ────────────────────────
CREATE OR REPLACE FUNCTION inventario.fn_anular_traslado(
    p_id_movimiento INT,
    p_motivo        TEXT,
    p_anulado_por   INT,
    p_id_empresa    INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    v_movimiento             inventario.mov_movimientos%ROWTYPE;
    v_comprobante            comprobantes.cfg_comprobante%ROWTYPE;
    v_id_mov_comp            INT;
    v_id_estado_anulado_inv  INT;
    v_id_estado_anulado_comp INT;
    v_linea                  inventario.mov_movimientos_detalle%ROWTYPE;
    v_costo_origen           NUMERIC(14,4);
    v_stock_destino          NUMERIC(14,4);
    v_nuevo_costo_origen     NUMERIC(14,4);
    v_stock_origen_actual    NUMERIC(14,4);
    v_costo_origen_actual    NUMERIC(14,4);
BEGIN

    -- ── 1. Leer y bloquear movimiento ─────────────────────────────
    SELECT * INTO v_movimiento
    FROM inventario.mov_movimientos
    WHERE id = p_id_movimiento
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'No se encontró el movimiento con id=%', p_id_movimiento;
    END IF;

    -- ── 2. Verificar que es un Traslado ───────────────────────────
    IF (SELECT nombre FROM inventario.ref_tipo_movimiento
        WHERE id = v_movimiento.id_tipo_movimiento) != 'Traslados' THEN
        RAISE EXCEPTION 'Solo se pueden anular movimientos de tipo Traslado';
    END IF;

    -- ── 3. Verificar que no esté ya anulado ───────────────────────
    IF (SELECT codigo FROM inventario.ref_estado_movimiento
        WHERE id = v_movimiento.id_estado) = '05' THEN
        RAISE EXCEPTION 'Este movimiento ya fue anulado';
    END IF;

    -- ── 4. Verificar que el comprobante permite anulación ─────────
    v_id_mov_comp := v_movimiento.id_mov_comprobante;

    SELECT c.* INTO v_comprobante
    FROM comprobantes.cfg_comprobante c
    INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id_comprobante = c.id
    WHERE mg.id = v_id_mov_comp
    FOR UPDATE;

    IF v_comprobante.permite_anulacion = FALSE THEN
        RAISE EXCEPTION 'El comprobante "%" no permite anulación',
            v_comprobante.nombre_comprobante;
    END IF;

    -- ── 5. Obtener IDs de estado Anulado ──────────────────────────
    SELECT id INTO v_id_estado_anulado_inv
    FROM inventario.ref_estado_movimiento WHERE codigo = '05';

    SELECT id INTO v_id_estado_anulado_comp
    FROM comprobantes.ref_estado_comprobante WHERE nombre = 'ANULADO';

    -- ── 6. Revertir existencias por cada línea ────────────────────
    FOR v_linea IN
        SELECT * FROM inventario.mov_movimientos_detalle
        WHERE id_movimiento = p_id_movimiento
    LOOP
        -- Leer costo en destino (el que se usó para el traslado)
        SELECT COALESCE(costo_promedio, 0)
        INTO v_costo_origen
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = v_movimiento.id_bodega_destino
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        -- 6.1 Restar en bodega DESTINO (devolver lo trasladado)
        UPDATE inventario.inv_existencias_articulos
        SET stock_actual = stock_actual - v_linea.cantidad,
            updated_at   = NOW()
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = v_movimiento.id_bodega_destino
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        -- 6.2 Sumar en bodega ORIGEN (recuperar lo que salió)
        SELECT COALESCE(stock_actual, 0), COALESCE(costo_promedio, 0)
        INTO v_stock_origen_actual, v_costo_origen_actual
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = v_movimiento.id_bodega_origen
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        IF NOT FOUND THEN
            v_stock_origen_actual := 0;
            v_costo_origen_actual := 0;
        END IF;

        -- Recalcular costo promedio en origen al reincorporar el stock
        IF (v_stock_origen_actual + v_linea.cantidad) > 0 THEN
            v_nuevo_costo_origen := (
                (v_stock_origen_actual * v_costo_origen_actual)
                + (v_linea.cantidad * v_costo_origen)
            ) / (v_stock_origen_actual + v_linea.cantidad);
        ELSE
            v_nuevo_costo_origen := v_costo_origen;
        END IF;

        INSERT INTO inventario.inv_existencias_articulos (
            id_articulo, id_bodega, lote,
            stock_actual, costo_promedio, id_empresa, updated_at
        ) VALUES (
            v_linea.id_articulo,
            v_movimiento.id_bodega_origen,
            v_linea.lote,
            v_linea.cantidad,
            v_nuevo_costo_origen,
            p_id_empresa,
            NOW()
        )
        ON CONFLICT ON CONSTRAINT uix_existencias_articulo
        DO UPDATE SET
            stock_actual   = inventario.inv_existencias_articulos.stock_actual + v_linea.cantidad,
            costo_promedio = v_nuevo_costo_origen,
            updated_at     = NOW();

    END LOOP;

    -- ── 7. Actualizar estado del movimiento ───────────────────────
    UPDATE inventario.mov_movimientos
    SET id_estado        = v_id_estado_anulado_inv,
        anulado_por      = p_anulado_por,
        anulado_at       = NOW(),
        motivo_anulacion = p_motivo,
        updated_at       = NOW(),
        updated_by       = p_anulado_por
    WHERE id = p_id_movimiento;

    -- ── 8. Actualizar estado del comprobante ──────────────────────
    UPDATE comprobantes.mov_gestion_comprobantes
    SET id_estado_comprobante = v_id_estado_anulado_comp,
        id_user_anulo         = p_anulado_por,
        fecha_anulacion       = NOW(),
        motivo_anulacion      = p_motivo,
        updated_at            = NOW(),
        updated_by            = p_anulado_por
    WHERE id = v_id_mov_comp;

    -- ── 9. Revertir asientos contables si aplica ──────────────────
    IF v_comprobante.aplica_mov_contable = TRUE THEN
        DELETE FROM contabilidad.mov_movimientos_contables
        WHERE id_mov_comprobante = v_id_mov_comp;
    END IF;

END;
$$;

COMMENT ON FUNCTION inventario.fn_anular_traslado IS
'Anula un traslado revirtiendo existencias: resta en destino, suma en origen.
Recalcula costo promedio ponderado en origen al reincorporar el stock.
Elimina asientos contables si aplica_mov_contable = true.
ROLLBACK automático si cualquier paso falla.';