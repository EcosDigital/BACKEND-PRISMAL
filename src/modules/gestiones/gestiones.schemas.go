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

// TicketFiltros recibe los query params del listado.
// FIX: se agrega Pagina y PerPagina para evitar devolver todos los registros sin límite.
type TicketFiltros struct {
	IDEstado      int    `query:"id_estado"`
	IDNivel       int    `query:"id_nivel"`
	IDOrigen      int    `query:"id_origen"`
	IDColaborador int    `query:"id_colaborador"`
	Numero        string `query:"numero"`
	SoloAsignados bool   `query:"solo_asignados"`
	// Paginación
	Pagina    int `query:"pagina"`     // default 1
	PerPagina int `query:"per_pagina"` // default 25, máx 100
}

// PaginatedTickets envuelve el resultado con metadatos de paginación.
type PaginatedTickets struct {
	Data      []TicketResponse `json:"data"`
	Total     int64            `json:"total"`
	Pagina    int              `json:"pagina"`
	PerPagina int              `json:"per_pagina"`
	Paginas   int              `json:"paginas"`
}

type GestionRequest struct {
	Comentario    string `json:"comentario"         validate:"required,min=3"`
	IDEstadoNuevo int    `json:"id_estado_nuevo"    validate:"required,gt=0"`
	UserID        int64  `json:"user_id"            validate:"omitempty"`
}

type GestionResponse struct {
	ID         int64  `json:"id"`
	IDTicket   int64  `json:"id_ticket"`
	Comentario string `json:"comentario"`
	// Estado que tenía el ticket al momento de esta gestión
	IDEstadoAnterior int    `json:"id_estado_anterior"`
	EstadoAnterior   string `json:"estado_anterior"`
	ColorAnterior    string `json:"color_anterior"`
	// Estado al que se cambió
	IDEstadoNuevo int    `json:"id_estado_nuevo"`
	EstadoNuevo   string `json:"estado_nuevo"`
	ColorNuevo    string `json:"color_nuevo"`
	// Autor
	IDAutor  int64  `json:"id_autor"`
	Autor    string `json:"autor"`
	CreadoEn string `json:"created_at"`
}

// ─── Asignación ───────────────────────────────────────────────────────────────

// AsignacionRequest body para POST /ticket/:id/asignar
type AsignacionRequest struct {
	IDColaboradores []int64 `json:"id_colaboradores" validate:"required,min=1"`
	Tarea           string  `json:"tarea"            validate:"omitempty,max=1000"`
	UserID          int64   `json:"user_id"          validate:"omitempty"`
}

// AsignadoResponse representa un colaborador asignado a un ticket
type AsignadoResponse struct {
	ID            int64  `json:"id"`
	IDTicket      int64  `json:"id_ticket"`
	IDColaborador int64  `json:"id_colaborador"`
	Nombre        string `json:"nombre"`
	Tarea         string `json:"tarea"`
	AsignadoEn    string `json:"created_at"`
}

// ─── Estadísticas globales ────────────────────────────────────────────────────

// StatGlobalResponse es un contador por estado, independiente de filtros.
type StatGlobalResponse struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	ColorHex string `json:"color_hex"`
	Total    int    `json:"total"`
}

// StatsResponse agrupa los contadores globales y los indicadores operativos.
type StatsResponse struct {
	PorEstado  []StatGlobalResponse `json:"por_estado"`
	Vencidos   int                  `json:"vencidos"`
	SinAsignar int                  `json:"sin_asignar"`
	CreadosHoy int                  `json:"creados_hoy"`
}
