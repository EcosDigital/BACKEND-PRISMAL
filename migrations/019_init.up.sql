-- ============================================================
-- 019_init.up.sql
-- 1) Logo opcional por sede en configuracion.cfg_sedes.
--    Nota: horario_atencion, geolocalizacion_lat y geolocalizacion_lon
--    YA EXISTEN en cfg_sedes desde 006_init.up.sql — no se tocan aquí,
--    solo se agrega la imagen, que sí es un campo nuevo.
-- 2) Horario de atención estructurado por día de la semana
--    (configuracion.cfg_sede_horarios) — para negocios que no
--    atienden todos los días.
-- ============================================================

ALTER TABLE configuracion.cfg_sedes
    ADD COLUMN IF NOT EXISTS imagen_url VARCHAR(300) NULL;

COMMENT ON COLUMN configuracion.cfg_sedes.imagen_url IS
'Logo/foto propio de la sede (opcional). Si es NULL, se puede caer al logo general de la
empresa (cfg_empresas.logo_url) para presentación en canales externos como Mi Llave.';

-- ─── Horario de atención por día de la semana ─────────────────
CREATE TABLE IF NOT EXISTS configuracion.cfg_sede_horarios(
    id SERIAL PRIMARY KEY,
    id_sede INT NOT NULL REFERENCES configuracion.cfg_sedes(id),
    dia_semana SMALLINT NOT NULL CHECK (dia_semana BETWEEN 0 AND 6),
    abierto BOOLEAN NOT NULL DEFAULT TRUE,
    hora_apertura TIME NULL,
    hora_cierre TIME NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,
    CONSTRAINT uq_sede_horario_dia UNIQUE (id_sede, dia_semana),
    CONSTRAINT chk_sede_horario_rango CHECK (
        hora_cierre IS NULL OR hora_apertura IS NULL OR hora_cierre > hora_apertura
    )
);

COMMENT ON TABLE configuracion.cfg_sede_horarios IS
'Horario de atención por sede, un registro por día de la semana que el negocio opera.
No incluye múltiples turnos por día (ej. mañana/tarde separados) — un solo rango por día.';

COMMENT ON COLUMN configuracion.cfg_sede_horarios.dia_semana IS
'0=domingo, 1=lunes, 2=martes, 3=miércoles, 4=jueves, 5=viernes, 6=sábado
(mismo orden que EXTRACT(DOW FROM fecha) en Postgres).';

COMMENT ON COLUMN configuracion.cfg_sede_horarios.abierto IS
'Permite mostrar los 7 días en una pantalla de configuración con un switch por día,
en vez de que la ausencia de la fila sea la única forma de representar un día cerrado.';
