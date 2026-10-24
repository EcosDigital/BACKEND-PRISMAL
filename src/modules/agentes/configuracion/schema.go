package configuracion

type AgenteRequest struct {
	Nombre        string `json:"nombre"        validate:"required,min=2,max=100"`
	Instrucciones string `json:"instrucciones" validate:"required,min=10"`
	Modelo        string `json:"modelo"        validate:"omitempty,max=100"`
	IsActive      *bool  `json:"is_active"     validate:"required"`
	IDModuloRef   *int64 `json:"id_modulo_ref" validate:"omitempty,gt=0"`
	UserID        int64  `json:"user_id"       validate:"omitempty"`
	EmpresaID     int64  `json:"empresa_id"    validate:"omitempty"`
	SedeID        int64  `json:"sede_id"       validate:"omitempty"`
}

type AgenteResponse struct {
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	Instrucciones string `json:"instrucciones"`
	Modelo        string `json:"modelo"`
	IsActive      bool   `json:"is_active"`
	IDModuloRef   *int   `json:"id_modulo_ref"`
	Modulo        string `json:"modulo"`
}

type ModuloResponse struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}
