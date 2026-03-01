package comprobantes

import "time"

type ComprobanteRequest struct {
	IDModulo          int        `json:"id_modulo" validate:"gt=0"`
	IDTipoOperacion   int        `json:"id_tipo_operacion" validate:"gt=0"`
	Nombre            string     `json:"nombre" validate:"required,min=3"`
	Prefijo           string     `json:"prefijo" validate:"required,min=1"`
	ConsecutivoInicio int        `json:"consecutivo_inicio" validate:"required,min=1"`
	ConsecutivoFinal  *int       `json:"consecutivo_final"`
	PermiteAnular     *bool      `json:"permite_anular" validate:"required"`
	FechaInicio       time.Time  `json:"fecha_inicio" validate:"required"`
	FechaFinal        *time.Time `json:"fecha_final" validate:"required"`
	Token             string     `json:"token" validate:"required,min=5"`
	Resolucion        string     `json:"resolucion" validate:"required,min=5"`
	IsActive          *bool      `json:"is_active" validate:"required"`
	UserID            int64      `json:"user_id" validate:"omitempty"`
	EmpresaID         int64      `json:"empresa_id" validate:"omitempty"`
	SedeID            int64      `json:"sede_id" validate:"omitempty"`
}

type ComprobanteResponse struct {
	ID                int        `json:"id"`
	IDModulo          int        `json:"id_modulo"`
	Modulo            string     `json:"modulo"`
	IDTipoOperacion   int        `json:"id_tipo_operacion"`
	Operacion         string     `json:"operacion"`
	Nombre            string     `json:"nombre"`
	Prefijo           string     `json:"prefijo"`
	ConsecutivoInicio int        `json:"consecutivo_inicio"`
	ConsecutivoFinal  *int       `json:"consecutivo_final"`
	PermiteAnular     *bool      `json:"permite_anular"`
	FechaInicio       time.Time  `json:"fecha_inicio"`
	FechaFinal        *time.Time `json:"fecha_final"`
	Token             string     `json:"token"`
	Resolucion        string     `json:"resolucion"`
	IsActive          *bool      `json:"is_active"`
}

type ComprobanteUpdateRequest struct {
	Nombre            string     `json:"nombre" validate:"required,min=3"`
	Prefijo           string     `json:"prefijo" validate:"required,min=1"`
	ConsecutivoInicio int        `json:"consecutivo_inicio" validate:"required,min=1"`
	ConsecutivoFinal  *int       `json:"consecutivo_final"`
	PermiteAnular     *bool      `json:"permite_anular" validate:"required"`
	FechaInicio       time.Time  `json:"fecha_inicio" validate:"required"`
	FechaFinal        *time.Time `json:"fecha_final" validate:"required"`
	Token             string     `json:"token" validate:"required,min=5"`
	Resolucion        string     `json:"resolucion" validate:"required,min=5"`
	IsActive          *bool      `json:"is_active" validate:"required"`
	UserID            int64      `json:"user_id" validate:"omitempty"`
	EmpresaID         int64      `json:"empresa_id" validate:"omitempty"`
	SedeID            int64      `json:"sede_id" validate:"omitempty"`
}
