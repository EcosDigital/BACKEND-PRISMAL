package articulos

type ArticuloRequest struct {
	Codigo              string  `json:"codigo"               validate:"required,min=1,max=50"`
	Nombre              string  `json:"nombre"               validate:"required,min=3"`
	Descripcion         string  `json:"descripcion"          validate:"omitempty"`
	IDTipoArticulo      int     `json:"id_tipo_articulo"     validate:"gt=0"`
	IDUnidadMedida      int     `json:"id_unidad_medida"     validate:"gt=0"`
	FactorUnidadBase    float64 `json:"factor_unidad_base"   validate:"omitempty"`
	IDFormaFarmaceutica *int    `json:"id_forma_farmaceutica" validate:"omitempty"`
	IDPresentacion      *int    `json:"id_presentacion"      validate:"omitempty"`
	IsActive            *bool   `json:"is_active"            validate:"required"`
	UserID              int64   `json:"user_id"              validate:"omitempty"`
	EmpresaID           int64   `json:"empresa_id"           validate:"omitempty"`
	SedeID              int64   `json:"sede_id"              validate:"omitempty"`
}

type ArticuloUpdateRequest struct {
	Nombre              string  `json:"nombre"               validate:"required,min=3"`
	Descripcion         string  `json:"descripcion"          validate:"omitempty"`
	IDTipoArticulo      int     `json:"id_tipo_articulo"     validate:"gt=0"`
	IDUnidadMedida      int     `json:"id_unidad_medida"     validate:"gt=0"`
	FactorUnidadBase    float64 `json:"factor_unidad_base"   validate:"omitempty"`
	IDFormaFarmaceutica *int    `json:"id_forma_farmaceutica" validate:"omitempty"`
	IDPresentacion      *int    `json:"id_presentacion"      validate:"omitempty"`
	IsActive            *bool   `json:"is_active"            validate:"required"`
	UserID              int64   `json:"user_id"              validate:"omitempty"`
	EmpresaID           int64   `json:"empresa_id"           validate:"omitempty"`
	SedeID              int64   `json:"sede_id"              validate:"omitempty"`
}

type ArticuloResponse struct {
	ID            int     `json:"id"`
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	GrupoArticulo string  `json:"grupo_articulo"`
	UnidadMedida  string  `json:"unidad_medida"`
	FactorUnidad  float64 `json:"factor_unidad_base" gorm:"column:factor_unidad_base"`
	IsActive      *bool   `json:"is_active"`
}

type ArticuloResponseFull struct {
	ID                  int     `json:"id"`
	Codigo              string  `json:"codigo"`
	Nombre              string  `json:"nombre"`
	Descripcion         string  `json:"descripcion"`
	IDTipoArticulo      int     `json:"id_tipo_articulo"`
	GrupoArticulo       string  `json:"grupo_articulo"`
	IDUnidadMedida      int     `json:"id_unidad_medida"`
	UnidadMedida        string  `json:"unidad_medida"`
	FactorUnidadBase    float64 `json:"factor_unidad_base"    gorm:"column:factor_unidad_base"`
	IDFormaFarmaceutica *int    `json:"id_forma_farmaceutica" gorm:"column:id_forma_farmaceutica"`
	IDPresentacion      *int    `json:"id_presentacion"       gorm:"column:id_presentacion"`
	IsActive            *bool   `json:"is_active"             gorm:"column:is_active"`
}

type RefUnidadMedida struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type RefGrupoArticulo struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

// ─── Carga masiva ─────────────────────────────────────────────────────────────

// PlanoArticuloRow representa una fila del Excel recibida en el request.
// Los campos de referencia vienen como NOMBRE (string) igual que en el plano.
// El backend resuelve el ID haciendo lookup en las tablas de referencia.
type PlanoArticuloRow struct {
	// Columnas del plano (tal cual el Excel)
	Codigo               string  `json:"CODIGO"`
	Nombre               string  `json:"NOMBRE"`
	Descripcion          string  `json:"DESCRIPCION"`
	IDTipoArticuloNombre string  `json:"ID_TIPO_ARTICULO"` // nombre del grupo  → se resuelve a ID
	IDUnidadMediaNombre  string  `json:"ID_UNIDAD_MEDIA"`  // nombre de unidad  → se resuelve a ID
	FactorBase           float64 `json:"FACTOR_BASE"`
	IDFormaNombre        string  `json:"ID_FORMA_FARMACEUTICA"` // nombre forma farmacéutica (opcional)
	IDPresentacionNombre string  `json:"ID_PRESETACION"`        // nombre presentación → se resuelve a ID
	Estado               string  `json:"ESTADO"`                // "Activo" | "Inactivo"
}

// ImportArticulosRequest es el body del endpoint POST /articulos/import
type ImportArticulosRequest struct {
	Filas []PlanoArticuloRow `json:"filas" validate:"required,min=1"`
}

// ImportRowResult es el resultado por cada fila procesada
type ImportRowResult struct {
	Fila   int    `json:"fila"`
	Codigo string `json:"codigo"`
	Accion string `json:"accion"` // "creado" | "actualizado" | "error"
	Error  string `json:"error,omitempty"`
}

// ImportArticulosResponse es la respuesta final del endpoint
type ImportArticulosResponse struct {
	Creados      int               `json:"creados"`
	Actualizados int               `json:"actualizados"`
	Errores      int               `json:"errores"`
	Detalle      []ImportRowResult `json:"detalle"`
}

// refCache agrupa los mapas de lookup nombre→id para no consultar la BD por cada fila
type refCache struct {
	Grupos         map[string]int // nombre_lower → id
	Unidades       map[string]int
	Presentaciones map[string]int
	FormasFarma    map[string]int
}
