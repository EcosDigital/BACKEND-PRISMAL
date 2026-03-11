-- ==========================================================
-- MÓDULO GESTIONES — Tabla Principal de Tickets
-- ==========================================================

-- ----------------------------------------------------------
-- REFERENCIAL: Origen del ticket
-- Define si el caso viene desde el cliente (externo/tenant)
-- o fue creado internamente por un colaborador o el sistema.
-- ----------------------------------------------------------

CREATE TABLE gestiones.ref_origenes_ticket (
    id     SMALLINT    PRIMARY KEY,
    codigo VARCHAR(20) NOT NULL,
    nombre VARCHAR(60) NOT NULL,
    CONSTRAINT uq_ref_origenes_ticket_codigo UNIQUE (codigo)
);

INSERT INTO gestiones.ref_origenes_ticket (id, codigo, nombre)
VALUES
    (1, '001', 'Interno'),   -- Creado por un colaborador o admin
    (2, '002', 'Externo')    -- Reportado desde la plataforma del cliente (tenant)
ON CONFLICT (id) DO NOTHING;

-- ----------------------------------------------------------
-- TRANSACCIONAL: Tickets
-- ----------------------------------------------------------

CREATE TABLE gestiones.cfg_tickets_soporte (
    id          BIGSERIAL    PRIMARY KEY,
    titulo      VARCHAR(255) NOT NULL,
    descripcion TEXT         NOT NULL,
    id_nivel    INT          NOT NULL REFERENCES gestiones.cfg_niveles_caso(id),
    id_estado   INT          NOT NULL REFERENCES gestiones.cfg_estados_ticket(id),
    id_origen   SMALLINT     NOT NULL DEFAULT 1 REFERENCES gestiones.ref_origenes_ticket(id),
    id_tenant   INT          NULL,
    fecha_entrega  TIMESTAMP NULL DEFAULT NOW(),
    -- Datos libres del contacto que reporta
    contacto_nombre   VARCHAR(200) NOT NULL,
    contacto_email    VARCHAR(200) NULL,
    contacto_telefono VARCHAR(30)  NULL,
    created_at  TIMESTAMP  NOT NULL DEFAULT NOW(),
    created_by  INT          NOT NULL
);

COMMENT ON TABLE  gestiones.cfg_tickets_soporte                  IS 'Tabla principal de casos o tickets reportados al equipo de soporte';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.titulo           IS 'Título corto y descriptivo del caso';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.descripcion      IS 'Detalle completo del problema o solicitud';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.id_nivel         IS 'FK → cfg_niveles_caso. Prioridad e impacto del caso';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.id_estado        IS 'FK → cfg_estados_ticket. Estado actual del ticket en el flujo';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.id_origen        IS 'FK → ref_origenes_ticket. INTERNO = colaborador/admin, EXTERNO = cliente/tenant';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.id_tenant        IS 'ID de la empresa (tenant) que reporta el caso desde su plataforma. NULL si es interno';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.fecha_entrega     IS 'Fecha comprometida de solución o entrega del caso al cliente';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.contacto_nombre  IS 'Nombre de la persona que reporta el caso';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.contacto_email   IS 'Email de contacto del reportante';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.contacto_telefono IS 'Teléfono de contacto del reportante';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.created_at       IS 'Fecha y hora en que se creó el ticket';
COMMENT ON COLUMN gestiones.cfg_tickets_soporte.created_by       IS 'ID del usuario que registró el ticket';

-- ==========================================================
-- MÓDULO GESTIONES — Colaboradores Asignados al Ticket
-- ==========================================================
-- Un ticket puede tener varios colaboradores trabajando
-- en conjunto. Cada asignación puede llevar una tarea
-- o meta específica a cumplir dentro del caso.
-- La reasignación se gestiona eliminando el registro
-- actual e insertando el nuevo colaborador.
-- ==========================================================

CREATE TABLE gestiones.cfg_ticket_colaboradores (
    id             SERIAL       PRIMARY KEY,
    id_ticket      BIGINT       NOT NULL REFERENCES gestiones.cfg_tickets_soporte(id),
    id_colaborador INT          NOT NULL REFERENCES configuracion.cfg_terceros(id),
    tarea          TEXT         NULL,  -- Meta o tarea específica asignada a este colaborador
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by     INT          NOT NULL,
    CONSTRAINT uq_ticket_tercero UNIQUE (id_ticket, id_colaborador)
);

COMMENT ON TABLE  gestiones.cfg_ticket_colaboradores           IS 'Colaboradores asignados a un ticket. Cada uno puede tener una tarea específica dentro del caso';
COMMENT ON COLUMN gestiones.cfg_ticket_colaboradores.id_ticket  IS 'FK → cfg_tickets_soporte';
COMMENT ON COLUMN gestiones.cfg_ticket_colaboradores.id_colaborador IS 'FK → configuracion.cfg_terceros. Colaborador responsable';
COMMENT ON COLUMN gestiones.cfg_ticket_colaboradores.tarea      IS 'Tarea o meta específica asignada a este colaborador dentro del ticket';
COMMENT ON COLUMN gestiones.cfg_ticket_colaboradores.created_at IS 'Fecha y hora en que se realizó la asignación';
COMMENT ON COLUMN gestiones.cfg_ticket_colaboradores.created_by IS 'ID del usuario que realizó la asignación';

-- ==========================================================
-- MÓDULO GESTIONES — Historial de Actividad del Ticket
-- ==========================================================
-- Registra cronológicamente cada acción realizada sobre
-- un ticket: cambios de estado, comentarios, reasignaciones,
-- etc. Es el log completo de vida del caso.
-- ==========================================================

CREATE TABLE gestiones.mov_ticket_historial (
    id              BIGSERIAL   PRIMARY KEY,
    id_ticket       BIGINT      NOT NULL REFERENCES gestiones.cfg_tickets_soporte(id),
    -- Contenido del registro
    comentario      TEXT        NULL,
    -- Snapshot de cambio de estado (se puebla solo si id_tipo_accion = CAMBIO_ESTADO)
    id_estado_anterior  INT     NULL REFERENCES gestiones.cfg_estados_ticket(id),
    id_estado_nuevo     INT     NULL REFERENCES gestiones.cfg_estados_ticket(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      INT         NOT NULL
);

