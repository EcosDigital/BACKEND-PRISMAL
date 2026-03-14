CREATE SCHEMA IF NOT EXISTS contabilidad;

CREATE TABLE IF NOT EXISTS contabilidad.ref_naturaleza_contable(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO contabilidad.ref_naturaleza_contable (nombre) values 
    ('Activo'),
    ('Pasivo'),
    ('Patrimonio'),
    ('Ingreso'),
    ('Gasto'),
    ('Costo')
ON CONFLICT (nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS contabilidad.cfg_tipo_cuentas(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(20) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true
);

INSERT INTO contabilidad.cfg_tipo_cuentas (nombre) values 
    ('Inventario'),
    ('Caja'),
    ('Bancos'),
    ('Cuentas por Cobrar'),
    ('Cuentas por Pagar'),
    ('Costo Venta')
ON CONFLICT (nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS contabilidad.ref_nivel_cuenta(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO contabilidad.ref_nivel_cuenta (nombre) values 
    ('Nivel 1'),
    ('Nivel 2'),
    ('Nivel 3'),
    ('Nivel 4')
ON CONFLICT (nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS contabilidad.cfg_cuentas_contables(
    id SERIAL PRIMARY KEY,
    codigo_cuenta VARCHAR(20) NOT NULL,
    nombre_cuenta VARCHAR(250) NOT NULL,
    id_cuenta_padre INT NULL REFERENCES contabilidad.cfg_cuentas_contables(id),
    id_naturaleza INT NOT NULL REFERENCES contabilidad.ref_naturaleza_contable(id),
    id_tipo_cuenta INT NULL REFERENCES contabilidad.cfg_tipo_cuentas(id),
    id_nivel_cuenta INT NOT NULL REFERENCES contabilidad.ref_nivel_cuenta(id),
    permite_movimientos BOOLEAN DEFAULT TRUE,
    requiere_tercero BOOLEAN DEFAULT FALSE,
    requiere_centro_costo BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,
    id_empresa INT NOT NULL,
    id_sede INT NULL,
    CONSTRAINT uq_cuenta_empresa
        UNIQUE (id_empresa, codigo_cuenta)
);


CREATE TABLE IF NOT EXISTS contabilidad.ref_conceptos_articulos(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(450) NOT NULL unique
);

INSERT INTO contabilidad.ref_conceptos_articulos (nombre) VALUES
    ('Cuenta Inventario'),
    ('Cuenta Costo'),
    ('Cuenta Gasto')
ON CONFLICT (nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS contabilidad.cfg_cuentas_grupo_articulos (
    id SERIAL PRIMARY KEY,
    id_sede INT NOT NULL REFERENCES configuracion.cfg_sedes(id),
    id_bodega INT NOT NULL,
    id_tipo_articulo INT NOT NULL,
    id_concepto_articulo INT NOT NULL REFERENCES contabilidad.ref_conceptos_articulos(id),
    id_cuenta_contable INT NOT NULL REFERENCES contabilidad.cfg_cuentas_contables(id),
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,
    id_empresa INT NOT NULL,
    CONSTRAINT uq_cuentas_grupo_articulo
        UNIQUE (id_empresa, id_tipo_articulo, id_concepto_articulo)
);

CREATE TABLE IF NOT EXISTS contabilidad.ref_conceptos_generales(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(450) NOT NULL unique
);

INSERT INTO contabilidad.ref_conceptos_generales (nombre) VALUES
    ('Cuenta Cierre de Cajas'),
    ('Cuenta Ajuste Inicial')
ON CONFLICT (nombre) DO NOTHING;

CREATE TABLE contabilidad.cfg_cuentas_generales(
    id SERIAL PRIMARY KEY,
    id_concepto_general INT NOT NULL REFERENCES contabilidad.ref_conceptos_generales(id),
    id_cuenta_contable INT NOT NULL REFERENCES contabilidad.cfg_cuentas_contables(id),
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,
    id_sede INT NULL,
    id_empresa INT NOT NULL,
    CONSTRAINT uq_cuentas_generales
        UNIQUE (id_empresa, id_concepto_general)
);

CREATE TABLE contabilidad.mov_movimientos_contables(
    id SERIAL PRIMARY KEY,
    id_mov_comprobante INT NOT NULL REFERENCES comprobantes.mov_gestion_comprobantes(id),
    id_cuenta_contable INT NOT NULL REFERENCES contabilidad.cfg_cuentas_contables(id),
    id_tercero INT NULL REFERENCES configuracion.cfg_terceros(id),
    id_centro_costo INT NULL REFERENCES costos.cfg_centros_costo(id)
    valor NUMERIC(14,2) NOT NULL,
    Naturaleza CHAR(1) NOT NULL CHECK (naturaleza IN ('D','C')) -- Débito o Crédito
);
