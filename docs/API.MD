# Documentación de la API

## 🔹 Endpoints disponibles

> Nota: Todas las peticiones están sujetas al middleware CORS.  
> Solo se aceptan solicitudes desde el dominio definido en `FRONTEND_URL` (.env).  

### 0.0.1 Rutas Generales
**GET** `/api/general/TipoPersona`

#### Descripción
Responde con un array de json enlistando cada registro desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 1,
    "nombre": "Natural"
  },
  {
    "id": 2,
    "nombre": "Juridica"
  }
]
```

**GET** `/api/general/TipoDocument`
#### Descripción
Responde con un array de json enlistando cada registro desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 1,
    "nombre": "Tarjeta de identidad",
    "codigo": "TI",
    "tipo": 1
  },
  {
    "id": 2,
    "nombre": "Cédula de ciudadania",
    "codigo": "CC",
    "tipo": 1
  },
  {
    "id": 3,
    "nombre": "Pasaporte",
    "codigo": "PP",
    "tipo": 1
  },
  {
    "id": 4,
    "nombre": "Cédula de extranjería",
    "codigo": "CE",
    "tipo": 1
  },
  {
    "id": 5,
    "nombre": "Número de Identificación Tributaria",
    "codigo": "NIT",
    "tipo": 2
  }
]
```

**GET** `/api/general/Generos`
#### Descripción
Responde con un array de json enlistando cada registro desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 1,
    "nombre": "Masculino"
  },
  {
    "id": 2,
    "nombre": "Femenino"
  },
  {
    "id": 3,
    "nombre": "No Aplica"
  }
]
```

**GET** `/api/general/Departamentos`
#### Descripción
Responde con un array de json enlistando cada registro de los departamentos nacionales
desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 2,
    "codigo": 8,
    "nombre": "ATLANTICO"
  },
  {
    "id": 3,
    "codigo": 11,
    "nombre": "BOGOTA"
  },
  {
    "id": 4,
    "codigo": 13,
    "nombre": "BOLIVAR"
  },
]
```

**GET** `/api/general/Municipios/:id`
#### Descripción
Responde con un array de json enlistando cada registro de los municipios nacionales por departamento
desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 126,
    "codigo_departamento": 8,
    "codigo_municipio": 8001,
    "nombre_municipio": "BARRANQUILLA"
  },
  {
    "id": 127,
    "codigo_departamento": 8,
    "codigo_municipio": 8078,
    "nombre_municipio": "BARANOA"
  },
  {
    "id": 128,
    "codigo_departamento": 8,
    "codigo_municipio": 8137,
    "nombre_municipio": "CAMPO DE LA CRUZ"
  },
]
```

**GET** `/api/general/ZonasRurales`
#### Descripción
Responde con un array de json enlistando cada registro desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 1,
    "codigo": "U",
    "nombre": "Urbana"
  },
  {
    "id": 2,
    "codigo": "R",
    "nombre": "Rural"
  }
]
```

**GET** `/api/general/Ambitos`
#### Descripción
Responde con un array de json enlistando cada registro desde la tabla de referencia.

#### Respuesta exitosa (200)
```json
[
  {
    "id": 1,
    "nombre": "Municipal"
  },
  {
    "id": 2,
    "nombre": "Departamental"
  },
  {
    "id": 3,
    "nombre": "Nacional"
  },
  {
    "id": 4,
    "nombre": "Internacional"
  }
]
```

