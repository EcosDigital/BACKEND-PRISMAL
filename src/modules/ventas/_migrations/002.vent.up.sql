-- ============================================================
-- 002_fiados_up.sql
-- Trigger de saldo + funciones transaccionales
-- fn_registrar_fiado | fn_registrar_abono | fn_anular_mov_fiado
-- ============================================================

-- ─────────────────────────────────────────────────────────────
-- TRIGGER — actualizar saldo de la cuenta tras cada movimiento
-- BEFORE INSERT: calcula saldo_anterior, saldo_movimiento,
--               saldo_despues y actualiza cfg_cuentas_fiado
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION fiados.trg_fn_saldo()
RETURNS TRIGGER
LANGUAGE plpgsql AS
$$
DECLARE
    v_tipo_codigo  VARCHAR(20);
    v_delta        NUMERIC(14,2);
    v_saldo_ant    NUMERIC(14,2);
    v_saldo_nuevo  NUMERIC(14,2);
BEGIN
    -- Resolver el código del tipo de movimiento
    SELECT codigo INTO v_tipo_codigo
    FROM fiados.ref_tipo_movimiento
    WHERE id = NEW.id_tipo;

    -- Fiado suma al saldo, abono resta
    v_delta := CASE v_tipo_codigo WHEN '01' THEN NEW.valor ELSE -NEW.valor END;

    -- Leer saldo actual de la cuenta ANTES del movimiento
    SELECT saldo_actual INTO v_saldo_ant
    FROM fiados.cfg_cuentas_fiado
    WHERE id = NEW.id_cuenta;

    v_saldo_nuevo := v_saldo_ant + v_delta;

    -- Completar los tres campos de saldo en el movimiento
    NEW.saldo_anterior    := v_saldo_ant;
    NEW.saldo_movimiento  := v_delta;
    NEW.saldo_despues     := v_saldo_nuevo;

    -- Actualizar saldo en la cuenta
    UPDATE fiados.cfg_cuentas_fiado
    SET saldo_actual = v_saldo_nuevo,
        updated_at   = NOW()
    WHERE id = NEW.id_cuenta;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_saldo_fiado
    BEFORE INSERT ON fiados.mov_fiados
    FOR EACH ROW
    EXECUTE FUNCTION fiados.trg_fn_saldo();

COMMENT ON FUNCTION fiados.trg_fn_saldo IS
'Trigger BEFORE INSERT en mov_fiados.
Calcula saldo_anterior, saldo_movimiento y saldo_despues antes de insertar.
Actualiza saldo_actual en cfg_cuentas_fiado automáticamente.';


-- ─────────────────────────────────────────────────────────────
-- fn_registrar_fiado (CORREGIDA)
-- Registra una nueva deuda sobre la cuenta del tercero.
-- Si la cuenta no existe la crea. Si estaba cerrada la reabre.
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION fiados.fn_registrar_fiado(
    p_id_tercero        INT,
    p_id_comprobante    INT,
    p_valor             NUMERIC(14,2),
    p_observaciones     TEXT,
    p_fecha             DATE,
    p_created_by        INT,
    p_id_empresa        INT,
    p_id_sede           INT
)
RETURNS TABLE (
    id_movimiento       INT,
    id_cuenta           INT,
    id_mov_comprobante  INT,
    consecutivo         INT,
    prefijo             VARCHAR(10),
    saldo_resultante    NUMERIC(14,2)
)
LANGUAGE plpgsql AS
$$
DECLARE
    v_comp              comprobantes.cfg_comprobante%ROWTYPE;
    v_consecutivo       INT;
    v_id_mov_comp       INT;
    v_id_estado_comp    INT;
    v_id_cuenta         INT;
    v_id_estado_pend    INT;
    v_id_tipo_fiado     INT;
    v_saldo             NUMERIC(14,2);
    v_id_movimiento     INT;
    v_estado_cuenta     VARCHAR(20);
    v_limite            NUMERIC(14,2);
