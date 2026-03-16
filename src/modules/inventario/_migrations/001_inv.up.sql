CREATE SCHEMA IF NOT EXISTS inventario;

CREATE TABLE IF NOT EXISTS inventario.cfg_grupo_articulos(
    id SERIAL PRIMARY KEY NOT NULL,
    codigo VARCHAR(10) NOT NULL UNIQUE,
    nombre VARCHAR(150) NOT NULL,
    descripcion TEXT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre, descripcion) VALUES ('001', 'Medicamentos', 'Fármacos y productos regulados (con lote y vencimiento)') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre, descripcion) VALUES ('002', 'Alimentos y bebidas', 'Productos comestibles, ingredientes y bebidas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre, descripcion) VALUES ('003', 'Insumos y materiales', 'Suministros de uso o venta (empaques, tornillos, químicos, repuestos).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre, descripcion) VALUES ('004', 'Equipos y dispositivos', 'Herramientas, máquinas, dispositivos médicos o eléctricos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO inventario.cfg_grupo_articulos (codigo, nombre, descripcion) VALUES ('005', 'Servicios', 'Prestaciones no físicas (domicilios, instalaciones, recargas).') ON CONFLICT(codigo) DO NOTHING;

CREATE TABLE IF NOT EXISTS inventario.ref_unidad_medidas(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL UNIQUE
);

INSERT INTO inventario.ref_unidad_medidas (nombre) 
VALUES ('Unidad'), ('Kilogramo'), ('Libra'), ('Gramo'), ('Onza')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS inventario.ref_forma_farmaceutica(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(50) NOT NULL,
    nombre VARCHAR(250) NOT NULL
);

CREATE TABLE IF NOT EXISTS inventario.ref_presentacion_articulo(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(250) NOT NULL UNIQUE
);

INSERT INTO inventario.ref_presentacion_articulo (nombre) VALUES
  ('Unidad'),
  ('Caja'),
  ('Paquete'),
  ('Bolsa'),
  ('Frasco'),
  ('Ampolla'),
  ('Rollo'),
  ('Litro'),
  ('Kilo'),
  ('A granel'),
  ('Empacado al vacío'),
  ('Porción'),
  ('Par'),
  ('Juego'),
  ('Kit')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS inventario.cfg_articulos(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(50) NOT NULL,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    id_tipo_articulo INT NOT NULL REFERENCES inventario.cfg_grupo_articulos(id),
    id_unidad_medida INT NOT NULL REFERENCES inventario.ref_unidad_medidas(id),
    factor_unidad_base NUMERIC(14,4) DEFAULT 1,
    id_forma_farmaceutica INT NULL REFERENCES inventario.ref_forma_farmaceutica(id),
    id_presentacion INT NULL REFERENCES inventario.ref_presentacion_articulo(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS inventario.cfg_articulo_proveedores(
    id SERIAL PRIMARY KEY,
    id_articulo INT NOT NULL REFERENCES inventario.cfg_articulos(id),
    id_proveedor INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    marca VARCHAR(150) NULL,
    codigo_barras VARCHAR(150) NULL,
    codigo_cum VARCHAR(150) NULL,
    registro_invima VARCHAR(150) NULL,
    referencia_proveedor VARCHAR(250) NULL,
    created_by INT NOT NULL
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
    id_sede INT REFERENCES configuracion.cfg_sedes(id),
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

CREATE TABLE IF NOT EXISTS inventario.ref_estado_movimiento(
    id       SERIAL PRIMARY KEY,
    codigo   VARCHAR(20)  NOT NULL UNIQUE,
    nombre   VARCHAR(100) NOT NULL,
    es_final BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO inventario.ref_estado_movimiento (codigo, nombre, es_final)
VALUES
    ('01', 'Borrador', False),
    ('02', 'Pendiente', False),
    ('03', 'Aprobado', False),
    ('04', 'ejecutado', False),
    ('05', 'Anulado', False)
ON CONFLICT (codigo) DO NOTHING;

CREATE TABLE IF NOT EXISTS inventario.mov_movimientos(
    id SERIAL PRIMARY KEY,
    id_tipo_movimiento INT NOT NULL REFERENCES inventario.ref_tipo_movimiento(id),
    id_estado INT NOT NULL REFERENCES inventario.ref_estado_movimiento(id),
    fecha_movimiento DATE         NOT NULL DEFAULT CURRENT_DATE,
    id_bodega_origen INT NULL REFERENCES inventario.cfg_bodegas(id),
    id_bodega_destino INT NULL REFERENCES inventario.cfg_bodegas(id),
    id_mov_comprobante INT NOT NULL REFERENCES comprobantes.mov_gestion_comprobantes(id),
    id_orden_compra INT NULL,
    --documento soporte externo (factura orden de compra, remision. etc)
    referencia_externa VARCHAR(100) NULL,
    observaciones TEXT NULL,
    --APROBACION SOLO APLICA CUANDO EL TIPO MOV LO REQ
    aprobado_por INT NULL REFERENCES seguridad.cfg_usuarios(id),
    aprobado_at TIMESTAMP,
    --anulacion
    anulado_por INT NULL REFERENCES seguridad.cfg_usuarios(id),
    anulado_at TIMESTAMP,
    motivo_anulacion TEXT NULL,
    
    id_empresa INT NOT NULL,
    id_sede INT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,

    --validaciones bodegas distintas para el mov
    CONSTRAINT chk_traslado_bodegas_distintas CHECK (
        id_bodega_origen IS NULL
        OR id_bodega_destino IS NULL
        OR id_bodega_origen <> id_bodega_destino
    )
);

CREATE TABLE IF NOT EXISTS inventario.mov_movimientos_detalle (
    id SERIAL PRIMARY KEY,
    id_movimiento INT NOT NULL REFERENCES inventario.mov_movimientos(id),
    id_articulo INT NOT NULL REFERENCES inventario.cfg_articulos(id),
    id_proveedor INT NULL REFERENCES configuracion.cfg_terceros(id),
    marca VARCHAR(150) NULL,
    lote VARCHAR(100) NULL,
    cantidad NUMERIC(14,4) NOT NULL CHECK (cantidad > 0),
    valor_unitario  NUMERIC(14,4) NOT NULL DEFAULT 0,
    valor_impuesto  NUMERIC(14,2) NOT NULL DEFAULT 0,
    valor_total     NUMERIC(14,2) GENERATED ALWAYS AS
        (ROUND(cantidad * valor_unitario + valor_impuesto, 2)) STORED,
    observacion TEXT  NULL,
    created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    created_by INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventario.inv_existencias_articulos(
    id SERIAL PRIMARY KEY,
    id_articulo INT NOT NULL REFERENCES inventario.cfg_articulos(id),
    id_bodega INT NULL REFERENCES inventario.cfg_bodegas(id),
    lote VARCHAR(100) NULL,
    stock_actual NUMERIC(14,4) NOT NULL DEFAULT 0,
    costo_promedio NUMERIC(14,4) NULL,
    id_empresa INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uix_existencias_articulo
    UNIQUE (id_articulo, id_bodega, lote, id_empresa)
);