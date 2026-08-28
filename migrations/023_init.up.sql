-- Registro de "Espacio Comercial" en el catálogo central de módulos.
-- Vive únicamente en la base admin (igual que 011/015/017/018/020/022);
-- se excluye de las migraciones que se replican en cada tenant
-- (ver excludedFiles en onboarding/database_helper.go).
--
-- El código MD-XXX se calcula dinámicamente (siguiente consecutivo libre)
-- para no chocar con módulos creados manualmente desde la UI que no
-- quedaron registrados en ninguna migración anterior.

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
SELECT
    1,
    (SELECT id FROM configuracion.cfg_categorias WHERE codigo = 'CT-005'),
    3,
    'MD-' || LPAD((
        COALESCE(MAX(NULLIF(regexp_replace(codigo, '\D', '', 'g'), '')::INT), 0) + 1
    )::TEXT, 3, '0'),
    'Espacio Comercial',
    'Pantalla de videos musicales con publicidad intercalada',
    (SELECT COALESCE(MAX(orden_lista), 0) + 1 FROM configuracion.cfg_modulos),
    false,
    'bg-purple-500',
    'bg-purple-50',
    'bg-purple-50',
    'Tv',
    'src/modules/espacio_comercial/_migrations',
    true,
    NOW(),
    1,
    1,
    1
FROM configuracion.cfg_modulos
WHERE NOT EXISTS (
    SELECT 1 FROM configuracion.cfg_modulos
    WHERE migration_path = 'src/modules/espacio_comercial/_migrations'
);
