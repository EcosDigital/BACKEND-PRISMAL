package documentacion_sop

import "time"

// ─── Constantes de dominio ────────────────────────────────────────────────────

const (
	EstadoDraft    = 1 //borrador
	EstadoReview   = 2 // en revision
	EstadoApproved = 3 // aprobado
	EstadoObsolete = 4 // obsoleto

	// Categorías de SOP (ref_categoria_sop.id)
	CategoriaOperacion      = 1
	CategoriaeSoporte       = 2
	CategoriaTecnologia     = 3
	CategoriaPagos          = 4
	CategoriaComercial      = 5
	CategoriaMarketing      = 6
	CategoriaAdministracion = 7
)

// ─── Modelos de dominio ───────────────────────────────────────────────────────

// RefCategoriaSOP representa una fila de la tabla catálogo de categorías.
type RefCategoriaSOP struct {
	ID       int    `json:"id"        gorm:"column:id"`
	Nombre   string `json:"nombre"    gorm:"column:nombre"`
	IsActive bool   `json:"is_active" gorm:"column:is_active"`
}

// TableName le indica a GORM el nombre exacto de la tabla catálogo.
func (RefCategoriaSOP) TableName() string {
	return "documentacion.ref_categoria_sop"
}

// RefEstadoSOP representa una fila de la tabla catálogo de estados.
type RefEstadoSOP struct {
	ID       int    `json:"id"        gorm:"column:id"`
	Nombre   string `json:"nombre"    gorm:"column:nombre"`
	Color    string `json:"color"     gorm:"column:color"`
	IsActive bool   `json:"is_active" gorm:"column:is_active"`
}

// TableName le indica a GORM el nombre exacto de la tabla catálogo.
func (RefEstadoSOP) TableName() string {
	return "documentacion.ref_estado_sop"
}

// SOP es el modelo principal del módulo.
type SOP struct {
	ID            int64      `json:"id"             gorm:"column:id;primaryKey"`
	Code          string     `json:"code"           gorm:"column:code"`
	Title         string     `json:"title"          gorm:"column:title"`
	IDCategoria   int        `json:"id_categoria"   gorm:"column:id_categoria"`
	Area          string     `json:"area"           gorm:"column:area"`
	IDEstado      int        `json:"id_estado"      gorm:"column:id_estado"`
	Content       string     `json:"content"        gorm:"column:content"`
	Version       string     `json:"version"        gorm:"column:version"`
	EffectiveDate time.Time  `json:"effective_date" gorm:"column:effective_date"`
	ReviewDate    *time.Time `json:"review_date"   gorm:"column:review_date"`
	IsActive      bool       `json:"is_active"      gorm:"column:is_active"`
	EmpresaID     int64      `json:"empresa_id"     gorm:"column:empresa_id"`
	CreatedBy     int64      `json:"created_by"     gorm:"column:created_by"`
	UpdatedBy     int64      `json:"updated_by"     gorm:"column:updated_by"`
	CreatedAt     time.Time  `json:"created_at"     gorm:"column:created_at"`
	UpdatedAt     time.Time  `json:"updated_at"     gorm:"column:updated_at"`
}

// TableName le indica a GORM el nombre exacto de la tabla principal.
// Sin esto, GORM inferiría un nombre incorrecto para el esquema documentacion.
func (SOP) TableName() string {
	return "documentacion.cfg_sop"
}

// ─── Requests ─────────────────────────────────────────────────────────────────

// CreateSOPRequest contiene los campos necesarios para registrar un nuevo SOP.
type CreateSOPRequest struct {
	Code          string  `json:"code"           validate:"required,max=50"`
	Title         string  `json:"title"          validate:"required,max=500"`
	IDCategoria   int     `json:"id_categoria"   validate:"required,min=1"`
	Area          string  `json:"area"           validate:"required,max=200"`
	Content       string  `json:"content"        validate:"required"`
	Version       string  `json:"version"        validate:"required,max=20"`
	EffectiveDate string  `json:"effective_date" validate:"required"`
	ReviewDate    *string `json:"review_date"`

	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
}

// UpdateSOPRequest contiene los campos actualizables de un SOP existente.
type UpdateSOPRequest struct {
	Title         string  `json:"title"          validate:"required,max=500"`
	IDCategoria   int     `json:"id_categoria"   validate:"required,min=1"`
	Area          string  `json:"area"           validate:"required,max=200"`
	IDEstado      int     `json:"id_estado"      validate:"required,min=1"`
	Content       string  `json:"content"        validate:"required"`
	Version       string  `json:"version"        validate:"required,max=20"`
	EffectiveDate string  `json:"effective_date" validate:"required"` // formato: YYYY-MM-DD
	ReviewDate    *string `json:"review_date"`                        // opcional: YYYY-MM-DD

	// Inyectados desde JWT — no vienen del body
	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
}

// SOPFilters agrupa los filtros disponibles para el listado de SOPs.
type SOPFilters struct {
	IDCategoria int    // 0 = sin filtro
	IDEstado    int    // 0 = sin filtro
	Area        string // "" = sin filtro
	Search      string // búsqueda por título o código, "" = sin filtro
}

// ─── Responses ────────────────────────────────────────────────────────────────

// SOPResponse es el DTO público con los catálogos resueltos por JOIN.
type SOPResponse struct {
	ID            int64   `json:"id"`
	Code          string  `json:"code"`
	Title         string  `json:"title"`
	IDCategoria   int     `json:"id_categoria"`
	Categoria     string  `json:"categoria"`
	Area          string  `json:"area"`
	IDEstado      int     `json:"id_estado"`
	Estado        string  `json:"estado"`
	EstadoColor   string  `json:"estado_color"`
	Content       string  `json:"content"`
	Version       string  `json:"version"`
	EffectiveDate string  `json:"effective_date"`
	ReviewDate    *string `json:"review_date"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// CatalogoResponse agrupa los catálogos de referencia para poblar selects en el frontend.
type CatalogoResponse struct {
	Categorias []RefCategoriaSOP `json:"categorias"`
	Estados    []RefEstadoSOP    `json:"estados"`
}

// ─── Envoltorio estándar de respuesta HTTP ────────────────────────────────────

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
