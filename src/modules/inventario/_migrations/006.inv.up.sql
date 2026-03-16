-- ============================================================
-- 008_inv_up.sql
-- Tabla general ref_subtipo_movimiento asociada a ref_tipo_movimiento
-- Reemplaza ref_tipo_baja con un esquema más robusto y reutilizable
-- ============================================================

-- ─── 1. Tabla referencial general ────────────────────────────
CREATE TABLE IF NOT EXISTS inventario.ref_subtipo_movimiento (
    id                  SERIAL PRIMARY KEY,
    id_tipo_movimiento  INT          NOT NULL
        REFERENCES inventario.ref_tipo_movimiento(id),
    codigo              VARCHAR(10)  NOT NULL,
    nombre              VARCHAR(150) NOT NULL,
    descripcion         TEXT         NULL,
    is_active           BOOLEAN      NOT NULL DEFAULT TRUE,
    CONSTRAINT uq_subtipo_movimiento UNIQUE (id_tipo_movimiento, codigo)
);

COMMENT ON TABLE inventario.ref_subtipo_movimiento IS
'Subtipos de movimiento asociados al tipo de movimiento padre.
Permite clasificar el origen o motivo de cualquier movimiento de inventario.
Ejemplos:
  - Entradas (id=2): Factura proveedor, Saldo inicial, Ajuste positivo
  - Bajas (id=3):    Vencimiento, Daño, Ajuste negativo, Venta, Robo, Devolución
  - Traslados (id=4): Traslado interno, Traslado entre sedes';

-- ─── 2. Subtipos para Entradas (id_tipo_movimiento = 2) ──────
INSERT INTO inventario.ref_subtipo_movimiento
    (id_tipo_movimiento, codigo, nombre, descripcion)
SELECT t.id, v.codigo, v.nombre, v.descripcion
FROM inventario.ref_tipo_movimiento t,
(VALUES
    ('FAC',  'Factura de proveedor',   'Entrada por compra a proveedor con factura'),
    ('SI',   'Saldo inicial',          'Cargue inicial de inventario al sistema'),
    ('AJP',  'Ajuste positivo',        'Ajuste de inventario al alza por conteo físico'),
    ('RET',  'Devolución de cliente',  'Retorno de mercancía por parte del cliente'),
    ('OTE',  'Otro tipo de entrada',   NULL)
) AS v(codigo, nombre, descripcion)
WHERE t.nombre = 'Entradas'
ON CONFLICT (id_tipo_movimiento, codigo) DO NOTHING;

-- ─── 3. Subtipos para Bajas (id_tipo_movimiento = 3) ─────────
INSERT INTO inventario.ref_subtipo_movimiento
    (id_tipo_movimiento, codigo, nombre, descripcion)
SELECT t.id, v.codigo, v.nombre, v.descripcion
FROM inventario.ref_tipo_movimiento t,
(VALUES
    ('VEN',  'Vencimiento',            'Artículos vencidos o con fecha expirada'),
    ('DAN',  'Daño o deterioro',       'Artículos dañados, rotos o en mal estado'),
    ('AJN',  'Ajuste negativo',        'Ajuste de inventario a la baja por conteo físico'),
    ('VTA',  'Venta directa',          'Salida por venta sin documento de despacho'),
    ('ROB',  'Robo o pérdida',         'Artículos perdidos o robados'),
    ('DEVP', 'Devolución a proveedor', 'Retorno de mercancía al proveedor'),
    ('OTB',  'Otro motivo de baja',    NULL)
) AS v(codigo, nombre, descripcion)
WHERE t.nombre = 'Bajas'
ON CONFLICT (id_tipo_movimiento, codigo) DO NOTHING;

-- ─── 4. Subtipos para Traslados (id_tipo_movimiento = 4) ─────
INSERT INTO inventario.ref_subtipo_movimiento
    (id_tipo_movimiento, codigo, nombre, descripcion)
