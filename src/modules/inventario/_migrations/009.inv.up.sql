-- ============================================================
-- 009.inv.up.sql
-- Marca de sincronización con la base central (Mi Llave).
-- TRUE por defecto para no marcar como "fallido" nada existente;
-- se pone en FALSE cuando la escritura event-driven a la base
-- central falla, para que la pantalla avise al negocio.
-- ============================================================

ALTER TABLE inventario.cfg_catalogo_articulos
    ADD COLUMN IF NOT EXISTS sync_ok BOOLEAN NOT NULL DEFAULT TRUE;

COMMENT ON COLUMN inventario.cfg_catalogo_articulos.sync_ok IS
'FALSE = la última escritura a la base central (integraciones.mi_llave_articulos) falló.
El artículo sigue vigente en el tenant, pero puede no estar reflejado en Mi Llave todavía.
Se reintenta automáticamente la próxima vez que se edite/guarde esta fila.';
