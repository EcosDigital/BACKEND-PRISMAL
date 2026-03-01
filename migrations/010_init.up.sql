CREATE TABLE IF NOT EXISTS configuracion.ref_tipo_licencia (
    id SERIAL PRIMARY KEY,
    nombre TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    duracion_dias int,
    es_gratis BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true
);

INSERT INTO configuracion.ref_tipo_licencia (nombre, descripcion, duracion_dias, es_gratis)
VALUES ('TRIAL_15', 'Licencia de prueba 15 días', 15, true);

CREATE TABLE IF NOT EXISTS configuracion.ref_estado_licencia (
    id SERIAL PRIMARY KEY,
    nombre TEXT NOT NULL
);

INSERT INTO configuracion.ref_estado_licencia (nombre)
VALUES ('Activa');

INSERT INTO configuracion.ref_estado_licencia (nombre)
VALUES ('Vencida');

INSERT INTO configuracion.ref_estado_licencia (nombre)
VALUES ('Suspendida');


CREATE TABLE IF NOT EXISTS configuracion.cfg_tenants (
    id SERIAL PRIMARY KEY,
    id_tercero INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    nombre TEXT NOT NULL,
    dominio TEXT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS configuracion.cfg_licencias (
    id SERIAL PRIMARY KEY,
    codigo_licencia TEXT UNIQUE NOT NULL,
    id_tenant INT NOT NULL REFERENCES configuracion.cfg_tenants(id),
    id_cliente INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    id_producto INT NOT NULL REFERENCES configuracion.cfg_productos_software(id),
    id_tipo_licencia INT NOT NULL REFERENCES configuracion.ref_tipo_licencia(id),
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    id_estado INT NOT NULL REFERENCES configuracion.ref_estado_licencia(id),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS configuracion.mov_licencias_modulos (
    id SERIAL PRIMARY KEY,
    id_licencia INT NOT NULL REFERENCES configuracion.cfg_licencias(id),
    id_modulo INT NOT NULL REFERENCES configuracion.cfg_modulos(id),
    instalado BOOLEAN DEFAULT true
);