SELECT t.id, v.codigo, v.nombre, v.descripcion
FROM inventario.ref_tipo_movimiento t,
(VALUES
    ('TI',   'Traslado interno',       'Movimiento entre bodegas de la misma sede'),
    ('TS',   'Traslado entre sedes',   'Movimiento entre sedes distintas'),
    ('TC',   'Traslado por consigna',  'Envío en consignación'),
    ('OTT',  'Otro tipo de traslado',  NULL)
) AS v(codigo, nombre, descripcion)
WHERE t.nombre = 'Traslados'
ON CONFLICT (id_tipo_movimiento, codigo) DO NOTHING;

-- ─── 5. Agregar id_subtipo_movimiento a mov_movimientos ───────
-- Esta columna permite clasificar el MOTIVO u ORIGEN del movimiento
-- aplica para todos los tipos: entradas, bajas, traslados, etc.
ALTER TABLE inventario.mov_movimientos
    ADD COLUMN IF NOT EXISTS id_subtipo_movimiento INT NULL
        REFERENCES inventario.ref_subtipo_movimiento(id);

COMMENT ON COLUMN inventario.mov_movimientos.id_subtipo_movimiento IS
'Clasificación del motivo u origen del movimiento.
Para entradas: Factura, Saldo inicial, Ajuste positivo...
Para bajas: Vencimiento, Daño, Ajuste negativo, Venta...
Para traslados: Interno, Entre sedes...';

-- ─── 6. Eliminar columna id_tipo_baja de mov_movimientos_detalle ──
-- El subtipo va en el encabezado (mov_movimientos), no por línea,
-- excepto cuando líneas del mismo movimiento tienen motivos distintos.
-- Para bajas mantenemos la posibilidad de especificarlo por línea también.
ALTER TABLE inventario.mov_movimientos_detalle
    ADD COLUMN IF NOT EXISTS id_subtipo_movimiento INT NULL
        REFERENCES inventario.ref_subtipo_movimiento(id);

COMMENT ON COLUMN inventario.mov_movimientos_detalle.id_subtipo_movimiento IS
'Subtipo específico de esta línea. Permite que en una misma baja
un artículo sea por vencimiento y otro por daño.
Si es NULL se hereda el subtipo del encabezado (mov_movimientos).';

-- ─── 7. Actualizar tipo compuesto t_baja_detalle ─────────────
-- Reemplazar id_tipo_baja por id_subtipo_movimiento
DROP TYPE IF EXISTS inventario.t_baja_detalle CASCADE;

CREATE TYPE inventario.t_baja_detalle AS (
    id_articulo           INT,
    lote                  VARCHAR(100),
    cantidad              NUMERIC(14,4),
    id_subtipo_movimiento INT,     -- subtipo específico de esta línea (ej: Vencimiento)
    observacion           TEXT
);

