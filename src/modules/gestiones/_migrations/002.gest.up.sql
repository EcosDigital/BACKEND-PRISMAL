-- ==========================================================
-- MÓDULO GESTIONES — Políticas ANS y Niveles del Caso
-- ==========================================================

-- ----------------------------------------------------------
-- REFERENCIAL: Unidad de tiempo
-- Usada por cfg_politicas_ans para expresar los tiempos.
-- ----------------------------------------------------------
CREATE TABLE gestiones.ref_unidades_tiempo (
    id             SMALLINT    PRIMARY KEY,
    codigo         VARCHAR(10) NOT NULL,
    nombre         VARCHAR(40) NOT NULL,
    factor_minutos INT         NOT NULL,  -- Conversión a minutos para cálculos
    CONSTRAINT uq_ref_unidades_tiempo_codigo UNIQUE (codigo)
);

INSERT INTO gestiones.ref_unidades_tiempo (id, codigo, nombre, factor_minutos)
VALUES
    (1, '001', 'Minuto',  1),
    (2, '002',   'Hora',    60),
    (3, '003',    'Día',     1440),
    (4, '004', 'Semana',  10080)
ON CONFLICT (id) DO NOTHING;

-- ----------------------------------------------------------
-- CONFIGURACIÓN: Políticas ANS
-- Define los tiempos máximos de atención.
-- Se asigna a un nivel del caso.
-- ----------------------------------------------------------
CREATE TABLE gestiones.cfg_politicas_ans (
    id          SERIAL       PRIMARY KEY,
    codigo      VARCHAR(30)  NOT NULL,
    nombre      VARCHAR(100) NOT NULL,
    descripcion TEXT         NULL,
    -- Tiempo máximo de primera respuesta
    tiempo_respuesta    INT      NOT NULL CHECK (tiempo_respuesta > 0),
    id_unidad_respuesta SMALLINT NOT NULL REFERENCES gestiones.ref_unidades_tiempo(id),
    -- Tiempo máximo de resolución
    tiempo_resolucion    INT   NOT NULL CHECK (tiempo_resolucion > 0),
    id_unidad_resolucion SMALLINT NOT NULL REFERENCES gestiones.ref_unidades_tiempo(id),
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by  INT         NOT NULL,
    updated_at  TIMESTAMP NULL,
    updated_by  INT         NULL,

    CONSTRAINT uq_cfg_politicas_ans_codigo UNIQUE (codigo)
);

COMMENT ON TABLE  gestiones.cfg_politicas_ans                    IS 'Políticas de Acuerdo de Nivel de Servicio (ANS/SLA)';
COMMENT ON COLUMN gestiones.cfg_politicas_ans.tiempo_respuesta   IS 'Tiempo máximo para dar la primera respuesta al cliente';
COMMENT ON COLUMN gestiones.cfg_politicas_ans.id_unidad_respuesta IS 'FK → ref_unidades_tiempo. Unidad del tiempo_respuesta';
COMMENT ON COLUMN gestiones.cfg_politicas_ans.tiempo_resolucion  IS 'Tiempo máximo para resolver el caso completamente';
COMMENT ON COLUMN gestiones.cfg_politicas_ans.id_unidad_resolucion IS 'FK → ref_unidades_tiempo. Unidad del tiempo_resolucion';

INSERT INTO gestiones.cfg_politicas_ans
    (codigo, nombre, descripcion, tiempo_respuesta, id_unidad_respuesta, tiempo_resolucion, id_unidad_resolucion, created_by)
VALUES
    ('001',       'ANS Bajo',       'Para casos de impacto mínimo',              48, 2,  30, 3, 1),
    ('002',       'ANS Medio',      'Para casos de impacto moderado',             8, 2,  10, 3, 1),
    ('003',       'ANS Alto',       'Para casos de impacto significativo',        4, 2,   5, 3, 1),
    ('004',       'ANS Prioritario','Para casos críticos que afectan operaciones', 1, 2,   1, 3, 1)
ON CONFLICT (codigo) DO NOTHING;

-- ----------------------------------------------------------
-- CONFIGURACIÓN: Niveles del caso
-- Cada nivel lleva su política ANS asociada.
-- ----------------------------------------------------------

CREATE TABLE gestiones.cfg_niveles_caso (
    id          SERIAL      PRIMARY KEY,
    codigo      VARCHAR(20) NOT NULL,
    nombre      VARCHAR(80) NOT NULL,
    descripcion TEXT        NULL,
    color_hex   CHAR(10)    NOT NULL DEFAULT '#6B7280',
    orden       SMALLINT    NOT NULL DEFAULT 0, 
    id_politica_ans INT     NULL REFERENCES gestiones.cfg_politicas_ans(id),
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  INT         NOT NULL,
    updated_at  TIMESTAMPTZ NULL,
    updated_by  INT         NULL,

    CONSTRAINT uq_cfg_niveles_caso_codigo UNIQUE (codigo)
);

COMMENT ON TABLE  gestiones.cfg_niveles_caso                 IS 'Niveles de prioridad o impacto de un ticket';
COMMENT ON COLUMN gestiones.cfg_niveles_caso.color_hex       IS 'Color hexadecimal para badges en la UI. Ej: #EF4444 para crítico';
COMMENT ON COLUMN gestiones.cfg_niveles_caso.orden           IS 'Orden ascendente de criticidad: 1=más bajo, N=más crítico';
COMMENT ON COLUMN gestiones.cfg_niveles_caso.id_politica_ans IS 'FK → cfg_politicas_ans. ANS aplicable a los tickets de este nivel';


-- ----------------------------------------------------------
-- DATOS INICIALES
--
--  Nivel       Color     Orden  ANS asociado
--  BAJO        #22C55E   1      48h respuesta / 30d resolución
--  MEDIO       #F59E0B   2       8h respuesta / 10d resolución
--  ALTO        #F97316   3       4h respuesta /  5d resolución
--  PRIORITARIO #EF4444   4       1h respuesta /  1d resolución
-- ----------------------------------------------------------

INSERT INTO gestiones.cfg_niveles_caso
    (codigo, nombre, descripcion, color_hex, orden, id_politica_ans, created_by)
VALUES
    ('N001',        'Bajo',        'Impacto mínimo, no interrumpe operaciones',             '#22C55E', 1, 1, 1),
    ('N002',        'Medio',       'Impacto moderado, operación parcialmente afectada',      '#F59E0B', 2, 2, 1),
    ('N003',        'Alto',        'Impacto significativo, operación considerablemente afectada', '#F97316', 3, 3, 1),
    ('N004',        'Prioritario', 'Impacto crítico, operación completamente detenida',      '#EF4444', 4, 4, 1)
ON CONFLICT (codigo) DO NOTHING;