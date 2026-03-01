# Sprint Plan - Version 0.0.1

**Objetivo del Sprint**  
Implementacion del **módulo de terceros** dentro del sistema de gestión de inventario y facturación. 
Los terceros son los registros de informacion basica para personas naturales y juridicas
Permitiendo la **Centralizacion** de informacion entre clientes proveedores e.t.c.

---

## Requerimientos Funcionales
1. Registrar una persona natural/juridica con la informacion basica (contactos, ubicacion, inf Tributaria)
2. Consultar Terceros de forma individual o listada, con filtros por estado rol y tipo
3. Actualizar informacion del registro. incluyendo cambios de numero identificacion y tipo identificacion
4. Guardar historial de modificaciones al registro (tabla historica)

--
## Alcanze Técnico

### Endpoints Nuevos

## Endpoint Generales
- `GET     /general/TipoPersona`          → Enlista tipo de personas
- `GET     /general/TipoDocument`         → Enlista tipo de documentos
- `GET     /general/Generos`              → Enlista los generos
- `GET     /general/Departamentos`        → Enlista los departamentos nacionales
- `GET     /general/Municipios/:id`       → Enlista los municipios nacionales por id departamento
- `GET     /general/ZonasResidenciales`   → Enlista las zonas residenciales
- `GET     /general/ActividadesEconomic`  → Enlista las actividades economicas nacionales
- `GET     /general/Ambitos`              → Enlista las actividades economicas nacionales
- `GET     /general/Centralizacion`       → Enlista las centralizaciones por contabilidad
- `GET     /general/ResponsabilidadDian`  → Enlista las responsabilidades ante la dia
- `GET     /general/TipoContribuyente`    → Enlista los tipos de contribuyentes
- `GET     /general/RegimenIva`           → Enlista los regimen de iva
- `GET     /general/RegimenDian`          → Enlista los regimen ante la DIAN
- `GET     /general/ClaseTerceros`        → Enlista las clases por tercero

### Endpoints Terceros
- `POST     /third-parties` → Registro del tercero
- `GET      /third-parties` → Búsuqeda de registros recientes
- `GET      /third-parties/:id` → Obtener detalle del registro por ID
- `PUT      /third-parties/:id` → Actualizar registro del tercero
- `PATCH    /third-parties/:id` → cambiar estado (activo/inactivo).

### Base de datos (Migración)
- Crear tabla `cfg_terceros` (entidad principal de terceros).
  - Campos minimos: id, nombre completo, identificacion, direcion, datos tributarios, rol_id, tercero_id, estado, created_by,            created_at.
-Crear tabla `ref_tipo_persona`           (catalogo de personas naturales/juridicas)
-Crear tabla `ref_tipo_documento`         (catalogo de documento ti/cc)
-Crear tabla `ref_tipo_persona`           (catalogo de personas naturales/juridicas)
-Crear tabla `ref_genero`                 (catalogo de generos)
-Crear tabla `ref_departamentos`          (catalogo de departamentos nacionales `tabla general`)
-Crear tabla `ref_municipios`             (catalogo de municipios nacionales `tabla general`)
-Crear tabla `ref_zona_residencial`       (catalogo de zonas rurles `tabla general`)
-Crear tabla `ref_actividades_economicas` (catalogo de actividades economicas nacionales `tabla general`)
-Crear tabla `ref_ambitos_terceros`       (catalogo de ambitos por terceros)
-Crear tabla `ref_centralizado`           (catalogo de centralizaciones por contabilidad)
-Crear tabla `ref_responsabilidades_dian` (catalogo de responsabilidades dian nacionales `tabla general`)
-Crear tabla `ref_tipo_contribuyente`     (catalogo de tipo contribuyentes segun dian)
-Crear tabla `ref_regimen_iva`            (catalogo de regimen iva nacional )
-Crear tabla `ref_regimen_dian`           (catalogo de regimen dian)
-Crear tabla `ref_clase_personas`         (catalogo de clases para terceros `tabla general`)
-Crear tabla `cfg_clase_terceros`         (tabla de configuracion entre terceros y sus clases relacionadas)
-Crear tabla `hco_terceros`               (tabla de movimientos historicos por registro `tabla historica de respaldo`)


## Criterios de activacion
- No se puede registrar un tercero con numero identificacion previamente registrado
- Debe peritir seleccionar multiples clases para el tercero
- El Tercero debe estar activo para lograr ser ubicados en futuras busquedas

---

## Observaciones
- Este sprint sienta las bases para la **gestión de terceros y control de informacion** en el sistema.
con relacion a la informacion basica de proveedores/clientes/empleado e.tc



