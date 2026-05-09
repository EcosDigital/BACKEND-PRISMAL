package ano_fiscal

import "time"

// ─── Constantes de dominio ────────────────────────────────────────────────────

const (
	EstadoAbierto = 1 // ref_estado_año_fiscal.id = 1 → 'Abierto'
	EstadoCerrado = 2 // ref_estado_año_fiscal.id = 2 → 'Cerrado'
)

// ─── Modelos de dominio ───────────────────────────────────────────────────────

// RefEstadoAñoFiscal representa una fila de la tabla catálogo.
type RefEstadoAñoFiscal struct {
	ID       int    `json:"id"        gorm:"column:id"`
	Nombre   string `json:"nombre"    gorm:"column:nombre"`
	IsActive bool   `json:"is_active" gorm:"column:is_active"`
}

// TableName le indica a GORM el nombre exacto de la tabla del catálogo.
// Necesario porque el nombre contiene caracteres especiales (ñ).
func (RefEstadoAñoFiscal) TableName() string {
	return "contabilidad.ref_estado_año_fiscal"
}

// AñoFiscal es el modelo principal.
type AñoFiscal struct {
	ID        int64     `json:"id"         gorm:"column:id;primaryKey"`
	Year      int       `json:"year"       gorm:"column:year"`
	IDEstado  int       `json:"id_estado"  gorm:"column:id_estado"`
	EmpresaID int64     `json:"empresa_id" gorm:"column:empresa_id"`
	CreatedBy int64     `json:"created_by" gorm:"column:created_by"`
	UpdatedBy int64     `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName le indica a GORM el nombre exacto de la tabla principal.
// Sin esto, GORM infiere "año_fiscals" y genera SQL inválido en Updates/Create.
func (AñoFiscal) TableName() string {
	return "contabilidad.cfg_años_fiscales"
}

// ─── Requests ─────────────────────────────────────────────────────────────────

type CreateAñoFiscalRequest struct {
	Year int `json:"year" validate:"required,min=1900,max=2200"`

	// Inyectados desde JWT — no vienen del body
	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
	SedeID    int64 `json:"-"`
}

type UpdateEstadoAñoFiscalRequest struct {
	IDEstado int `json:"id_estado" validate:"required,min=1"`

	// Inyectado desde JWT
	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
}

// ─── Responses ────────────────────────────────────────────────────────────────

// AñoFiscalResponse es el DTO público con el nombre del estado resuelto por JOIN.
type AñoFiscalResponse struct {
	ID        int64  `json:"id"         gorm:"column:id"`
	Year      int    `json:"year"       gorm:"column:year"`
	IDEstado  int    `json:"id_estado"  gorm:"column:id_estado"`
	Estado    string `json:"estado"     gorm:"column:estado"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
}

// ─── Envoltorio estándar de respuesta HTTP ────────────────────────────────────

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
