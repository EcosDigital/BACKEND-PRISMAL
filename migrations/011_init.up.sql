BEGIN;

-- PRODUCTO
INSERT INTO configuracion.cfg_productos_software 
    (codigo, nombre, descripcion, logo, verssion, is_active, created_by, id_empresa) 
VALUES 
    ('ECO-002', 'Ecoss', 'Suite de modulos operativos y comerciales', '', '0.0.1', true, 1, 1);

-- CATEGORIA 1
INSERT INTO configuracion.cfg_categorias
    (id_producto, codigo, nombre, descripcion, is_active, created_by, id_empresa)
VALUES 
    (1, 'CT-001', 'Parametrizaciones', 'Ajustes y configuraciones de funcionamiento', true, 1, 1);

-- MODULO 1
INSERT INTO configuracion.cfg_modulos 
    (id_producto, id_categoria, id_estado, codigo, nombre, descripcion, orden_lista, es_interno, color, bg_color, border_color, icono, is_active, created_by, id_empresa)
VALUES
    (1, 1, 3, 'MD-001', 'Inicio', 'Estadisticas y metricas generales', 1, true, 'bg-green-500', 'bg-green-50', 'bg-green-50', 'Home', true, 1, 1);

-- MODULO 2
INSERT INTO configuracion.cfg_modulos 
    (id_producto, id_categoria, id_estado, codigo, nombre, descripcion, orden_lista, es_interno, color, bg_color, border_color, icono, is_active, created_by, id_empresa)
VALUES
    (1, 1, 3, 'MD-002', 'Configuraciones', 'Ajustes y parametrizaciones', 2, true, 'bg-blue-500', 'bg-blue-50', 'bg-blue-50', 'Settings', true, 1, 1);


-- CATEGORIA 2
INSERT INTO configuracion.cfg_categorias
    (id_producto, codigo, nombre, descripcion, is_active, created_by, id_empresa)
VALUES
    (1, 'CT-003', 'Finanzas', 'Contabilidad y cartera', true, 1, 1);

-- MODULO 3
INSERT INTO configuracion.cfg_modulos 
    (id_producto, id_categoria, id_estado, codigo, nombre, descripcion, orden_lista, es_interno, color, bg_color, border_color, icono, is_active, created_by, id_empresa)
VALUES
    (1, 2, 3, 'MD-003', 'Contabilidad', 'Cuentas y balances', 3, false, 'bg-blue-600', 'bg-blue-100', 'bg-blue-100', 'Calculator', true, 1, 1);


-- CATEGORIA 3
INSERT INTO configuracion.cfg_categorias
    (id_producto, codigo, nombre, descripcion, is_active, created_by, id_empresa)
VALUES
    (1, 'CT-002', 'Gestion', 'Operaciones y logistica', true, 1, 1);


-- MODULO 4
INSERT INTO configuracion.cfg_modulos 
    (id_producto, id_categoria, id_estado, codigo, nombre, descripcion, orden_lista, es_interno, color, bg_color, border_color, icono, is_active, created_by, id_empresa)
VALUES
    (1, 3, 3, 'MD-004', 'Inventario', 'Stock y movimientos', 3, false, 'bg-green-500', 'bg-green-50', 'bg-green-50', 'Package', true, 1, 1);

COMMIT;
