-- ============================================================
-- 030_init.up.sql
-- Función "Notificaciones" dentro del módulo Configuración (MD-002):
-- panel para enviar notificaciones push por usuario, rol o módulo.
--
-- Vive únicamente en la base admin (igual que 011/015/017/018/020/
-- 022/023/024/025/026/028): el menú/permisos por rol se resuelven contra
-- la BD admin, no contra cada tenant. Se excluye de las migraciones
-- que se replican en cada tenant (excludedFiles en
-- onboarding/database_helper.go).
-- ============================================================

INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto, m.id, 'Notificaciones',
    COALESCE((SELECT MAX(f.orden_lista) FROM configuracion.cfg_funciones f WHERE f.id_modulo = m.id), 0) + 1,
    '/config/notificaciones', true, 1, 1
FROM configuracion.cfg_modulos m
WHERE m.codigo = 'MD-002'
AND NOT EXISTS (
    SELECT 1 FROM configuracion.cfg_funciones f
    WHERE f.id_modulo = m.id AND f.ruta_acceso = '/config/notificaciones'
);
