CREATE SCHEMA inventario;

CREATE TABLE IF NOT EXISTS inventario.cfg_grupo_articulos(
    id SERIAL PRIMARY KEY NOT NULL,
    codigo VARCHAR(10) NOT NULL UNIQUE,
    nombre VARCHAR(150) NOT NULL,
    descripcion TEXT NULL
    is_active BOOLEAN DEFAULT TRUE,
);

INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre, descripcion) VALUES ('001', 'Medicamentos', 'Fármacos y productos regulados (con lote y vencimiento)') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre) VALUES ('002', 'Alimentos y bebidas', 'Productos comestibles, ingredientes y bebidas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre) VALUES ('003', 'Insumos y materiales', 'Suministros de uso o venta (empaques, tornillos, químicos, repuestos).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre) VALUES ('004', 'Equipos y dispositivos', 'Herramientas, máquinas, dispositivos médicos o eléctricos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre) VALUES ('005', 'Servicios', 'Prestaciones no físicas (domicilios, instalaciones, recargas).') ON CONFLICT(codigo) DO NOTHING;

CREATE TABLE inventario.ref_unidad_medidas(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL
);

INSERT INTO inventario.ref_unidad_medidas (nombre) VALUES ('Unidad');

CREATE TABLE IF NOT EXISTS inventario.cfg_articulos(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(50) NOT NULL,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NOT NULL,
    id_tipo_articulo INT NOT NULL REFERENCES inventario.cfg_grupo_articulos(id),
    id_unidad_medida INT NOT NULL REFERENCES inventario.ref_unidad_medidas(id),
    concentracion VARCHAR(50) NULL,
    id_forma_farmaceutica INT NULL REFERENCES inventario.ref_forma_farmaceutica(id),
    id_presentacion INT NOT NULL REFERENCES inventario.ref_presentacion_articulo(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    update_by INT,
    update_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS inventario.cfg_articulo_proveedores(
    id SERIAL PRIMARY KEY,
    id_articulo INT NOT NULL REFERENCES inventario.cfg_articulos(id),
    id_proveedor INT NOT NULL REFERENCES configuracion.cfg_articulos(id),
    marca VARCHAR(150) NULL,
    codigo_barras VARCHAR(150) NULL,
    codigo_cum VARCHAR(150) NULL,
    registro_invima VARCHAR(150) NULL,
    referencia_proveedor VARCHAR(250) NULL,
    created_by INT NOT NULL,
);

CREATE TABLE IF NOT EXISTS inventario.ref_tipo_bodega(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(250) NOT NULL
);

INSERT INTO inventario.ref_tipo_bodega (nombre) VALUES ('Propia');
INSERT INTO inventario.ref_tipo_bodega (nombre) VALUES ('Externa');
INSERT INTO inventario.ref_tipo_bodega (nombre) VALUES ('Tercerizada');
INSERT INTO inventario.ref_tipo_bodega (nombre) VALUES ('Virtual');

CREATE TABLE IF NOT EXISTS inventario.cfg_bodegas(
    id SERIAL PRIMARY KEY,
    id_tipo_bodega INT NOT NULL REFERENCES inventario.ref_tipo_bodega(id),
    codigo VARCHAR(10) NOT NULL UNIQUE,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    is_central BOOLEAN,
    permite_ventas BOOLEAN DEFAULT TRUE,
    stock_minimo_global NUMERIC(12,2),
    stock_maximo_global NUMERIC(12,2),
    alerta_stock BOOLEAN DEFAULT TRUE,
    aplica_mov_contable BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT TRUE,
    id_empresa INT NOT NULL REFERENCES configuracion.cfg_empresas(id),
    id_sede INT REFERENCES concentracion.cfg_sedes(id),
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inventario.ref_tipo_movimiento(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(250) NOT NULL
);

INSERT INTO inventario.ref_tipo_movimiento (nombre) VALUES ('Saldo Inicial');
INSERT INTO inventario.ref_tipo_movimiento (nombre) VALUES ('Entradas');
INSERT INTO inventario.ref_tipo_movimiento (nombre) VALUES ('Bajas');
INSERT INTO inventario.ref_tipo_movimiento (nombre) VALUES ('Traslados');
INSERT INTO inventario.ref_tipo_movimiento (nombre) VALUES ('Despachos');

CREATE TABLE IF NOT EXISTS inventario.mov_movimientos(
    id SERIAL PRIMARY KEY
)