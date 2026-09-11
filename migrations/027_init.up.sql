-- ============================================================
-- 027_init.up.sql
-- Catálogo de Mesas (módulo Configuración, MD-002).
-- A diferencia de "integraciones.mov_pedidos" (compartida en la
-- BD admin), esta tabla es un catálogo operativo del negocio: debe
-- existir en la base propia de cada tenant. Por eso NO se agrega a
-- excludedFiles en onboarding/database_helper.go — se aplica igual
-- en la BD admin y en cada prismar_<tenant> al aprovisionarse.
-- ============================================================

CREATE TABLE IF NOT EXISTS configuracion.ref_estado_mesa(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(20) NOT NULL,
    descripcion TEXT NULL
);

INSERT INTO configuracion.ref_estado_mesa (nombre, descripcion) VALUES ('Disponible', 'Mesa libre, lista para asignar un nuevo pedido');
INSERT INTO configuracion.ref_estado_mesa (nombre, descripcion) VALUES ('Ocupada', 'Mesa con un pedido en curso');

CREATE TABLE IF NOT EXISTS configuracion.cfg_mesas(
    id SERIAL PRIMARY KEY NOT NULL,
    id_estado INT NOT NULL REFERENCES configuracion.ref_estado_mesa(id),
    codigo VARCHAR(50) NOT NULL,
    nombre VARCHAR(250) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

-- Mesas por defecto (estado 1 = Disponible)
INSERT INTO configuracion.cfg_mesas (id_estado, codigo, nombre, is_active, created_by, id_empresa, id_sede) VALUES (1, 'MS-01', 'Mesa 1', true, 1, 1, 1);
INSERT INTO configuracion.cfg_mesas (id_estado, codigo, nombre, is_active, created_by, id_empresa, id_sede) VALUES (1, 'MS-02', 'Mesa 2', true, 1, 1, 1);
INSERT INTO configuracion.cfg_mesas (id_estado, codigo, nombre, is_active, created_by, id_empresa, id_sede) VALUES (1, 'MS-03', 'Mesa 3', true, 1, 1, 1);
INSERT INTO configuracion.cfg_mesas (id_estado, codigo, nombre, is_active, created_by, id_empresa, id_sede) VALUES (1, 'MS-04', 'Mesa 4', true, 1, 1, 1);
INSERT INTO configuracion.cfg_mesas (id_estado, codigo, nombre, is_active, created_by, id_empresa, id_sede) VALUES (1, 'MS-05', 'Mesa 5', true, 1, 1, 1);
INSERT INTO configuracion.cfg_mesas (id_estado, codigo, nombre, is_active, created_by, id_empresa, id_sede) VALUES (1, 'MS-06', 'Mesa 6', true, 1, 1, 1);
