-- ============================================================
-- 002_prospeccion_up.sql
-- Schema: prospeccion
-- Tablas: cfg_scraping_jobs, prospeccion_leads
-- Índices: detección de duplicados por teléfono y nombre+ciudad
-- ============================================================

-- ─── Historial de jobs de scraping ───────────────────────────
-- Registra cada ejecución: quién la pidió, qué buscó,
-- cuántos leads resultaron, tiempo de ejecución y estado final.
CREATE TABLE IF NOT EXISTS prospeccion.cfg_scraping_jobs (
    id              SERIAL          PRIMARY KEY,
    keyword         VARCHAR(250)    NOT NULL,
    ciudad          VARCHAR(150)    NOT NULL,
    limite          INT             NOT NULL DEFAULT 50,
    id_estado       INT             NOT NULL REFERENCES prospeccion.ref_estado_job(id),

    -- Contadores resultado
    total_encontrados   INT         NOT NULL DEFAULT 0,
    total_guardados     INT         NOT NULL DEFAULT 0,
    total_duplicados    INT         NOT NULL DEFAULT 0,
    total_errores       INT         NOT NULL DEFAULT 0,

    -- Mensaje de error global si el job falla completamente
    error_mensaje   TEXT            NULL,

    -- Trazabilidad
    id_empresa      INT             NOT NULL,
    id_sede         INT             NULL,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    created_by      INT             NOT NULL,
    finished_at     TIMESTAMP       NULL
);

-- ─── Leads capturados desde Google Maps ──────────────────────
CREATE TABLE IF NOT EXISTS prospeccion.prospeccion_leads (
    id              SERIAL          PRIMARY KEY,

    -- Vínculo al job que lo originó
    id_job          INT             NOT NULL REFERENCES prospeccion.cfg_scraping_jobs(id),

    -- Datos extraídos de Google Maps
    nombre          VARCHAR(250)    NOT NULL,
    tipo_negocio    VARCHAR(250)    NULL,
    ciudad          VARCHAR(150)    NOT NULL,
    direccion       TEXT            NULL,
    telefono        VARCHAR(50)     NULL,
    sitio_web       VARCHAR(500)    NULL,
    rating          NUMERIC(3,1)    NULL,
    total_reviews   INT             NULL,

    -- URL de la ficha en Google Maps (útil para auditoría y re-scraping)
    maps_url        TEXT            NULL,

    -- Estado del lead dentro del CRM
    id_estado_lead  INT             NOT NULL REFERENCES prospeccion.ref_estado_lead(id)
                                    DEFAULT 1,

    -- Trazabilidad
    id_empresa      INT             NOT NULL,
    id_sede         INT             NULL,
    created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
    created_by      INT             NOT NULL,
    updated_at      TIMESTAMP       NULL,
    updated_by      INT             NULL
);

-- ─── Índices de búsqueda y rendimiento ───────────────────────

-- Filtrado frecuente por empresa
CREATE INDEX IF NOT EXISTS idx_leads_empresa
    ON prospeccion.prospeccion_leads (id_empresa);

-- Filtrado por job
CREATE INDEX IF NOT EXISTS idx_leads_job
    ON prospeccion.prospeccion_leads (id_job);

-- Búsqueda de duplicados por teléfono (dentro de la misma empresa)
-- Un teléfono NULL nunca se considera duplicado → se filtra en app
CREATE UNIQUE INDEX IF NOT EXISTS uix_leads_telefono_empresa
    ON prospeccion.prospeccion_leads (telefono, id_empresa)
    WHERE telefono IS NOT NULL AND telefono <> '';

-- Búsqueda de duplicados por nombre + ciudad (dentro de la misma empresa)
CREATE UNIQUE INDEX IF NOT EXISTS uix_leads_nombre_ciudad_empresa
    ON prospeccion.prospeccion_leads (nombre, ciudad, id_empresa);

-- ─── Comentarios de documentación ────────────────────────────
COMMENT ON TABLE prospeccion.cfg_scraping_jobs IS
'Historial de ejecuciones de scraping. Cada fila representa un job
iniciado desde el endpoint POST /scraping/run. Registra keyword,
ciudad, límite de resultados, contadores de éxito/duplicados/errores
y quién lo ejecutó.';

COMMENT ON TABLE prospeccion.prospeccion_leads IS
'Leads capturados desde Google Maps. La detección de duplicados opera
en dos niveles: (1) teléfono único por empresa y (2) nombre+ciudad
únicos por empresa. Ambas restricciones son UNIQUE parciales para
permitir NULLs sin falsos positivos.';

COMMENT ON COLUMN prospeccion.prospeccion_leads.maps_url IS
'URL directa a la ficha de Google Maps. Permite auditar la fuente
del dato y facilita re-scraping futuro si se necesita actualizar.';