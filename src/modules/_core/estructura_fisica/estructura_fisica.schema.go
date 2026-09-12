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

// CambiarEstadoMesaRequest body del PUT /core/mesas/:id/estado — cambia
// solo el estado (por id, no por nombre: el frontend ya conoce los ids de
// configuracion.ref_estado_mesa) sin reenviar codigo/nombre/is_active.
type CambiarEstadoMesaRequest struct {
	IDEstado int `json:"id_estado" validate:"required,gt=0"`
}

type MesaResponse struct {
	ID       int    `json:"id"`
	Codigo   string `json:"codigo"`
	Nombre   string `json:"nombre"`
	IDEstado int    `json:"id_estado"`
	Estado   string `json:"estado"`
	IsActive bool   `json:"is_active"`
}
