-- ============================================================
-- 031_init.up.sql
-- Parámetros generales (módulo Configuración MD-002).
--
-- 1) Catálogo maestro de eventos que disparan notificaciones
--    automáticas. Cada tenant copia a su base solo los eventos de los
--    módulos que tiene contratados, al abrir Parámetros generales. En
--    el tenant se reconocen por id_ref/codigo, nunca por su id local.
-- 2) Función "Parámetros generales" dentro de Configuración.
--
-- Vive únicamente en la base admin (igual que 011/015/017/018/020/
-- 022/023/024/025/026/028/030). Se excluye de las migraciones que se
-- replican en cada tenant (excludedFiles en onboarding/database_helper.go).
-- ============================================================

CREATE TABLE IF NOT EXISTS configuracion.cfg_eventos_notificacion(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(20) NOT NULL UNIQUE,
    nombre VARCHAR(150) NOT NULL,
    descripcion TEXT NULL,
    id_modulo INT NOT NULL REFERENCES configuracion.cfg_modulos(id),
    id_tipo_notificacion INT NOT NULL REFERENCES notificaciones.ref_tipo_notificacion(id),
    orden_lista INT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ─── Eventos del módulo Ventas (MD-008) ──
INSERT INTO configuracion.cfg_eventos_notificacion
    (codigo, nombre, descripcion, id_modulo, id_tipo_notificacion, orden_lista)
SELECT
    'EV-001', 'Orden de mesa lista para servir',
    'Se envía cuando una orden de mesa pasa a "Listo para servir".',
    m.id, t.id, 1
FROM configuracion.cfg_modulos m
CROSS JOIN notificaciones.ref_tipo_notificacion t
WHERE m.codigo = 'MD-008'
AND t.nombre = 'Informativa'
ON CONFLICT (codigo) DO NOTHING;

-- ─── Función "Parámetros generales" dentro de Configuración (MD-002) ──
INSERT INTO configuracion.cfg_funciones
    (id_producto, id_modulo, nombre, orden_lista, ruta_acceso, is_active, created_by, id_empresa)
SELECT
    m.id_producto, m.id, 'Parámetros generales',
    COALESCE((SELECT MAX(f.orden_lista) FROM configuracion.cfg_funciones f WHERE f.id_modulo = m.id), 0) + 1,
    '/config/parametros-generales', true, 1, 1
FROM configuracion.cfg_modulos m
WHERE m.codigo = 'MD-002'
AND NOT EXISTS (
    SELECT 1 FROM configuracion.cfg_funciones f
    WHERE f.id_modulo = m.id AND f.ruta_acceso = '/config/parametros-generales'
);
