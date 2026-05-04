CREATE SCHEMA configuracion;

CREATE TABLE IF NOT EXISTS configuracion.ref_tipo_persona(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO configuracion.ref_tipo_persona (nombre) VALUES ('Natural') ON CONFLICT(nombre) DO NOTHING;
INSERT INTO configuracion.ref_tipo_persona (nombre) VALUES ('Juridica') ON CONFLICT(nombre) DO NOTHING;

CREATE TABLE IF NOT EXISTS configuracion.ref_tipo_documento (
    id SERIAL PRIMARY KEY,
    nombre varchar(50) NOT NULL,
    codigo varchar(10) NOT NULL UNIQUE,
    tipo INT NOT NULL
);

INSERT INTO configuracion.ref_tipo_documento (nombre, codigo, tipo) 
VALUES ('Tarjeta de identidad', 'TI', 1) ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_tipo_documento (nombre, codigo, tipo) 
VALUES ('Cédula de ciudadania', 'CC', 1) ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_tipo_documento (nombre, codigo, tipo)
 VALUES ('Pasaporte', 'PP', 1) ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_tipo_documento (nombre, codigo, tipo) 
VALUES ('Cédula de extranjería', 'CE', 1) ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_tipo_documento (nombre, codigo, tipo)
VALUES ('Número de Identificación Tributaria', 'NIT', 2) ON CONFLICT(codigo) DO NOTHING;

CREATE TABLE IF NOT EXISTS configuracion.ref_genero (
    id SERIAL PRIMARY KEY NOT NULL,
    nombre varchar(10) NOT NULl
);

INSERT INTO configuracion.ref_genero (nombre)
VALUES ('Masculino'), ('Femenino'), ('No Aplica');

--tabla departamentos
CREATE TABLE IF NOT EXISTS configuracion.ref_departamentos(
    id SERIAL PRIMARY KEY NOT NULL,
    codigo INT NOT NULL UNIQUE,
    nombre VARCHAR(150) NOT NULL
);

--tabla municipios
CREATE TABLE IF NOT EXISTS configuracion.ref_municipios(
    id SERIAL PRIMARY KEY NOT NULL,
    codigo_departamento INT NOT NULL,
    codigo_municipio INT NOT NULL UNIQUE,
    nombre_municipio varchar(150) NOT NULL
);

CREATE TABLE IF NOT EXISTS configuracion.ref_zona_residencial(
    id SERIAL PRIMARY KEY NOT NULL,
    codigo varchar(5) NOT NULL,
    nombre varchar(150) NOT NULL
);

INSERT INTO configuracion.ref_zona_residencial (codigo, nombre) 
VALUES ('U', 'Urbana');
INSERT INTO configuracion.ref_zona_residencial (codigo, nombre) 
VALUES ('R', 'Rural');

CREATE TABLE IF NOT EXISTS configuracion.ref_actividades_economicas(
    id SERIAL PRIMARY KEY,
    codigo VARCHAR(12) NOT NULL UNIQUE,
    nombre TEXT NOT NULL
);


INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0111', 'Cultivo de cereales (excepto arroz), legumbres y semillas oleaginosas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0112', 'Cultivo de arroz.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0113', 'Cultivo de hortalizas, raíces y tubérculos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0114', 'Cultivo de tabaco.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0115', 'Cultivo de plantas textiles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0119', 'Otros cultivos transitorios n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0121', 'Cultivo de frutas tropicales y subtropicales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0122', 'Cultivo de plátano y banano.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0123', 'Cultivo de café.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0124', 'Cultivo de caña de azúcar.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0125', 'Cultivo de flor de corte.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0129', 'Otros cultivos permanentes n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0130', 'Propagación de plantas (actividades de los viveros, excepto viveros forestales).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0150', 'Explotación mixta (agrícola y pecuaria).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1020', 'Procesamiento y conservación de frutas, legumbres, hortalizas y tubérculos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1030', 'Elaboración de aceites y grasas de origen vegetal y animal.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1200', 'Elaboración de productos de tabaco.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1511', 'Curtido y recurtido de cueros; recurtido y teñido de pieles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1512', 'Fabricación de artículos de viaje, bolsos de mano y artículos similares elaborados en cuero, y fabricación de artículos de talabartería y guarnicionería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1521', 'Fabricación de calzado de cuero y piel, con cualquier tipo de suela.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1522', 'Fabricación de otros tipos de calzado, excepto calzado de cuero y piel.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1820', 'Producción de copias a partir de grabaciones originales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1910', 'Fabricación de productos de hornos de coque.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1921', 'Fabricación de productos de la refinación del petróleo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1922', 'Actividad de mezcla de combustibles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2030', 'Fabricación de fibras sintéticas y artificiales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2211', 'Fabricación de llantas y neumáticos de caucho') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2212', 'Reencauche de llantas usadas') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2219', 'Fabricación de formas básicas de caucho y otros productos de caucho n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2310', 'Fabricación de vidrio y productos de vidrio.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2421', 'Industrias básicas de metales preciosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2429', 'Industrias básicas de otros metales no ferrosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2511', 'Fabricación de productos metálicos para uso estructural.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2512', 'Fabricación de tanques, depósitos y recipientes de metal, excepto los utilizados para el envase o transporte de mercancías.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2513', 'Fabricación de generadores de vapor, excepto calderas de agua caliente para calefacción central.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2610', 'Fabricación de componentes y tableros electrónicos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2731', 'Fabricación de hilos y cables eléctricos y de fibra óptica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2732', 'Fabricación de dispositivos de cableado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2811', 'Fabricación de motores, turbinas, y partes para motores de combustión interna.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2812', 'Fabricación de equipos de potencia hidráulica y neumática.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2813', 'Fabricación de otras bombas, compresores, grifos y válvulas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2930', 'Fabricación de partes, piezas (autopartes) y accesorios (lujos) para vehículos automotores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3120', 'Fabricación de colchones y somieres.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3210', 'Fabricación de joyas, bisutería y artículos conexos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3220', 'Fabricación de instrumentos musicales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3230', 'Fabricación de artículos y equipo para la práctica del deporte.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3311', 'Mantenimiento y reparación especializado de productos elaborados en metal.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3312', 'Mantenimiento y reparación especializado de maquinaria y equipo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3313', 'Mantenimiento y reparación especializado de equipo electrónico y óptico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3320', 'Instalación especializada de maquinaria y equipo industrial.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3511', 'Generación de energía eléctrica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3512', 'Transmisión de energía eléctrica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3520', 'Producción de gas; distribución de combustibles gaseosos por tuberías.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3530', 'Suministro de vapor y aire acondicionado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4511', 'Comercio de vehículos automotores nuevos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4512', 'Comercio de vehículos automotores usados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4530', 'Comercio de partes, piezas (autopartes) y accesorios (lujos) para vehículos automotores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4541', 'Comercio de motocicletas y de sus partes, piezas y accesorios.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4542', 'Mantenimiento y reparación de motocicletas y de sus partes y piezas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5011', 'Transporte de pasajeros marítimo y de cabotaje.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5012', 'Transporte de carga marítimo y de cabotaje.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5111', 'Transporte aéreo nacional de pasajeros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5112', 'Transporte aéreo internacional de pasajeros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5121', 'Transporte aéreo nacional de carga.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5122', 'Transporte aéreo internacional de carga.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5221', 'Actividades de estaciones, vías y servicios complementarios para el transporte terrestre.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5222', 'Actividades de puertos y servicios complementarios para el transporte acuático.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5223', 'Actividades de aeropuertos, servicios de navegación aérea y demás actividades conexas al transporte aéreo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5224', 'Manipulación de carga.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5229', 'Otras actividades complementarias al transporte.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5511', 'Alojamiento en hoteles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5512', 'Alojamiento en apartahoteles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5513', 'Alojamiento en centros vacacionales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5519', 'Otros tipos de alojamientos para visitantes.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5530', 'Servicio por horas') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6010', 'Actividades de programación y transmisión en el servicio de radiodifusión sonora.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6120', 'Actividades de telecomunicaciones inalámbricas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6411', 'Banco Central.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6412', 'Bancos comerciales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6421', 'Actividades de las corporaciones financieras.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6422', 'Actividades de las compañías de financiamiento.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6423', 'Banca de segundo piso.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6424', 'Actividades de las cooperativas financieras.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6511', 'Seguros generales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6512', 'Seguros de vida.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6513', 'Reaseguros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6514', 'Capitalización.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7010', 'Actividades de administración empresarial.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7020', 'Actividades de consultoría de gestión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7210', 'Investigaciones y desarrollo experimental en el campo de las ciencias naturales y la ingeniería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7220', 'Investigaciones y desarrollo experimental en el campo de las ciencias sociales y las humanidades.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7310', 'Publicidad.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7320', 'Estudios de mercado y realización de encuestas de opinión pública.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8030', 'Actividades de detectives e investigadores privados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8511', 'Educación de la primera infancia.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8512', 'Educación preescolar.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8513', 'Educación básica primaria.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9900', 'Actividades de organizaciones y entidades extraterritoriales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0010', 'Asalariados') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0081', 'Personas Naturales sin Actividad Económica') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0082', 'Personas Naturales Subsidiadas por Terceros') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0090', 'Rentistas de Capital, solo para personas naturales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0126', 'Cultivo de palma para aceite (palma africana) y otros frutos oleaginosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0127', 'Cultivo de plantas con las que se preparan bebidas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0128', 'Cultivo de especias y de plantas aromáticas y medicinales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0141', 'Cría de ganado bovino y bufalino.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0142', 'Cría de caballos y otros equinos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0143', 'Cría de ovejas y cabras.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0144', 'Cría de ganado porcino.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0145', 'Cría de aves de corral.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0149', 'Cría de otros animales n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0161', 'Actividades de apoyo a la agricultura.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0162', 'Actividades de apoyo a la ganadería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0163', 'Actividades posteriores a la cosecha.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0164', 'Tratamiento de semillas para propagación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0170', 'Caza ordinaria y mediante trampas y actividades de servicios conexas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0210', 'Silvicultura y otras actividades forestales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0220', 'Extracción de madera.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0230', 'Recolección de productos forestales diferentes a la madera.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0240', 'Servicios de apoyo a la silvicultura.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0311', 'Pesca marítima.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0312', 'Pesca de agua dulce.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0321', 'Acuicultura marítima.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0322', 'Acuicultura de agua dulce.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0510', 'Extracción de hulla (carbón de piedra).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0520', 'Extracción de carbón lignito.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0610', 'Extracción de petróleo crudo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0620', 'Extracción de gas natural.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0710', 'Extracción de minerales de hierro.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0721', 'Extracción de minerales de uranio y de torio.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0722', 'Extracción de oro y otros metales preciosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0723', 'Extracción de minerales de níquel.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0729', 'Extracción de otros minerales metalíferos no ferrosos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0811', 'Extracción de piedra, arena, arcillas comunes, yeso y anhidrita.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0812', 'Extracción de arcillas de uso industrial, caliza, caolín y bentonitas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0820', 'Extracción de esmeraldas, piedras preciosas y semipreciosas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0891', 'Extracción de minerales para la fabricación de abonos y productos químicos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0892', 'Extracción de halita (sal).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0899', 'Extracción de otros minerales no metálicos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0910', 'Actividades de apoyo para la extracción de petróleo y de gas natural.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('0990', 'Actividades de apoyo para otras actividades de explotación de minas y canteras.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1011', 'Procesamiento y conservación de carne y productos cárnicos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1012', 'Procesamiento y conservación de pescados, crustáceos y moluscos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1040', 'Elaboración de productos lácteos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1051', 'Elaboración de productos de molinería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1052', 'Elaboración de almidones y productos derivados del almidón.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1061', 'Trilla de café.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1062', 'Descafeinado, tostión y molienda del café.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1063', 'Otros derivados del café.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1071', 'Elaboración y refinación de azúcar.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1072', 'Elaboración de panela.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1081', 'Elaboración de productos de panadería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1082', 'Elaboración de cacao, chocolate y productos de confitería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1083', 'Elaboración de macarrones, fideos, alcuzcuz y productos farináceos similares.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1084', 'Elaboración de comidas y platos preparados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1089', 'Elaboración de otros productos alimenticios n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1090', 'Elaboración de alimentos preparados para animales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1101', 'Destilación, rectificación y mezcla de bebidas alcohólicas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1102', 'Elaboración de bebidas fermentadas no destiladas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1103', 'Producción de malta, elaboración de cervezas y otras bebidas malteadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1104', 'Elaboración de bebidas no alcohólicas, producción de aguas minerales y de otras aguas embotelladas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1311', 'Preparación e hilatura de fibras textiles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1312', 'Tejeduría de productos textiles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1313', 'Acabado de productos textiles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1391', 'Fabricación de tejidos de punto y ganchillo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1392', 'Confección de artículos con materiales textiles, excepto prendas de vestir.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1393', 'Fabricación de tapetes y alfombras para pisos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1394', 'Fabricación de cuerdas, cordeles, cables, bramantes y redes.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1399', 'Fabricación de otros artículos textiles n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1410', 'Confección de prendas de vestir, excepto prendas de piel.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1420', 'Fabricación de artículos de piel.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1430', 'Fabricación de artículos de punto y ganchillo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1513', 'Fabricación de artículos de viaje, bolsos de mano y artículos similares; artículos de talabartería y guarnicionería elaborados en otros materiales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1523', 'Fabricación de partes del calzado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1610', 'Aserrado, acepillado e impregnación de la madera.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1620', 'Fabricación de hojas de madera para enchapado; fabricación de tableros contrachapados, tableros laminados, tableros de partículas y otros tableros y paneles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1630', 'Fabricación de partes y piezas de madera, de carpintería y ebanistería para la construcción.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1640', 'Fabricación de recipientes de madera.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1690', 'Fabricación de otros productos de madera; fabricación de artículos de corcho, cestería y espartería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1701', 'Fabricación de pulpas (pastas) celulósicas; papel y cartón.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1702', 'Fabricación de papel y cartón ondulado (corrugado); fabricación de envases, empaques y de embalajes de papel y cartón.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1709', 'Fabricación de otros artículos de papel y cartón.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1811', 'Actividades de impresión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('1812', 'Actividades de servicios relacionados con la impresión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2011', 'Fabricación de sustancias y productos químicos básicos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2012', 'Fabricación de abonos y compuestos inorgánicos nitrogenados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2013', 'Fabricación de plásticos en formas primarias.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2014', 'Fabricación de caucho sintético en formas primarias.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2021', 'Fabricación de plaguicidas y otros productos químicos de uso agropecuario.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2022', 'Fabricación de pinturas, barnices y revestimientos similares, tintas para impresión y masillas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2023', 'Fabricación de jabones y detergentes, preparados para limpiar y pulir; perfumes y preparados de tocador.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2029', 'Fabricación de otros productos químicos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2100', 'Fabricación de productos farmacéuticos, sustancias químicas medicinales y productos botánicos de uso farmacéutico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2221', 'Fabricación de formas básicas de plástico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2229', 'Fabricación de artículos de plástico n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2391', 'Fabricación de productos refractarios.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2392', 'Fabricación de materiales de arcilla para la construcción.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2393', 'Fabricación de otros productos de cerámica y porcelana.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2394', 'Fabricación de cemento, cal y yeso.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2395', 'Fabricación de artículos de hormigón, cemento y yeso.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2396', 'Corte, tallado y acabado de la piedra.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2399', 'Fabricación de otros productos minerales no metálicos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2410', 'Industrias básicas de hierro y de acero.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2431', 'Fundición de hierro y de acero.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2432', 'Fundición de metales no ferrosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2520', 'Fabricación de armas y municiones.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2591', 'Forja, prensado, estampado y laminado de metal; pulvimetalurgia.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2592', 'Tratamiento y revestimiento de metales; mecanizado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2593', 'Fabricación de artículos de cuchillería, herramientas de mano y artículos de ferretería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2599', 'Fabricación de otros productos elaborados de metal n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2620', 'Fabricación de computadoras y de equipo periférico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2630', 'Fabricación de equipos de comunicación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2640', 'Fabricación de aparatos electrónicos de consumo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2651', 'Fabricación de equipo de medición, prueba, navegación y control.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2652', 'Fabricación de relojes.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2660', 'Fabricación de equipo de irradiación y equipo electrónico de uso médico y terapéutico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2670', 'Fabricación de instrumentos ópticos y equipo fotográfico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2680', 'Fabricación de medios magnéticos y ópticos para almacenamiento de datos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2711', 'Fabricación de motores, generadores y transformadores eléctricos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2712', 'Fabricación de aparatos de distribución y control de la energía eléctrica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2720', 'Fabricación de pilas, baterías y acumuladores eléctricos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2740', 'Fabricación de equipos eléctricos de iluminación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2750', 'Fabricación de aparatos de uso doméstico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2790', 'Fabricación de otros tipos de equipo eléctrico n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2814', 'Fabricación de cojinetes, engranajes, trenes de engranajes y piezas de transmisión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2815', 'Fabricación de hornos, hogares y quemadores industriales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2816', 'Fabricación de equipo de elevación y manipulación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2817', 'Fabricación de maquinaria y equipo de oficina (excepto computadoras.)') ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2818', 'Fabricación de herramientas manuales con motor.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2819', 'Fabricación de otros tipos de maquinaria y equipo de uso general n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2821', 'Fabricación de maquinaria agropecuaria y forestal.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2822', 'Fabricación de máquinas formadoras de metal y de máquinas herramienta.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2823', 'Fabricación de maquinaria para la metalurgia.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2824', 'Fabricación de maquinaria para explotación de minas y canteras y para obras de construcción.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2825', 'Fabricación de maquinaria para la elaboración de alimentos, bebidas y tabaco.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2826', 'Fabricación de maquinaria para la elaboración de productos textiles, prendas de vestir y cueros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2829', 'Fabricación de otros tipos de maquinaria y equipo de uso especial n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2910', 'Fabricación de vehículos automotores y sus motores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('2920', 'Fabricación de carrocerías para vehículos automotores; fabricación de remolques y semirremolques.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3011', 'Construcción de barcos y de estructuras flotantes.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3012', 'Construcción de embarcaciones de recreo y deporte.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3020', 'Fabricación de locomotoras y de material rodante para ferrocarriles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3030', 'Fabricación de aeronaves, naves espaciales y de maquinaria conexa.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3040', 'Fabricación de vehículos militares de combate.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3091', 'Fabricación de motocicletas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3092', 'Fabricación de bicicletas y de sillas de ruedas para personas con discapacidad.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3099', 'Fabricación de otros tipos de equipo de transporte n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3110', 'Fabricación de muebles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3240', 'Fabricación de juegos, juguetes y rompecabezas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3250', 'Fabricación de instrumentos, aparatos y materiales médicos y odontológicos (incluido mobiliario).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3290', 'Otras industrias manufactureras n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3314', 'Mantenimiento y reparación especializado de equipo eléctrico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3315', 'Mantenimiento y reparación especializado de equipo de transporte, excepto los vehículos automotores, motocicletas y bicicletas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3319', 'Mantenimiento y reparación de otros tipos de equipos y sus componentes n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3513', 'Distribución de energía eléctrica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3514', 'Comercialización de energía eléctrica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3600', 'Captación, tratamiento y distribución de agua.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3700', 'Evacuación y tratamiento de aguas residuales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3811', 'Recolección de desechos no peligrosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3812', 'Recolección de desechos peligrosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3821', 'Tratamiento y disposición de desechos no peligrosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3822', 'Tratamiento y disposición de desechos peligrosos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3830', 'Recuperación de materiales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('3900', 'Actividades de saneamiento ambiental y otros servicios de gestión de desechos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4111', 'Construcción de edificios residenciales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4112', 'Construcción de edificios no residenciales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4210', 'Construcción de carreteras y vías de ferrocarril.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4220', 'Construcción de proyectos de servicio público.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4290', 'Construcción de otras obras de ingeniería civil.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4311', 'Demolición.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4312', 'Preparación del terreno.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4321', 'Instalaciones eléctricas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4322', 'Instalaciones de fontanería, calefacción y aire acondicionado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4329', 'Otras instalaciones especializadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4330', 'Terminación y acabado de edificios y obras de ingeniería civil.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4390', 'Otras actividades especializadas para la construcción de edificios y obras de ingeniería civil.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4520', 'Mantenimiento y reparación de vehículos automotores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4610', 'Comercio al por mayor a cambio de una retribución o por contrata.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4620', 'Comercio al por mayor de materias primas agropecuarias; animales vivos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4631', 'Comercio al por mayor de productos alimenticios.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4632', 'Comercio al por mayor de bebidas y tabaco.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4641', 'Comercio al por mayor de productos textiles, productos confeccionados para uso doméstico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4642', 'Comercio al por mayor de prendas de vestir.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4643', 'Comercio al por mayor de calzado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4644', 'Comercio al por mayor de aparatos y equipo de uso doméstico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4645', 'Comercio al por mayor de productos farmacéuticos, medicinales, cosméticos y de tocador.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4649', 'Comercio al por mayor de otros utensilios domésticos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4651', 'Comercio al por mayor de computadores, equipo periférico y programas de informática.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4652', 'Comercio al por mayor de equipo, partes y piezas electrónicos y de telecomunicaciones.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4653', 'Comercio al por mayor de maquinaria y equipo agropecuarios.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4659', 'Comercio al por mayor de otros tipos de maquinaria y equipo n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4661', 'Comercio al por mayor de combustibles sólidos, líquidos, gaseosos y productos conexos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4662', 'Comercio al por mayor de metales y productos metalíferos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4663', 'Comercio al por mayor de materiales de construcción, artículos de ferretería, pinturas, productos de vidrio, equipo y materiales de fontanería y calefacción.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4664', 'Comercio al por mayor de productos químicos básicos, cauchos y plásticos en formas primarias y productos químicos de uso agropecuario.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4665', 'Comercio al por mayor de desperdicios, desechos y chatarra.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4669', 'Comercio al por mayor de otros productos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4690', 'Comercio al por mayor no especializado.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4711', 'Comercio al por menor en establecimientos no especializados con surtido compuesto principalmente por alimentos, bebidas o tabaco.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4719', 'Comercio al por menor en establecimientos no especializados, con surtido compuesto principalmente por productos diferentes de alimentos (víveres en general), bebidas y tabaco.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4721', 'Comercio al por menor de productos agrícolas para el consumo en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4722', 'Comercio al por menor de leche, productos lácteos y huevos, en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4723', 'Comercio al por menor de carnes (incluye aves de corral), productos cárnicos, pescados y productos de mar, en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4724', 'Comercio al por menor de bebidas y productos del tabaco, en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4729', 'Comercio al por menor de otros productos alimenticios n.c.p., en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4731', 'Comercio al por menor de combustible para automotores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4732', 'Comercio al por menor de lubricantes (aceites, grasas), aditivos y productos de limpieza para vehículos automotores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4741', 'Comercio al por menor de computadores, equipos periféricos, programas de informática y equipos de telecomunicaciones en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4742', 'Comercio al por menor de equipos y aparatos de sonido y de video, en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4751', 'Comercio al por menor de productos textiles en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4752', 'Comercio al por menor de artículos de ferretería, pinturas y productos de vidrio en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4753', 'Comercio al por menor de tapices, alfombras y cubrimientos para paredes y pisos en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4754', 'Comercio al por menor de electrodomésticos y gasodomésticos de uso doméstico, muebles y equipos de iluminación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4755', 'Comercio al por menor de artículos y utensilios de uso doméstico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4759', 'Comercio al por menor de otros artículos domésticos en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4761', 'Comercio al por menor de libros, periódicos, materiales y artículos de papelería y escritorio, en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4762', 'Comercio al por menor de artículos deportivos, en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4769', 'Comercio al por menor de otros artículos culturales y de entretenimiento n.c.p. en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4771', 'Comercio al por menor de prendas de vestir y sus accesorios (incluye artículos de piel) en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4772', 'Comercio al por menor de todo tipo de calzado y artículos de cuero y sucedáneos del cuero en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4773', 'Comercio al por menor de productos farmacéuticos y medicinales, cosméticos y artículos de tocador en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4774', 'Comercio al por menor de otros productos nuevos en establecimientos especializados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4775', 'Comercio al por menor de artículos de segunda mano.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4781', 'Comercio al por menor de alimentos, bebidas y tabaco, en puestos de venta móviles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4782', 'Comercio al por menor de productos textiles, prendas de vestir y calzado, en puestos de venta móviles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4789', 'Comercio al por menor de otros productos en puestos de venta móviles.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4791', 'Comercio al por menor realizado a través de internet.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4792', 'Comercio al por menor realizado a través de casas de venta o por correo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4799', 'Otros tipos de comercio al por menor no realizado en establecimientos, puestos de venta o mercados.') ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4911', 'Transporte férreo de pasajeros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4912', 'Transporte férreo de carga.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4921', 'Transporte de pasajeros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4922', 'Transporte mixto.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4923', 'Transporte de carga por carretera.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('4930', 'Transporte por tuberías.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5021', 'Transporte fluvial de pasajeros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5022', 'Transporte fluvial de carga.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5210', 'Almacenamiento y depósito.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5310', 'Actividades postales nacionales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5320', 'Actividades de mensajería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5514', 'Alojamiento rural.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5520', 'Actividades de zonas de camping y parques para vehículos recreacionales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5590', 'Otros tipos de alojamiento n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5611', 'Expendio a la mesa de comidas preparadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5612', 'Expendio por autoservicio de comidas preparadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5613', 'Expendio de comidas preparadas en cafeterías.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5619', 'Otros tipos de expendio de comidas preparadas n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5621', 'Catering para eventos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5629', 'Actividades de otros servicios de comidas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5630', 'Expendio de bebidas alcohólicas para el consumo dentro del establecimiento.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5811', 'Edición de libros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5812', 'Edición de directorios y listas de correo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5813', 'Edición de periódicos, revistas y otras publicaciones periódicas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5819', 'Otros trabajos de edición.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5820', 'Edición de programas de informática (software).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5911', 'Actividades de producción de películas cinematográficas, videos, programas, anuncios y comerciales de televisión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5912', 'Actividades de posproducción de películas cinematográficas, videos, programas, anuncios y comerciales de televisión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5913', 'Actividades de distribución de películas cinematográficas, videos, programas, anuncios y comerciales de televisión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5914', 'Actividades de exhibición de películas cinematográficas y videos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('5920', 'Actividades de grabación de sonido y edición de música.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6020', 'Actividades de programación y transmisión de televisión.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6110', 'Actividades de telecomunicaciones alámbricas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6130', 'Actividades de telecomunicación satelital.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6190', 'Otras actividades de telecomunicaciones.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6201', 'Actividades de desarrollo de sistemas informáticos (planificación, análisis, diseño, programación, pruebas).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6202', 'Actividades de consultoría informática y actividades de administración de instalaciones informáticas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6209', 'Otras actividades de tecnologías de información y actividades de servicios informáticos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6311', 'Procesamiento de datos, alojamiento (hosting) y actividades relacionadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6312', 'Portales web.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6391', 'Actividades de agencias de noticias.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6399', 'Otras actividades de servicio de información n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6431', 'Fideicomisos, fondos y entidades financieras similares.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6432', 'Fondos de cesantías.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6491', 'Leasing financiero (arrendamiento financiero).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6492', 'Actividades financieras de fondos de empleados y otras formas asociativas del sector solidario.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6493', 'Actividades de compra de cartera o factoring.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6494', 'Otras actividades de distribución de fondos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6495', 'Instituciones especiales oficiales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6499', 'Otras actividades de servicio financiero, excepto las de seguros y pensiones n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6521', 'Servicios de seguros sociales de salud.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6522', 'Servicios de seguros sociales de riesgos profesionales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6531', 'Régimen de prima media con prestación definida (RPM).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6532', 'Régimen de ahorro individual (RAI).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6611', 'Administración de mercados financieros.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6612', 'Corretaje de valores y de contratos de productos básicos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6613', 'Otras actividades relacionadas con el mercado de valores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6614', 'Actividades de las casas de cambio.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6615', 'Actividades de los profesionales de compra y venta de divisas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6619', 'Otras actividades auxiliares de las actividades de servicios financieros n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6621', 'Actividades de agentes y corredores de seguros') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6629', 'Evaluación de riesgos y daños, y otras actividades de servicios auxiliares') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6630', 'Actividades de administración de fondos.') ON CONFLICT(codigo) DO NOTHING;

INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6810', 'Actividades inmobiliarias realizadas con bienes propios o arrendados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6820', 'Actividades inmobiliarias realizadas a cambio de una retribución o por contrata.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6910', 'Actividades jurídicas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('6920', 'Actividades de contabilidad, teneduría de libros, auditoría financiera y asesoría tributaria.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7110', 'Actividades de arquitectura e ingeniería y otras actividades conexas de consultoría técnica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7120', 'Ensayos y análisis técnicos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7410', 'Actividades especializadas de diseño.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7420', 'Actividades de fotografía.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7490', 'Otras actividades profesionales, científicas y técnicas n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7500', 'Actividades veterinarias.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7710', 'Alquiler y arrendamiento de vehículos automotores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7721', 'Alquiler y arrendamiento de equipo recreativo y deportivo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7722', 'Alquiler de videos y discos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7729', 'Alquiler y arrendamiento de otros efectos personales y enseres domésticos n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7730', 'Alquiler y arrendamiento de otros tipos de maquinaria, equipo y bienes tangibles n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7740', 'Arrendamiento de propiedad intelectual y productos similares, excepto obras protegidas por derechos de autor.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7810', 'Actividades de agencias de empleo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7820', 'Actividades de agencias de empleo temporal.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7830', 'Otras actividades de suministro de recurso humano.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7911', 'Actividades de las agencias de viaje.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7912', 'Actividades de operadores turísticos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('7990', 'Otros servicios de reserva y actividades relacionadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8010', 'Actividades de seguridad privada.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8020', 'Actividades de servicios de sistemas de seguridad.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8110', 'Actividades combinadas de apoyo a instalaciones.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8121', 'Limpieza general interior de edificios.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8129', 'Otras actividades de limpieza de edificios e instalaciones industriales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8130', 'Actividades de paisajismo y servicios de mantenimiento conexos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8211', 'Actividades combinadas de servicios administrativos de oficina.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8219', 'Fotocopiado, preparación de documentos y otras actividades especializadas de apoyo a oficina.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8220', 'Actividades de centros de llamadas (Call center).') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8230', 'Organización de convenciones y eventos comerciales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8291', 'Actividades de agencias de cobranza y oficinas de calificación crediticia.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8292', 'Actividades de envase y empaque.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8299', 'Otras actividades de servicio de apoyo a las empresas n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8411', 'Actividades legislativas de la administración pública.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8412', 'Actividades ejecutivas de la administración pública.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8413', 'Regulación de las actividades de organismos que prestan servicios de salud, educativos, culturales y otros servicios sociales, excepto servicios de seguridad social.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8414', 'Actividades reguladoras y facilitadoras de la actividad económica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8415', 'Actividades de los otros órganos de control.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8421', 'Relaciones exteriores.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8422', 'Actividades de defensa.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8423', 'Orden público y actividades de seguridad.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8424', 'Administración de justicia.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8430', 'Actividades de planes de seguridad social de afiliación obligatoria.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8521', 'Educación básica secundaria.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8522', 'Educación media académica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8523', 'Educación media técnica y de formación laboral.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8530', 'Establecimientos que combinan diferentes niveles de educación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8541', 'Educación técnica profesional.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8542', 'Educación tecnológica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8543', 'Educación de instituciones universitarias o de escuelas tecnológicas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8544', 'Educación de universidades.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8551', 'Formación académica no formal.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8552', 'Enseñanza deportiva y recreativa.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8553', 'Enseñanza cultural.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8559', 'Otros tipos de educación n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8560', 'Actividades de apoyo a la educación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8610', 'Actividades de hospitales y clínicas, con internación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8621', 'Actividades de la práctica médica, sin internación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8622', 'Actividades de la práctica odontológica.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8691', 'Actividades de apoyo diagnóstico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8692', 'Actividades de apoyo terapéutico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8699', 'Otras actividades de atención de la salud humana.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8710', 'Actividades de atención residencial medicalizada de tipo general.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8720', 'Actividades de atención residencial, para el cuidado de pacientes con retardo mental, enfermedad mental y consumo de sustancias psicoactivas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8730', 'Actividades de atención en instituciones para el cuidado de personas mayores y/o discapacitadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8790', 'Otras actividades de atención en instituciones con alojamiento') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8810', 'Actividades de asistencia social sin alojamiento para personas mayores y discapacitadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('8890', 'Otras actividades de asistencia social sin alojamiento.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9001', 'Creación literaria.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9002', 'Creación musical.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9003', 'Creación teatral.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9004', 'Creación audiovisual.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9005', 'Artes plásticas y visuales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9006', 'Actividades teatrales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9007', 'Actividades de espectáculos musicales en vivo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9008', 'Otras actividades de espectáculos en vivo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9101', 'Actividades de bibliotecas y archivos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9102', 'Actividades y funcionamiento de museos, conservación de edificios y sitios históricos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9103', 'Actividades de jardines botánicos, zoológicos y reservas naturales.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9200', 'Actividades de juegos de azar y apuestas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9311', 'Gestión de instalaciones deportivas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9312', 'Actividades de clubes deportivos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9319', 'Otras actividades deportivas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9321', 'Actividades de parques de atracciones y parques temáticos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9329', 'Otras actividades recreativas y de esparcimiento n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9411', 'Actividades de asociaciones empresariales y de empleadores') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9412', 'Actividades de asociaciones profesionales') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9420', 'Actividades de sindicatos de empleados.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9491', 'Actividades de asociaciones religiosas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9492', 'Actividades de asociaciones políticas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9499', 'Actividades de otras asociaciones n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9511', 'Mantenimiento y reparación de computadores y de equipo periférico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9512', 'Mantenimiento y reparación de equipos de comunicación.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9521', 'Mantenimiento y reparación de aparatos electrónicos de consumo.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9522', 'Mantenimiento y reparación de aparatos y equipos domésticos y de jardinería.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9523', 'Reparación de calzado y artículos de cuero.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9524', 'Reparación de muebles y accesorios para el hogar.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9529', 'Mantenimiento y reparación de otros efectos personales y enseres domésticos.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9601', 'Lavado y limpieza, incluso la limpieza en seco, de productos textiles y de piel.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9602', 'Peluquería y otros tratamientos de belleza.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9603', 'Pompas fúnebres y actividades relacionadas.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9609', 'Otras actividades de servicios personales n.c.p.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9700', 'Actividades de los hogares individuales como empleadores de personal doméstico.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9810', 'Actividades no diferenciadas de los hogares individuales como productores de bienes para uso propio.') ON CONFLICT(codigo) DO NOTHING;
INSERT INTO configuracion.ref_actividades_economicas (codigo, nombre) VALUES ('9820', 'Actividades no diferenciadas de los hogares individuales como productores de servicios para uso propio.') ON CONFLICT(codigo) DO NOTHING;


