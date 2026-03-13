package comprobantes

type ComprobanteRequest struct {
	IDModulo          int    `json:"id_modulo" validate:"gt=0"`
	IDTipoOperacion   int    `json:"id_tipo_operacion" validate:"gt=0"`
	Nombre            string `json:"nombre" validate:"required,min=3"`
	Prefijo           string `json:"prefijo" validate:"required,min=1"`
	ConsecutivoInicio int    `json:"consecutivo_inicio" validate:"required,min=1"`
	ConsecutivoFinal  *int   `json:"consecutivo_final"`
	PermiteAnular     *bool  `json:"permite_anular" validate:"required"`
	FechaInicio       string `json:"fecha_inicio" validate:"required"`
	FechaFinal        string `json:"fecha_final" validate:"required"`
	Token             string `json:"token" validate:"omitempty,min=5"`
	Resolucion        string `json:"resolucion" validate:"omitempty,min=5"`
	IsActive          *bool  `json:"is_active" validate:"required"`
	UserID            int64  `json:"user_id" validate:"omitempty"`
	EmpresaID         int64  `json:"empresa_id" validate:"omitempty"`
	SedeID            int64  `json:"sede_id" validate:"omitempty"`
}

type ComprobanteResponse struct {
	ID          int    `json:"id"`
	IDModulo    int    `json:"id_modulo"`
	Modulo      string `json:"modulo"`
	Operacion   string `json:"operacion"`
	Nombre      string `json:"nombre"`
	Consecutivo int    `json:"consecutivo"`
	IsActive    *bool  `json:"is_active"`
}

type ComprobanteResponseFull struct {
	ID                int    `json:"id"`
	IDModulo          int    `json:"id_modulo"`
	IDTipoOperacion   int    `json:"id_tipo_operacion"`
	Prefijo           string `json:"prefijo"`
	Nombre            string `json:"nombre"`
	ConsecutivoInicio int    `json:"consecutivo_inicio"`
	Consecutivo       int    `json:"consecutivo_actual"`
	ConsecutivoFinal  *int   `json:"consecutivo_final"`
	FechaInicio       string `json:"fecha_inicio"`
	FechaFinal        string `json:"fecha_final"`
	PermiteAnulacion  *bool  `json:"permite_anulacion" gorm:"column:permite_anulacion"`
	Token             string `json:"token"`
	Resolucion        string `json:"resolucion"`
	IsActive          *bool  `json:"is_active"`
}

type ComprobanteUpdateRequest struct {
	Nombre            string `json:"nombre" validate:"required,min=3"`
	Prefijo           string `json:"prefijo" validate:"required,min=1"`
	ConsecutivoInicio int    `json:"consecutivo_inicio" validate:"required,min=1"`
	ConsecutivoFinal  *int   `json:"consecutivo_final"`
	PermiteAnular     *bool  `json:"permite_anular" validate:"required"`
	FechaInicio       string `json:"fecha_inicio" validate:"required"`
	FechaFinal        string `json:"fecha_final" validate:"required"`
	Token             string `json:"token" validate:"omitempty,min=5"`
	Resolucion        string `json:"resolucion" validate:"omitempty,min=5"`
	IsActive          *bool  `json:"is_active" validate:"required"`
	UserID            int64  `json:"user_id" validate:"omitempty"`
	EmpresaID         int64  `json:"empresa_id" validate:"omitempty"`
	SedeID            int64  `json:"sede_id" validate:"omitempty"`
}
