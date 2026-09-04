-- ============================================================
-- 025_init.up.sql
-- Pedidos de domicilio (Mi Llave → Prismar). Vive en la base
-- admin, en el mismo schema `integraciones` de las tablas del
-- catálogo — es la fuente única: el negocio consulta y actualiza
-- el pedido directo aquí, no hay copia en la base del tenant.
-- ============================================================

-- ─── Estados del pedido ────────────────────────────────────────
CREATE TABLE IF NOT EXISTS integraciones.ref_estado_pedido(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(10) NOT NULL UNIQUE,
    nombre VARCHAR(100) NOT NULL,
    es_final BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO integraciones.ref_estado_pedido (codigo, nombre, es_final) VALUES
    ('001', 'Solicitado', false),
    ('002', 'Aceptado', false),
    ('003', 'En preparación', false),
    ('004', 'Listo para entregar', false),
    ('005', 'En camino', false),
    ('006', 'Llegando', false),
    ('007', 'Entregado', true),
    ('008', 'Rechazado', true),
    ('009', 'Cancelado', true)
ON CONFLICT (codigo) DO NOTHING;

-- ─── Cabecera del pedido ───────────────────────────────────────
CREATE TABLE IF NOT EXISTS integraciones.mov_pedidos(
    id SERIAL PRIMARY KEY,
    id_negocio INT NOT NULL REFERENCES integraciones.mi_llave_negocios(id),
    id_tenant INT NOT NULL REFERENCES configuracion.cfg_tenants(id),
    id_referencia_millave VARCHAR(100) NULL,
    id_estado INT NOT NULL REFERENCES integraciones.ref_estado_pedido(id),
    valor_subtotal NUMERIC(14,2) NOT NULL DEFAULT 0,
    valor_domicilio NUMERIC(14,2) NOT NULL DEFAULT 0,
    valor_total NUMERIC(14,2) NOT NULL DEFAULT 0,
    observaciones TEXT NULL,
    motivo_rechazo TEXT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

COMMENT ON TABLE integraciones.mov_pedidos IS
'Registro permanente del pedido, desde que se solicita hasta que se entrega o cancela.
No es una tabla temporal ni de staging. id_referencia_millave es el ID que usa Mi Llave
de su propio lado, para correlacionar el mismo pedido entre los dos sistemas.';

-- ─── Detalle: artículos del pedido ─────────────────────────────
CREATE TABLE IF NOT EXISTS integraciones.mov_pedidos_detalle(
    id SERIAL PRIMARY KEY,
    id_pedido INT NOT NULL REFERENCES integraciones.mov_pedidos(id) ON DELETE CASCADE,
    id_articulo_origen INT NOT NULL,
    codigo VARCHAR(50) NULL,
    nombre VARCHAR(250) NOT NULL,
    cantidad INT NOT NULL CHECK (cantidad > 0),
    precio_unitario NUMERIC(14,2) NOT NULL CHECK (precio_unitario >= 0),
    subtotal NUMERIC(14,2) GENERATED ALWAYS AS (cantidad * precio_unitario) STORED,
    observacion TEXT NULL
);

COMMENT ON COLUMN integraciones.mov_pedidos_detalle.codigo IS
'Copia congelada del código del artículo al momento del pedido — si el negocio luego
cambia el código o el nombre real, este pedido no debe cambiar retroactivamente.';

-- ─── Persona y ubicación de entrega ─────────────────────────────
CREATE TABLE IF NOT EXISTS integraciones.mov_pedidos_entrega(
    id SERIAL PRIMARY KEY,
    id_pedido INT NOT NULL UNIQUE REFERENCES integraciones.mov_pedidos(id) ON DELETE CASCADE,
    cliente_nombre VARCHAR(150) NOT NULL,
    cliente_telefono VARCHAR(20) NOT NULL,
    tipo_documento VARCHAR(10) NOT NULL,
    cliente_documento VARCHAR(30) NOT NULL,
    direccion_entrega VARCHAR(300) NOT NULL,
    geolocalizacion_lat NUMERIC(10,7) NULL,
    geolocalizacion_lon NUMERIC(10,7) NULL,
    referencia TEXT NULL
);

COMMENT ON COLUMN integraciones.mov_pedidos_entrega.tipo_documento IS
'Texto libre (CC, CE, PP, TI, NIT — mismos códigos de configuracion.ref_tipo_documento),
NO es FK al id numérico de esa tabla: el dato viene de Mi Llave, un sistema externo, y
los IDs internos de una tabla no son un contrato seguro para cruzar entre dos sistemas.';
