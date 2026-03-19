package kardex

// ─── Request ──────────────────────────────────────────────────

type KardexRequest struct {
	IDArticulo  int64  `query:"id_articulo"  validate:"gt=0"`
	IDBodega    int64  `query:"id_bodega"    validate:"gt=0"`
	Lote        string `query:"lote"         validate:"omitempty"`
	FechaInicio string `query:"fecha_inicio" validate:"required"`
	FechaFin    string `query:"fecha_fin"    validate:"required"`
}

// ─── Fila del Kardex ──────────────────────────────────────────

type KardexFila struct {
	FilaOrden      int    `json:"fila_orden"`
	Fecha          string `json:"fecha"`
	TipoMovimiento string `json:"tipo_movimiento"`
	// Clase: derivada en Go por derivarClase() — no viene de la función SQL.
	// Valores: entrada | baja | traslado_entrada | traslado_salida |
	//          despacho | saldo_inicial
	Clase             string  `json:"clase"`
	Documento         string  `json:"documento"`
	ReferenciaExterna string  `json:"referencia_externa"`
	Observaciones     string  `json:"observaciones"`
	CantidadEntrada   float64 `json:"cantidad_entrada"`
	CantidadSalida    float64 `json:"cantidad_salida"`
	CostoUnitario     float64 `json:"costo_unitario"`
	SaldoCantidad     float64 `json:"saldo_cantidad"`
	EsAnulado         bool    `json:"es_anulado"`
	MotivoAnulacion   string  `json:"motivo_anulacion"`
	IDMovimiento      int     `json:"id_movimiento"`
}

// ─── Cabecera del reporte ─────────────────────────────────────

type KardexCabecera struct {
	IDArticulo   int64   `json:"id_articulo"`
	CodigoArt    string  `json:"codigo_articulo"`
	NombreArt    string  `json:"nombre_articulo"`
	UnidadMedida string  `json:"unidad_medida"`
	IDBodega     int64   `json:"id_bodega"`
	NombreBodega string  `json:"nombre_bodega"`
	Lote         string  `json:"lote"`
	FechaInicio  string  `json:"fecha_inicio"`
	FechaFin     string  `json:"fecha_fin"`
	StockFinal   float64 `json:"stock_final"`
}

// ─── Empresa ──────────────────────────────────────────────────

type KardexEmpresa struct {
	RazonSocial  string `json:"razon_social"`
	Nit          string `json:"nit"`
	Dv           string `json:"dv"`
	Direccion    string `json:"direccion"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Telefono     string `json:"telefono"`
	LogoURL      string `json:"logo_url"`
}

// ─── Response completo ────────────────────────────────────────

type KardexReporte struct {
	Empresa  KardexEmpresa  `json:"empresa"`
	Cabecera KardexCabecera `json:"cabecera"`
	Filas    []KardexFila   `json:"filas"`
}

// ─── Refs para el formulario ──────────────────────────────────

type KardexRefBodega struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

type KardexRefLote struct {
	Lote string `json:"lote"`
}
