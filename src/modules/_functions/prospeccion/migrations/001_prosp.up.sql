-- ============================================================
-- 001_prospeccion_up.sql
-- Schema: prospeccion
-- Módulo: Scraping / Prospección de leads desde Google Maps
-- Tablas: ref_estado_lead, ref_estado_job, cfg_scraping_jobs,
--         prospeccion_leads
-- ============================================================

-- ─── Schema ──────────────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS prospeccion;

-- ─── Estados posibles de un lead ─────────────────────────────
CREATE TABLE IF NOT EXISTS prospeccion.ref_estado_lead (
    id     SERIAL       PRIMARY KEY,
    codigo VARCHAR(20)  NOT NULL UNIQUE,
    nombre VARCHAR(100) NOT NULL
);

INSERT INTO prospeccion.ref_estado_lead (codigo, nombre) VALUES
    ('001',      'Nuevo'),
    ('002', 'Contactado'),
    ('003', 'Calificado'),
    ('004', 'Descartado')
ON CONFLICT (codigo) DO NOTHING;

-- ─── Estados posibles de un job de scraping ──────────────────
CREATE TABLE IF NOT EXISTS prospeccion.ref_estado_job (
    id     SERIAL       PRIMARY KEY,
    codigo VARCHAR(20)  NOT NULL UNIQUE,
    nombre VARCHAR(100) NOT NULL
);

INSERT INTO prospeccion.ref_estado_job (codigo, nombre) VALUES
    ('001',  'Pendiente'),
    ('002', 'Ejecutando'),
    ('003', 'Completado'),
    ('004',    'Fallido')
ON CONFLICT (codigo) DO NOTHING;