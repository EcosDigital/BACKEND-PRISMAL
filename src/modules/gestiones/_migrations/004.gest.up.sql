-- ============================================================
-- PRISMAR ERP - MÓDULO GESTIONES
-- Migración: 004_gest_up.sql
-- Descripción: Índices de rendimiento, notas internas,
--              adjuntos y tabla de encuestas de satisfacción.
-- Versión: 1.1
-- ============================================================

-- ============================================================
-- SECCIÓN 1: ÍNDICES DE RENDIMIENTO
-- Previenen degradación cuando el volumen de tickets crece.
-- ============================================================

-- Filtros más frecuentes en el listado de tickets
CREATE INDEX IF NOT EXISTS idx_tickets_estado
    ON gestiones.cfg_tickets_soporte(id_estado);
 
CREATE INDEX IF NOT EXISTS idx_tickets_nivel
    ON gestiones.cfg_tickets_soporte(id_nivel);
 
CREATE INDEX IF NOT EXISTS idx_tickets_origen
    ON gestiones.cfg_tickets_soporte(id_origen);
 
-- Orden por defecto en el listado (DESC por id)
CREATE INDEX IF NOT EXISTS idx_tickets_id_desc
    ON gestiones.cfg_tickets_soporte(id DESC);
 
-- Filtro de creados hoy / rangos de fecha
CREATE INDEX IF NOT EXISTS idx_tickets_created_at
    ON gestiones.cfg_tickets_soporte(created_at DESC);
 
-- Filtro de vencidos: fecha_entrega < NOW() AND estado abierto
CREATE INDEX IF NOT EXISTS idx_tickets_vencidos
    ON gestiones.cfg_tickets_soporte(fecha_entrega)
    WHERE fecha_entrega IS NOT NULL AND id_estado NOT IN (6, 7);
 
-- JOIN en asignación de colaboradores
CREATE INDEX IF NOT EXISTS idx_colab_ticket
    ON gestiones.cfg_ticket_colaboradores(id_ticket);
 
CREATE INDEX IF NOT EXISTS idx_colab_colaborador
    ON gestiones.cfg_ticket_colaboradores(id_colaborador);
 
-- Composite para filtro "asignados a mí" + listado
CREATE INDEX IF NOT EXISTS idx_colab_ticket_colab
    ON gestiones.cfg_ticket_colaboradores(id_ticket, id_colaborador);
 
-- Timeline de gestiones: ORDER BY created_at ASC por ticket
CREATE INDEX IF NOT EXISTS idx_historial_ticket_fecha
    ON gestiones.mov_ticket_historial(id_ticket, created_at ASC);

-- ============================================================
-- SECCIÓN 2: NOTAS INTERNAS PARA COLABORADORES
-- Visibles SOLO para roles internos. Nunca exponer a clientes.
-- tipo = 'NOTE' → comentario interno libre
-- tipo = 'TASK' → tarea asignable con estado de completado
-- ============================================================


