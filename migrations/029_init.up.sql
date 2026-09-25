CREATE SCHEMA IF NOT EXISTS notificaciones;

CREATE TABLE IF NOT EXISTS notificaciones.ref_tipo_notificacion(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

INSERT INTO notificaciones.ref_tipo_notificacion (nombre) VALUES
    ('Informativa'), ('Alerta'), ('Urgente')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS notificaciones.ref_origen_destinatario(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO notificaciones.ref_origen_destinatario (nombre) VALUES
    ('Usuario'), ('Rol'), ('Modulo')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS notificaciones.cfg_dispositivos_push(
    id SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
    token TEXT NOT NULL,
    navegador VARCHAR(100) NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    CONSTRAINT uix_dispositivo_token UNIQUE (token)
);

-- enviado_por NULL = notificación automática del sistema
CREATE TABLE IF NOT EXISTS notificaciones.mov_notificaciones(
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(200) NOT NULL,
    mensaje TEXT NOT NULL,
    id_modulo_ref INT NULL REFERENCES configuracion.ref_modulos_tenant(id),
    id_tipo_notificacion INT NOT NULL REFERENCES notificaciones.ref_tipo_notificacion(id),
    enviado_por INT NULL REFERENCES seguridad.cfg_usuarios(id),
    id_empresa INT NOT NULL DEFAULT 1,
    id_sede INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notificaciones.mov_destinatarios(
    id SERIAL PRIMARY KEY,
    id_notificacion INT NOT NULL REFERENCES notificaciones.mov_notificaciones(id),
    id_usuario INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
    id_origen INT NOT NULL REFERENCES notificaciones.ref_origen_destinatario(id),
    leido BOOLEAN NOT NULL DEFAULT FALSE,
    leido_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uix_notificacion_usuario UNIQUE (id_notificacion, id_usuario)
);

CREATE INDEX IF NOT EXISTS idx_mov_destinatarios_usuario_leido ON notificaciones.mov_destinatarios(id_usuario, leido);
CREATE INDEX IF NOT EXISTS idx_cfg_dispositivos_usuario ON notificaciones.cfg_dispositivos_push(id_usuario);

-- ─── Parámetros de notificación ──
-- copia local de los eventos del catálogo maestro (configuracion.
-- cfg_eventos_notificacion en la base admin). se llena sola al abrir
-- Parámetros generales, solo con los módulos contratados. id_ref es el
-- id del maestro: el id local cambia de un tenant a otro.
CREATE TABLE IF NOT EXISTS notificaciones.ref_evento_notificacion(
    id SERIAL PRIMARY KEY,
    id_ref INT NOT NULL UNIQUE,
    codigo VARCHAR(20) NOT NULL UNIQUE,
    nombre VARCHAR(150) NOT NULL,
    descripcion TEXT NULL,
    id_modulo_ref INT NOT NULL REFERENCES configuracion.ref_modulos_tenant(id),
    id_tipo_notificacion INT NOT NULL REFERENCES notificaciones.ref_tipo_notificacion(id),
    orden_lista INT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at TIMESTAMP NULL
);

-- a quién le llega cada evento en cada sede. id_origen dice si id_destino
-- es un usuario o un rol (notificaciones.ref_origen_destinatario).
CREATE TABLE IF NOT EXISTS notificaciones.cfg_evento_destinatarios(
    id SERIAL PRIMARY KEY,
    id_evento INT NOT NULL REFERENCES notificaciones.ref_evento_notificacion(id),
    id_origen INT NOT NULL REFERENCES notificaciones.ref_origen_destinatario(id),
    id_destino INT NOT NULL,
    id_empresa INT NOT NULL DEFAULT 1,
    id_sede INT NOT NULL,
    created_by INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uix_evento_destinatario UNIQUE (id_evento, id_origen, id_destino, id_sede)
);

CREATE INDEX IF NOT EXISTS idx_cfg_evento_destinatarios_evento_sede
    ON notificaciones.cfg_evento_destinatarios(id_evento, id_sede);