-- ─── 8. Actualizar fn_registrar_baja ─────────────────────────
CREATE OR REPLACE FUNCTION inventario.fn_registrar_baja(
    p_id_comprobante          INT,
    p_id_bodega               INT,
    p_id_subtipo_movimiento   INT,     -- subtipo general del documento
    p_referencia_externa      VARCHAR(100),
    p_observaciones           TEXT,
    p_fecha_movimiento        DATE,
    p_created_by              INT,
    p_id_empresa              INT,
    p_id_sede                 INT,
    p_detalle                 inventario.t_baja_detalle[]
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
    v_id_tipo_baja_mov   INT;

    v_linea              inventario.t_baja_detalle;
    v_stock_actual       NUMERIC(14,4);
    v_costo_actual       NUMERIC(14,4);
    v_nuevo_stock        NUMERIC(14,4);
    v_nuevo_costo        NUMERIC(14,4);
    v_total_costo        NUMERIC(14,2) := 0;

    v_id_cuenta_inventario INT;
    v_id_cuenta_gasto      INT;
    v_valor_linea          NUMERIC(14,2);

    -- Subtipo efectivo de cada línea (línea o heredado del encabezado)
    v_subtipo_linea        INT;
BEGIN

    -- ── 1. Validar comprobante ────────────────────────────────────
    SELECT * INTO v_comprobante
    FROM comprobantes.cfg_comprobante
    WHERE id = p_id_comprobante FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Comprobante no encontrado (id=%)', p_id_comprobante;
    END IF;
    IF v_comprobante.is_active = FALSE THEN
        RAISE EXCEPTION 'El comprobante "%" está inactivo', v_comprobante.nombre_comprobante;
    END IF;

    -- ── 2. Validar stock antes de ejecutar ────────────────────────
    FOREACH v_linea IN ARRAY p_detalle LOOP
        SELECT COALESCE(stock_actual, 0)
        INTO v_stock_actual
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo
          AND id_bodega   = p_id_bodega
          AND id_empresa  = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        IF NOT FOUND THEN v_stock_actual := 0; END IF;

        IF v_stock_actual < v_linea.cantidad THEN
            RAISE EXCEPTION
                'Stock insuficiente para artículo id=% lote=%. Disponible: %, solicitado: %',
                v_linea.id_articulo, COALESCE(v_linea.lote, 'Sin lote'),
                v_stock_actual, v_linea.cantidad;
        END IF;
    END LOOP;

    -- ── 3. Consecutivo ────────────────────────────────────────────
    IF v_comprobante.consecutivo_actual = 0 THEN
        v_consecutivo := v_comprobante.consecutivo_inicial;
    ELSE
        v_consecutivo := v_comprobante.consecutivo_actual + 1;
    END IF;

    IF v_comprobante.consecutivo_fin IS NOT NULL
       AND v_consecutivo > v_comprobante.consecutivo_fin THEN
        RAISE EXCEPTION 'El comprobante "%" ha agotado su rango (máx: %)',
            v_comprobante.nombre_comprobante, v_comprobante.consecutivo_fin;
    END IF;

    -- ── 4. Costo total ────────────────────────────────────────────
    FOREACH v_linea IN ARRAY p_detalle LOOP
        SELECT COALESCE(costo_promedio, 0) INTO v_costo_actual
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo AND id_bodega = p_id_bodega
          AND id_empresa = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');
        IF NOT FOUND THEN v_costo_actual := 0; END IF;
        v_total_costo := v_total_costo + ROUND(v_linea.cantidad * v_costo_actual, 2);
    END LOOP;

    -- ── 5. Estado APLICADO ────────────────────────────────────────
    SELECT id INTO v_id_estado_aplicado
    FROM comprobantes.ref_estado_comprobante WHERE nombre = 'APLICADO';

    -- ── 6. mov_gestion_comprobantes ───────────────────────────────
    INSERT INTO comprobantes.mov_gestion_comprobantes (
        id_comprobante, id_estado_comprobante, autorizado,
        documento_soporte, fecha_movimiento, fecha_creacion,
        valor_impuesto, valor_total_comprobante,
        consecutivo_comprobante, prefijo_comprobante,
        observaciones, created_at, created_by, id_empresa, id_sede
    ) VALUES (
        p_id_comprobante, v_id_estado_aplicado, TRUE,
        p_referencia_externa, p_fecha_movimiento, NOW(),
        0, v_total_costo, v_consecutivo, v_comprobante.prefijo_comprobante,
        p_observaciones, NOW(), p_created_by, p_id_empresa, p_id_sede
    )
    RETURNING id INTO v_id_mov_comprobante;

    UPDATE comprobantes.cfg_comprobante
    SET consecutivo_actual = v_consecutivo, updated_at = NOW(), updated_by = p_created_by
    WHERE id = p_id_comprobante;

    -- ── 7. mov_movimientos ────────────────────────────────────────
    SELECT id INTO v_id_tipo_baja_mov
    FROM inventario.ref_tipo_movimiento WHERE nombre = 'Bajas';

    INSERT INTO inventario.mov_movimientos (
        id_tipo_movimiento, id_estado, id_subtipo_movimiento,
        fecha_movimiento, id_bodega_origen,
        id_mov_comprobante, referencia_externa, observaciones,
        id_empresa, id_sede, created_at, created_by
    ) VALUES (
        v_id_tipo_baja_mov,
        (SELECT id FROM inventario.ref_estado_movimiento WHERE codigo = '04'),
        p_id_subtipo_movimiento,                -- subtipo general del documento
        p_fecha_movimiento, p_id_bodega,
        v_id_mov_comprobante, p_referencia_externa, p_observaciones,
        p_id_empresa, p_id_sede, NOW(), p_created_by
    )
    RETURNING id INTO v_id_movimiento;

    -- ── 8. Procesar líneas ────────────────────────────────────────
    FOREACH v_linea IN ARRAY p_detalle LOOP

        SELECT stock_actual, COALESCE(costo_promedio, 0)
        INTO v_stock_actual, v_costo_actual
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo AND id_bodega = p_id_bodega
          AND id_empresa = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        v_nuevo_stock := v_stock_actual - v_linea.cantidad;
        v_nuevo_costo := CASE WHEN v_nuevo_stock > 0 THEN v_costo_actual ELSE 0 END;

        UPDATE inventario.inv_existencias_articulos
        SET stock_actual = v_nuevo_stock, costo_promedio = v_nuevo_costo, updated_at = NOW()
        WHERE id_articulo = v_linea.id_articulo AND id_bodega = p_id_bodega
          AND id_empresa = p_id_empresa
          AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');

        v_valor_linea := ROUND(v_linea.cantidad * v_costo_actual, 2);

        -- Subtipo efectivo: el de la línea si viene, sino el del encabezado
        v_subtipo_linea := COALESCE(v_linea.id_subtipo_movimiento, p_id_subtipo_movimiento);

        INSERT INTO inventario.mov_movimientos_detalle (
            id_movimiento, id_articulo, lote, cantidad,
            valor_unitario, valor_impuesto,
            id_subtipo_movimiento, observacion,
            created_at, created_by
        ) VALUES (
            v_id_movimiento, v_linea.id_articulo, v_linea.lote, v_linea.cantidad,
            v_costo_actual, 0,
            v_subtipo_linea, v_linea.observacion,
            NOW(), p_created_by
        );

        -- Contabilidad
        IF v_comprobante.aplica_mov_contable = TRUE THEN
            SELECT cc.id_cuenta_contable INTO v_id_cuenta_inventario
            FROM contabilidad.cfg_cuentas_grupo_articulos cc
            INNER JOIN inventario.cfg_articulos a ON a.id_tipo_articulo = cc.id_tipo_articulo
            WHERE a.id = v_linea.id_articulo AND cc.id_bodega = p_id_bodega
              AND cc.id_empresa = p_id_empresa
              AND cc.id_concepto_articulo = (SELECT id FROM contabilidad.ref_conceptos_articulos WHERE nombre = 'Cuenta Inventario')
            LIMIT 1;

            SELECT cc.id_cuenta_contable INTO v_id_cuenta_gasto
            FROM contabilidad.cfg_cuentas_grupo_articulos cc
            INNER JOIN inventario.cfg_articulos a ON a.id_tipo_articulo = cc.id_tipo_articulo
            WHERE a.id = v_linea.id_articulo AND cc.id_bodega = p_id_bodega
              AND cc.id_empresa = p_id_empresa
              AND cc.id_concepto_articulo = (SELECT id FROM contabilidad.ref_conceptos_articulos WHERE nombre = 'Cuenta Gasto')
            LIMIT 1;

            IF v_id_cuenta_inventario IS NOT NULL THEN
                INSERT INTO contabilidad.mov_movimientos_contables (id_mov_comprobante, id_cuenta_contable, naturaleza, valor)
                VALUES (v_id_mov_comprobante, v_id_cuenta_inventario, 'C', v_valor_linea);
            END IF;
            IF v_id_cuenta_gasto IS NOT NULL THEN
                INSERT INTO contabilidad.mov_movimientos_contables (id_mov_comprobante, id_cuenta_contable, naturaleza, valor)
                VALUES (v_id_mov_comprobante, v_id_cuenta_gasto, 'D', v_valor_linea);
            END IF;
        END IF;

    END LOOP;

    RETURN QUERY SELECT v_id_movimiento, v_id_mov_comprobante,
                        v_consecutivo, v_comprobante.prefijo_comprobante;
END;
$$;

-- ─── 9. Actualizar fn_anular_baja ────────────────────────────
CREATE OR REPLACE FUNCTION inventario.fn_anular_baja(
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
    v_stock_actual           NUMERIC(14,4);
    v_costo_actual           NUMERIC(14,4);
    v_nuevo_costo            NUMERIC(14,4);
BEGIN
    SELECT * INTO v_movimiento FROM inventario.mov_movimientos WHERE id = p_id_movimiento FOR UPDATE;
    IF NOT FOUND THEN RAISE EXCEPTION 'No se encontró el movimiento id=%', p_id_movimiento; END IF;

    IF (SELECT nombre FROM inventario.ref_tipo_movimiento WHERE id = v_movimiento.id_tipo_movimiento) != 'Bajas' THEN
        RAISE EXCEPTION 'Solo se pueden anular movimientos de tipo Baja';
    END IF;
    IF (SELECT codigo FROM inventario.ref_estado_movimiento WHERE id = v_movimiento.id_estado) = '05' THEN
        RAISE EXCEPTION 'Este movimiento ya fue anulado';
    END IF;

    v_id_mov_comp := v_movimiento.id_mov_comprobante;
    SELECT c.* INTO v_comprobante
    FROM comprobantes.cfg_comprobante c
    INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id_comprobante = c.id
    WHERE mg.id = v_id_mov_comp FOR UPDATE;

    IF v_comprobante.permite_anulacion = FALSE THEN
        RAISE EXCEPTION 'El comprobante "%" no permite anulación', v_comprobante.nombre_comprobante;
    END IF;

    SELECT id INTO v_id_estado_anulado_inv FROM inventario.ref_estado_movimiento WHERE codigo = '05';
    SELECT id INTO v_id_estado_anulado_comp FROM comprobantes.ref_estado_comprobante WHERE nombre = 'ANULADO';

    FOR v_linea IN SELECT * FROM inventario.mov_movimientos_detalle WHERE id_movimiento = p_id_movimiento LOOP
        SELECT COALESCE(stock_actual, 0), COALESCE(costo_promedio, 0)
        INTO v_stock_actual, v_costo_actual
        FROM inventario.inv_existencias_articulos
        WHERE id_articulo = v_linea.id_articulo AND id_bodega = v_movimiento.id_bodega_origen
          AND id_empresa = p_id_empresa AND COALESCE(lote, '') = COALESCE(v_linea.lote, '');
        IF NOT FOUND THEN v_stock_actual := 0; v_costo_actual := 0; END IF;

        IF (v_stock_actual + v_linea.cantidad) > 0 THEN
            v_nuevo_costo := ((v_stock_actual * v_costo_actual) + (v_linea.cantidad * v_linea.valor_unitario))
                             / (v_stock_actual + v_linea.cantidad);
        ELSE
            v_nuevo_costo := v_linea.valor_unitario;
        END IF;

        INSERT INTO inventario.inv_existencias_articulos (
            id_articulo, id_bodega, lote, stock_actual, costo_promedio, id_empresa, updated_at
        ) VALUES (
            v_linea.id_articulo, v_movimiento.id_bodega_origen, v_linea.lote,
            v_linea.cantidad, v_nuevo_costo, p_id_empresa, NOW()
        )
        ON CONFLICT ON CONSTRAINT uix_existencias_articulo
        DO UPDATE SET
            stock_actual   = inventario.inv_existencias_articulos.stock_actual + v_linea.cantidad,
            costo_promedio = v_nuevo_costo,
            updated_at     = NOW();
    END LOOP;

    UPDATE inventario.mov_movimientos
    SET id_estado = v_id_estado_anulado_inv, anulado_por = p_anulado_por,
        anulado_at = NOW(), motivo_anulacion = p_motivo,
        updated_at = NOW(), updated_by = p_anulado_por
    WHERE id = p_id_movimiento;

    UPDATE comprobantes.mov_gestion_comprobantes
    SET id_estado_comprobante = v_id_estado_anulado_comp,
        id_user_anulo = p_anulado_por, fecha_anulacion = NOW(),
        motivo_anulacion = p_motivo, updated_at = NOW(), updated_by = p_anulado_por
    WHERE id = v_id_mov_comp;

    IF v_comprobante.aplica_mov_contable = TRUE THEN
        DELETE FROM contabilidad.mov_movimientos_contables WHERE id_mov_comprobante = v_id_mov_comp;
    END IF;
END;
$$;

COMMENT ON FUNCTION inventario.fn_registrar_baja IS
'Registra baja con subtipo en encabezado y por línea (hereda si no se especifica).
Bloquea si stock insuficiente. Genera asientos si aplica_mov_contable=true.';