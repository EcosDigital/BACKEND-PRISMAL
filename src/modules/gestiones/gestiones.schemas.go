package gestiones

type TicketRequest struct {
	Titulo           string `json:"titulo"            validate:"required,min=5,max=255"`
	Descripcion      string `json:"descripcion"       validate:"required,min=10"`
	IDNivel          int    `json:"id_nivel"          validate:"required,gt=0"`
	IDOrigen         int    `json:"id_origen"         validate:"required,gt=0"`
	ContactoNombre   string `json:"contacto_nombre"   validate:"omitempty,min=3,max=200"`
	ContactoEmail    string `json:"contacto_email"    validate:"omitempty,email,max=200"`
	ContactoTelefono string `json:"contacto_telefono" validate:"omitempty,max=30"`
	UserID           int64  `json:"user_id"           validate:"omitempty"`
}

type TicketResponse struct {
	ID               int64  `json:"id"`
	Titulo           string `json:"titulo"`
	Descripcion      string `json:"descripcion"`
	IDNivel          int    `json:"id_nivel"`
	Nivel            string `json:"nivel"`
	ColorNivel       string `json:"color_nivel"`
	IDEstado         int    `json:"id_estado"`
	Estado           string `json:"estado"`
	ColorEstado      string `json:"color_estado"`
	IDOrigen         int    `json:"id_origen"`
	Origen           string `json:"origen"`
	IDTenant         *int   `json:"id_tenant"`
	FechaEntrega     string `json:"fecha_entrega"`
	ContactoNombre   string `json:"contacto_nombre"`
	ContactoEmail    string `json:"contacto_email"`
	ContactoTelefono string `json:"contacto_telefono"`
	CreadoEn         string `json:"created_at"`
}

// TicketFiltros recibe los query params del listado
type TicketFiltros struct {
	IDEstado      int    `query:"id_estado"`
	IDNivel       int    `query:"id_nivel"`
	IDOrigen      int    `query:"id_origen"`
	IDColaborador int    `query:"id_colaborador"`
	Numero        string `query:"numero"`
	SoloAsignados bool   `query:"solo_asignados"`
}
