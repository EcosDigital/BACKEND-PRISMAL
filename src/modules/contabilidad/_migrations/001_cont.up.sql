CREATE SCHEMA comprobantes;

CREATE TABLE IF NOT EXISTS comprobantes.ref_tipo_operacion(
    id SERIAL PRIMARY KEY,
    id_modulo INT NOT NULL,
    nombre VARCHAR(150) NOT NULL,
    CONSTRAINT uq_ref_tipo_operacion UNIQUE (id_modulo, nombre)
);

INSERT INTO comprobantes.ref_tipo_operacion (id_modulo, nombre)
SELECT id, 'Saldo Inicial'
FROM configuracion.cfg_modulos
WHERE codigo = 'MD-004'
ON CONFLICT (id_modulo, nombre) DO NOTHING;

INSERT INTO comprobantes.ref_tipo_operacion (id_modulo, nombre)
SELECT id, 'Entrada'
FROM configuracion.cfg_modulos
WHERE codigo = 'MD-004'
ON CONFLICT (id_modulo, nombre) DO NOTHING;


INSERT INTO comprobantes.ref_tipo_operacion (id_modulo, nombre)
SELECT id, 'Baja'
FROM configuracion.cfg_modulos
WHERE codigo = 'MD-004'
ON CONFLICT (id_modulo, nombre) DO NOTHING;


INSERT INTO comprobantes.ref_tipo_operacion (id_modulo, nombre)
SELECT id, 'Traslado'
FROM configuracion.cfg_modulos
WHERE codigo = 'MD-004'
ON CONFLICT (id_modulo, nombre) DO NOTHING;


INSERT INTO comprobantes.ref_tipo_operacion (id_modulo, nombre)
SELECT id, 'Despacho'
FROM configuracion.cfg_modulos
WHERE codigo = 'MD-004'
ON CONFLICT (id_modulo, nombre) DO NOTHING;


CREATE TABLE IF NOT EXISTS comprobantes.cfg_comprobante(
    id SERIAL PRIMARY KEY NOT NULL,
    id_modulo INT NOT NULL REFERENCES configuracion.cfg_modulos(id),
    id_tipo_operacion INT NOT NULL REFERENCES comprobantes.ref_tipo_operacion(id),
    nombre_comprobante VARCHAR(250) NOT NULL,
    prefijo_comprobante VARCHAR(10) NOT NULL,
    consecutivo_inicial INT NOT NULL,
    consecutivo_actual INT NOT NULL,
    consecutivo_fin INT NULL,
    permite_anulacion BOOLEAN DEFAULT TRUE,
    fecha_inicio DATE NOT NULL,
    fecha_final DATE NULL,
    token VARCHAR(450) NULL,
    resolucion VARCHAR(450) NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL,
    CHECK (fecha_final IS NULL OR fecha_final >= fecha_inicio)
);

CREATE TABLE IF NOT EXISTS comprobantes.cfg_usuarios_comprobante(
    id SERIAL PRIMARY KEY NOT NULL,
    id_comprobante INT NOT NULL REFERENCES comprobantes.cfg_comprobante(id),
    id_usuario INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL
);

CREATE TABLE IF NOT EXISTS comprobantes.ref_estado_comprobante(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL,
    CONSTRAINT uq_ref_estado_comprobante UNIQUE (nombre)
);

INSERT INTO comprobantes.ref_estado_comprobante (nombre) VALUES
    ('BORRADOR'),
    ('PENDIENTE_AUTORIZACION'),
    ('AUTORIZADO'),
    ('APLICADO'),
    ('ANULADO'),
    ('RECHAZADO')
ON CONFLICT (nombre) DO NOTHING;


CREATE TABLE IF NOT EXISTS comprobantes.mov_gestion_comprobantes(
    id SERIAL PRIMARY KEY NOT NULL,
    id_comprobante INT NOT NULL REFERENCES comprobantes.cfg_comprobante(id),
    id_estado_comprobante INT REFERENCES comprobantes.ref_estado_comprobante(id),
    autorizado BOOLEAN DEFAULT FALSE,
    fecha_autorizacion DATE,
    documento_soporte VARCHAR(250),
    user_autoriza INT,
    id_tercero INT REFERENCES configuracion.cfg_terceros(id),
    fecha_movimiento DATE,
    fecha_creacion TIMESTAMP,
    valor_anterior NUMERIC(14,2),
    valor_movimiento NUMERIC(14,2),
    valor_descuento NUMERIC(14,2),
    valor_impuesto NUMERIC(14,2),
    valor_total_comprobante NUMERIC(14,2),
    consecutivo_comprobante INT,
    prefijo_comprobante VARCHAR(10),
    id_user_anulo INT,
    fecha_anulacion TIMESTAMP,
    motivo_anulacion TEXT,
    observaciones TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);