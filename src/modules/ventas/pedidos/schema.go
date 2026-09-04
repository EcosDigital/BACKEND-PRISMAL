package pedidos

// ─── Responses ───────────────────────────────────────────────────────────────

// PedidoListItem fila del listado de pedidos del tenant.
type PedidoListItem struct {
	ID                 int64   `json:"id"                    gorm:"column:id"`
	IDReferenciaMillave *string `json:"id_referencia_millave" gorm:"column:id_referencia_millave"`
	EstadoCodigo       string  `json:"estado_codigo"         gorm:"column:estado_codigo"`
	EstadoNombre       string  `json:"estado_nombre"         gorm:"column:estado_nombre"`
	ClienteNombre      string  `json:"cliente_nombre"        gorm:"column:cliente_nombre"`
	ValorTotal         float64 `json:"valor_total"           gorm:"column:valor_total"`
	CreatedAt          string  `json:"created_at"            gorm:"column:created_at"`
}

// PedidoDetalleItem fila de integraciones.mov_pedidos_detalle.
type PedidoDetalleItem struct {
	ID              int64   `json:"id"               gorm:"column:id"`
	Codigo          *string `json:"codigo"           gorm:"column:codigo"`
	Nombre          string  `json:"nombre"           gorm:"column:nombre"`
	Cantidad        int     `json:"cantidad"         gorm:"column:cantidad"`
	PrecioUnitario  float64 `json:"precio_unitario"  gorm:"column:precio_unitario"`
	Subtotal        float64 `json:"subtotal"         gorm:"column:subtotal"`
	Observacion     *string `json:"observacion"      gorm:"column:observacion"`
}

// PedidoEntrega fila de integraciones.mov_pedidos_entrega.
type PedidoEntrega struct {
	ClienteNombre      string   `json:"cliente_nombre"       gorm:"column:cliente_nombre"`
	ClienteTelefono    string   `json:"cliente_telefono"     gorm:"column:cliente_telefono"`
	TipoDocumento      string   `json:"tipo_documento"       gorm:"column:tipo_documento"`
	ClienteDocumento   string   `json:"cliente_documento"    gorm:"column:cliente_documento"`
	DireccionEntrega   string   `json:"direccion_entrega"    gorm:"column:direccion_entrega"`
	GeolocalizacionLat *float64 `json:"geolocalizacion_lat"  gorm:"column:geolocalizacion_lat"`
	GeolocalizacionLon *float64 `json:"geolocalizacion_lon"  gorm:"column:geolocalizacion_lon"`
	Referencia         *string  `json:"referencia"           gorm:"column:referencia"`
}

// PedidoDetalle respuesta completa de GET /ventas/pedidos/:id.
type PedidoDetalle struct {
	ID                   int64               `json:"id"                    gorm:"column:id"`
	IDReferenciaMillave  *string             `json:"id_referencia_millave" gorm:"column:id_referencia_millave"`
	EstadoCodigo         string              `json:"estado_codigo"         gorm:"column:estado_codigo"`
	EstadoNombre         string              `json:"estado_nombre"         gorm:"column:estado_nombre"`
	EstadoEsFinal        bool                `json:"estado_es_final"       gorm:"column:estado_es_final"`
	ValorSubtotal        float64             `json:"valor_subtotal"        gorm:"column:valor_subtotal"`
	ValorDomicilio       float64             `json:"valor_domicilio"       gorm:"column:valor_domicilio"`
	ValorTotal           float64             `json:"valor_total"           gorm:"column:valor_total"`
	Observaciones        *string             `json:"observaciones"         gorm:"column:observaciones"`
	MotivoRechazo        *string             `json:"motivo_rechazo"        gorm:"column:motivo_rechazo"`
	CreatedAt            string              `json:"created_at"            gorm:"column:created_at"`
	UpdatedAt            *string             `json:"updated_at"            gorm:"column:updated_at"`
	Detalle              []PedidoDetalleItem `json:"detalle" gorm:"-"`
	Entrega              PedidoEntrega       `json:"entrega" gorm:"-"`
}

// ─── Requests ────────────────────────────────────────────────────────────────

// CambiarEstadoRequest body del PUT /ventas/pedidos/:id/estado
type CambiarEstadoRequest struct {
	CodigoEstado  string `json:"codigo_estado"  validate:"required"`
	MotivoRechazo string `json:"motivo_rechazo" validate:"omitempty"`
}
