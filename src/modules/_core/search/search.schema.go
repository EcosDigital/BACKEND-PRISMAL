package search

type SearchDynamics struct {
	NombreCodigo string  `json:"nombre_codigo" validate:"omitempty,max=200"`
	Filters      Filters `json:"filters" validate:"omitempty"`
	Operacion    int     `json:"operacion" validate:"required,min=1"`
	UserID       int64   `json:"user_id" validate:"omitempty"`
	EmpresaID    int64   `json:"empresa_id" validate:"omitempty"`
	SedeID       int64   `json:"sede_id" validate:"omitempty"`
}

type Filters struct {
	IdTipoRegistro   *int    `json:"id_tipo_registro"`
	FechaInicial     *string `json:"fecha_inicial"`
	FechaFinal       *string `json:"fecha_final"`
	IdEstadoRegistro *int    `json:"id_estado_registro"`
	Active           *bool   `json:"estado"`
	AplicarLimit     *bool   `json:"aplicar_limit"`
}
