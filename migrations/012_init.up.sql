CREATE TABLE IF NOT EXISTS configuracion.ref_modulos_tenant(
    id SERIAL PRIMARY KEY,
    id_ref INT NOT NULL UNIQUE,
    codigo VARCHAR(20) NOT NULL,
    nombre VARCHAR(150) NOT NULL
);