CREATE TABLE IF NOT EXISTS configuracion.ref_ambitos_terceros(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(150) NOT NULL
);

INSERT INTO  configuracion.ref_ambitos_terceros (nombre) VALUES ('Municipal');
INSERT INTO  configuracion.ref_ambitos_terceros (nombre) VALUES ('Departamental');
INSERT INTO  configuracion.ref_ambitos_terceros (nombre) VALUES ('Nacional');
INSERT INTO  configuracion.ref_ambitos_terceros (nombre) VALUES ('Internacional');

CREATE TABLE IF NOT EXISTS configuracion.ref_centralizacion(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(150) NOT NULL
);
INSERT INTO  configuracion.ref_centralizacion (nombre) VALUES ('Centralizado');
INSERT INTO  configuracion.ref_centralizacion (nombre) VALUES ('Descentralizado');
INSERT INTO  configuracion.ref_centralizacion (nombre) VALUES ('Mixto');

CREATE TABLE IF NOT EXISTS configuracion.ref_responsabilidades_dian(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(150) NOT NULL
);

INSERT INTO configuracion.ref_responsabilidades_dian (nombre) VALUES ('Gran contribuyente');
INSERT INTO configuracion.ref_responsabilidades_dian (nombre) VALUES ('Autorretenedor');
INSERT INTO configuracion.ref_responsabilidades_dian (nombre) VALUES ('Responsable de IVA');
INSERT INTO configuracion.ref_responsabilidades_dian (nombre) VALUES ('Régimen simple');
INSERT INTO configuracion.ref_responsabilidades_dian (nombre) VALUES ('Facturador electrónico');
INSERT INTO configuracion.ref_responsabilidades_dian (nombre) VALUES ('No responsable');

