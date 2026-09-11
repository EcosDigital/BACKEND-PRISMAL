package estructura_fisica

type MesaRequest struct {
	Codigo    string `json:"codigo" validate:"required"`
	Nombre    string `json:"nombre" validate:"required,min=1"`
	IDEstado  int    `json:"id_estado" validate:"required,gt=0"`
	IsActive  *bool  `json:"is_active" validate:"required"`
	UserID    int64  `json:"user_id" validate:"omitempty"`
	EmpresaID int64  `json:"empresa_id" validate:"omitempty"`
	SedeID    int64  `json:"sede_id" validate:"omitempty"`
}

type MesaResponse struct {
	ID       int    `json:"id"`
	Codigo   string `json:"codigo"`
	Nombre   string `json:"nombre"`
	IDEstado int    `json:"id_estado"`
	Estado   string `json:"estado"`
	IsActive bool   `json:"is_active"`
}
