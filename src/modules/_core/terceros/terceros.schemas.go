package terceros

import (
	"time"
)

type FechaJSON struct {
	time.Time
}

func (f *FechaJSON) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // quitar comillas
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s) // acepta solo YYYY-MM-DD
	if err != nil {
		return err
	}
	f.Time = t
	return nil
}

type TerceroRequest struct {
	IdTipoPersona      int    `json:"id_tipo_persona" validate:"required,gt=0"`
	IdTipoDocumento    int    `json:"id_tipo_documento" validate:"required,gt=0"`
	NumeroDocumento    string `json:"numero_documento" validate:"required,min=5,max=20"`
	Dv                 string `json:"dv" validate:"omitempty,len=1"`
	PrimerNombre       string `json:"primer_nombre" validate:"omitempty,alphaunicode,min=2,max=30"`
	SegundoNombre      string `json:"segundo_nombre" validate:"omitempty,alphaunicode,max=30"`
	PrimerApellido     string `json:"primer_apellido" validate:"omitempty,alphaunicode,min=2,max=30"`
	SegundoApellido    string `json:"segundo_apellido" validate:"omitempty,alphaunicode,min=2,max=30"`
	RazonSocial        string `json:"razon_social" validate:"omitempty,min=2,max=150"`
	RepresentanteLegal string `json:"representante_legal" validate:"omitempty,min=2,max=200"`
	IdGenero           int    `json:"id_genero" validate:"required,gt=0"`
	Telefono           string `json:"telefono" validate:"required,numeric,min=7,max=15"`
	Telefono_2         string `json:"telefono_2" validate:"omitempty,numeric,min=7,max=15"`
	Email              string `json:"email" validate:"required,email,max=100"`
	PaginaWeb          string `json:"pagina_web" validate:"omitempty,url,max=300"`
	Direccion          string `json:"direccion" validate:"required,min=5,max=150"`
	IdPais             int    `json:"id_pais" valdate:"required,gt=0"`
	IdDepartamento     int    `json:"id_departamento" validate:"required,gt=0"`
	IdCiudad           int    `json:"id_ciudad"`
	IdZona             int    `json:"id_zona" valdate:"required,gt=0"`
	Estado             *bool  `json:"estado" validate:"required"`
	ClaseTercero       []int  `json:"clase_tercero" validate:"omitempty,dive,gt=0"`
	UserID             int64  `json:"user_id" validate:"omitempty"`
	EmpresaID          int64  `json:"empresa_id" validate:"omitempty"`
	SedeID             int64  `json:"sede_id" validate:"omitempty"`
}

type TerceroResponse struct {
	ID              int    `json:"id"`
	NumeroDocumento string `json:"numero_documento"`
	NombreCompleto  string `json:"nombre_completo"`
	Email           string `json:"email"`
	Telefono        string `json:"telefono"`
	TipoDocumento   string `json:"tipo_documento"`
	Estado          bool   `json:"estado"`
}

type TerceroResponseFull struct {
	IdTipoPersona      int    `json:"id_tipo_persona"`
	IdTipoDocumento    int    `json:"id_tipo_documento"`
	NumeroDocumento    string `json:"numero_documento"`
	Dv                 string `json:"dv"`
	PrimerNombre       string `json:"primer_nombre"`
	SegundoNombre      string `json:"segundo_nombre"`
	PrimerApellido     string `json:"primer_apellido"`
	SegundoApellido    string `json:"segundo_apellido"`
	RazonSocial        string `json:"razon_social"`
	RepresentanteLegal string `json:"representante_legal"`
	IdGenero           int    `json:"id_genero"`
	Telefono           string `json:"telefono"`
	Telefono_2         string `json:"telefono_2"`
	Email              string `json:"email"`
	PaginaWeb          string `json:"pagina_web"`
	Direccion          string `json:"direccion"`
	IdPais             int    `json:"id_pais"`
	IdDepartamento     int    `json:"id_departamento"`
	IdCiudad           int    `json:"id_ciudad"`
	IdZona             int    `json:"id_zona"`
	Estado             *bool  `json:"estado"`
	ClaseTercero       []int  `json:"clase_tercero" gorm:"-"`
	UserID             int64  `json:"user_id"`
	EmpresaID          int64  `json:"empresa_id"`
	SedeID             int64  `json:"sede_id"`
}