CREATE TABLE IF NOT EXISTS configuracion.ref_tipo_contribuyente(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(150) NOT NULL
);

INSERT INTO configuracion.ref_tipo_contribuyente (nombre) VALUES ('Persona Natural');
INSERT INTO configuracion.ref_tipo_contribuyente (nombre) VALUES ('Persona Jurídica');
INSERT INTO configuracion.ref_tipo_contribuyente (nombre) VALUES ('Gran Contribuyente');
INSERT INTO configuracion.ref_tipo_contribuyente (nombre) VALUES ('Régimen Simple');
INSERT INTO configuracion.ref_tipo_contribuyente (nombre) VALUES ('ESAL');

CREATE TABLE IF NOT EXISTS configuracion.ref_regimen_iva(
   id SERIAL PRIMARY KEY NOT NULL,
   nombre VARCHAR(150) NOT NULL 
);

INSERT INTO configuracion.ref_regimen_iva (nombre) VALUES ('Responsable de IVA');
INSERT INTO configuracion.ref_regimen_iva (nombre) VALUES ('No responsable de IVA');
INSERT INTO configuracion.ref_regimen_iva (nombre) VALUES ('Responsable INC');
INSERT INTO configuracion.ref_regimen_iva (nombre) VALUES ('Exento');
INSERT INTO configuracion.ref_regimen_iva (nombre) VALUES ('Excluido');

