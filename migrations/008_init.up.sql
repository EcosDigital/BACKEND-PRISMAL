CREATE TABLE IF NOT EXISTS configuracion.cfg_productos_software(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(12) NOT NULL UNIQUE,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    logo VARCHAR(300) NULL,
    verssion VARCHAR(300) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.cfg_categorias(
    id SERIAL PRIMARY KEY NOT NULL,
    id_producto INT NOT NULL REFERENCES configuracion.cfg_productos_software(id),
    codigo VARCHAR(50) NOT NULL UNIQUE,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.ref_estado_modulo(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(20) NOT NULL,
    descripcion TEXT NULL
);

INSERT INTO configuracion.ref_estado_modulo (nombre, descripcion) VALUES ('Desarrollo', 'En construcción');
INSERT INTO configuracion.ref_estado_modulo (nombre, descripcion) VALUES ('Pruebas',    'Validación interna');
INSERT INTO configuracion.ref_estado_modulo (nombre, descripcion) VALUES ('Producción', 'Disponible para clientes');

CREATE TABLE IF NOT EXISTS configuracion.cfg_modulos(
    id SERIAL PRIMARY KEY NOT NULL,
    id_producto INT NOT NULL REFERENCES configuracion.cfg_productos_software(id),
    id_categoria INT NOT NULL REFERENCES configuracion.cfg_categorias(id),
    id_estado INT NOT NULL REFERENCES configuracion.ref_estado_modulo(id),
    codigo VARCHAR(50) NOT NULL UNIQUE,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    orden_lista INT NOT NULL,
    es_interno BOOLEAN,
    color VARCHAR(50) NULL,
    bg_color VARCHAR(50) NULL,
    border_color VARCHAR(50) NULL,
    icono VARCHAR(150) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.cfg_funciones(
    id SERIAL PRIMARY KEY NOT NULL,
    id_producto INT NOT NULL REFERENCES configuracion.cfg_productos_software(id),
    id_modulo INT NOT NULL REFERENCES configuracion.cfg_modulos(id),
    nombre VARCHAR(250) NOT NULL,
    orden_lista INT NOT NULL,
    ruta_acceso TEXT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.cfg_subfunciones(
    id SERIAL PRIMARY KEY NOT NULL,
    id_funcion INT NOT NULL REFERENCES configuracion.cfg_funciones(id),
    nombre VARCHAR(250) NOT NULL,
    ruta_acceso TEXT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.cfg_detalle_subfunciones(
    id SERIAL PRIMARY KEY NOT NULL,
    id_modulo INT NOT NULL REFERENCES configuracion.cfg_modulos(id),
    id_subfuncion INT NOT NULL REFERENCES configuracion.cfg_subfunciones(id),
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    ruta_acceso TEXT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

INSERT INTO seguridad.cfg_empresas_roles (id_empresa, id_rol, created_by, created_at)
VALUES (1, 1, 1, NOW());