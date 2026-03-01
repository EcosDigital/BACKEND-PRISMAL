WITH inserted_tercero AS (
    INSERT INTO configuracion.cfg_terceros (
        id_tipo_persona, id_tipo_documento, numero_documento, primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, id_genero, telefono, telefono_2, email, direccion, id_pais, id_departamento, id_ciudad, id_zona, estado, created_by, update_by, created_at, update_at
    )
    VALUES (
        1, 2, '1003081009', 'Andrés', 'Felipe', 'Pastrana', 'Llanez', 1, '3217815529', '', 'administrador@ecosistemas.com',
        'calle 33 n 33-41', 1, 10, 12, 1, true, 1, null, NOW(), null
    )
    RETURNING id
)
INSERT INTO seguridad.cfg_usuarios (
    email, password, id_rol, id_tercero, image_profile, activo, created_at, update_at, created_by, update_by
)
SELECT 
    'administrador@ecosistemas.com',
    '$2a$12$xTgnPYFonbG.ztMzX09vC.oSiytiAD2i/FWqvhjFTtsm8lN/M7vsK',
    1,
    id,
    null,
    true,
    NOW(),
    null,
    1,
    null
FROM inserted_tercero;
