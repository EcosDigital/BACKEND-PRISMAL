CREATE OR REPLACE FUNCTION configuracion.qry_generales(
    operacion INTEGER = null,
    codigo_dep INTEGER = null
)
RETURNS JSONB
LANGUAGE plpgsql
AS $function$
DECLARE
    id_registro INTEGER;
    resultado JSONB;
BEGIN
    CASE
        WHEN operacion = 1 THEN 
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT 
                id,
                nombre
                FROM configuracion.ref_tipo_persona
            ) t;
            RETURN resultado;

        WHEN operacion = 2 THEN 
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT 
                id,
                nombre,
                codigo,
                tipo
                FROM configuracion.ref_tipo_documento
            ) t;
            RETURN resultado;

        WHEN operacion = 3 THEN 
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT 
                id,
                nombre
                FROM configuracion.ref_genero
            ) t;
            RETURN resultado;

        WHEN operacion = 4 THEN 
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT *
                FROM configuracion.ref_departamentos
            ) t;
            RETURN resultado;

        
         WHEN operacion = 5 THEN 
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
               SELECT *
                FROM configuracion.ref_municipios
                WHERE codigo_departamento = (select codigo from configuracion.ref_departamentos WHERE id = codigo_dep) 
            ) t;
            RETURN resultado;

        WHEN operacion = 6 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_zona_residencial
            ) t;
            RETURN resultado;

        WHEN operacion = 7 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_actividades_economicas
            ) t;
            RETURN resultado;

         WHEN operacion = 8 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_ambitos_terceros
            ) t;
            RETURN resultado;

        WHEN operacion = 9 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_centralizacion
            ) t;
            RETURN resultado;

         WHEN operacion = 10 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_responsabilidades_dian
            ) t;
            RETURN resultado;

         WHEN operacion = 11 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_tipo_contribuyente
            ) t;
            RETURN resultado;

        WHEN operacion = 12 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_regimen_iva
            ) t;
            RETURN resultado;

        WHEN operacion = 13 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_regimen_dian
            ) t;
            RETURN resultado;

        WHEN operacion = 14 THEN
            SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT * 
                    FROM configuracion.ref_clase_personas
            ) t;
            RETURN resultado;


    END CASE;
END;    
$function$