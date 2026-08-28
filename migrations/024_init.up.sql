-- Función de acceso al módulo Espacio Comercial (registrado en 023).
-- Vive únicamente en la base admin (igual que 011/015/017/018/020/022/023);
-- se excluye de las migraciones que se replican en cada tenant
-- (ver excludedFiles en onboarding/database_helper.go).

INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto,
    m.id,
    'Configuración',
    1,
    '/espacio-comercial',
    true,
    1,
    1
FROM configuracion.cfg_modulos m
WHERE m.migration_path = 'src/modules/espacio_comercial/_migrations'
  AND NOT EXISTS (
      SELECT 1 FROM configuracion.cfg_funciones f
      WHERE f.id_modulo = m.id AND f.ruta_acceso = '/espacio-comercial'
  );
