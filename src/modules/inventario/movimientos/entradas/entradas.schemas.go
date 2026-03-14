package entradas

import "time"

type EntradaDetalleRequest struct {
	IDArticulo    int     `json:"id_articulo"   validate:"gt=0"`
	IDProveedor   *int    `json:"id_proveedor"  validate:"omitempty"`
	Marca         string  `json:"marca"         validate:"omitempty"`
	Lote          string  `json:"lote"          validate:"omitempty"`
	Cantidad      float64 `json:"cantidad"      validate:"gt=0"`
	ValorUnitario float64 `json:"valor_unitario" validate:"gte=0"`
	IDImpuesto    *int    `json:"id_impuesto"   validate:"omitempty"`
	Observacion   string  `json:"observacion"   validate:"omitempty"`
}

type EntradaRequest struct {
	IDComprobante     int                     `json:"id_comprobante"    validate:"gt=0"`
	IDBodega          int                     `json:"id_bodega"         validate:"gt=0"`
	IDTercero         *int                    `json:"id_tercero"        validate:"omitempty"`
	ReferenciaExterna string                  `json:"referencia_externa" validate:"omitempty"`
	Observaciones     string                  `json:"observaciones"     validate:"omitempty"`
	FechaMovimiento   string                  `json:"fecha_movimiento"  validate:"required"`
	Detalle           []EntradaDetalleRequest `json:"detalle"           validate:"required,min=1,dive"`
	// Inyectados desde el contexto JWT
	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
	SedeID    int64 `json:"sede_id"    validate:"omitempty"`
}

type EntradaResponse struct {
	IDMovimiento     int    `json:"id_movimiento"`
	IDMovComprobante int    `json:"id_mov_comprobante"`
	Consecutivo      int    `json:"consecutivo"`
	Prefijo          string `json:"prefijo"`
}

type EntradaListResponse struct {
	ID                int     `json:"id"`
	Consecutivo       int     `json:"consecutivo"`
	Prefijo           string  `json:"prefijo"`
	FechaMovimiento   string  `json:"fecha_movimiento"`
	Proveedor         string  `json:"proveedor"`
	ReferenciaExterna string  `json:"referencia_externa"`
	Bodega            string  `json:"bodega"`
	ValorTotal        float64 `json:"valor_total"`
	Estado            string  `json:"estado"`
}

type EntradaDetalleResponse struct {
	IDArticulo    int     `json:"id_articulo"`
	Codigo        string  `json:"codigo"`
	Nombre        string  `json:"nombre"`
	UnidadMedida  string  `json:"unidad_medida"`
	Marca         string  `json:"marca"`
	Lote          string  `json:"lote"`
	Cantidad      float64 `json:"cantidad"`
	ValorUnitario float64 `json:"valor_unitario"`
	ValorImpuesto float64 `json:"valor_impuesto"`
	ValorTotal    float64 `json:"valor_total"`
	Observacion   string  `json:"observacion"`
}

type EntradaFullResponse struct {
	ID                int                      `json:"id"`
	Consecutivo       int                      `json:"consecutivo"`
	Prefijo           string                   `json:"prefijo"`
	Comprobante       string                   `json:"comprobante"`
	FechaMovimiento   string                   `json:"fecha_movimiento"`
	Bodega            string                   `json:"bodega"`
	Proveedor         string                   `json:"proveedor"`
	ReferenciaExterna string                   `json:"referencia_externa"`
	Observaciones     string                   `json:"observaciones"`
	Estado            string                   `json:"estado"`
	ValorSubtotal     float64                  `json:"valor_subtotal"`
	ValorImpuesto     float64                  `json:"valor_impuesto"`
	ValorTotal        float64                  `json:"valor_total"`
	Detalle           []EntradaDetalleResponse `json:"detalle"`
}

// ─── Refs para el formulario ──────────────────────────────────────────────────

type RefComprobante struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Prefijo string `json:"prefijo"`
}

type RefImpuesto struct {
	ID         int     `json:"id"`
	Nombre     string  `json:"nombre"`
	Porcentaje float64 `json:"porcentaje"`
}

// ─── Existencia para mostrar en el formulario ─────────────────────────────────

type ArticuloExistencia struct {
	IDArticulo    int       `json:"id_articulo"`
	Codigo        string    `json:"codigo"`
	Nombre        string    `json:"nombre"`
	UnidadMedida  string    `json:"unidad_medida"`
	StockActual   float64   `json:"stock_actual"`
	CostoPromedio float64   `json:"costo_promedio"`
	_             time.Time // evita import no usado
}
