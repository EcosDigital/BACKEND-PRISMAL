-- ============================================================
-- 004_agen.up.sql
-- Datos de conexión de cada agente con canales externos
-- (por ahora WhatsApp). Solo aplica a bases de tenants operativos.
-- ============================================================

-- "details" guarda, para los agentes conectados a WhatsApp, p. ej.:
--   {"whatsapp": {"telefono": "+57...", "id_telefono_meta": "...",
--                 "id_cuenta_meta": "...", "llave_cifrada": "...",
--                 "conectado": true}}
ALTER TABLE agentes.cfg_agentes
    ADD COLUMN IF NOT EXISTS details JSONB NOT NULL DEFAULT '{}'::jsonb;
