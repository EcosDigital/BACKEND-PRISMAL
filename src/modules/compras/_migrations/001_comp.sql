-- ==========================================================
-- MIGRACIÓN: Módulo de compras
-- Archivo: 001_compras.up.sql
-- Depende de: 001_inv.up.sql, 004_cont_impuestos.up.sql
--
-- DECISIONES DE DISEÑO:
--   - Estado y numeración los hereda del comprobante
--   - Solicitud y OC no generan movimiento contable
--   - La entrada al inventario (EM-01) es quien genera
--     el comprobante contable y la CxP automáticamente
--   - No existe tabla de factura proveedor: el número
--     físico de la factura se registra como referencia_externa
--     en inventario.mov_movimientos
-- ==========================================================

CREATE SCHEMA IF NOT EXISTS compras;

-- ----------------------------------------------------------
-- 1. SOLICITUD DE COMPRA
--    Documento interno. No genera movimiento contable.
--    Estado y numeración los hereda del comprobante.
-- ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS compras.mov_solicitudes_compra (
    id                SERIAL PRIMARY KEY,
    -- Comprobante padre (ej: tipo "SC-01")
    id_mov_comprobante    INT  NOT NULL REFERENCES comprobantes.mov_gestion_comprobantes(id),

    fecha_requerida   DATE NULL,
    id_bodega_destino INT  NOT NULL REFERENCES inventario.cfg_bodegas(id),
    observaciones     TEXT NULL,

    id_empresa        INT       NOT NULL,
    id_sede           INT       NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by        INT       NOT NULL,
    updated_at        TIMESTAMP NULL,
    updated_by        INT       NULL
);

-- ----------------------------------------------------------
-- 2. DETALLE DE SOLICITUD DE COMPRA
-- ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS compras.mov_solicitudes_compra_detalle (
    id                  SERIAL PRIMARY KEY,
    id_solicitud        INT           NOT NULL REFERENCES compras.mov_solicitudes_compra(id),
    id_articulo         INT           NOT NULL REFERENCES inventario.cfg_articulos(id),
    cantidad_solicitada NUMERIC(12,2) NOT NULL CHECK (cantidad_solicitada > 0),
    observacion         TEXT          NULL,
    created_at          TIMESTAMP     NOT NULL DEFAULT NOW(),
    created_by          INT           NOT NULL
);

-- ----------------------------------------------------------
-- 3. ORDEN DE COMPRA
--    Puede originarse desde una solicitud aprobada o
--    crearse directamente. No genera movimiento contable.
--    Estado y numeración los hereda del comprobante.
-- ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS compras.mov_ordenes_compra (
    id                SERIAL PRIMARY KEY,

    -- Comprobante padre (ej: tipo "OC-01")
    id_mov_comprobante    INT  NOT NULL REFERENCES comprobantes.mov_gestion_comprobantes(id),

    fecha_entrega_est DATE NULL,
    id_proveedor      INT  NOT NULL REFERENCES configuracion.cfg_terceros(id),
    id_bodega_destino INT  NOT NULL REFERENCES inventario.cfg_bodegas(id),

    -- Si viene de una solicitud aprobada (opcional)
    id_solicitud      INT  NULL REFERENCES compras.mov_solicitudes_compra(id),

    -- Flete pactado con el proveedor en la OC
    valor_flete       NUMERIC(14,2) NOT NULL DEFAULT 0,
    id_impuesto_flete INT           NULL REFERENCES contabilidad.cfg_impuestos(id),

    observaciones     TEXT NULL,

    id_empresa        INT       NOT NULL,
    id_sede           INT       NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by        INT       NOT NULL,
    updated_at        TIMESTAMP NULL,
    updated_by        INT       NULL
);

-- ----------------------------------------------------------
-- 4. DETALLE DE ORDEN DE COMPRA
-- ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS compras.mov_ordenes_compra_detalle (
    id                    SERIAL PRIMARY KEY,
    id_orden_compra       INT           NOT NULL REFERENCES compras.mov_ordenes_compra(id),
    id_articulo           INT           NOT NULL REFERENCES inventario.cfg_articulos(id),
    id_proveedor INT           NULL REFERENCES configuracion.cfg_terceros(id),
    cantidad_ordenada     NUMERIC(12,2) NOT NULL CHECK (cantidad_ordenada > 0),
    -- Se actualiza cada vez que el usuario registra una entrada
    -- en inventario asociada a esta OC
    cantidad_recibida     NUMERIC(12,2) NOT NULL DEFAULT 0,

    precio_unitario       NUMERIC(14,4) NOT NULL DEFAULT 0,
    descuento_porcentaje  NUMERIC(7,4)  NOT NULL DEFAULT 0,

    -- IVA opcional por artículo (solo IVA, sin retenciones)
    id_impuesto_iva       INT           NULL REFERENCES contabilidad.cfg_impuestos(id),

    observacion           TEXT          NULL,
    created_at            TIMESTAMP     NOT NULL DEFAULT NOW(),
    created_by            INT           NOT NULL
);

-- ----------------------------------------------------------
-- 5. CUENTAS POR PAGAR
--    Se genera automáticamente cuando se aplica la entrada
--    al inventario (comprobante tipo EM-01).
--    El comprobante de la entrada es el documento oficial
--    que respalda la deuda con el proveedor.
--    Los anticipos reducen el valor_pendiente.
-- ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS compras.mov_cuentas_por_pagar (
    id                SERIAL PRIMARY KEY,

    -- Comprobante de la entrada que originó la CxP (ej: EM-01)
    id_mov_comprobante    INT           NOT NULL REFERENCES comprobantes.mov_gestion_comprobantes(id),

    id_proveedor      INT           NOT NULL REFERENCES configuracion.cfg_terceros(id),
    fecha_vencimiento DATE          NOT NULL,
    valor_original    NUMERIC(16,2) NOT NULL,
    valor_pagado      NUMERIC(16,2) NOT NULL DEFAULT 0,
    valor_pendiente   NUMERIC(16,2) NOT NULL,

    id_empresa        INT       NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by        INT       NOT NULL,
    updated_at        TIMESTAMP NULL,
    updated_by        INT       NULL,

    CONSTRAINT uq_cxp_comprobante UNIQUE (id_mov_comprobante)
);