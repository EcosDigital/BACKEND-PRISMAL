package catalogo

// ─── Cabecera del catálogo ──────────────────────────────────────────────────

type CatalogoRequest struct {
	Nombre      string `json:"nombre"      validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"omitempty"`
	IDSede      int64  `json:"id_sede"     validate:"gt=0"`
	UserID      int64  `json:"user_id"     validate:"omitempty"`
	EmpresaID   int64  `json:"empresa_id"  validate:"omitempty"`
}

type CatalogoUpdateRequest struct {
	Nombre      string `json:"nombre"      validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"omitempty"`
	IsActive    *bool  `json:"is_active"   validate:"required"`
	UserID      int64  `json:"user_id"     validate:"omitempty"`
}

type CatalogoResponse struct {
	ID       int    `json:"id"`
	Codigo   string `json:"codigo"`
	Nombre   string `json:"nombre"`
	Sede     string `json:"sede"`
	IsActive *bool  `json:"is_active"`
}

type CatalogoResponseFull struct {
	ID          int    `json:"id"`
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	IDSede      int    `json:"id_sede"`
	Sede        string `json:"sede"`
	IsActive    *bool  `json:"is_active" gorm:"column:is_active"`
}

// ─── Detalle: artículos publicados dentro del catálogo ──────────────────────

type CatalogoArticuloRequest struct {
	IDArticulo      int64   `json:"id_articulo"      validate:"gt=0"`
	PrecioPublico   float64 `json:"precio_publico"   validate:"gte=0"`
	AplicaDomicilio *bool   `json:"aplica_domicilio" validate:"omitempty"`
	UserID          int64   `json:"user_id"          validate:"omitempty"`
}

type CatalogoArticuloUpdateRequest struct {
	PrecioPublico   float64 `json:"precio_publico"   validate:"gte=0"`
	AplicaDomicilio *bool   `json:"aplica_domicilio" validate:"required"`
	UserID          int64   `json:"user_id"          validate:"omitempty"`
}

type CatalogoArticuloResponse struct {
	ID              int     `json:"id"`
	IDArticulo      int     `json:"id_articulo"`
	CodigoArticulo  string  `json:"codigo_articulo"`
	NombreArticulo  string  `json:"nombre_articulo"`
	ImagenURL       string  `json:"imagen_url"`
	PrecioPublico   float64 `json:"precio_publico"`
	AplicaDomicilio *bool   `json:"aplica_domicilio"`
	SyncOk          bool    `json:"sync_ok" gorm:"column:sync_ok"`
}
