CREATE OR REPLACE FUNCTION configuracion.qry_busquedas_dinamicas(
    operacion INTEGER,
    nombre_codigo TEXT DEFAULT NULL,
    id_tipo_registro INT DEFAULT NULL,
    fecha_inicial DATE DEFAULT NULL,
    fecha_final DATE DEFAULT NULL,
    id_estado_registro INT DEFAULT NULL,
    active BOOLEAN DEFAULT NULL,
    aplicar_limit BOOLEAN DEFAULT FALSE
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
                    p.id,
                    p.codigo,
                    p.logo,
                    p.nombre,
                    p.descripcion,
                    p.verssion,
                    p.is_active
                FROM configuracion.cfg_productos_software p
                WHERE (nombre_codigo IS NULL OR p.codigo = nombre_codigo OR nombre = nombre_codigo)
                AND (active IS NULL OR p.is_active = active)
                AND (fecha_inicial IS NULL OR p.created_at >= fecha_inicial)
                AND (fecha_final IS NULL OR p.created_at <= fecha_final)
                ORDER BY p.id ASC
            ) t;
            RETURN resultado;

        WHEN operacion = 2 THEN
             SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT
                    ct.id,
                    ct.id_producto,
                    p.nombre as producto,
                    ct.codigo,
                    ct.nombre,
                    ct.descripcion,
                    ct.is_active
                FROM configuracion.cfg_categorias ct
                JOIN configuracion.cfg_productos_software p ON p.id = ct.id_producto
                WHERE (nombre_codigo IS NULL OR ct.codigo = nombre_codigo OR ct.nombre ILIKE '%' || nombre_codigo || '%')
                AND (active IS NULL OR ct.is_active = active)
                AND (id_tipo_registro IS NULL OR ct.id_producto = id_tipo_registro )
                AND (fecha_inicial IS NULL OR ct.created_at >= fecha_inicial)
                AND (fecha_final IS NULL OR ct.created_at <= fecha_final)
                ORDER BY ct.id ASC
            ) t;
            RETURN resultado;

        WHEN operacion = 3 THEN
             SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT
                    m.id,
                    m.id_producto,
                    p.nombre as producto,
                    m.id_categoria,
                    ct.nombre as categoria,
                    m.codigo,
                    m.id_estado,
                    m.nombre,
                    m.descripcion,
                    m.color,
                    m.bg_color,
                    m.border_color,
                    m.icono,
                    m.is_active,
                    m.es_interno
                FROM configuracion.cfg_modulos m
                JOIN configuracion.cfg_productos_software p ON p.id = m.id_producto
                JOIN configuracion.cfg_categorias ct ON ct.id = m.id_categoria
                WHERE (nombre_codigo IS NULL OR m.codigo = nombre_codigo OR m.nombre ILIKE '%' || nombre_codigo || '%')
                AND (active IS NULL OR m.is_active = active)
                AND (id_tipo_registro IS NULL OR m.id_categoria = id_tipo_registro )
                AND (fecha_inicial IS NULL OR m.created_at >= fecha_inicial)
                AND (fecha_final IS NULL OR m.created_at <= fecha_final)
                ORDER BY m.id ASC
            ) t;
            RETURN resultado;


        WHEN operacion = 4 THEN
             SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                SELECT
                    f.id,
                    f.id_producto,
                    p.nombre as producto,
                    f.id_modulo,
                    m.nombre as modulo,
                    f.nombre,
                    f.orden_lista,
                    f.ruta_acceso,
                    f.is_active
                FROM configuracion.cfg_funciones f
                JOIN configuracion.cfg_productos_software p ON p.id = f.id_producto
                JOIN configuracion.cfg_modulos m ON m.id = f.id_modulo
                WHERE (nombre_codigo IS NULL OR f.nombre ILIKE '%' || nombre_codigo || '%')
                AND (active IS NULL OR f.is_active = active)
                AND (id_tipo_registro IS NULL OR f.id_modulo = id_tipo_registro )
                AND (fecha_inicial IS NULL OR f.created_at >= fecha_inicial)
                AND (fecha_final IS NULL OR f.created_at <= fecha_final)
                ORDER BY f.id ASC
            ) t;
            RETURN resultado;


        WHEN operacion = 5 THEN
             SELECT to_jsonb(array_agg(t))
            INTO resultado
            FROM (
                
                SELECT
        t.id,
        tp.nombre AS tipo_persona,
        t.numero_documento,
        CASE
            WHEN t.id_tipo_persona = 1 THEN
                ti.codigo || ' - ' ||
                COALESCE(t.primer_nombre, '') || ' ' ||
                COALESCE(t.segundo_nombre, '') || ' ' ||
                COALESCE(t.primer_apellido, '') || ' ' ||
                COALESCE(t.segundo_apellido, '')
            ELSE
                t.razon_social
            END AS nombre,
                g.nombre AS genero
            FROM configuracion.cfg_terceros t
            JOIN configuracion.ref_tipo_persona tp 
                ON tp.id = t.id_tipo_persona
            LEFT JOIN configuracion.ref_genero g 
                ON g.id = t.id_genero
            LEFT JOIN configuracion.ref_tipo_documento ti 
                ON ti.id = t.id_tipo_documento

             WHERE (
                nombre_codigo IS NULL
                OR t.numero_documento ILIKE '%' || nombre_codigo || '%'
                OR t.razon_social ILIKE '%' || nombre_codigo || '%'
                OR (
                    COALESCE(t.primer_nombre, '') || ' ' ||
                    COALESCE(t.segundo_nombre, '') || ' ' ||
                    COALESCE(t.primer_apellido, '') || ' ' ||
                    COALESCE(t.segundo_apellido, '')
                ) ILIKE '%' || nombre_codigo || '%'
            )
            ORDER BY t.id ASC

            ) t;
            RETURN resultado;
        

    END CASE;
    RETURN '[]'::jsonb;
END;    
$function$