type TerceroSearchBD struct {
	ID              int    `json:"id"`
	TipoPersona     string `json:"tipo_persona"`
	TipoDocumento   string `json:"tipo_documento"`
	NumeroDocumento string `json:"numero_documento"`
	NombreCompeto   string `json:"nombre_completo"`
	Direccion       string `json:"direccion"`
	Telefono        string `json:"telefono"`
	Genero          string `json:"genero"`
	Email           string `json:"email"`
	Departamento    string `json:"departamento"`
	Ciudad          string `json:"ciudad"`
}

type TerceroUpdateRequest struct {
	PrimerNombre       string `json:"primer_nombre" validate:"omitempty,alphaunicode,min=2,max=30"`
	SegundoNombre      string `json:"segundo_nombre" validate:"omitempty,alphaunicode,max=30"`
	PrimerApellido     string `json:"primer_apellido" validate:"omitempty,alphaunicode,min=2,max=30"`
	SegundoApellido    string `json:"segundo_apellido" validate:"omitempty,alphaunicode,min=2,max=30"`
	RazonSocial        string `json:"razon_social" validate:"omitempty,min=2,max=150"`
	RepresentanteLegal string `json:"representante_legal" validate:"omitempty,min=2,max=200"`
	IdGenero           int    `json:"id_genero" validate:"required,gt=0"`
	Telefono           string `json:"telefono" validate:"required,numeric,min=7,max=15"`
	Telefono_2         string `json:"telefono_2" validate:"omitempty,numeric,min=7,max=15"`
	Email              string `json:"email" validate:"required,email,max=100"`
	PaginaWeb          string `json:"pagina_web" validate:"omitempty,url,max=300"`
	Direccion          string `json:"direccion" validate:"required,min=5,max=150"`
	IdPais             int    `json:"id_pais" valdate:"required,gt=0"`
	IdDepartamento     int    `json:"id_departamento" validate:"required,gt=0"`
	Idciudad           int    `json:"id_ciudad" validate:"required,gt=0"`
	IdZona             int    `json:"id_zona" valdate:"required,gt=0"`
	FechaNacimiento    string `json:"fecha_nacimiento"`
	Estado             *bool  `json:"estado" validate:"required"`
	IdActividadEco     int    `json:"id_actividad_economica" validate:"omitempty,gt=0"`
	IdAmbito           int    `json:"id_ambito" validate:"omitempty,gt=0"`
	IdCentralizacion   int    `json:"id_centralizacion" validate:"omitempty,gt=0"`
	NumeroAcciones     int    `json:"numero_acciones" validate:"omitempty,max=7"`
	CodigoContable     int    `json:"codigo_contable" validate:"omitempty,max=5"`
	IdRespDian         int    `json:"id_responsabilidad_dian" validate:"omitempty,gt=0"`
	IdTipoCont         int    `json:"id_tipo_contribuyente" validate:"omitempty,gt=0"`
	IdRegimenIva       int    `json:"id_regimen_iva" validate:"omitempty,gt=0"`
	IdRegimenDian      int    `json:"id_regimen_dian" validate:"omitempty,gt=0"`
	ClaseTercero       []int  `json:"clase_tercero" validate:"omitempty,dive,gt=0"`
	UserID             int64  `json:"user_id" validate:"omitempty"`
	EmpresaID          int64  `json:"empresa_id" validate:"omitempty"`
	SedeID             int64  `json:"sede_id" validate:"omitempty"`
}

type TerceroUpdateIndenty struct {
	IdTipoPersona   int    `json:"id_tipo_persona" validate:"required,gt=0"`
	IdTipoDocumento int    `json:"id_tipo_documento" validate:"required,gt=0"`
	NumeroDocumento string `json:"numero_documento" validate:"required,min=5,max=20"`
	Dv              string `json:"dv" validate:"omitempty,len=1"`
}
