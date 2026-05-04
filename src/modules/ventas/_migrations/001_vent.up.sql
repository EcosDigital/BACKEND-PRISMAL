-- ============================================================
-- 001_fiados_up.sql
-- Módulo de Fiados — Schema, tablas y comprobantes
-- Se apoya en configuracion.cfg_terceros para los clientes
-- ============================================================

CREATE SCHEMA IF NOT EXISTS fiados;

-- ─────────────────────────────────────────────────────────────
-- COMPROBANTES del módulo
-- Sigue el mismo patrón de 002_inv_up.sql
-- cod_modulo = MD-008 (Fiados), tipos: 1=Fiado, 2=Abono
-- ─────────────────────────────────────────────────────────────

INSERT INTO comprobantes.cfg_comprobante (
    id_modulo,
    id_tipo_operacion,
    nombre_comprobante,
    prefijo_comprobante,
    consecutivo_inicial,
    consecutivo_actual,
    consecutivo_fin,
    permite_anulacion,
    fecha_inicio,
    fecha_final,
    is_active,
    created_at,
    created_by,
    id_empresa,
    id_sede
)
VALUES
(
    8,
    1,
    'Registro de Fiado',
    'FI', 1, 0, 9999999,
    true, NOW(), '2030-12-31',
    true, NOW(), 1, 1, 1
),
(
    8,
    2,
    'Abono a Fiado',
    'AB', 1, 0, 9999999,
    true, NOW(), '2030-12-31',
    true, NOW(), 1, 1, 1
)
ON CONFLICT (prefijo_comprobante, id_modulo) DO NOTHING;

-- ─────────────────────────────────────────────────────────────
-- TABLAS DE REFERENCIA
-- estado de fiado (atrasado, pendiente, al dia) 
-- tipo movimiento (fiado, abono)
-- ───────────────────────────────────────────────────────────── 

CREATE TABLE IF NOT EXISTS fiados.ref_estado_credito (
    id  SERIAL PRIMARY KEY,
    codigo VARCHAR(20)  NOT NULL UNIQUE,
    nombre   VARCHAR(100) NOT NULL
);

INSERT INTO fiados.ref_estado_credito (codigo, nombre)
VALUES
    ('01', 'Atrasado'),
    ('02', 'Pendiente'),
    ('03', 'Al día')
ON CONFLICT (codigo) DO NOTHING;

CREATE TABLE IF NOT EXISTS fiados.ref_tipo_movimiento (
    id  SERIAL PRIMARY KEY,
    codigo VARCHAR(20)  NOT NULL UNIQUE,
    nombre   VARCHAR(100) NOT NULL
);

INSERT INTO fiados.ref_tipo_movimiento (codigo, nombre)
VALUES
    ('01', 'Fiado'),
    ('02', 'Abono')
ON CONFLICT (codigo) DO NOTHING;


-- ─────────────────────────────────────────────────────────────
-- CUENTA DE FIADO
-- Una cuenta por cliente (tercero). Mientras estado = 'ABIERTA'
-- el cliente puede seguir acumulando deuda o hacer abonos.
-- El tercero ya contiene nombre, cédula, teléfono y dirección.
-- ───────────────────────────────────────────────────────────── 

CREATE TABLE IF NOT EXISTS fiados.cfg_cuentas_fiado (
    id  SERIAL PRIMARY KEY,
    id_tercero  INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    saldo_actual    NUMERIC(14,2)  NOT NULL DEFAULT 0,
    -- Límite de crédito que el tendero asigna (NULL = sin límite)
    limite_credito  NUMERIC(14,2)  NULL,
    -- ABIERTA: puede recibir movimientos
    -- CERRADA: saldo en 0, cerrada manualmente o por abono total
    -- CASTIGADA: deuda incobrable, no acepta más movimientos
    id_estado  INT NOT NULL REFERENCES fiados.ref_estado_credito,
    observaciones   TEXT           NULL,
    -- Fechas de apertura y cierre
    fecha_apertura  DATE           NOT NULL DEFAULT CURRENT_DATE,
    fecha_cierre    DATE           NULL,
    -- Multi-empresa
    id_empresa      INT            NOT NULL,
    id_sede         INT            NULL,
    -- Auditoría
    created_at      TIMESTAMP      NOT NULL DEFAULT NOW(),
    created_by      INT            NOT NULL,
    updated_at      TIMESTAMP      NULL,
    updated_by      INT            NULL,
    -- Un mismo tercero solo puede tener UNA cuenta abierta por empresa
    CONSTRAINT uq_cuenta_abierta_tercero
        UNIQUE (id_tercero, id_empresa)
);

