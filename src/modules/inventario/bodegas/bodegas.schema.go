package bodega

type BodegaRequest struct {
	IDTipoBodega      int     `json:"id_tipo_bodega" validate:"gt=0"`
	IDSede            int     `json:"id_sede" validate:"gt=0"`
	Codigo            string  `json:"codigo" validate:"required,min=1,max=10"`
	Nombre            string  `json:"nombre"               validate:"required,min=3"`
	Descripcion       string  `json:"descripcion"          validate:"omitempty"`
	IsCentral         *bool   `json:"is_central"           validate:"required"`
	PermiteVentas     *bool   `json:"permite_ventas"       validate:"required"`
	StockMinimoGlobal float64 `json:"stock_minimo_global"  validate:"omitempty"`
	StockMaximoGlobal float64 `json:"stock_maximo_global"  validate:"omitempty"`
	AlertaStock       *bool   `json:"alerta_stock"         validate:"required"`
	AplicaMovContable *bool   `json:"aplica_mov_contable"  validate:"required"`
	IsActive          *bool   `json:"is_active"            validate:"required"`
	UserID            int64   `json:"user_id"              validate:"omitempty"`
	EmpresaID         int64   `json:"empresa_id"           validate:"omitempty"`
	SedeID            int64   `json:"sede_id"              validate:"omitempty"`
}

type BodegaUpdateRequest struct {
	IDTipoBodega      int     `json:"id_tipo_bodega"       validate:"gt=0"`
	Nombre            string  `json:"nombre"               validate:"required,min=3"`
	Descripcion       string  `json:"descripcion"          validate:"omitempty"`
	IsCentral         *bool   `json:"is_central"           validate:"required"`
	PermiteVentas     *bool   `json:"permite_ventas"       validate:"required"`
	StockMinimoGlobal float64 `json:"stock_minimo_global"  validate:"omitempty"`
	StockMaximoGlobal float64 `json:"stock_maximo_global"  validate:"omitempty"`
	AlertaStock       *bool   `json:"alerta_stock"         validate:"required"`
	AplicaMovContable *bool   `json:"aplica_mov_contable"  validate:"required"`
	IsActive          *bool   `json:"is_active"            validate:"required"`
	UserID            int64   `json:"user_id"              validate:"omitempty"`
	EmpresaID         int64   `json:"empresa_id"           validate:"omitempty"`
	SedeID            int64   `json:"sede_id"              validate:"omitempty"`
}

type BodegaResponse struct {
	ID            int    `json:"id"`
	Codigo        string `json:"codigo"`
	Nombre        string `json:"nombre"`
	TipoBodega    string `json:"tipo_bodega"`
	IsCentral     *bool  `json:"is_central"`
	PermiteVentas *bool  `json:"permite_ventas"`
	IsActive      *bool  `json:"is_active"`
}

type BodegaResponseFull struct {
	ID                int     `json:"id"`
	IDTipoBodega      int     `json:"id_tipo_bodega"`
	TipoBodega        string  `json:"tipo_bodega"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
	Descripcion       string  `json:"descripcion"`
	IsCentral         *bool   `json:"is_central"          gorm:"column:is_central"`
	PermiteVentas     *bool   `json:"permite_ventas"      gorm:"column:permite_ventas"`
	StockMinimoGlobal float64 `json:"stock_minimo_global" gorm:"column:stock_minimo_global"`
	StockMaximoGlobal float64 `json:"stock_maximo_global" gorm:"column:stock_maximo_global"`
	AlertaStock       *bool   `json:"alerta_stock"        gorm:"column:alerta_stock"`
	AplicaMovContable *bool   `json:"aplica_mov_contable" gorm:"column:aplica_mov_contable"`
	IsActive          *bool   `json:"is_active"           gorm:"column:is_active"`
}
