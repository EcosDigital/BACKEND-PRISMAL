WITH inserted_empresa AS (
    INSERT INTO configuracion.cfg_empresas(
        id_tipo_empresa, nit, dv, razon_social, descripcion, direccion, id_pais, id_departamento, id_ciudad, id_zona, telefono, telefono_2,
        email, fax, pagina_web, logo_url, usa_sedes, numero_sedes, id_naturaleza, id_actividad_economica, id_tipo_contribuyente,
        id_regimen_tributario, tipo_documento_rep, numero_documento, representante_legal, codigo_licencia, is_active, created_by, created_at, update_at
    )
    VALUES (
         3, '7777', '7', 'Ecosistemas Digital', 'Software amigable', 'Montelibano - Cordoba', 1, 10, 429, 1, '3217815529', '',
        'administracion@ecosistemasdigital.com', '', 'https://ecosistemasdigital.com/', '', false, 0, 1, 355, 1, 1, 2,
        '1003081009', 'Andrés Pastrana', 'ECO00001', true, 1, now(), null
    )
    RETURNING id
)