BEGIN
    -- ── 1. Validaciones básicas ───────────────────────────────
    IF p_valor IS NULL OR p_valor <= 0 THEN
        RAISE EXCEPTION 'El valor del fiado debe ser mayor a cero';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM configuracion.cfg_terceros WHERE id = p_id_tercero
    ) THEN
        RAISE EXCEPTION 'Tercero id=% no encontrado en cfg_terceros', p_id_tercero;
    END IF;

    -- ── 2. Comprobante: leer, bloquear y calcular consecutivo ─
    SELECT * INTO v_comp
    FROM comprobantes.cfg_comprobante
    WHERE id = p_id_comprobante
    FOR UPDATE;

    IF NOT FOUND OR v_comp.is_active = FALSE THEN
        RAISE EXCEPTION 'Comprobante id=% no encontrado o inactivo', p_id_comprobante;
    END IF;

    v_consecutivo := CASE
        WHEN v_comp.consecutivo_actual = 0 THEN v_comp.consecutivo_inicial
        ELSE v_comp.consecutivo_actual + 1
    END;

    IF v_comp.consecutivo_fin IS NOT NULL
       AND v_consecutivo > v_comp.consecutivo_fin THEN
        RAISE EXCEPTION 'Comprobante "%" agotó su rango de consecutivos (máx: %)',
            v_comp.nombre_comprobante, v_comp.consecutivo_fin;
    END IF;

    -- ── 3. Insertar comprobante ───────────────────────────────
    SELECT id INTO v_id_estado_comp
    FROM comprobantes.ref_estado_comprobante WHERE nombre = 'APLICADO';

    INSERT INTO comprobantes.mov_gestion_comprobantes (
        id_comprobante, id_estado_comprobante, autorizado,
        id_tercero, fecha_movimiento, fecha_creacion,
        valor_impuesto, valor_total_comprobante,
        consecutivo_comprobante, prefijo_comprobante,
        observaciones, created_at, created_by,
        id_empresa, id_sede
    ) VALUES (
        p_id_comprobante, v_id_estado_comp, TRUE,
        p_id_tercero, COALESCE(p_fecha, CURRENT_DATE), NOW(),
        0, p_valor,
        v_consecutivo, v_comp.prefijo_comprobante,
        p_observaciones, NOW(), p_created_by,
        p_id_empresa, p_id_sede
    )
    RETURNING id INTO v_id_mov_comp;

    UPDATE comprobantes.cfg_comprobante
    SET consecutivo_actual = v_consecutivo,
        updated_at = NOW(), updated_by = p_created_by
    WHERE id = p_id_comprobante;

    -- ── 4. Resolver o crear cuenta (sin bloque anidado extra) ─
    SELECT id, saldo_actual
    INTO v_id_cuenta, v_saldo
    FROM fiados.cfg_cuentas_fiado
    WHERE id_tercero = p_id_tercero
      AND id_empresa = p_id_empresa
    FOR UPDATE;

    -- Estado "Pendiente" para cuenta nueva o reapertura
    SELECT id INTO v_id_estado_pend
    FROM fiados.ref_estado_credito WHERE codigo = '02';

    IF NOT FOUND OR v_id_cuenta IS NULL THEN
        -- Cuenta nueva
        INSERT INTO fiados.cfg_cuentas_fiado (
            id_tercero, saldo_actual, limite_credito,
            id_estado, fecha_apertura,
            id_empresa, id_sede,
            created_at, created_by
        ) VALUES (
            p_id_tercero, 0, NULL,
            v_id_estado_pend,
            COALESCE(p_fecha, CURRENT_DATE),
            p_id_empresa, p_id_sede,
            NOW(), p_created_by
        )
        RETURNING id, saldo_actual INTO v_id_cuenta, v_saldo;
    ELSE
        -- Verificar estado de cuenta existente
        SELECT re.codigo
        INTO v_estado_cuenta
        FROM fiados.cfg_cuentas_fiado cf
        INNER JOIN fiados.ref_estado_credito re ON re.id = cf.id_estado
        WHERE cf.id = v_id_cuenta;

        IF v_estado_cuenta = 'CASTIGADA' THEN
            RAISE EXCEPTION
                'La cuenta del cliente está castigada y no acepta nuevos movimientos';
        END IF;
    END IF;

    -- Validar límite de crédito si existe (sin bloque anidado)
    SELECT cf.limite_credito INTO v_limite
    FROM fiados.cfg_cuentas_fiado cf
    WHERE cf.id = v_id_cuenta;

    IF v_limite IS NOT NULL AND (v_saldo + p_valor) > v_limite THEN
        RAISE EXCEPTION
            'El fiado de % supera el límite de crédito asignado (saldo: %, límite: %)',
            p_valor, v_saldo, v_limite;
    END IF;

    -- ── 5. Tipo de movimiento FIADO ───────────────────────────
    SELECT id INTO v_id_tipo_fiado
    FROM fiados.ref_tipo_movimiento WHERE codigo = '01';

    -- ── 6. Insertar movimiento (trigger calcula los saldos) ───
    INSERT INTO fiados.mov_fiados (
        id_cuenta, id_mov_comprobante, id_tipo,
        fecha_movimiento, valor, observaciones,
        saldo_anterior, saldo_movimiento, saldo_despues,
        id_empresa, id_sede,
        created_at, created_by
    ) VALUES (
        v_id_cuenta, v_id_mov_comp, v_id_tipo_fiado,
        COALESCE(p_fecha, CURRENT_DATE), p_valor, p_observaciones,
        0, 0, 0,   -- el trigger BEFORE los sobreescribe
        p_id_empresa, p_id_sede,
        NOW(), p_created_by
    )
    RETURNING id INTO v_id_movimiento;

    -- ── 7. Actualizar estado cuenta a Pendiente ───────────────
    UPDATE fiados.cfg_cuentas_fiado
    SET id_estado  = v_id_estado_pend,
        updated_at = NOW(),
        updated_by = p_created_by
    WHERE id = v_id_cuenta;

    -- ── 8. Saldo final ────────────────────────────────────────
    SELECT saldo_actual INTO v_saldo
    FROM fiados.cfg_cuentas_fiado WHERE id = v_id_cuenta;

    RETURN QUERY
    SELECT v_id_movimiento, v_id_cuenta, v_id_mov_comp,
           v_consecutivo, v_comp.prefijo_comprobante, v_saldo;
