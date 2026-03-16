package reportes_inv

// ─── Response para el reporte de entrada ─────────────────────

type EmpresaReporte struct {
	RazonSocial  string `json:"razon_social"`
	Nit          string `json:"nit"`
	Dv           string `json:"dv"`
	Direccion    string `json:"direccion"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Telefono     string `json:"telefono"`
	Email        string `json:"email"`
	LogoURL      string `json:"logo_url"`
}

type EntradaReporteDetalle struct {
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	Presentacion  string  `json:"presentacion"`
	UnidadMedida  string  `json:"unidad_medida"`
	Cantidad      float64 `json:"cantidad"`
	ValorUnitario float64 `json:"valor_unitario"`
	PorcentajeIva float64 `json:"porcentaje_iva"`
	ValorIva      float64 `json:"valor_iva"`
	Subtotal      float64 `json:"subtotal"`
}

type EntradaReporte struct {
	// Encabezado empresa — ignorado por GORM, se llena manualmente
	Empresa EmpresaReporte `json:"empresa" gorm:"-"`

	// Encabezado del documento
	ID                int    `json:"id"`
	Consecutivo       int    `json:"consecutivo"`
	Prefijo           string `json:"prefijo"`
	Comprobante       string `json:"comprobante"`
	FechaMovimiento   string `json:"fecha_movimiento"`
	Estado            string `json:"estado"`
	ReferenciaExterna string `json:"referencia_externa"`
	Observaciones     string `json:"observaciones"`

	// Proveedor
	Proveedor          string `json:"proveedor"`
	NitProveedor       string `json:"nit_proveedor"`
	DireccionProveedor string `json:"direccion_proveedor"`
	TelefonoProveedor  string `json:"telefono_proveedor"`

	// Bodega
	Bodega string `json:"bodega"`

	// Usuario que elaboró
	Elaboro string `json:"elaboro"`

	// Detalle artículos — ignorado por GORM, se llena manualmente
	Detalle []EntradaReporteDetalle `json:"detalle" gorm:"-"`

	// Totales
	Subtotal       float64 `json:"subtotal"`
	TotalIva       float64 `json:"total_iva"`
	TotalDescuento float64 `json:"total_descuento"`
	TotalNeto      float64 `json:"total_neto"`
}

// ─── Response para el reporte de traslado ─────────────────────

type TrasladoReporteDetalle struct {
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	UnidadMedida  string  `json:"unidad_medida"`
	Lote          string  `json:"lote"`
	Cantidad      float64 `json:"cantidad"`
	CostoUnitario float64 `json:"costo_unitario"` // costo promedio origen
	CostoTotal    float64 `json:"costo_total"`
}

type TrasladoReporte struct {
	// Empresa — ignorado por GORM
	Empresa EmpresaReporteTraslado `json:"empresa" gorm:"-"`

	// Encabezado
	ID              int    `json:"id"`
	Consecutivo     int    `json:"consecutivo"`
	Prefijo         string `json:"prefijo"`
	Comprobante     string `json:"comprobante"`
	FechaMovimiento string `json:"fecha_movimiento"`
	Estado          string `json:"estado"`
	Observaciones   string `json:"observaciones"`

	// Bodegas
	BodegaOrigen  string `json:"bodega_origen"`
	BodegaDestino string `json:"bodega_destino"`

	// Usuario
	Elaboro string `json:"elaboro"`

	// Detalle — ignorado por GORM
	Detalle []TrasladoReporteDetalle `json:"detalle" gorm:"-"`

	// Totales
	TotalUnidades float64 `json:"total_unidades"`
	TotalCosto    float64 `json:"total_costo"`
}

type EmpresaReporteTraslado struct {
	RazonSocial  string `json:"razon_social"`
	Nit          string `json:"nit"`
	Dv           string `json:"dv"`
	Direccion    string `json:"direccion"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Telefono     string `json:"telefono"`
	LogoURL      string `json:"logo_url"`
}