CREATE TABLE configuracion.ref_regimen_dian(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(150) NOT NULL 
);

INSERT INTO configuracion.ref_regimen_dian (nombre) VALUES ('Régimen Ordinario');
INSERT INTO configuracion.ref_regimen_dian (nombre) VALUES ('Régimen Simple');
INSERT INTO configuracion.ref_regimen_dian (nombre) VALUES ('Régimen Especial');
INSERT INTO configuracion.ref_regimen_dian (nombre) VALUES ('Gran Contribuyente');
INSERT INTO configuracion.ref_regimen_dian (nombre) VALUES ('No Contribuyente');

CREATE TABLE IF NOT EXISTS configuracion.ref_clase_personas(
    id SERIAL PRIMARY KEY NOT NULL,
    nombre VARCHAR(300) NOT NULL
);

INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Proveedor');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Cliente');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Empleado');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Contratista');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Distribuidor');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Laboratorio');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Entidad educativa');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Entidad de salud');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Operador logístico');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Aliado comercial');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Personal administrativo');
INSERT INTO configuracion.ref_clase_personas (nombre) VALUES ('Otro');

CREATE TABLE IF NOT EXISTS configuracion.cfg_terceros(
    id SERIAL PRIMARY KEY NOT NULL,
    id_tipo_persona INT NOT NULL REFERENCES configuracion.ref_tipo_persona(id),
    id_tipo_documento INT NOT NULL REFERENCES configuracion.ref_tipo_documento(id),
    numero_documento VARCHAR(20) NOT NULL UNIQUE,
    dv CHAR(1),
    primer_nombre VARCHAR(100),
    segundo_nombre VARCHAR(100),
    primer_apellido VARCHAR(100),
    segundo_apellido VARCHAR(100),
    razon_social VARCHAR(150),
    representante_legal VARCHAR(200),
    id_genero INT NOT NULL REFERENCES configuracion.ref_genero(id),
    fecha_nacimiento DATE NULL,
    telefono VARCHAR(20) NOT NULL,
    telefono_2 VARCHAR(20),
    email VARCHAR(100) NULL,
    pagina_web VARCHAR(300),
    direccion VARCHAR(350) NOT NULL,
    id_pais INT NOT NULL,
    id_departamento INT NOT NULL REFERENCES configuracion.ref_departamentos(id),
    id_ciudad INT NOT NULL REFERENCES configuracion.ref_municipios(id),
    id_zona INT NOT NULL REFERENCES configuracion.ref_zona_residencial(id),
    estado BOOLEAN DEFAULT TRUE,
    id_actividad_economica INT NULL REFERENCES configuracion.ref_actividades_economicas(id),
    id_ambito INT  NULL REFERENCES configuracion.ref_ambitos_terceros(id),
    id_centralizacion INT NULL REFERENCES configuracion.ref_centralizacion(id),
    numero_acciones VARCHAR(350) NULL,
    codigo_contable VARCHAR(450) NULL,
    id_responsabilidad_dian INT NULL REFERENCES configuracion.ref_responsabilidades_dian(id),
    id_tipo_contribuyente INT NULL REFERENCES configuracion.ref_tipo_contribuyente(id),
    id_regimen_iva INT NULL REFERENCES configuracion.ref_regimen_iva(id),
    id_regimen_dian INT NULL REFERENCES configuracion.ref_regimen_dian(id),
    created_by INT DEFAULT NULL,
    update_by INT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS configuracion.cfg_clase_terceros(
    id SERIAL PRIMARY KEY,
    id_tercero INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    id_clase INT NOT NULL REFERENCES configuracion.ref_clase_personas(id),
    created_by INT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS configuracion.hco_modificacion_documentos(
    id SERIAL PRIMARY KEY,
    id_tercero INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    tipo_persona_ant INT,
    tipo_persona_new INT,
    tipo_documento_ant INT,
    tipo_documento_new INT,
    numero_documento_ant VARCHAR(20),
    numero_documento_new VARCHAR(20),
    created_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

--schema de trabajo
CREATE SCHEMA IF NOT EXISTS seguridad;

CREATE TABLE IF NOT EXISTS seguridad.ref_tipo_rol(
   id SERIAL PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO seguridad.ref_tipo_rol (nombre) VALUES
     ('Administrativo'),
     ('Operativo'),
     ('Gerencial'),
     ('Financiero'),
     ('Soporte')
 ON CONFLICT(nombre) DO NOTHING;


--tabla de roles
CREATE TABLE IF NOT EXISTS seguridad.cfg_roles_usuario (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL,
    id_tipo INT NOT NULL REFERENCES seguridad.ref_tipo_rol(id),
    intentos_login INT NOT NULL,
    is_multisession BOOLEAN DEFAULT FALSE,
    is_externo BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT,
    id_empresa INT NOT NULL,
    id_sede INT NULL
);


INSERT INTO seguridad.cfg_roles_usuario (nombre, id_tipo, intentos_login, is_multisession, is_externo, created_by, id_empresa) 
VALUES ('Administrador', 1, 3, false, false, 1, 1);

--tabla principal de usuarios
CREATE TABLE IF NOT EXISTS seguridad.cfg_usuarios(
    id SERIAL PRIMARY KEY,
    email VARCHAR(150) NOT NULL,
    password VARCHAR(300) NOT NULL,
    id_rol INT NOT NULL REFERENCES seguridad.cfg_roles_usuario(id),
    id_tercero INT NOT NULL REFERENCES configuracion.cfg_terceros(id),
    image_profile text,
    activo BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP,
    created_by INT NOT NULL,
    update_by INT
);

--tabla historico de sessiones 
CREATE TABLE IF NOT EXISTS seguridad.hco_inicios_usuarios(
    id SERIAL PRIMARY KEY,
    id_usuario INT NOT NULL REFERENCES seguridad.cfg_usuarios(id),
    fecha_inicio TIMESTAMP DEFAULT NOW(),
    hora_inicio TIME NOT NULL,
    fecha_cierre TIMESTAMP,
    hora_cierre TIME,
    ip_maquina VARCHAR(150) NOT NULL
);
