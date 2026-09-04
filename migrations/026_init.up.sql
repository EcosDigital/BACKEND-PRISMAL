-- ============================================================
-- 026_init.up.sql
-- Pedidos de domicilio (Mi Llave → Prismar): aviso en tiempo real
-- + registro de la función en el módulo Ventas. Admin-only, igual
-- que 025_init.up.sql (agregar a excludedFiles en database_helper.go).
-- ============================================================

-- ─── Aviso en tiempo real: pg_notify al llegar un pedido nuevo ──
CREATE OR REPLACE FUNCTION integraciones.fn_notify_pedido_nuevo() RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify('pedidos_nuevos',
        json_build_object('id_tenant', NEW.id_tenant, 'id_pedido', NEW.id)::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_notify_pedido_nuevo ON integraciones.mov_pedidos;
CREATE TRIGGER trg_notify_pedido_nuevo
    AFTER INSERT ON integraciones.mov_pedidos
    FOR EACH ROW EXECUTE FUNCTION integraciones.fn_notify_pedido_nuevo();

-- ─── Función "Pedidos a Domicilio" dentro del módulo Ventas (MD-008) ──
-- Idempotente: resuelve id_modulo/id_producto por código (no hardcodea IDs)
-- y solo inserta si esa ruta no existe todavía para ese módulo.
INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto, m.id, 'Pedidos a Domicilio',
    COALESCE((SELECT MAX(f.orden_lista) FROM configuracion.cfg_funciones f WHERE f.id_modulo = m.id), 0) + 1,
    '/ventas/search-pedidos', true, 1, 1
FROM configuracion.cfg_modulos m
WHERE m.codigo = 'MD-008'
AND NOT EXISTS (
    SELECT 1 FROM configuracion.cfg_funciones f
    WHERE f.id_modulo = m.id AND f.ruta_acceso = '/ventas/search-pedidos'
);
