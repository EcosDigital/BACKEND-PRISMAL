-- ============================================================
-- 003_agen.up.sql
-- Módulo del que es encargado cada agente (entre los instalados
-- en el tenant). Solo aplica a bases de tenants operativos.
-- ============================================================

ALTER TABLE agentes.cfg_agentes
    ADD COLUMN IF NOT EXISTS id_modulo_ref INT NULL REFERENCES configuracion.ref_modulos_tenant(id);
