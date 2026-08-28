-- El tenant ya está aislado por su propia base de datos (resuelta por
-- subdominio); no se necesita un token extra para identificar la pantalla
-- pública. Se elimina la columna.
ALTER TABLE espacio_comercial.cfg_espacio_comercial DROP COLUMN IF EXISTS public_token;