END;
$$;

COMMENT ON FUNCTION fiados.fn_registrar_fiado IS
'Registra una nueva deuda (fiado) en la cuenta del cliente.
Crea la cuenta si no existe. Si estaba castigada lanza excepción.
Genera comprobante FI y actualiza consecutivo.
El trigger BEFORE INSERT calcula saldo_anterior, saldo_movimiento y saldo_despues.
ROLLBACK automático si cualquier paso falla.';


-- ─────────────────────────────────────────────────────────────
-- fn_registrar_abono (corregida para consistencia)
-- Registra un pago parcial o total sobre la cuenta.
-- Si el saldo queda en 0 marca la cuenta como Al día.
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION fiados.fn_registrar_abono(
    p_id_cuenta         INT,
    p_id_comprobante    INT,
    p_valor             NUMERIC(14,2),
    p_observaciones     TEXT,
    p_fecha             DATE,
    p_created_by        INT,
    p_id_empresa        INT,
    p_id_sede           INT
)
RETURNS TABLE (
    id_movimiento       INT,
    id_mov_comprobante  INT,
    consecutivo         INT,
    prefijo             VARCHAR(10),
    saldo_resultante    NUMERIC(14,2),
    cuenta_saldada      BOOLEAN
)
LANGUAGE plpgsql AS
$$
DECLARE
    v_comp              comprobantes.cfg_comprobante%ROWTYPE;
    v_consecutivo       INT;
    v_id_mov_comp       INT;
    v_id_estado_comp    INT;
    v_id_tercero        INT;
    v_id_tipo_abono     INT;
    v_saldo_nuevo       NUMERIC(14,2);
    v_id_movimiento     INT;
    v_cuenta_saldada    BOOLEAN := FALSE;
    v_id_estado_al_dia  INT;
    v_id_estado_pend    INT;
