-- ============================================================
-- 002_agen.up.sql
-- Mensajes entre los agentes y quienes les escriben.
-- Solo aplica a bases de tenants operativos — NO se ejecuta en la
-- BD admin.
-- ============================================================

-- "detalle" guarda con qué agente y con quién fue el mensaje, p. ej.:
--   {"id_agente": 1, "origen": "prismar", "id_usuario": 5}
--   {"id_agente": 1, "origen": "whatsapp", "telefono": "+57..."}   (más adelante)
CREATE TABLE IF NOT EXISTS agentes.mov_mensajes(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mensaje TEXT NOT NULL,
    es_ia BOOLEAN NOT NULL,
    detalle JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT NOW(),
    id_empresa INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mov_mensajes_detalle
    ON agentes.mov_mensajes USING GIN (detalle);
