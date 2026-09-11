-- ============================================================
-- 004_vent.up.sql
-- Schema "ventas": pedidos de mesa (Órdenes).
-- Solo aplica a bases de tenants operativos — NO se ejecuta en la
-- BD admin. Vive en backend/src/modules/ventas/_migrations, que
-- ExecuteTenantMigrations aplica al aprovisionar cada
-- prismar_<tenant>; RunMigrations (BD admin) nunca lee esta carpeta.
-- ============================================================

CREATE SCHEMA IF NOT EXISTS ventas;

CREATE TABLE IF NOT EXISTS ventas.ref_estado_orden(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(10) NOT NULL UNIQUE,
    nombre VARCHAR(100) NOT NULL,
    es_final BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO ventas.ref_estado_orden (codigo, nombre, es_final) VALUES ('001', 'Solicitado', false) ON CONFLICT (codigo) DO NOTHING;
INSERT INTO ventas.ref_estado_orden (codigo, nombre, es_final) VALUES ('002', 'En preparación', false) ON CONFLICT (codigo) DO NOTHING;
INSERT INTO ventas.ref_estado_orden (codigo, nombre, es_final) VALUES ('003', 'Listo para servir', false) ON CONFLICT (codigo) DO NOTHING;
INSERT INTO ventas.ref_estado_orden (codigo, nombre, es_final) VALUES ('004', 'Entregado', true) ON CONFLICT (codigo) DO NOTHING;
INSERT INTO ventas.ref_estado_orden (codigo, nombre, es_final) VALUES ('005', 'Cancelado', true) ON CONFLICT (codigo) DO NOTHING;

-- id_mesa es NULL cuando todas las mesas físicas están ocupadas y aun así
-- hay que tomar el pedido (mesa "ficticia"); identificador_mesa siempre
-- se llena, con el código/nombre real o con lo que el mesero escriba.
CREATE TABLE IF NOT EXISTS ventas.mov_ordenes(
    id SERIAL PRIMARY KEY,
    id_estado INT NOT NULL REFERENCES ventas.ref_estado_orden(id),
    id_mesa INT NULL REFERENCES configuracion.cfg_mesas(id),
    identificador_mesa VARCHAR(250) NOT NULL,
    valor_subtotal NUMERIC(14,2) NOT NULL DEFAULT 0,
    valor_total NUMERIC(14,2) NOT NULL DEFAULT 0,
    observaciones TEXT NULL,
    motivo_rechazo TEXT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS ventas.mov_ordenes_detalle(
    id SERIAL PRIMARY KEY,
    id_orden INT NOT NULL REFERENCES ventas.mov_ordenes(id) ON DELETE CASCADE,
    id_articulo_origen INT NOT NULL,
    codigo VARCHAR(50) NULL,
    nombre VARCHAR(250) NOT NULL,
    cantidad INT NOT NULL CHECK (cantidad > 0),
    precio_unitario NUMERIC(14,2) NOT NULL CHECK (precio_unitario >= 0),
    subtotal NUMERIC(14,2) GENERATED ALWAYS AS (cantidad * precio_unitario) STORED,
    observacion TEXT NULL
);

-- Aviso en tiempo real: pg_notify al llegar una orden nueva. A diferencia
-- del trigger de "integraciones.mov_pedidos" (una sola BD admin, un solo
-- listener), este trigger corre en la BD propia de cada tenant, así que
-- no necesita id_tenant en el payload.
CREATE OR REPLACE FUNCTION ventas.fn_notify_orden_nueva() RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify('ordenes_nuevas', json_build_object('id_orden', NEW.id)::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_notify_orden_nueva ON ventas.mov_ordenes;
CREATE TRIGGER trg_notify_orden_nueva
    AFTER INSERT ON ventas.mov_ordenes
    FOR EACH ROW EXECUTE FUNCTION ventas.fn_notify_orden_nueva();