BEGIN
    -- ── 1. Validaciones ───────────────────────────────────────
    IF p_valor IS NULL OR p_valor <= 0 THEN
        RAISE EXCEPTION 'El valor del abono debe ser mayor a cero';
    END IF;

    SELECT id_tercero INTO v_id_tercero
    FROM fiados.cfg_cuentas_fiado
    WHERE id = p_id_cuenta AND id_empresa = p_id_empresa
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Cuenta id=% no encontrada en empresa=%',
            p_id_cuenta, p_id_empresa;
    END IF;

    -- Verificar que la cuenta no esté castigada
    IF EXISTS (
        SELECT 1
        FROM fiados.cfg_cuentas_fiado cf
        INNER JOIN fiados.ref_estado_credito re ON re.id = cf.id_estado
        WHERE cf.id = p_id_cuenta AND re.codigo = '00'
    ) THEN
        RAISE EXCEPTION 'La cuenta id=% está castigada', p_id_cuenta;
    END IF;

    -- ── 2. Comprobante ────────────────────────────────────────
    SELECT * INTO v_comp
    FROM comprobantes.cfg_comprobante
    WHERE id = p_id_comprobante FOR UPDATE;

    IF NOT FOUND OR v_comp.is_active = FALSE THEN
        RAISE EXCEPTION 'Comprobante id=% no encontrado o inactivo', p_id_comprobante;
    END IF;

    v_consecutivo := CASE
        WHEN v_comp.consecutivo_actual = 0 THEN v_comp.consecutivo_inicial
        ELSE v_comp.consecutivo_actual + 1
    END;

    IF v_comp.consecutivo_fin IS NOT NULL
       AND v_consecutivo > v_comp.consecutivo_fin THEN
        RAISE EXCEPTION 'Comprobante "%" agotó su rango de consecutivos',
            v_comp.nombre_comprobante;
    END IF;

    SELECT id INTO v_id_estado_comp
    FROM comprobantes.ref_estado_comprobante WHERE nombre = 'APLICADO';

    INSERT INTO comprobantes.mov_gestion_comprobantes (
        id_comprobante, id_estado_comprobante, autorizado,
        id_tercero, fecha_movimiento, fecha_creacion,
        valor_impuesto, valor_total_comprobante,
        consecutivo_comprobante, prefijo_comprobante,
        observaciones, created_at, created_by,
        id_empresa, id_sede
    ) VALUES (
        p_id_comprobante, v_id_estado_comp, TRUE,
        v_id_tercero, COALESCE(p_fecha, CURRENT_DATE), NOW(),
        0, p_valor,
        v_consecutivo, v_comp.prefijo_comprobante,
        p_observaciones, NOW(), p_created_by,
        p_id_empresa, p_id_sede
    )
    RETURNING id INTO v_id_mov_comp;

    UPDATE comprobantes.cfg_comprobante
    SET consecutivo_actual = v_consecutivo,
        updated_at = NOW(), updated_by = p_created_by
    WHERE id = p_id_comprobante;

    -- ── 3. Tipo de movimiento ABONO ───────────────────────────
    SELECT id INTO v_id_tipo_abono
    FROM fiados.ref_tipo_movimiento WHERE codigo = '02';

    -- ── 4. Insertar movimiento (trigger calcula los saldos) ───
    INSERT INTO fiados.mov_fiados (
        id_cuenta, id_mov_comprobante, id_tipo,
        fecha_movimiento, valor, observaciones,
        saldo_anterior, saldo_movimiento, saldo_despues,
        id_empresa, id_sede,
        created_at, created_by
    ) VALUES (
        p_id_cuenta, v_id_mov_comp, v_id_tipo_abono,
        COALESCE(p_fecha, CURRENT_DATE), p_valor, p_observaciones,
        0, 0, 0,
        p_id_empresa, p_id_sede,
        NOW(), p_created_by
    )
    RETURNING id INTO v_id_movimiento;

    -- ── 5. Leer saldo final y actualizar estado ───────────────
    SELECT saldo_actual INTO v_saldo_nuevo
    FROM fiados.cfg_cuentas_fiado WHERE id = p_id_cuenta;

    SELECT id INTO v_id_estado_al_dia
    FROM fiados.ref_estado_credito WHERE codigo = '03';  -- Al día

    SELECT id INTO v_id_estado_pend
    FROM fiados.ref_estado_credito WHERE codigo = '02';  -- Pendiente

    IF v_saldo_nuevo <= 0 THEN
        -- Saldo saldado: estado Al día
        UPDATE fiados.cfg_cuentas_fiado
        SET id_estado    = v_id_estado_al_dia,
            fecha_cierre = CURRENT_DATE,
            updated_at   = NOW(),
            updated_by   = p_created_by
        WHERE id = p_id_cuenta;

        v_cuenta_saldada := TRUE;
    ELSE
        -- Abono parcial: puede seguir en Pendiente o Atrasado
        UPDATE fiados.cfg_cuentas_fiado
        SET id_estado  = v_id_estado_pend,
            updated_at = NOW(),
            updated_by = p_created_by
        WHERE id = p_id_cuenta
          AND id_estado = (SELECT id FROM fiados.ref_estado_credito WHERE codigo = '02');
    END IF;

    RETURN QUERY
    SELECT v_id_movimiento, v_id_mov_comp,
           v_consecutivo, v_comp.prefijo_comprobante,
           v_saldo_nuevo, v_cuenta_saldada;
