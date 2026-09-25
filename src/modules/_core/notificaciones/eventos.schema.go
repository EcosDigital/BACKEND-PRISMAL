package notificaciones

// modulo contratado por el tenant (configuracion.ref_modulos_tenant).
// id_ref es el id del modulo en el catalogo de la base admin
type moduloTenant struct {
	ID    int64 `gorm:"column:id"`
	IDRef int64 `gorm:"column:id_ref"`
}

// evento del catalogo maestro (configuracion.cfg_eventos_notificacion en la
// base admin). el tipo viene por nombre para traducirlo al id local
type eventoCatalogo struct {
	ID          int64   `gorm:"column:id"`
	Codigo      string  `gorm:"column:codigo"`
	Nombre      string  `gorm:"column:nombre"`
	Descripcion *string `gorm:"column:descripcion"`
	IDModulo    int64   `gorm:"column:id_modulo"`
	Tipo        string  `gorm:"column:tipo"`
	OrdenLista  int     `gorm:"column:orden_lista"`
}

// evento ya copiado en el tenant, para compararlo contra el catalogo
type eventoLocal struct {
	IDRef              int64   `gorm:"column:id_ref"`
	Codigo             string  `gorm:"column:codigo"`
	Nombre             string  `gorm:"column:nombre"`
	Descripcion        *string `gorm:"column:descripcion"`
	IDModuloRef        int64   `gorm:"column:id_modulo_ref"`
	IDTipoNotificacion int     `gorm:"column:id_tipo_notificacion"`
	OrdenLista         int     `gorm:"column:orden_lista"`
	IsActive           bool    `gorm:"column:is_active"`
}

// evento local con su modulo, antes de agruparlo por modulo
type eventoModuloRow struct {
	IDModulo     int64  `gorm:"column:id_modulo"`
	CodigoModulo string `gorm:"column:codigo_modulo"`
	NombreModulo string `gorm:"column:nombre_modulo"`
	Codigo       string `gorm:"column:codigo"`
	Nombre       string `gorm:"column:nombre"`
	Descripcion  string `gorm:"column:descripcion"`
	Tipo         string `gorm:"column:tipo"`
}

type EventoResponse struct {
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Tipo        string `json:"tipo"`
}

// eventos de un modulo: cada modulo es una pestaña de Parametros generales
type ModuloEventosResponse struct {
	IDModulo int64            `json:"id_modulo"`
	Codigo   string           `json:"codigo"`
	Nombre   string           `json:"nombre"`
	Eventos  []EventoResponse `json:"eventos"`
}

// body del PUT /core/notificaciones/eventos/:codigo/destinatarios. se envia
// la lista completa: reemplaza lo que habia configurado en la sede
type DestinatariosEventoRequest struct {
	Usuarios  []int64 `json:"usuarios"   validate:"omitempty,dive,gt=0"`
	Roles     []int64 `json:"roles"      validate:"omitempty,dive,gt=0"`
	UserID    int64   `json:"user_id"    validate:"omitempty"`
	EmpresaID int64   `json:"empresa_id" validate:"omitempty"`
	SedeID    int64   `json:"sede_id"    validate:"omitempty"`
}

// a quien le llega un evento en la sede, separado en usuarios y roles
type DestinatariosEventoResponse struct {
	Usuarios []int64 `json:"usuarios"`
	Roles    []int64 `json:"roles"`
}

// fila de notificaciones.cfg_evento_destinatarios
type eventoDestinoRow struct {
	IDOrigen  int   `gorm:"column:id_origen"`
	IDDestino int64 `gorm:"column:id_destino"`
}

// usuario para la pantalla de Parametros generales: se filtra por rol y se
// muestra en la tabla de seleccionados, por eso el correo y el rol van
// aparte y no pegados en el nombre
type UsuarioDestinoDetalle struct {
	ID     int64  `json:"id"     gorm:"column:id"`
	Nombre string `json:"nombre" gorm:"column:nombre"`
	Email  string `json:"email"  gorm:"column:email"`
	IDRol  int64  `json:"id_rol" gorm:"column:id_rol"`
	Rol    string `json:"rol"    gorm:"column:rol"`
}

// evento activo del tenant con lo necesario para disparar su notificacion
type eventoActivo struct {
	ID                 int64 `gorm:"column:id"`
	IDModuloRef        int64 `gorm:"column:id_modulo_ref"`
	IDTipoNotificacion int   `gorm:"column:id_tipo_notificacion"`
}
