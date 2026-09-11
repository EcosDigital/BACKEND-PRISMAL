-- ============================================================
-- 028_init.up.sql
-- Función "Estructura Física" dentro del módulo Configuración (MD-002).
-- Empieza gestionando el catálogo de Mesas (configuracion.cfg_mesas),
-- pero el nombre queda abierto para extenderse más adelante a la
-- estructura física completa del negocio (zonas, áreas, posible
-- conexión con centros de costo en contabilidad).
--
-- Vive únicamente en la base admin (igual que 011/015/017/018/020/
-- 022/023/024/025/026): el menú/permisos por rol se resuelven contra
-- la BD admin, no contra cada tenant. Se excluye de las migraciones
-- que se replican en cada tenant (agregar a excludedFiles en
-- onboarding/database_helper.go).
-- ============================================================

INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto, m.id, 'Estructura Física',
    COALESCE((SELECT MAX(f.orden_lista) FROM configuracion.cfg_funciones f WHERE f.id_modulo = m.id), 0) + 1,
    '/config/estructura-fisica', true, 1, 1
FROM configuracion.cfg_modulos m
WHERE m.codigo = 'MD-002'
AND NOT EXISTS (
    SELECT 1 FROM configuracion.cfg_funciones f
    WHERE f.id_modulo = m.id AND f.ruta_acceso = '/config/estructura-fisica'
);

-- ─── Función "Órdenes" dentro del módulo Ventas (MD-008) ──
-- El mesero arma la orden en la mesa (artículos + cantidades) y la
-- manda a preparación. Distinta de "Pedidos a Domicilio" (Mi Llave).
INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto, m.id, 'Órdenes',
    COALESCE((SELECT MAX(f.orden_lista) FROM configuracion.cfg_funciones f WHERE f.id_modulo = m.id), 0) + 1,
    '/ventas/ordenes', true, 1, 1
FROM configuracion.cfg_modulos m
WHERE m.codigo = 'MD-008'
AND NOT EXISTS (
    SELECT 1 FROM configuracion.cfg_funciones f
    WHERE f.id_modulo = m.id AND f.ruta_acceso = '/ventas/ordenes'
);

-- ─── Renombrar "Pedidos a Domicilio" → "Pedidos" ──
-- La pantalla unifica domicilios (Mi Llave) y pedidos de mesa (Órdenes),
-- así que el nombre ya no debe hacer referencia solo a domicilio.
UPDATE configuracion.cfg_funciones
SET nombre = 'Pedidos'
WHERE ruta_acceso = '/ventas/search-pedidos' AND nombre = 'Pedidos a Domicilio';
