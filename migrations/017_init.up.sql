-- Habilitación manual por tenant para el canal de delivery externo.
-- Regla: solo si delivery_app = TRUE el tenant puede tener catálogo
-- visible hacia ese canal. El onboarding SIEMPRE lo deja en FALSE; se
-- habilita manualmente en base de datos, negocio por negocio.
ALTER TABLE configuracion.cfg_tenants
    ADD COLUMN IF NOT EXISTS delivery_app BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN configuracion.cfg_tenants.delivery_app IS
'TRUE = el tenant fue habilitado manualmente para publicar catálogo hacia el canal de delivery externo.
Por defecto FALSE en el onboarding; se activa a mano en BD por tenant, no hay proceso automático.';

-- ─── Ítem de menú "Catálogo" bajo el módulo Inventario (MD-004) ──────────────
-- Idempotente: usa el id_modulo/id_producto resueltos por código (no hardcodea
-- IDs) y solo inserta si esa ruta no existe todavía para ese módulo. En bases
-- de tenant nuevas (donde MD-004 aún no está sembrado) simplemente no inserta
-- nada — no falla.
INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto,
    m.id,
    'Catálogo',
    4,
    '/inventario/search-catalogos',
    true,
    1,
    1
FROM configuracion.cfg_modulos m
WHERE m.codigo = 'MD-004'
AND NOT EXISTS (
    SELECT 1
    FROM configuracion.cfg_funciones f
    WHERE f.id_modulo = m.id
      AND f.ruta_acceso = '/inventario/search-catalogos'
);
