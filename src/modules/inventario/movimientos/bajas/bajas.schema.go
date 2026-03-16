package bajas

// ─── Requests ─────────────────────────────────────────────────

type BajaDetalleRequest struct {
	IDArticulo          int     `json:"id_articulo"            validate:"gt=0"`
	Lote                string  `json:"lote"                   validate:"omitempty"`
	Cantidad            float64 `json:"cantidad"               validate:"gt=0"`
	IDSubtipoMovimiento int     `json:"id_subtipo_movimiento"  validate:"omitempty"` // si 0 hereda del encabezado
	Observacion         string  `json:"observacion"            validate:"omitempty"`
}

type BajaRequest struct {
	IDComprobante       int                  `json:"id_comprobante"         validate:"gt=0"`
	IDBodega            int                  `json:"id_bodega"              validate:"gt=0"`
	IDSubtipoMovimiento int                  `json:"id_subtipo_movimiento"  validate:"gt=0"` // subtipo general del documento
	ReferenciaExterna   string               `json:"referencia_externa"     validate:"omitempty"`
	Observaciones       string               `json:"observaciones"          validate:"omitempty"`
	FechaMovimiento     string               `json:"fecha_movimiento"       validate:"required"`
	Detalle             []BajaDetalleRequest `json:"detalle"                validate:"required,min=1,dive"`
	UserID              int64                `json:"user_id"                validate:"omitempty"`
	EmpresaID           int64                `json:"empresa_id"             validate:"omitempty"`
	SedeID              int64                `json:"sede_id"                validate:"omitempty"`
}

type AnulacionBajaRequest struct {
	Motivo    string `json:"motivo"     validate:"required,min=5"`
	UserID    int64  `json:"user_id"    validate:"omitempty"`
	EmpresaID int64  `json:"empresa_id" validate:"omitempty"`
}

// ─── Responses ────────────────────────────────────────────────

type BajaResponse struct {
	IDMovimiento     int    `json:"id_movimiento"`
	IDMovComprobante int    `json:"id_mov_comprobante"`
	Consecutivo      int    `json:"consecutivo"`
	Prefijo          string `json:"prefijo"`
}

type BajaListResponse struct {
	ID                int     `json:"id"`
	Consecutivo       int     `json:"consecutivo"`
	Prefijo           string  `json:"prefijo"`
	FechaMovimiento   string  `json:"fecha_movimiento"`
	Bodega            string  `json:"bodega"`
	ReferenciaExterna string  `json:"referencia_externa"`
	Estado            string  `json:"estado"`
	TotalCosto        float64 `json:"total_costo"`
}

type BajaDetalleResponse struct {
	IDArticulo    int     `json:"id_articulo"`
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	UnidadMedida  string  `json:"unidad_medida"`
	Lote          string  `json:"lote"`
	Cantidad      float64 `json:"cantidad"`
	TipoBaja      string  `json:"tipo_baja"`
	CostoUnitario float64 `json:"costo_unitario"`
	CostoTotal    float64 `json:"costo_total"`
	Observacion   string  `json:"observacion"`
}

type BajaFullResponse struct {
	ID                int                   `json:"id"`
	Consecutivo       int                   `json:"consecutivo"`
	Prefijo           string                `json:"prefijo"`
	Comprobante       string                `json:"comprobante"`
	FechaMovimiento   string                `json:"fecha_movimiento"`
	Bodega            string                `json:"bodega"`
	ReferenciaExterna string                `json:"referencia_externa"`
	Observaciones     string                `json:"observaciones"`
	Estado            string                `json:"estado"`
	TotalCosto        float64               `json:"total_costo"`
	Detalle           []BajaDetalleResponse `json:"detalle" gorm:"-"`
}

// ─── Refs ─────────────────────────────────────────────────────

type RefSubtipoMovimiento struct {
	ID               int    `json:"id"`
	IDTipoMovimiento int    `json:"id_tipo_movimiento"`
	Codigo           string `json:"codigo"`
	Nombre           string `json:"nombre"`
}

type RefComprobanteBaja struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Prefijo string `json:"prefijo"`
}

type RefLoteBaja struct {
	Lote        string  `json:"lote"`
	StockActual float64 `json:"stock_actual"`
}

// ─── Reporte ──────────────────────────────────────────────────

type BajaReporteDetalle struct {
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	UnidadMedida  string  `json:"unidad_medida"`
	Lote          string  `json:"lote"`
	TipoBaja      string  `json:"tipo_baja"`
	Cantidad      float64 `json:"cantidad"`
	CostoUnitario float64 `json:"costo_unitario"`
	CostoTotal    float64 `json:"costo_total"`
	Observacion   string  `json:"observacion"`
}

type EmpresaReporteBaja struct {
	RazonSocial  string `json:"razon_social"`
	Nit          string `json:"nit"`
	Dv           string `json:"dv"`
	Direccion    string `json:"direccion"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Telefono     string `json:"telefono"`
	LogoURL      string `json:"logo_url"`
}

type BajaReporte struct {
	Empresa           EmpresaReporteBaja   `json:"empresa"  gorm:"-"`
	ID                int                  `json:"id"`
	Consecutivo       int                  `json:"consecutivo"`
	Prefijo           string               `json:"prefijo"`
	Comprobante       string               `json:"comprobante"`
	FechaMovimiento   string               `json:"fecha_movimiento"`
	Estado            string               `json:"estado"`
	Bodega            string               `json:"bodega"`
	ReferenciaExterna string               `json:"referencia_externa"`
	Observaciones     string               `json:"observaciones"`
	Elaboro           string               `json:"elaboro"`
	Detalle           []BajaReporteDetalle `json:"detalle"  gorm:"-"`
	TotalUnidades     float64              `json:"total_unidades"`
	TotalCosto        float64              `json:"total_costo"`
}
