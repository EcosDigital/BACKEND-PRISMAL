CREATE SCHEMA IF NOT EXISTS reportes;

CREATE TABLE IF NOT EXISTS reportes.cfg_informes_generales(
    id SERIAL PRIMARY KEY,
    id_modulo INT NOT NULL REFERENCES configuracion.cfg_modulos(id),
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT,
    funcion_sql VARCHAR(250) NOT NULL,
    is_active BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL
);

CREATE TABLE IF NOT EXISTS reportes.ref_tipo_filtros(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO reportes.ref_tipo_filtros (nombre) VALUES 
('date'),
('select'),
('number'),
('text'),
('search')
ON CONFLICT(nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS reportes.ref_filtros(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(50) UNIQUE NOT NULL,
    label VARCHAR(100) NOT NULL,
    id_tipo INT NOT NULL REFERENCES reportes.ref_tipo_filtros(id)
);

INSERT INTO reportes.ref_filtros (codigo, label, id_tipo)
VALUES 
('fecha_inicial', 'Fecha Inicial', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='date')),
('fecha_final',   'Fecha Final',   (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='date')),
('id_empresa', 'Empresa', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select')),
('id_sede', 'Sede', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select')),
('id_usuario', 'Usuario', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select')),
('id_tercero', 'Tercero', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='search')),
('id_bodega', 'Bodega', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select')),
('id_articulo', 'Articulo', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select')),
('id_cuenta', 'Cuenta Contable', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='search')),
('id_centro_costo', 'Centro Costo', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select')),
('is_active', 'Activo', (SELECT id FROM reportes.ref_tipo_filtros WHERE nombre='select'))
ON CONFLICT (codigo) DO NOTHING;


CREATE TABLE IF NOT EXISTS reportes.cfg_informes_filtros(
    id SERIAL PRIMARY KEY,
    id_informe INT NOT NULL REFERENCES reportes.cfg_informes_generales(id),
    id_filtro INT NOT NULL REFERENCES reportes.ref_filtros(id),
    requerido BOOLEAN DEFAULT FALSE,
    CONSTRAINT uq_reporte_filtro UNIQUE (id_informe, id_filtro)
);