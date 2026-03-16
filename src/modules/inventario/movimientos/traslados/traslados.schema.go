package traslados

// ─── Requests ─────────────────────────────────────────────────

type TrasladoDetalleRequest struct {
	IDArticulo  int     `json:"id_articulo"  validate:"gt=0"`
	Lote        string  `json:"lote"         validate:"omitempty"`
	Cantidad    float64 `json:"cantidad"     validate:"gt=0"`
	Observacion string  `json:"observacion"  validate:"omitempty"`
}

type TrasladoRequest struct {
	IDComprobante     int                      `json:"id_comprobante"      validate:"gt=0"`
	IDBodegaOrigen    int                      `json:"id_bodega_origen"    validate:"gt=0"`
	IDBodegaDestino   int                      `json:"id_bodega_destino"   validate:"gt=0"`
	ReferenciaExterna string                   `json:"referencia_externa"  validate:"omitempty"`
	Observaciones     string                   `json:"observaciones"       validate:"omitempty"`
	FechaMovimiento   string                   `json:"fecha_movimiento"    validate:"required"`
	Detalle           []TrasladoDetalleRequest `json:"detalle"             validate:"required,min=1,dive"`
	UserID            int64                    `json:"user_id"             validate:"omitempty"`
	EmpresaID         int64                    `json:"empresa_id"          validate:"omitempty"`
	SedeID            int64                    `json:"sede_id"             validate:"omitempty"`
}

// ─── Responses ────────────────────────────────────────────────

type TrasladoResponse struct {
	IDMovimiento     int    `json:"id_movimiento"`
	IDMovComprobante int    `json:"id_mov_comprobante"`
	Consecutivo      int    `json:"consecutivo"`
	Prefijo          string `json:"prefijo"`
}

type TrasladoListResponse struct {
	ID                int    `json:"id"`
	Consecutivo       int    `json:"consecutivo"`
	Prefijo           string `json:"prefijo"`
	FechaMovimiento   string `json:"fecha_movimiento"`
	BodegaOrigen      string `json:"bodega_origen"`
	BodegaDestino     string `json:"bodega_destino"`
	ReferenciaExterna string `json:"referencia_externa"`
	Estado            string `json:"estado"`
}

type TrasladoDetalleResponse struct {
	IDArticulo   int     `json:"id_articulo"`
	Codigo       string  `json:"codigo"`
	Nombre       string  `json:"nombre"`
	UnidadMedida string  `json:"unidad_medida"`
	Lote         string  `json:"lote"`
	Cantidad     float64 `json:"cantidad"`
	Observacion  string  `json:"observacion"`
}

type TrasladoFullResponse struct {
	ID                int                       `json:"id"`
	Consecutivo       int                       `json:"consecutivo"`
	Prefijo           string                    `json:"prefijo"`
	Comprobante       string                    `json:"comprobante"`
	FechaMovimiento   string                    `json:"fecha_movimiento"`
	BodegaOrigen      string                    `json:"bodega_origen"`
	BodegaDestino     string                    `json:"bodega_destino"`
	ReferenciaExterna string                    `json:"referencia_externa"`
	Observaciones     string                    `json:"observaciones"`
	Estado            string                    `json:"estado"`
	Detalle           []TrasladoDetalleResponse `json:"detalle"    gorm:"-"`
}

// ─── Refs ─────────────────────────────────────────────────────

type RefBodega struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

type RefLote struct {
	Lote        string  `json:"lote"` // NULL = sin lote
	StockActual float64 `json:"stock_actual"`
}

type RefComprobanteTraslado struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Prefijo string `json:"prefijo"`
}

// ─── Anulación ────────────────────────────────────────────────

type AnulacionTrasladoRequest struct {
	Motivo    string `json:"motivo"     validate:"required,min=5"`
	UserID    int64  `json:"user_id"    validate:"omitempty"`
	EmpresaID int64  `json:"empresa_id" validate:"omitempty"`
}
