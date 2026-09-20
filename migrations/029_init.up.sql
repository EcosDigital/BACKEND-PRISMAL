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

CREATE TABLE IF NOT EXISTS notificaciones.mov_notificaciones(
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(200) NOT NULL,
    mensaje TEXT NOT NULL,
    id_modulo_ref INT NULL REFERENCES configuracion.ref_modulos_tenant(id),
    id_tipo_notificacion INT NOT NULL REFERENCES notificaciones.ref_tipo_notificacion(id),
    enviado_por INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
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
