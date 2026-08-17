-- Configuración opcional de integración con el asistente de IA, por empresa.
-- Ambos campos son opcionales porque el modelo/endpoint puede cambiar por
-- cliente o con el tiempo, y no toda empresa tiene por qué tener uno asignado.
ALTER TABLE configuracion.cfg_empresas
    ADD COLUMN IF NOT EXISTS ia_endpoint VARCHAR(300) NULL,
    ADD COLUMN IF NOT EXISTS ia_token TEXT NULL;
