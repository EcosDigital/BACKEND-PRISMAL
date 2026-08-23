-- ============================================================
-- 021_init.up.sql
-- Permite horarios que cruzan medianoche (ej. abre 5pm, cierra
-- 2am del día siguiente) en configuracion.cfg_sede_horarios.
-- hora_cierre < hora_apertura ya NO es un error — significa que
-- el negocio cierra al día siguiente. Solo se exige que, si el
-- día está marcado como abierto, ambas horas estén presentes
-- (eso se valida en el backend, no aquí).
-- ============================================================

ALTER TABLE configuracion.cfg_sede_horarios
    DROP CONSTRAINT IF EXISTS chk_sede_horario_rango;