CREATE TABLE IF NOT EXISTS gestiones.mov_ticket_notas_internas (
    id            BIGSERIAL    PRIMARY KEY,
    id_ticket     BIGINT       NOT NULL REFERENCES gestiones.cfg_tickets_soporte(id),
    -- Tipo de nota
    tipo          VARCHAR(10)  NOT NULL DEFAULT 'NOTE'
                               CONSTRAINT chk_nota_tipo CHECK (tipo IN ('NOTE', 'TASK')),
    contenido     TEXT         NOT NULL,
    -- Campos exclusivos de TASK
    is_completada  BOOLEAN      NULL,
    completada_en  TIMESTAMPTZ  NULL,
    completada_por INT          NULL,
    -- Colaborador al que se asigna la tarea (NULL = sin asignar específico)
    id_asignado    INT          NULL REFERENCES configuracion.cfg_terceros(id),
    -- Control de visibilidad: INTERNAL = todos los colaboradores, ADMIN = solo admins
    visibilidad   VARCHAR(10)  NOT NULL DEFAULT 'INTERNAL'
                               CONSTRAINT chk_nota_visib CHECK (visibilidad IN ('INTERNAL', 'ADMIN')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by    INT          NOT NULL,
    updated_at    TIMESTAMPTZ  NULL,
    updated_by    INT          NULL
);

COMMENT ON TABLE  gestiones.mov_ticket_notas_internas               IS 'Notas y tareas internas de un ticket. NUNCA exponer a clientes/tenants externos.';
COMMENT ON COLUMN gestiones.mov_ticket_notas_internas.tipo          IS 'NOTE = comentario interno libre | TASK = tarea con estado de completado';
COMMENT ON COLUMN gestiones.mov_ticket_notas_internas.visibilidad   IS 'INTERNAL = visible a todos los colaboradores | ADMIN = solo administradores';
COMMENT ON COLUMN gestiones.mov_ticket_notas_internas.id_asignado   IS 'FK → cfg_terceros. Colaborador responsable de la tarea. NULL si no está asignada.';
COMMENT ON COLUMN gestiones.mov_ticket_notas_internas.is_completada IS 'Solo aplica cuando tipo = TASK. TRUE cuando el colaborador marca la tarea como hecha.';
 
CREATE INDEX IF NOT EXISTS idx_notas_ticket
    ON gestiones.mov_ticket_notas_internas(id_ticket);
 
CREATE INDEX IF NOT EXISTS idx_notas_asignado
    ON gestiones.mov_ticket_notas_internas(id_asignado)
    WHERE id_asignado IS NOT NULL;
 
CREATE INDEX IF NOT EXISTS idx_notas_pendientes
    ON gestiones.mov_ticket_notas_internas(id_ticket, is_completada)
    WHERE tipo = 'TASK' AND is_completada = FALSE;

-- ============================================================
-- SECCIÓN 3: ADJUNTOS
-- Archivos adjuntos a tickets o gestiones específicas.
-- El almacenamiento físico (S3/R2) se gestiona en la aplicación.
-- ============================================================

CREATE TABLE IF NOT EXISTS gestiones.mov_ticket_adjuntos (
    id               BIGSERIAL    PRIMARY KEY,
    id_ticket        BIGINT       NOT NULL REFERENCES gestiones.cfg_tickets_soporte(id),
    -- Referencia opcional a una gestión específica del historial
    id_gestion       BIGINT       NULL REFERENCES gestiones.mov_ticket_historial(id),
    -- Referencia opcional a una nota interna
    id_nota_interna  BIGINT       NULL REFERENCES gestiones.mov_ticket_notas_internas(id),
    nombre_original  VARCHAR(255) NOT NULL,
    nombre_storage   VARCHAR(500) NOT NULL, -- Ruta/key en el bucket (S3/R2/local)
    mime_type        VARCHAR(100) NOT NULL,
    tamano_bytes     BIGINT       NOT NULL DEFAULT 0,
    url_publica      TEXT         NULL,     -- URL firmada o pública según política del bucket
    -- Visibilidad: PUBLIC = visible al cliente, INTERNAL = solo colaboradores
    visibilidad      VARCHAR(10)  NOT NULL DEFAULT 'PUBLIC'
                                  CONSTRAINT chk_adj_visib CHECK (visibilidad IN ('PUBLIC', 'INTERNAL')),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by       INT          NOT NULL
);

COMMENT ON TABLE  gestiones.mov_ticket_adjuntos                    IS 'Archivos adjuntos a tickets, gestiones o notas internas';
COMMENT ON COLUMN gestiones.mov_ticket_adjuntos.nombre_storage     IS 'Clave/ruta del archivo en el bucket de almacenamiento (S3, R2, local)';
COMMENT ON COLUMN gestiones.mov_ticket_adjuntos.visibilidad        IS 'PUBLIC = visible al cliente reportante | INTERNAL = solo colaboradores internos';

CREATE INDEX IF NOT EXISTS idx_adjuntos_ticket
    ON gestiones.mov_ticket_adjuntos(id_ticket);

-- ============================================================
-- SECCIÓN 4: ENCUESTAS DE SATISFACCIÓN (CSAT)
-- Se activa cuando el ticket llega al estado con is_encuesta=TRUE.
-- Un ticket solo puede tener una encuesta (UNIQUE sobre id_ticket).
-- ============================================================

CREATE TABLE IF NOT EXISTS gestiones.mov_ticket_encuestas (
    id          BIGSERIAL    PRIMARY KEY,
    id_ticket   BIGINT       NOT NULL UNIQUE REFERENCES gestiones.cfg_tickets_soporte(id),
 
    -- Puntuación CSAT: 1 (muy insatisfecho) → 5 (muy satisfecho)
    puntuacion  SMALLINT     NOT NULL CONSTRAINT chk_csat_rango CHECK (puntuacion BETWEEN 1 AND 5),
    comentario  TEXT         NULL,
 
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    -- Sin created_by: la llena el cliente/contacto externo (puede ser anónimo)
);

COMMENT ON TABLE  gestiones.mov_ticket_encuestas            IS 'Encuesta de satisfacción CSAT por ticket. Se habilita cuando el estado tiene is_encuesta=TRUE (ej: Resuelto).';
COMMENT ON COLUMN gestiones.mov_ticket_encuestas.puntuacion IS '1=Muy insatisfecho, 2=Insatisfecho, 3=Neutral, 4=Satisfecho, 5=Muy satisfecho';
