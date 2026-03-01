CREATE SCHEMA costos;

CREATE TABLE IF NOT EXISTS costos.cfg_area_costo(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(15) NOT NULL UNIQUE,
    nombre VARCHAR(350) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS costos.cfg_unidad_funcional(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(15) NOT NULL UNIQUE,
    nombre VARCHAR(350) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);

CREATE TABLE IF NOT EXISTS costos.cfg_centros_costo(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(15) NOT NULL UNIQUE,
    nombre VARCHAR(350) NOT NULL,
    id_area INT NULL REFERENCES costos.cfg_area_costo(id),
    id_unidad_funcional INT NULL REFERENCES costos.cfg_unidad_funcional(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_by INT,
    updated_at TIMESTAMP,
    id_empresa INT NOT NULL,
    id_sede INT NULL,
    CONSTRAINT uq_centro_costo_empresa
        UNIQUE (id_empresa, codigo)
);