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

-- =============================================
-- DATOS POR DEFECTO: ÁREAS DE COSTO
-- =============================================
INSERT INTO costos.cfg_area_costo (codigo, nombre, is_active, created_by, id_empresa, id_sede)
VALUES
    ('ADM', 'Administración',   TRUE, 1, 1, NULL),
    ('OPE', 'Operaciones',      TRUE, 1, 1, NULL),
    ('COM', 'Comercial',        TRUE, 1, 1, NULL),
    ('FIN', 'Finanzas',         TRUE, 1, 1, NULL),
    ('LOG', 'Logística',        TRUE, 1, 1, NULL),
    ('PRD', 'Producción',       TRUE, 1, 1, NULL),
    ('RHH', 'Recursos Humanos', TRUE, 1, 1, NULL),
    ('TIC', 'Tecnología',       TRUE, 1, 1, NULL);

-- =============================================
-- DATOS POR DEFECTO: UNIDADES FUNCIONALES
-- =============================================
INSERT INTO costos.cfg_unidad_funcional (codigo, nombre, is_active, created_by, id_empresa, id_sede)
VALUES
    ('GRL', 'Gerencia General',      TRUE, 1, 1, NULL),
    ('CTB', 'Contabilidad',          TRUE, 1, 1, NULL),
    ('TES', 'Tesorería',             TRUE, 1, 1, NULL),
    ('ALM', 'Almacén',               TRUE, 1, 1, NULL),
    ('INV', 'Inventario',            TRUE, 1, 1, NULL),
    ('VTA', 'Ventas',                TRUE, 1, 1, NULL),
    ('CMP', 'Compras',               TRUE, 1, 1, NULL),
    ('DSP', 'Despacho y Entregas',   TRUE, 1, 1, NULL),
    ('PRD', 'Producción',            TRUE, 1, 1, NULL),
    ('NOM', 'Nómina',                TRUE, 1, 1, NULL),
    ('SOP', 'Soporte y Sistemas',    TRUE, 1, 1, NULL),
    ('CAL', 'Calidad',               TRUE, 1, 1, NULL);