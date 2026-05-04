INSERT INTO configuracion.cfg_categorias (
    id_producto,
    codigo,
    nombre,
    descripcion,
    is_active,
    created_at,
    created_by,
    id_empresa,
    id_sede
)
VALUES (
    1,
    'CT-005',
    'Comercial',
    'Ventas, facturación y gestión de clientes',
    true,
    NOW(),
    1,
    1,
    1
)
ON CONFLICT (codigo) DO NOTHING;

INSERT INTO configuracion.cfg_modulos (
    id_producto,
    id_categoria,
    id_estado,
    codigo,
    nombre,
    descripcion,
    orden_lista,
    es_interno,
    color,
    bg_color,
    border_color,
    icono,
    migration_path,
    is_active,
    created_at,
    created_by,
    id_empresa,
    id_sede
)
VALUES (
    1,
    (SELECT id FROM configuracion.cfg_categorias WHERE codigo = 'CT-005'),
    3,
    'MD-008',
    'Ventas',
    'Gestión de ventas y facturación',
    4,
    false,
    'bg-red-500',
    'bg-red-50',
    'bg-red-50',
    'ShoppingCart',
    'src/modules/ventas/_migrations',
    true,
    NOW(),
    1,
    1,
    1
)
ON CONFLICT (codigo) DO NOTHING;