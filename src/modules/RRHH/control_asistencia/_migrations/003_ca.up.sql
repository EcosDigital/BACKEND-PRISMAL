-- ============================================================
-- MIGRACIÓN: 003_CA_UP.sql
-- MÓDULO: Control de Asistencia (CA)
-- SISTEMA: Prismar ERP
-- DESCRIPCIÓN: Creación de la tabla de movimientos de
--              marcaciones. Diseño basado en eventos (event-
--              driven): cada fila representa una marcación
--              individual. Sin cálculos ni lógica de negocio.
-- DEPENDE DE: 001_CA_UP.sql, 002_CA_UP.sql
-- AUTOR: Andres Pastra - Ingeniero de Software
-- FECHA: 2026-03-16
-- ============================================================

-- ============================================================
-- ESQUEMA
-- ============================================================
 
CREATE SCHEMA IF NOT EXISTS control_asistencia;

-- ============================================================
-- TABLAS DE MOVIMIENTO
-- ============================================================
 
-- ------------------------------------------------------------
-- MOV: Marcaciones
-- Registro inmutable de cada evento de marcación realizado
-- por un empleado. Cada fila representa un único evento.
-- No almacena cálculos; la lógica de negocio (tardanzas,
-- horas trabajadas, ausencias) se resuelve en el backend.
--
-- id_empleado no tiene FK declarada porque la tabla de
-- empleados pertenece al módulo de Recursos Humanos.
-- ------------------------------------------------------------

CREATE TABLE IF NOT EXISTS control_asistencia.ref_origen_marcacion(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(150)    NOT NULL,
    activo              BOOLEAN         NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS control_asistencia.mov_marcaciones (
    id                  SERIAL          PRIMARY KEY,
    id_empleado         INTEGER         NOT NULL REFERENCES configuracion.cfg_terceros(id),
    id_punto_marcacion  INTEGER         NOT NULL REFERENCES control_asistencia.cfg_puntos_marcacion(id),
    fecha_hora          TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Geolocalización capturada en el dispositivo al marcar
    latitud             NUMERIC(10, 7),
    longitud            NUMERIC(10, 7),
    precision_metros    INTEGER,
    -- Datos de red y dispositivo para trazabilidad y auditoría
    direccion_ip        VARCHAR(50),
    user_agent          TEXT,
    -- Clasificación del evento
    tipo_evento         VARCHAR(30)     NOT NULL DEFAULT 'ASISTENCIA',
    origen              INTEGER NOT NULL REFERENCES control_asistencia.ref_origen_marcacion(id),
    -- Campo libre para notas o contexto adicional del evento
    observacion         TEXT,
    -- Fecha de inserción del registro en la base de datos
    fecha_creacion      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_mov_marcaciones_punto_marcacion
        FOREIGN KEY (id_punto_marcacion)
        REFERENCES control_asistencia.cfg_puntos_marcacion (id)
);