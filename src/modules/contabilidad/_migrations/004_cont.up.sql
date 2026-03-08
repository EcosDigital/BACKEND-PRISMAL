CREATE TABLE IF NOT EXISTS contabilidad.ref_tipo_impuesto (
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(20)  NOT NULL UNIQUE,
    nombre VARCHAR(100) NOT NULL,
    suma BOOLEAN DEFAULT FALSE,
    resta BOOLEAN DEFAULT FALSE
);

INSERT INTO contabilidad.ref_tipo_impuesto (codigo, nombre, suma, resta) VALUES
    ('Iva', 'Impuesto al Valor Agregado', true, false),
    ('Rete_Renta', 'Retención en la Fuente', false, true),
    ('Rete_Iva', 'Retención sobre el IVA', false, true),
    ('Rete_Ica', 'Retención de Industria y Comercio', false, true)
ON CONFLICT (codigo) DO NOTHING;

CREATE TABLE IF NOT EXISTS contabilidad.cfg_impuestos (
    id  SERIAL PRIMARY KEY,
    id_tipo_impuesto INT  NOT NULL REFERENCES contabilidad.ref_tipo_impuesto(id),
    codigo  VARCHAR(20)  NOT NULL,
    nombre  VARCHAR(150) NOT NULL,
    porcentaje  NUMERIC(7,4)   NOT NULL CHECK (porcentaje >= 0 AND porcentaje <= 100),
    is_active BOOLEAN  NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,
    id_empresa INT NOT NULL,
    CONSTRAINT uq_impuesto_codigo_empresa UNIQUE (id_empresa, codigo)
);