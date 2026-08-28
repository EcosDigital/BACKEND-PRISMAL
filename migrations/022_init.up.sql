BEGIN;

-- Función "Instalar Módulos" dentro del módulo Configuraciones (MD-002).
-- Vive únicamente en el catálogo central (igual que 011/015/017/018/020),
-- por eso se excluye de las migraciones que se replican en cada tenant
-- (ver excludedFiles en onboarding/database_helper.go).
INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    p.id,
    m.id,
    'Instalar Módulos',
    COALESCE((SELECT MAX(f.orden_lista) FROM configuracion.cfg_funciones f WHERE f.id_modulo = m.id), 0) + 1,
    '/config/instalar-modulos',
    true,
    1,
    1
FROM configuracion.cfg_productos_software p
JOIN configuracion.cfg_modulos m ON m.id_producto = p.id AND m.codigo = 'MD-002'
WHERE p.codigo = 'ECO-002'
  AND NOT EXISTS (
      SELECT 1 FROM configuracion.cfg_funciones f2
      WHERE f2.id_modulo = m.id AND f2.ruta_acceso = '/config/instalar-modulos'
  );

COMMIT;
