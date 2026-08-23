-- ============================================================
-- 020_init.up.sql
-- Horario estructurado por negocio en la base central, espejo de
-- configuracion.cfg_sede_horarios en cada tenant (esa SÍ vive en
-- el tenant, esta vive solo aquí). Necesario para que Mi Llave
-- pueda calcular "¿está abierto ahora?" en el momento de la
-- consulta — no se puede resolver con un job periódico ni con
-- un campo de texto libre, porque el estado abierto/cerrado
-- cambia con el reloj, no con el guardado del catálogo.
-- ============================================================

CREATE TABLE IF NOT EXISTS integraciones.mi_llave_horarios(
    id SERIAL PRIMARY KEY,
    id_negocio INT NOT NULL REFERENCES integraciones.mi_llave_negocios(id) ON DELETE CASCADE,
    dia_semana SMALLINT NOT NULL CHECK (dia_semana BETWEEN 0 AND 6),
    abierto BOOLEAN NOT NULL DEFAULT TRUE,
    hora_apertura TIME NULL,
    hora_cierre TIME NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    CONSTRAINT uq_mi_llave_horario_dia UNIQUE (id_negocio, dia_semana)
);

COMMENT ON TABLE integraciones.mi_llave_horarios IS
'Espejo central de configuracion.cfg_sede_horarios del tenant, un registro por día que el
negocio atiende. Se escribe/borra por evento igual que mi_llave_negocios/mi_llave_articulos —
si un día no tiene fila aquí, ese negocio está cerrado ese día para Mi Llave.';

-- El horario de texto libre queda obsoleto ahora que existe la versión estructurada.
ALTER TABLE integraciones.mi_llave_negocios
    DROP COLUMN IF EXISTS horario_atencion;
