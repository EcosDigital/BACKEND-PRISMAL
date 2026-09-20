package notificaciones

import "time"

type TipoNotificacionResponse struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type DispositivoPushRequest struct {
	Token     string `json:"token" validate:"required"`
	Navegador string `json:"navegador" validate:"omitempty,max=100"`
	UserID    int64  `json:"user_id" validate:"omitempty"`
}

type BajaDispositivoRequest struct {
	Token  string `json:"token" validate:"required"`
	UserID int64  `json:"user_id" validate:"omitempty"`
}

// ids de notificaciones.ref_origen_destinatario
const (
	OrigenUsuario = 1
	OrigenRol     = 2
	OrigenModulo  = 3
)

type OrigenDestinatarioResponse struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// NotificacionRequest body del POST /core/notificaciones.
// id_origen define hacia donde va dirigida (usuario, rol o modulo, segun
// notificaciones.ref_origen_destinatario) e id_destino es el id de esa
// persona, rol o modulo. el frontend ya conoce los ids de los catalogos.
type NotificacionRequest struct {
	Titulo             string `json:"titulo" validate:"required,min=3,max=200"`
	Mensaje            string `json:"mensaje" validate:"required,min=3"`
	IDTipoNotificacion int    `json:"id_tipo_notificacion" validate:"required,gt=0"`
	IDOrigen           int    `json:"id_origen" validate:"required,gt=0"`
	IDDestino          int64  `json:"id_destino" validate:"required,gt=0"`
	IDModuloRef        *int64 `json:"id_modulo_ref" validate:"omitempty,gt=0"`
	UserID             int64  `json:"user_id" validate:"omitempty"`
	EmpresaID          int64  `json:"empresa_id" validate:"omitempty"`
	SedeID             int64  `json:"sede_id" validate:"omitempty"`
}

// fila de seguridad.cfg_sedes_roles, usada para resolver que roles tienen
// acceso a un modulo (el arbol de permisos vive como JSON en esa tabla)
type ConfigRolModulos struct {
	IDRol       int64  `json:"id_rol"`
	JsonModules string `json:"json_modules"`
}

// tope de notificaciones que devuelve el historial de la campanita
const maxHistorial = 50

type NotificacionUsuarioResponse struct {
	ID        int64  `json:"id"`
	Titulo    string `json:"titulo"`
	Mensaje   string `json:"mensaje"`
	Tipo      string `json:"tipo"`
	Modulo    string `json:"modulo"`
	Origen    string `json:"origen"`
	Leido     bool   `json:"leido"`
	CreatedAt string `json:"created_at"`
}

// opcion de los selectores de destino del panel de envio
type DestinoResponse struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

// fila de notificaciones.mov_notificaciones para el insert. se usa estructura
// y no mapa porque GORM completa el id del registro creado en el campo ID:
// ese id lo necesita el reparto a los destinatarios
type NotificacionRow struct {
	ID                 int64     `gorm:"column:id;primaryKey"`
	Titulo             string    `gorm:"column:titulo"`
	Mensaje            string    `gorm:"column:mensaje"`
	IDModuloRef        *int64    `gorm:"column:id_modulo_ref"`
	IDTipoNotificacion int       `gorm:"column:id_tipo_notificacion"`
	EnviadoPor         int64     `gorm:"column:enviado_por"`
	IDEmpresa          int64     `gorm:"column:id_empresa"`
	IDSede             int64     `gorm:"column:id_sede"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}
