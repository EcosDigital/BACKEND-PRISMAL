package costos

type CostosRequest struct {
	Codigo    string `json:"codigo" validate:"required,min=2,max=12"`
	Nombre    string `json:"nombre" validate:"required,min=3"`
	IDArea    *int   `json:"id_area" validate:"omitempty,gt=0"`
	IDUnidad  *int   `json:"id_unidad" validate:"omitempty,gt=0"`
	IsActive  *bool  `json:"is_active" validate:"required"`
	UserID    int64  `json:"user_id" validate:"omitempty"`
	EmpresaID int64  `json:"empresa_id" validate:"omitempty"`
	SedeID    int64  `json:"sede_id" validate:"omitempty"`
}

type CostosResponse struct {
	ID       int    `json:"id"`
	Codigo   string `json:"codigo"`
	Nombre   string `json:"nombre"`
	IDArea   *int   `json:"id_area" validate:"omitempty,gt=0"`
	Area     string `json:"area"`
	IDUnidad *int   `json:"id_unidad" validate:"omitempty,gt=0"`
	Unidad   string `json:"unidad"`
	IsActive *bool  `json:"is_active" validate:"required"`
}

type CostosUpdateRequest struct {
	Nombre    string `json:"nombre" validate:"required,min=3"`
	IDArea    *int   `json:"id_area" validate:"omitempty,gt=0"`
	IDUnidad  *int   `json:"id_unidad" validate:"omitempty,gt=0"`
	IsActive  *bool  `json:"is_active" validate:"required"`
	UserID    int64  `json:"user_id" validate:"omitempty"`
	EmpresaID int64  `json:"empresa_id" validate:"omitempty"`
	SedeID    int64  `json:"sede_id" validate:"omitempty"`
}