END;
$$;

COMMENT ON FUNCTION fiados.fn_registrar_abono IS
'Registra un pago (abono) sobre la cuenta de fiado.
Genera comprobante AB. El trigger calcula saldo_anterior/saldo_movimiento/saldo_despues.
Si el saldo queda <= 0 marca el estado Al día (código 03) y registra fecha_cierre.
ROLLBACK automático si cualquier paso falla.';


-- ─────────────────────────────────────────────────────────────
-- fn_anular_mov_fiado (sin cambios)
-- Anula el comprobante del movimiento y revierte el saldo
-- mediante un movimiento compensatorio inverso.
-- ─────────────────────────────────────────────────────────────
CREATE OR REPLACE FUNCTION fiados.fn_anular_mov_fiado(
    p_id_movimiento  INT,
    p_motivo         TEXT,
    p_anulado_por    INT,
    p_id_empresa     INT
)
RETURNS VOID
LANGUAGE plpgsql AS
$$
DECLARE
    v_mov               fiados.mov_fiados%ROWTYPE;
    v_comp_orig         comprobantes.cfg_comprobante%ROWTYPE;
    v_id_estado_anulado INT;
    v_tipo_codigo       VARCHAR(20);
    v_tipo_inv_id       INT;
    v_consecutivo       INT;
    v_id_mov_comp_anul  INT;
    v_id_estado_comp    INT;
    v_id_tercero        INT;
