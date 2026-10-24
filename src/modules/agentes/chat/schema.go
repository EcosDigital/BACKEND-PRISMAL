package chat

type ChatRequest struct {
	IDAgente   int64  `json:"id_agente" validate:"required,gt=0"`
	Mensaje    string `json:"mensaje"   validate:"required,min=1,max=2000"`
	IDSesion   string `json:"id_sesion" validate:"required,max=64"`
	UserID     int64  `json:"user_id"   validate:"omitempty"`
	EmpresaID  int64  `json:"empresa_id" validate:"omitempty"`
	SedeID     int64  `json:"sede_id"   validate:"omitempty"`
	TenantSlug string `json:"-"`
}

type ChatResponse struct {
	Respuesta string `json:"respuesta"`
}
