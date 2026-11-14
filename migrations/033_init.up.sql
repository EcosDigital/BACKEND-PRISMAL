-- Tokens de conexión por usuario (para MCP / agentes externos).
-- Se guarda solo el hash SHA-256; el token real se muestra una única vez al crearlo.
-- Vive en cada tenant junto a seguridad.cfg_usuarios.
CREATE TABLE IF NOT EXISTS seguridad.cfg_tokens_conexion (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_usuario INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
    nombre VARCHAR(100) NOT NULL,
    token_hash CHAR(64) NOT NULL UNIQUE,
    prefijo VARCHAR(12) NOT NULL,
    revocado BOOLEAN NOT NULL DEFAULT false,
    ultimo_uso TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL
);
