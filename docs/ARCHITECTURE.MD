## Arquitectura del Proyecto

Este proyecto sigue una arquitectura en capas simple que puede evolucionar hacia Clean Architecture.
El proyecto sigue una estructura en capas:


```bash
-cmd/: contiene el punto de entrada al servidor (`main.go`).
-migrations/: Archivos `.sql` que define la evolucion de la base de datos. (migraciones automaticas)
-src/config/: Manejo de la configuracion global (lectura de variables `.env`).
-src/database/: Gestiones directas en la base de datos como conexion y migraciones
-src/loggin/: Sistema de logs por consola y archivo (`app.log`) 
-src/routes/: define enpoints (rutas backend)
-src/controllers/: validan request y llaman servicios
-src/services/: contienen la logica de negocio
-src/repositorie/: abstraen acceso a base de datos
-src/midlwares: porciones de codigo reutilizables para validaciones
-src/schema/: estructuras de datos a validar para los request

## Estructura actual
core
├── README.MD
├── cmd
│   └── main.go
├── docs
│   ├── ARCHITECTURE.MD
│   └── SPRINT.MD
├── go.mod
├── go.sum
├── licence.md
├── logs
│   └── app.log
├── migrations
│   ├── 001_init.up.sql
│   ├── 002_init.up.sql
│   ├── 003_init.up.sql
│   ├── 004_init.up.sql
│   ├── 005_init.up.sql
│   ├── 006_init.up.sql
│   ├── 007_init.up.sql
│   ├── 008_init.up.sql
│   ├── 009_init.up.sql
│   ├── 010_init.up.sql
│   └── 011_init.up.sql
├── public
│   └── uploads
│       └── 136efea4-5b5c-4ad4-b284-c7f576a0da55.png
└── src
    ├── config
    │   ├── configuracion.go
    │   └── redis.go
    ├── controllers
    │   ├── configuracion
    │   │   ├── empresa.controller.go
    │   │   ├── terceros.controller.go
    │   │   └── upload.controller.go
    │   ├── general.controllers.go
    │   └── seguridad
    │       ├── auth.controller.go
    │       └── usuarios.controller.go
    ├── database
    │   ├── database.go
    │   └── migrates.go
    ├── logging
    │   └── logger.go
    ├── middlewares
    │   ├── auth_user.go
    │   ├── jwt.go
    │   └── validator.go
    ├── repository
    │   ├── configuracion
    │   │   ├── empresas.repository.go
    │   │   └── terceros.repository.go
    │   ├── general.reposiotory.go
    │   └── seguridad
    │       ├── auth.repository.go
    │       └── usuarios.repository.go
    ├── routes
    │   ├── configuracion
    │   │   ├── empresas.routes.go
    │   │   ├── general.routes.go
    │   │   ├── terceros.routes.go
    │   │   └── upload.routes.go
    │   ├── index.go
    │   └── seguridad
    │       ├── auth.routes.go
    │       └── usuarios.routes.go
    ├── schemas
    │   ├── configuracion
    │   │   ├── empresas.schemas.go
    │   │   ├── general.schema.go
    │   │   └── terceros.schemas.go
    │   └── seguridad
    │       ├── auth.schema.go
    │       └── usuarios.schema.go
    ├── services
    │   ├── configuracion
    │   │   ├── empresa.services.go
    │   │   └── terceros.services.go
    │   ├── general.services.go
    │   └── seguridad
    │       ├── autth.services.go
    │       └── usuarios.services.go
    └── utils
        ├── session_manager.go
        ├── token.go
        ├── upload.go
        └── utils.go

## 🔄 Flujo general de una petición
1. **Cliente →** solicitud HTTP.  
2. **Router →** determina qué controlador manejará la petición.  
3. **Midleware →** Valida los request en base a un schema
4. **Controller →** valida request y llama al service.  
5. **Service →** aplica la lógica de negocio y se comunica con el repositorio.  
6. **Repository →** ejecuta queries en PostgreSQL.  
7. **Logger →** registra información de la petición y respuesta.  
8. **Response →** se envía al cliente.


## Middlewares implementados
- `logger`: Registra cada request y código de respuesta en consola.
- `recover`: Captura errores inesperados (panic) y evita que el servidor se caiga.
- `CORS` : Restringe peticiones al dominio configurado en `FRONTEND_URL` dentro de `.env`.
- `logging`: registra informacion en consola y en `logs/app.log`.
    - `Info`: logs informativos (inicio de servidor, respuestas exitosas).
    - `Error`: logs de error en caso de fallos.
- `migrates.go` : Ejecuta migraciones con `golang-migrate` al iniciar el servidor
- `validator.go` : realiza validaciones de los request en base a las estrcturas en la carpeta schema

## 🗃️ Migraciones automáticas
- Al iniciar el servidor, se ejecutan las migraciones pendientes en la carpeta `migrations/`.  
- Si una migración ya fue aplicada, `golang-migrate` la ignora.  
- Esto asegura que todas las instancias del backend trabajen con la misma estructura de base de datos.


## Configuración
- Todas las variables de entorno se gestionan en `.env`.
- El archivo `src/config/configuracion.go` carga estas variables y expone constantes globales.
- Actualmente el servidor utiliza `BACKEND_PORT` para definir en qué puerto escuchar.

## Test Automatizados
- Elaboracion de test unitarios e integrados que permiten ejecutarlos antes de salir a produccion con el objetivo de garantizar la calidad del codigo desarrollado