COMMENT ON TABLE fiados.cfg_cuentas_fiado IS
'Cuenta de crédito informal (fiado) asociada a un tercero (cliente).
El tercero se registra en configuracion.cfg_terceros con clase CLIENTE.
saldo_actual se mantiene actualizado por el trigger trg_actualizar_saldo.';

CREATE INDEX IF NOT EXISTS idx_cuentas_fiado_empresa
    ON fiados.cfg_cuentas_fiado (id_empresa, id_estado);
 
CREATE INDEX IF NOT EXISTS idx_cuentas_fiado_tercero
    ON fiados.cfg_cuentas_fiado (id_tercero, id_empresa);

-- ─────────────────────────────────────────────────────────────
-- MOVIMIENTOS DE FIADO
-- Registra cada evento: nueva deuda (FIADO) o pago (ABONO).
-- Cada movimiento genera un comprobante en el sistema central.
-- Las observaciones capturan la lista de lo entregado al cliente
-- (arroz, aceite, panela...) o cualquier nota del tendero.
-- ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS fiados.mov_fiados (
    id SERIAL PRIMARY KEY,
    id_cuenta INT NOT NULL REFERENCES fiados.cfg_cuentas_fiado(id),
    -- Comprobante generado por el sistema central
    id_mov_comprobante  INT  NOT NULL REFERENCES comprobantes.mov_gestion_comprobantes(id),
    id_tipo INT NOT NULL REFERENCES fiados.ref_tipo_movimiento(id),
    fecha_movimiento    DATE           NOT NULL DEFAULT CURRENT_DATE,
    valor  NUMERIC(14,2)  NOT NULL CONSTRAINT chk_valor_positivo CHECK (valor > 0),
    -- Lista de artículos entregados o nota del tendero
    -- Ej: "Arroz 2kg, aceite, panela, gaseosa"
    -- Capturado por voz (Web Speech API) o escritura manual
    observaciones       TEXT           NULL,
    -- Saldo de la cuenta DESPUÉS de este movimiento
    -- (snapshot para historial sin recalcular)
    saldo_anterior       NUMERIC(14,2)  NOT NULL DEFAULT 0,
    saldo_movimiento       NUMERIC(14,2)  NOT NULL DEFAULT 0,
    saldo_despues       NUMERIC(14,2)  NOT NULL DEFAULT 0,
    -- Multi-empresa
    id_empresa          INT            NOT NULL,
    id_sede             INT            NULL,
    -- Auditoría
    created_at          TIMESTAMP      NOT NULL DEFAULT NOW(),
    created_by          INT            NOT NULL,
    updated_at          TIMESTAMP      NULL,
    updated_by          INT            NULL,
    CONSTRAINT chk_fecha_no_futura
        CHECK (fecha_movimiento <= CURRENT_DATE)
);

COMMENT ON TABLE fiados.mov_fiados IS
'Movimientos de fiado: FIADO (nueva deuda) o ABONO (pago).
Cada fila genera un comprobante en comprobantes.mov_gestion_comprobantes.
observaciones: lista de artículos entregados o nota libre del tendero.
saldo_despues: snapshot del saldo de la cuenta al posterior al movimiento.';

COMMENT ON COLUMN fiados.mov_fiados.observaciones IS
'Lista de artículos fiados o nota del tendero.
Puede capturarse por voz (Web Speech API en el frontend) o escritura.
Ejemplos: "Arroz 2kg, aceite 1L, panela" / "Mercado semanal" / "Gaseosas fiesta"';

CREATE INDEX IF NOT EXISTS idx_mov_fiados_cuenta
    ON fiados.mov_fiados (id_cuenta, fecha_movimiento DESC);
 
CREATE INDEX IF NOT EXISTS idx_mov_fiados_empresa
    ON fiados.mov_fiados (id_empresa, fecha_movimiento DESC);
 
CREATE INDEX IF NOT EXISTS idx_mov_fiados_comprobante
    ON fiados.mov_fiados (id_mov_comprobante);