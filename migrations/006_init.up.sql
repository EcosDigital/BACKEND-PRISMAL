CREATE TABLE IF NOT EXISTS configuracion.ref_paises(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(300) NOT NULL UNIQUE
);

INSERT INTO configuracion.ref_paises (nombre) VALUES ('Colombia') ON CONFLICT(nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS configuracion.ref_tipo_empresa(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(300) NOT NULL UNIQUE
);

INSERT INTO configuracion.ref_tipo_empresa (nombre)
VALUES ('Comercial') ON CONFLICT(nombre) DO NOTHING;

INSERT INTO configuracion.ref_tipo_empresa (nombre)
VALUES ('Industrial') ON CONFLICT(nombre) DO NOTHING;

INSERT INTO configuracion.ref_tipo_empresa (nombre)
VALUES ('Servicios') ON CONFLICT(nombre) DO NOTHING;

INSERT INTO configuracion.ref_tipo_empresa (nombre)
VALUES ('Distribuidora') ON CONFLICT(nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS configuracion.ref_regimen_tributario(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(300) NOT NULL UNIQUE
);

INSERT INTO configuracion.ref_regimen_tributario(nombre)
VALUES ('Régimen Ordinario') ON CONFLICT(nombre) DO NOTHING;

INSERT INTO configuracion.ref_regimen_tributario(nombre)
VALUES ('Régimen Especial') ON CONFLICT(nombre) DO NOTHING;

INSERT INTO configuracion.ref_regimen_tributario(nombre)
VALUES ('Régimen Simple de Tributación') ON CONFLICT(nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS configuracion.ref_actividad_economica(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(6) NOT NULL UNIQUE,
    nombre VARCHAR(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.ref_naturaleza_empresa(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(300) NOT NULL UNIQUE
);

INSERT INTO configuracion.ref_naturaleza_empresa (nombre) VALUES ('Privada') ON CONFLICT(nombre) DO NOTHING;
INSERT INTO configuracion.ref_naturaleza_empresa (nombre) VALUES ('Pública') ON CONFLICT(nombre) DO NOTHING;
INSERT INTO configuracion.ref_naturaleza_empresa (nombre) VALUES ('Mixta') ON CONFLICT(nombre) DO NOTHING;


CREATE TABLE IF NOT EXISTS configuracion.cfg_empresas(
    id SERIAL PRIMARY KEY,
    id_tipo_empresa INT NOT NULL REFERENCES configuracion.ref_tipo_empresa(id),
    nit VARCHAR(50) NOT NULL,
    dv char(1) NOT NULL,
    razon_social VARCHAR(300) NOT NULL,
    descripcion TEXT NULL,
    direccion VARCHAR(300) NOT NULL,
    id_pais INT NOT NULL REFERENCES configuracion.ref_paises(id),
    id_departamento INT NOT NULL REFERENCES configuracion.ref_departamentos(id),
    id_ciudad INT NOT NULL REFERENCES configuracion.ref_municipios(id),
    id_zona INT NOT NULL REFERENCES configuracion.ref_zona_residencial(id),
    telefono VARCHAR(20) NOT NULL,
    telefono_2 VARCHAR(20),
    email VARCHAR(150) NOT NULL,
    fax VARCHAR(150) NULL,
    pagina_web VARCHAR(300),
    logo_url VARCHAR(300),
    usa_sedes BOOLEAN DEFAULT FALSE,
    numero_sedes INT NOT NULL,
    id_naturaleza INT NOT NULL REFERENCES configuracion.ref_naturaleza_empresa(id),
    id_actividad_economica INT NOT NULL REFERENCES configuracion.ref_actividades_economicas(id),
    id_tipo_contribuyente INT NOT NULL REFERENCES configuracion.ref_tipo_contribuyente(id),
    id_regimen_tributario INT NOT NULL REFERENCES configuracion.ref_regimen_tributario(id),
    tipo_documento_rep INT NOT NULL REFERENCES configuracion.ref_tipo_documento(id),
    numero_documento VARCHAR(20),
    representante_legal VARCHAR(200),
    codigo_licencia VARCHAR NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_by INT NOT NULL,
    update_by INT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS configuracion.ref_tipo_sedes(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(300) NOT NULL UNIQUE
);

INSERT INTO configuracion.ref_tipo_sedes (nombre) VALUES ('Propia') ON CONFLICT(nombre) DO NOTHING;
INSERT INTO configuracion.ref_tipo_sedes (nombre) VALUES ('Externa') ON CONFLICT(nombre) DO NOTHING;
INSERT INTO configuracion.ref_tipo_sedes (nombre) VALUES ('Mixta') ON CONFLICT(nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS configuracion.cfg_sedes(
    id SERIAL PRIMARY KEY,
    id_empresa INT NOT NULL REFERENCES configuracion.cfg_empresas(id),
    nombre TEXT NOT NULL,
    codigo VARCHAR(7) NOT NULL UNIQUE,
    id_tipo_sede INT NOT NULL REFERENCES configuracion.ref_tipo_sedes(id),
    direccion VARCHAR(300) NOT NULL,
    id_pais INT NOT NULL REFERENCES configuracion.ref_paises(id),
    id_departamento INT NOT NULL REFERENCES configuracion.ref_departamentos(id),
    id_ciudad INT NOT NULL REFERENCES configuracion.ref_municipios(id),
    id_zona INT REFERENCES configuracion.ref_zona_residencial(id),
    telefono VARCHAR(20) NOT NULL,
    telefono_2 VARCHAR(20),
    email VARCHAR(150),
    fax VARCHAR(50),
    pagina_web VARCHAR(300),
    responsable_sede VARCHAR(150),
    cargo_responsable VARCHAR(100),
    horario_atencion VARCHAR(100), -- Horario de operación
    geolocalizacion_lat NUMERIC(10,7), -- Latitud GPS
    geolocalizacion_lon NUMERIC(10,7), -- Longitud GPS
    estado BOOLEAN DEFAULT TRUE, 
    created_by INT NOT NULL,
    update_by INT,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP
);
 
CREATE TABLE IF NOT EXISTS seguridad.cfg_empresas_roles(
    id SERIAL PRIMARY KEY,
    id_empresa INT NOT NULL REFERENCES configuracion.cfg_empresas(id),
    id_rol INT NOT NULL REFERENCES seguridad.cfg_roles_usuario(id),
    created_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS seguridad.cfg_sedes_roles(
    id SERIAL PRIMARY KEY,
    id_empresa INT NOT NULL REFERENCES configuracion.cfg_empresas(id),
    id_sede INT NOT NULL REFERENCES configuracion.cfg_sedes(id),
    id_rol INT NOT NULL REFERENCES seguridad.cfg_roles_usuario(id),
    json_modules JSONB NOT NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);