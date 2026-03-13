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
	FactorUnidad  float64 `json:"factor_unidad_base"`
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
