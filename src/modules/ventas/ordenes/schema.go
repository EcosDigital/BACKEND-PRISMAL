package ordenes

// ─── Requests ────────────────────────────────────────────────────────────────

// OrdenDetalleRequest una línea (artículo + cantidad) de la orden.
type OrdenDetalleRequest struct {
	IDArticuloOrigen int64   `json:"id_articulo_origen" validate:"required,gt=0"`
	Codigo           string  `json:"codigo"`
	Nombre           string  `json:"nombre"           validate:"required"`
	Cantidad         int     `json:"cantidad"         validate:"required,gt=0"`
	PrecioUnitario   float64 `json:"precio_unitario"  validate:"gte=0"`
	Observacion      string  `json:"observacion"`
}

// OrdenRequest body del POST /ventas/ordenes.
// IDMesa es nulo cuando no hay mesa física disponible (mesa "ficticia");
// IdentificadorMesa siempre se llena — con el código/nombre real o con lo
// que el mesero escriba.
type OrdenRequest struct {
	IDMesa            *int64                `json:"id_mesa"`
	IdentificadorMesa string                `json:"identificador_mesa" validate:"required"`
	Observaciones     string                `json:"observaciones"`
	Detalle           []OrdenDetalleRequest `json:"detalle" validate:"required,min=1,dive"`
	UserID            int64                 `json:"user_id"    validate:"omitempty"`
	EmpresaID         int64                 `json:"empresa_id" validate:"omitempty"`
	SedeID            int64                 `json:"sede_id"    validate:"omitempty"`
}

// CambiarEstadoOrdenRequest body del PUT /ventas/ordenes/:id/estado
type CambiarEstadoOrdenRequest struct {
	CodigoEstado  string `json:"codigo_estado"  validate:"required"`
	MotivoRechazo string `json:"motivo_rechazo" validate:"omitempty"`
}

// ─── Responses ───────────────────────────────────────────────────────────────

// OrdenListItem fila del listado de órdenes del tenant.
type OrdenListItem struct {
	ID                int64   `json:"id"                gorm:"column:id"`
	IdentificadorMesa string  `json:"identificador_mesa" gorm:"column:identificador_mesa"`
	EstadoCodigo      string  `json:"estado_codigo"      gorm:"column:estado_codigo"`
	EstadoNombre      string  `json:"estado_nombre"      gorm:"column:estado_nombre"`
	ValorTotal        float64 `json:"valor_total"        gorm:"column:valor_total"`
	CreatedAt         string  `json:"created_at"         gorm:"column:created_at"`
}

// OrdenDetalleItem fila de ventas.mov_ordenes_detalle.
type OrdenDetalleItem struct {
	ID             int64   `json:"id"              gorm:"column:id"`
	Codigo         *string `json:"codigo"          gorm:"column:codigo"`
	Nombre         string  `json:"nombre"          gorm:"column:nombre"`
	Cantidad       int     `json:"cantidad"        gorm:"column:cantidad"`
	PrecioUnitario float64 `json:"precio_unitario" gorm:"column:precio_unitario"`
	Subtotal       float64 `json:"subtotal"        gorm:"column:subtotal"`
	Observacion    *string `json:"observacion"     gorm:"column:observacion"`
}

// OrdenDetalle respuesta completa de GET /ventas/ordenes/:id.
type OrdenDetalle struct {
	ID                int64              `json:"id"                 gorm:"column:id"`
	IdentificadorMesa string             `json:"identificador_mesa" gorm:"column:identificador_mesa"`
	EstadoCodigo      string             `json:"estado_codigo"      gorm:"column:estado_codigo"`
	EstadoNombre      string             `json:"estado_nombre"      gorm:"column:estado_nombre"`
	EstadoEsFinal     bool               `json:"estado_es_final"    gorm:"column:estado_es_final"`
	ValorSubtotal     float64            `json:"valor_subtotal"     gorm:"column:valor_subtotal"`
	ValorTotal        float64            `json:"valor_total"        gorm:"column:valor_total"`
	Observaciones     *string            `json:"observaciones"      gorm:"column:observaciones"`
	MotivoRechazo     *string            `json:"motivo_rechazo"     gorm:"column:motivo_rechazo"`
	CreatedAt         string             `json:"created_at"         gorm:"column:created_at"`
	UpdatedAt         *string            `json:"updated_at"         gorm:"column:updated_at"`
	Detalle           []OrdenDetalleItem `json:"detalle" gorm:"-"`
}