BEGIN
    -- ── 1. Leer y validar el movimiento ───────────────────────
    SELECT * INTO v_mov
    FROM fiados.mov_fiados
    WHERE id = p_id_movimiento AND id_empresa = p_id_empresa
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Movimiento id=% no encontrado', p_id_movimiento;
    END IF;

    IF p_motivo IS NULL OR TRIM(p_motivo) = '' THEN
        RAISE EXCEPTION 'Debe indicar un motivo de anulación';
    END IF;

    -- Verificar que el comprobante no esté ya anulado
    IF (
        SELECT re.nombre
        FROM comprobantes.mov_gestion_comprobantes mg
        INNER JOIN comprobantes.ref_estado_comprobante re
            ON re.id = mg.id_estado_comprobante
        WHERE mg.id = v_mov.id_mov_comprobante
    ) = 'ANULADO' THEN
        RAISE EXCEPTION 'El movimiento id=% ya fue anulado', p_id_movimiento;
    END IF;

    -- ── 2. Verificar que el comprobante permite anulación ─────
    SELECT c.* INTO v_comp_orig
    FROM comprobantes.cfg_comprobante c
    INNER JOIN comprobantes.mov_gestion_comprobantes mg
        ON mg.id_comprobante = c.id
    WHERE mg.id = v_mov.id_mov_comprobante;

    IF v_comp_orig.permite_anulacion = FALSE THEN
        RAISE EXCEPTION 'El comprobante "%" no permite anulación',
            v_comp_orig.nombre_comprobante;
    END IF;

    -- ── 3. Anular el comprobante original ─────────────────────
    SELECT id INTO v_id_estado_anulado
    FROM comprobantes.ref_estado_comprobante WHERE nombre = 'ANULADO';

    UPDATE comprobantes.mov_gestion_comprobantes
    SET id_estado_comprobante = v_id_estado_anulado,
        id_user_anulo         = p_anulado_por,
        fecha_anulacion       = NOW(),
        motivo_anulacion      = TRIM(p_motivo),
        updated_at            = NOW(),
        updated_by            = p_anulado_por
    WHERE id = v_mov.id_mov_comprobante;

    -- ── 4. Revertir saldo con movimiento compensatorio ────────
    -- El tipo inverso: si era FIADO (01) → ABONO (02), y viceversa
    SELECT codigo INTO v_tipo_codigo
    FROM fiados.ref_tipo_movimiento WHERE id = v_mov.id_tipo;

    SELECT id INTO v_tipo_inv_id
    FROM fiados.ref_tipo_movimiento
    WHERE codigo = CASE v_tipo_codigo WHEN '01' THEN '02' ELSE '01' END;

    -- Generar nuevo comprobante de anulación sobre el mismo comprobante base
    v_consecutivo := v_comp_orig.consecutivo_actual + 1;

    SELECT id INTO v_id_estado_comp
    FROM comprobantes.ref_estado_comprobante WHERE nombre = 'APLICADO';

    SELECT id_tercero INTO v_id_tercero
    FROM fiados.cfg_cuentas_fiado WHERE id = v_mov.id_cuenta;

    INSERT INTO comprobantes.mov_gestion_comprobantes (
        id_comprobante, id_estado_comprobante, autorizado,
        id_tercero, fecha_movimiento, fecha_creacion,
        valor_impuesto, valor_total_comprobante,
        consecutivo_comprobante, prefijo_comprobante,
        observaciones, created_at, created_by,
        id_empresa, id_sede
    ) VALUES (
        v_comp_orig.id, v_id_estado_comp, TRUE,
        v_id_tercero, CURRENT_DATE, NOW(),
        0, v_mov.valor,
        v_consecutivo,
        v_comp_orig.prefijo_comprobante || '-AN',
        'ANULACIÓN: ' || TRIM(p_motivo),
        NOW(), p_anulado_por,
        p_id_empresa, v_mov.id_sede
    )
    RETURNING id INTO v_id_mov_comp_anul;

    UPDATE comprobantes.cfg_comprobante
    SET consecutivo_actual = v_consecutivo,
        updated_at = NOW(), updated_by = p_anulado_por
    WHERE id = v_comp_orig.id;

    -- Movimiento compensatorio (trigger revierte el saldo)
    INSERT INTO fiados.mov_fiados (
        id_cuenta, id_mov_comprobante, id_tipo,
        fecha_movimiento, valor, observaciones,
        saldo_anterior, saldo_movimiento, saldo_despues,
        id_empresa, id_sede,
        created_at, created_by
    ) VALUES (
        v_mov.id_cuenta, v_id_mov_comp_anul, v_tipo_inv_id,
        CURRENT_DATE, v_mov.valor,
        'ANULACIÓN MOV #' || p_id_movimiento || ': ' || TRIM(p_motivo),
        0, 0, 0,
        p_id_empresa, v_mov.id_sede,
        NOW(), p_anulado_por
    );

    -- ── 5. Marcar movimiento original como actualizado ────────
    UPDATE fiados.mov_fiados
    SET updated_at = NOW(),
        updated_by = p_anulado_por
    WHERE id = p_id_movimiento;

END;
$$;

COMMENT ON FUNCTION fiados.fn_anular_mov_fiado IS
'Anula un movimiento de fiado o abono.
Estrategia: anula el comprobante original e inserta un movimiento
compensatorio de tipo inverso (FIADO→ABONO o ABONO→FIADO).
El trigger BEFORE INSERT del compensatorio revierte el saldo automáticamente.
ROLLBACK automático si cualquier paso falla.';