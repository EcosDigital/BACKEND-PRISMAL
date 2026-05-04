-- ============================================================
-- 003_prospeccion_up.sql
-- Schema: prospeccion
-- Función: fn_check_lead_duplicado
-- Propósito: Verificar desde el backend si un lead es duplicado
--            antes de intentar el INSERT, retornando el motivo.
-- ============================================================

CREATE OR REPLACE FUNCTION prospeccion.fn_check_lead_duplicado(
    p_telefono   VARCHAR(50),
    p_nombre     VARCHAR(250),
    p_ciudad     VARCHAR(150),
    p_id_empresa INT
)
RETURNS TABLE (
    es_duplicado    BOOLEAN,
    motivo          VARCHAR(50),   -- 'TELEFONO' | 'NOMBRE_CIUDAD' | NULL
    id_lead_existente INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_id_por_telefono     INT := NULL;
    v_id_por_nombre_ciudad INT := NULL;
BEGIN

    -- ── Chequeo por teléfono ──────────────────────────────────
    IF p_telefono IS NOT NULL AND p_telefono <> '' THEN
        SELECT id INTO v_id_por_telefono
        FROM prospeccion.prospeccion_leads
        WHERE telefono   = p_telefono
          AND id_empresa = p_id_empresa
        LIMIT 1;
    END IF;

    -- Teléfono tiene prioridad: si hay match, retornar inmediatamente
    IF v_id_por_telefono IS NOT NULL THEN
        RETURN QUERY SELECT TRUE, 'TELEFONO'::VARCHAR(50), v_id_por_telefono;
        RETURN;
    END IF;

    -- ── Chequeo por nombre + ciudad ───────────────────────────
    SELECT id INTO v_id_por_nombre_ciudad
    FROM prospeccion.prospeccion_leads
    WHERE nombre     = p_nombre
      AND ciudad     = p_ciudad
      AND id_empresa = p_id_empresa
    LIMIT 1;

    IF v_id_por_nombre_ciudad IS NOT NULL THEN
        RETURN QUERY SELECT TRUE, 'NOMBRE_CIUDAD'::VARCHAR(50), v_id_por_nombre_ciudad;
        RETURN;
    END IF;

    -- ── No es duplicado ───────────────────────────────────────
    RETURN QUERY SELECT FALSE, NULL::VARCHAR(50), NULL::INT;

END;
$$;

COMMENT ON FUNCTION prospeccion.fn_check_lead_duplicado IS
'Verifica si un lead es duplicado antes del INSERT.
Chequea primero por teléfono (dentro de la empresa), luego por
nombre+ciudad. Retorna es_duplicado, el motivo y el ID del registro
existente. El backend llama esta función antes de cada INSERT para
acumular el contador de duplicados del job sin lanzar excepciones.';