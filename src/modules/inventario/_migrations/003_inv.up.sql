-- ============================================================
-- 003_inv_up.sql
-- Función principal para registro de entradas al inventario
-- Incluye: consecutivo, existencias, costo promedio y contabilidad
-- ============================================================

-- ─── Parámetro tipo para las líneas del detalle ──────────────
-- Se usa para pasar el arreglo de artículos desde el backend
CREATE TYPE inventario.t_entrada_detalle AS (
    id_articulo      INT,
    id_proveedor     INT,
    marca            VARCHAR(150),
    lote             VARCHAR(100),
    cantidad         NUMERIC(14,4),
    valor_unitario   NUMERIC(14,4),
    id_impuesto      INT,        -- NULL si no aplica impuesto
    observacion      TEXT
);

-- ─── Función principal ────────────────────────────────────────
CREATE OR REPLACE FUNCTION inventario.fn_registrar_entrada(
    p_id_comprobante    INT,
    p_id_bodega         INT,
    p_id_tercero        INT,        -- proveedor (puede ser NULL)
    p_referencia_externa VARCHAR(100),
    p_observaciones     TEXT,
    p_fecha_movimiento  DATE,
    p_created_by        INT,
    p_id_empresa        INT,
    p_id_sede           INT,
    p_detalle           inventario.t_entrada_detalle[]
)
RETURNS TABLE (
    id_movimiento           INT,
    id_mov_comprobante      INT,
    consecutivo_generado    INT,
    prefijo_comprobante     VARCHAR(10)
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_comprobante           comprobantes.cfg_comprobante%ROWTYPE;
    v_consecutivo           INT;
    v_id_mov_comprobante    INT;
    v_id_movimiento         INT;
    v_id_estado_aplicado    INT;
    v_id_tipo_entrada       INT;

    -- Para cada línea del detalle
    v_linea                 inventario.t_entrada_detalle;
    v_valor_unitario        NUMERIC(14,4);
    v_porcentaje_imp        NUMERIC(7,4)  := 0;
    v_valor_impuesto        NUMERIC(14,2) := 0;
    v_subtotal              NUMERIC(14,2) := 0;
    v_total_comprobante     NUMERIC(14,2) := 0;
    v_total_impuesto        NUMERIC(14,2) := 0;

    -- Para existencias y costo promedio
    v_stock_anterior        NUMERIC(14,4) := 0;
    v_costo_anterior        NUMERIC(14,4) := 0;
    v_nuevo_costo           NUMERIC(14,4) := 0;

    -- Para contabilidad
    v_id_cuenta_inventario  INT;
    v_id_cuenta_costo       INT;
    v_valor_total_linea     NUMERIC(14,2);

BEGIN

    -- ── 1. Bloquear y leer comprobante (FOR UPDATE evita concurrencia) ──
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

    -- ── 2. Calcular consecutivo ──────────────────────────────────────────
    IF v_comprobante.consecutivo_actual = 0 THEN
        v_consecutivo := v_comprobante.consecutivo_inicial;
    ELSE
        v_consecutivo := v_comprobante.consecutivo_actual + 1;
    END IF;

    -- Validar que no supere el consecutivo final
    IF v_comprobante.consecutivo_fin IS NOT NULL
       AND v_consecutivo > v_comprobante.consecutivo_fin THEN
        RAISE EXCEPTION 'El comprobante "%" ha agotado su rango de consecutivos (máx: %)',
            v_comprobante.nombre_comprobante,
            v_comprobante.consecutivo_fin;
    END IF;

    -- ── 3. Calcular totales del comprobante ──────────────────────────────
    FOREACH v_linea IN ARRAY p_detalle LOOP

        v_subtotal := v_linea.cantidad * v_linea.valor_unitario;

        -- Obtener porcentaje del impuesto si aplica
        IF v_linea.id_impuesto IS NOT NULL THEN
            SELECT porcentaje INTO v_porcentaje_imp
            FROM contabilidad.cfg_impuestos
            WHERE id = v_linea.id_impuesto;

            v_valor_impuesto := ROUND(v_subtotal * (v_porcentaje_imp / 100), 2);
        ELSE
            v_valor_impuesto := 0;
        END IF;

        v_total_impuesto    := v_total_impuesto + v_valor_impuesto;
        v_total_comprobante := v_total_comprobante + v_subtotal + v_valor_impuesto;

    END LOOP;

    -- ── 4. Obtener estado APLICADO ────────────────────────────────────────
    SELECT id INTO v_id_estado_aplicado
    FROM comprobantes.ref_estado_comprobante
    WHERE nombre = 'APLICADO';

    -- ── 5. Registrar en mov_gestion_comprobantes ──────────────────────────
    INSERT INTO comprobantes.mov_gestion_comprobantes (
        id_comprobante,
        id_estado_comprobante,
        autorizado,
        id_tercero,
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
        p_id_tercero,
        p_referencia_externa,
        p_fecha_movimiento,
        NOW(),
        v_total_impuesto,
        v_total_comprobante,
        v_consecutivo,
        v_comprobante.prefijo_comprobante,
        p_observaciones,
        NOW(),
        p_created_by,
        p_id_empresa,
        p_id_sede
    )
    RETURNING id INTO v_id_mov_comprobante;

    -- ── 6. Actualizar consecutivo en cfg_comprobante ──────────────────────
    UPDATE comprobantes.cfg_comprobante
    SET consecutivo_actual = v_consecutivo,
        updated_at         = NOW(),
        updated_by         = p_created_by
    WHERE id = p_id_comprobante;

    -- ── 7. Obtener id_tipo_movimiento = 'Entradas' ────────────────────────
    SELECT id INTO v_id_tipo_entrada
    FROM inventario.ref_tipo_movimiento
    WHERE nombre = 'Entradas';

    -- ── 8. Registrar en mov_movimientos ───────────────────────────────────
    INSERT INTO inventario.mov_movimientos (
        id_tipo_movimiento,
        id_estado,
        fecha_movimiento,
        id_bodega_destino,
        id_mov_comprobante,
        referencia_externa,
        observaciones,
        id_empresa,
        id_sede,
        created_at,
        created_by
    ) VALUES (
        v_id_tipo_entrada,
        (SELECT id FROM inventario.ref_estado_movimiento WHERE codigo = '04'), -- Ejecutado
        p_fecha_movimiento,
        p_id_bodega,
        v_id_mov_comprobante,
        p_referencia_externa,
        p_observaciones,
        p_id_empresa,
        p_id_sede,
        NOW(),
        p_created_by
    )
    RETURNING id INTO v_id_movimiento;

    -- ── 9. Procesar cada línea del detalle ────────────────────────────────
    FOREACH v_linea IN ARRAY p_detalle LOOP

        -- 9.1 Calcular impuesto de esta línea
        IF v_linea.id_impuesto IS NOT NULL THEN
            SELECT porcentaje INTO v_porcentaje_imp
            FROM contabilidad.cfg_impuestos
            WHERE id = v_linea.id_impuesto;
            v_valor_impuesto := ROUND(
                v_linea.cantidad * v_linea.valor_unitario * (v_porcentaje_imp / 100), 2
            );
        ELSE
            v_porcentaje_imp := 0;
            v_valor_impuesto := 0;
        END IF;

        -- 9.2 Insertar línea en mov_movimientos_detalle
        INSERT INTO inventario.mov_movimientos_detalle (
            id_movimiento,
            id_articulo,
            id_proveedor,
            marca,
            lote,
            cantidad,
            observacion,
            created_at,
            created_by
        ) VALUES (
            v_id_movimiento,
            v_linea.id_articulo,
            v_linea.id_proveedor,
            v_linea.marca,
            v_linea.lote,
            v_linea.cantidad,
            v_linea.observacion,
            NOW(),
            p_created_by
        );

        -- 9.3 Leer stock y costo promedio actuales
        SELECT stock_actual, COALESCE(costo_promedio, 0)
        INTO v_stock_anterior, v_costo_anterior
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = p_id_bodega
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        IF NOT FOUND THEN
            v_stock_anterior  := 0;
            v_costo_anterior  := 0;
        END IF;

        -- 9.4 Calcular nuevo costo promedio ponderado
        IF (v_stock_anterior + v_linea.cantidad) > 0 THEN
            v_nuevo_costo := (
                (v_stock_anterior * v_costo_anterior)
                + (v_linea.cantidad * v_linea.valor_unitario)
            ) / (v_stock_anterior + v_linea.cantidad);
        ELSE
            v_nuevo_costo := v_linea.valor_unitario;
        END IF;

        -- 9.5 Upsert en inv_existencias_articulos
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
            p_id_bodega,
            v_linea.lote,
            v_linea.cantidad,
            v_nuevo_costo,
            p_id_empresa,
            NOW()
        )
        ON CONFLICT ON CONSTRAINT uix_existencias_articulo
        DO UPDATE SET
            stock_actual   = inventario.inv_existencias_articulos.stock_actual + v_linea.cantidad,
            costo_promedio = v_nuevo_costo,
            updated_at     = NOW();

        -- ── 10. Movimientos contables (solo si el comprobante lo requiere) ──
        IF v_comprobante.aplica_mov_contable = TRUE THEN

            v_valor_total_linea := ROUND(
                v_linea.cantidad * v_linea.valor_unitario + v_valor_impuesto, 2
            );

            -- Buscar cuenta de inventario para el grupo del artículo
            SELECT cc.id_cuenta_contable INTO v_id_cuenta_inventario
            FROM contabilidad.cfg_cuentas_grupo_articulos cc
            INNER JOIN inventario.cfg_articulos a ON a.id_tipo_articulo = cc.id_tipo_articulo
            WHERE a.id            = v_linea.id_articulo
              AND cc.id_bodega    = p_id_bodega
              AND cc.id_empresa   = p_id_empresa
              AND cc.id_concepto_articulo = (
                  SELECT id FROM contabilidad.ref_conceptos_articulos
                  WHERE nombre = 'Cuenta Inventario'
              )
            LIMIT 1;

            -- Buscar cuenta de costo para el grupo del artículo
            SELECT cc.id_cuenta_contable INTO v_id_cuenta_costo
            FROM contabilidad.cfg_cuentas_grupo_articulos cc
            INNER JOIN inventario.cfg_articulos a ON a.id_tipo_articulo = cc.id_tipo_articulo
            WHERE a.id            = v_linea.id_articulo
              AND cc.id_bodega    = p_id_bodega
              AND cc.id_empresa   = p_id_empresa
              AND cc.id_concepto_articulo = (
                  SELECT id FROM contabilidad.ref_conceptos_articulos
                  WHERE nombre = 'Cuenta Costo'
              )
            LIMIT 1;

            -- Solo genera asientos si las cuentas están configuradas
            IF v_id_cuenta_inventario IS NOT NULL THEN

                -- DÉBITO → Cuenta Inventario (entra mercancía)
                INSERT INTO contabilidad.mov_movimientos_contables (
                    id_mov_comprobante,
                    id_cuenta_contable,
                    id_tercero,
                    naturaleza,
                    valor
                ) VALUES (
                    v_id_mov_comprobante,
                    v_id_cuenta_inventario,
                    p_id_tercero,
                    'D',
                    v_valor_total_linea
                );

            END IF;

            IF v_id_cuenta_costo IS NOT NULL THEN

                -- CRÉDITO → Cuenta Costo/Proveedor (sale el dinero)
                INSERT INTO contabilidad.mov_movimientos_contables (
                    id_mov_comprobante,
                    id_cuenta_contable,
                    id_tercero,
                    naturaleza,
                    valor
                ) VALUES (
                    v_id_mov_comprobante,
                    v_id_cuenta_costo,
                    p_id_tercero,
                    'C',
                    v_valor_total_linea
                );

            END IF;

        END IF; -- fin aplica_mov_contable

    END LOOP; -- fin FOREACH detalle

    -- ── 11. Retornar resultado ─────────────────────────────────────────────
    RETURN QUERY
    SELECT
        v_id_movimiento,
        v_id_mov_comprobante,
        v_consecutivo,
        v_comprobante.prefijo_comprobante;

END;
$$;

-- ─── Comentarios de documentación ────────────────────────────
COMMENT ON FUNCTION inventario.fn_registrar_entrada IS
'Registra una entrada al inventario dentro de una sola transacción.
Pasos: consecutivo → mov_gestion_comprobantes → mov_movimientos →
       mov_movimientos_detalle → inv_existencias_articulos (costo promedio) →
       mov_movimientos_contables (si aplica_mov_contable = true).
En caso de error en cualquier paso se hace ROLLBACK automático.';