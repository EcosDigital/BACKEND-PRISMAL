-- ============================================================
-- 001_esco_up.sql
-- Módulo Espacio Comercial — schema, tablas de configuración
-- Pantalla pública de reproducción: playlists de YouTube +
-- videos publicitarios propios, intercalados cada N videos.
-- ============================================================

CREATE SCHEMA IF NOT EXISTS espacio_comercial;

-- ─────────────────────────────────────────────────────────────
-- CONFIGURACIÓN GENERAL (una fila por empresa)
-- public_token identifica el link público de la pantalla de
-- reproducción (sin login). Se genera en Go con google/uuid.
-- ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS espacio_comercial.cfg_espacio_comercial (
    id              SERIAL PRIMARY KEY,
    id_empresa      INT           NOT NULL,
    frecuencia_ads  INT           NOT NULL DEFAULT 3
                        CONSTRAINT chk_frecuencia_ads_positiva CHECK (frecuencia_ads > 0),
    activo          BOOLEAN       NOT NULL DEFAULT TRUE,
    public_token    VARCHAR(36)   NOT NULL,
    created_at      TIMESTAMP     NOT NULL DEFAULT NOW(),
    created_by      INT           NOT NULL,
    updated_at      TIMESTAMP     NULL,
    updated_by      INT           NULL,
    CONSTRAINT uq_esco_config_empresa UNIQUE (id_empresa),
    CONSTRAINT uq_esco_config_token UNIQUE (public_token)
);

COMMENT ON TABLE espacio_comercial.cfg_espacio_comercial IS
'Configuración del módulo Espacio Comercial por empresa: cada cuántos
videos musicales se inserta un anuncio, si la pantalla pública está
activa, y el token que identifica su link público.';

-- ─────────────────────────────────────────────────────────────
-- PLAYLISTS DE VIDEOS MUSICALES (YouTube)
-- Varias playlists por empresa; "orden" define la secuencia en la
-- que se reproducen (al terminar una se pasa a la siguiente).
-- ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS espacio_comercial.cfg_playlists (
    id          SERIAL      PRIMARY KEY,
    id_empresa  INT         NOT NULL,
    nombre      VARCHAR(150) NOT NULL,
    orden       INT         NOT NULL DEFAULT 0,
    activo      BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    created_by  INT         NOT NULL,
    updated_at  TIMESTAMP   NULL,
    updated_by  INT         NULL
);

COMMENT ON TABLE espacio_comercial.cfg_playlists IS
'Playlists de videos musicales de YouTube configuradas por el tenant.
Dentro de cada playlist el orden de reproducción se baraja en el
cliente; "orden" solo define la secuencia entre playlists.';

CREATE INDEX IF NOT EXISTS idx_esco_playlists_empresa
    ON espacio_comercial.cfg_playlists (id_empresa, orden);

-- ─────────────────────────────────────────────────────────────
-- VIDEOS DE CADA PLAYLIST (URLs de YouTube)
-- ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS espacio_comercial.cfg_playlist_videos (
    id           SERIAL      PRIMARY KEY,
    playlist_id  INT         NOT NULL REFERENCES espacio_comercial.cfg_playlists(id) ON DELETE CASCADE,
    youtube_url  TEXT        NOT NULL,
    orden        INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE espacio_comercial.cfg_playlist_videos IS
'URLs de YouTube que componen una playlist. "orden" es solo de
referencia para el panel de administración.';

CREATE INDEX IF NOT EXISTS idx_esco_playlist_videos_playlist
    ON espacio_comercial.cfg_playlist_videos (playlist_id, orden);

-- ─────────────────────────────────────────────────────────────
-- VIDEOS PUBLICITARIOS (subidos por el tenant, archivo propio)
-- "orden" define la rotación secuencial entre anuncios.
-- ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS espacio_comercial.cfg_ads (
    id            SERIAL       PRIMARY KEY,
    id_empresa    INT          NOT NULL,
    nombre        VARCHAR(150) NOT NULL,
    archivo_path  TEXT         NOT NULL,
    orden         INT          NOT NULL DEFAULT 0,
    activo        BOOLEAN      NOT NULL DEFAULT TRUE,
    tamano_bytes  BIGINT       NOT NULL DEFAULT 0,
    mime_type     VARCHAR(100) NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    created_by    INT          NOT NULL
);

COMMENT ON TABLE espacio_comercial.cfg_ads IS
'Videos publicitarios propios del negocio (subidos como archivo),
insertados cada N videos musicales según cfg_espacio_comercial.
frecuencia_ads. La rotación entre ellos es secuencial por "orden".';

CREATE INDEX IF NOT EXISTS idx_esco_ads_empresa
    ON espacio_comercial.cfg_ads (id_empresa, orden);
