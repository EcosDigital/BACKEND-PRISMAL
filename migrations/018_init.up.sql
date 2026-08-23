-- ============================================================
-- 018_init.up.sql
-- Tabla central (base admin) que consolida lo que cada tenant
-- tiene publicado en su Catálogo público local — es la fuente
-- que va a leer la app Mi Llave. Se llena por escritura event-
-- driven desde el backend (no hay job periódico): cada alta,
-- edición o baja en el catálogo del tenant escribe aquí en la
-- misma operación.
-- Regla: si una fila existe aquí, está viva y es consumible por
-- Mi Llave ahora mismo. No hay columna de "activo/visible" — se
-- borra la fila cuando deja de cumplir las condiciones.
-- ============================================================

CREATE SCHEMA IF NOT EXISTS integraciones;

-- ─── Negocios (una fila por sede habilitada) ──────────────────
CREATE TABLE IF NOT EXISTS integraciones.mi_llave_negocios(
    id SERIAL PRIMARY KEY,
    id_tenant INT NOT NULL REFERENCES configuracion.cfg_tenants(id),
    id_sede_origen INT NOT NULL,
    id_catalogo_origen INT NOT NULL,
    nombre VARCHAR(150) NOT NULL,
    direccion VARCHAR(300) NULL,
    telefono VARCHAR(20) NULL,
    horario_atencion VARCHAR(100) NULL,
    geolocalizacion_lat NUMERIC(10,7) NULL,
    geolocalizacion_lon NUMERIC(10,7) NULL,
    logo_url VARCHAR(300) NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    CONSTRAINT uq_mi_llave_negocio_sede UNIQUE (id_tenant, id_sede_origen)
);

COMMENT ON TABLE integraciones.mi_llave_negocios IS
'Proyección central (una fila por sede habilitada) de lo que cada tenant publica localmente
en inventario.cfg_catalogos. Se escribe/borra por evento desde el backend, no hay sincronización
periódica. id_sede_origen / id_catalogo_origen son punteros de vuelta al registro real dentro
de la base del tenant — no tienen significado fuera de Prismar.';

-- ─── Artículos publicados (una fila por producto con domicilio activo) ─
CREATE TABLE IF NOT EXISTS integraciones.mi_llave_articulos(
    id SERIAL PRIMARY KEY,
    id_negocio INT NOT NULL REFERENCES integraciones.mi_llave_negocios(id) ON DELETE CASCADE,
    id_articulo_origen INT NOT NULL,
    codigo VARCHAR(50) NULL,
    nombre VARCHAR(250) NOT NULL,
    descripcion TEXT NULL,
    imagen_url VARCHAR(300) NULL,
    precio_publico NUMERIC(14,2) NOT NULL CHECK (precio_publico >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    CONSTRAINT uq_mi_llave_articulo_origen UNIQUE (id_negocio, id_articulo_origen)
);

COMMENT ON TABLE integraciones.mi_llave_articulos IS
'Proyección central (una fila por producto publicado y con aplica_domicilio=true) de
inventario.cfg_catalogo_articulos de cada tenant. id_articulo_origen es el id real de esa
publicación dentro de la base del tenant — se usa para encontrar qué fila actualizar/borrar
cuando llega un evento de edición o baja. ON DELETE CASCADE: al borrar un negocio se borran
solos todos sus productos.';
