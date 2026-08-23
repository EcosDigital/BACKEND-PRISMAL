-- ============================================================
-- 008.inv.up.sql
-- 1) Cabecera del catálogo público (por tenant, aislado — sin
--    tabla central ni sincronización todavía). Un catálogo por
--    sede.
-- 2) Imagen opcional de producto en cfg_articulos.
-- 3) Detalle del catálogo: artículos publicados dentro de cada
--    catálogo, con precio público y si aplica para domicilio.
-- ============================================================

CREATE TABLE IF NOT EXISTS inventario.cfg_catalogos(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(20) NOT NULL UNIQUE,
    nombre VARCHAR(150) NOT NULL,
    descripcion TEXT NULL,
    id_empresa INT NOT NULL REFERENCES configuracion.cfg_empresas(id),
    id_sede INT NOT NULL UNIQUE REFERENCES configuracion.cfg_sedes(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL
);

COMMENT ON TABLE inventario.cfg_catalogos IS
'Cabecera del catálogo público por sede. Un registro = un catálogo (un catálogo por sede,
forzado por el UNIQUE en id_sede). El código se genera automáticamente en el service de Go
(formato CAT-0001, consecutivo simple), no se digita ni se genera en la BD.';

COMMENT ON COLUMN inventario.cfg_catalogos.codigo IS
'Autogenerado por el backend al crear el catálogo (CAT-0001, CAT-0002...). No editable.';

COMMENT ON COLUMN inventario.cfg_catalogos.is_active IS
'Activa/desactiva el catálogo completo (todas sus publicaciones a la vez).';

-- ─── Imagen opcional de producto ──────────────────────────────
ALTER TABLE inventario.cfg_articulos
    ADD COLUMN IF NOT EXISTS imagen_url VARCHAR(300) NULL;

COMMENT ON COLUMN inventario.cfg_articulos.imagen_url IS
'URL pública de la imagen del artículo (opcional). Se sube/reemplaza a través de
POST /inventario/articulos/:id/imagen, que borra el archivo físico anterior al reemplazarla.
NULL si el artículo nunca tuvo imagen cargada.';

-- ─── Detalle del catálogo ──────────────────────────────────────
-- Artículos publicados dentro de un catálogo, con su precio público
-- y si aplica o no para el canal de delivery. Sencillo a propósito:
-- precio de oferta / vigencia / promociones quedan para una
-- iteración posterior, no se incluyen aquí todavía.
CREATE TABLE IF NOT EXISTS inventario.cfg_catalogo_articulos(
    id SERIAL PRIMARY KEY,
    id_catalogo INT NOT NULL REFERENCES inventario.cfg_catalogos(id),
    id_articulo INT NOT NULL REFERENCES inventario.cfg_articulos(id),
    precio_publico NUMERIC(14,2) NOT NULL CHECK (precio_publico >= 0),
    aplica_domicilio BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by INT NOT NULL,
    updated_at TIMESTAMP NULL,
    updated_by INT NULL,
    CONSTRAINT uq_catalogo_articulo UNIQUE (id_catalogo, id_articulo)
);

COMMENT ON TABLE inventario.cfg_catalogo_articulos IS
'Detalle del catálogo público: qué artículos están publicados dentro de cada catálogo (cabecera
por sede en cfg_catalogos), con su precio público. Un artículo no puede repetirse dentro del
mismo catálogo (UNIQUE id_catalogo + id_articulo) — "editar" es un UPDATE sobre la fila, no un
nuevo insert.';

COMMENT ON COLUMN inventario.cfg_catalogo_articulos.aplica_domicilio IS
'TRUE = este artículo publicado debe aparecer en la lista de la app de delivery.
Por defecto FALSE: publicar en el catálogo no implica automáticamente que aplique a domicilio,
es una decisión explícita adicional del negocio.';
