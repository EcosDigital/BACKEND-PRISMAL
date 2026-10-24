-- ============================================================
-- 001_agen.up.sql
-- Schema "agentes": configuración de los agentes de IA de cada negocio.
-- Solo aplica a bases de tenants operativos — NO se ejecuta en la
-- BD admin.
-- ============================================================

CREATE SCHEMA IF NOT EXISTS agentes;

CREATE TABLE IF NOT EXISTS agentes.cfg_agentes(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    instrucciones TEXT NOT NULL,
    modelo VARCHAR(100) NOT NULL DEFAULT 'gemini-3.7-flash',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);
