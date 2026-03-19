package control_asistencia

import "time"

type PuntoMarcacionRequest struct {
	IDTipoPunto int64   `json:"id_tipo_punto" validate:"gt=0"`
	Nombre      string  `json:"nombre"        validate:"required,min=3,max=150"`
	Descripcion string  `json:"descripcion"   validate:"omitempty"`
	Latitud     float64 `json:"latitud"       validate:"required"`
	Longitud    float64 `json:"longitud"      validate:"required"`
	RadioMetros int     `json:"radio_metros"  validate:"omitempty"`
	// Inyectado desde contexto JWT
	UserID    int64 `json:"user_id"                validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id"             validate:"omitempty"`
	SedeID    int64 `json:"sede_id"                validate:"omitempty"`
}

type UpdatePuntoMarcacionRequest struct {
	SedeID      int64   `json:"id_sede"       validate:"gt=0"`
	EmpresaID   int64   `json:"id_empresa"    validate:"gt=0"`
	IDTipoPunto int64   `json:"id_tipo_punto" validate:"omitempty,gt=0"`
	Nombre      string  `json:"nombre"        validate:"omitempty,min=3,max=150"`
	Descripcion string  `json:"descripcion"   validate:"omitempty"`
	Latitud     float64 `json:"latitud"       validate:"omitempty"`
	Longitud    float64 `json:"longitud"      validate:"omitempty"`
	RadioMetros int     `json:"radio_metros"  validate:"omitempty"`
	// Inyectado desde contexto JWT
	UserID int64 `json:"-"`
}

type CambioEstadoRequest struct {
	Activo bool `json:"activo"`
}

// ─── Response ─────────────────────────────────────────────────────────────────

type PuntoMarcacionResponse struct {
	ID            int64   `json:"id"`
	IDSede        int64   `json:"id_sede"`
	IDTipoPunto   int64   `json:"id_tipo_punto"`
	TipoPunto     string  `json:"tipo_punto"`
	Nombre        string  `json:"nombre"`
	Descripcion   string  `json:"descripcion"`
	TokenQR       string  `json:"token_qr"`
	Latitud       float64 `json:"latitud"`
	Longitud      float64 `json:"longitud"`
	RadioMetros   int     `json:"radio_metros"`
	Activo        bool    `json:"activo"`
	FechaCreacion string  `json:"fecha_creacion"`
}

type PuntoMarcacionListResponse struct {
	ID            int64   `json:"id"`
	IDSede        int64   `json:"id_sede"`
	TipoPunto     string  `json:"tipo_punto"`
	Nombre        string  `json:"nombre"`
	TokenQR       string  `json:"token_qr"`
	Latitud       float64 `json:"latitud"`
	Longitud      float64 `json:"longitud"`
	RadioMetros   int     `json:"radio_metros"`
	Activo        bool    `json:"activo"`
	FechaCreacion string  `json:"fecha_creacion"`
}

// ─── Refs ─────────────────────────────────────────────────────────────────────

type RefTipoPuntoMarcacion struct {
	ID          int    `json:"id"`
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// ─── Dummy para evitar import no usado ───────────────────────────────────────

var _ = time.Now

//=======  MARCACIONES EMPLEADO ================== ///

// MarcacionRequest es el body que envía el frontend al escanear el QR.
// Los campos de trazabilidad (IP, user_agent, id_empleado) se inyectan en backend.
type MarcacionRequest struct {
	Token       string  `json:"token"     validate:"required"`
	Latitud     float64 `json:"latitud"`
	Longitud    float64 `json:"longitud"`
	Precision   int     `json:"precision"`
	IDEmpleado  int64   `json:"-"`
	DireccionIP string  `json:"-"`
	UserAgent   string  `json:"-"`
}

// PuntoMarcacionValidacion contiene los datos del punto necesarios
// para ejecutar las validaciones de geofencing y estado.
type PuntoMarcacionValidacion struct {
	ID          int64
	Latitud     float64
	Longitud    float64
	RadioMetros int
	Activo      bool
}

type MarcacionResponse struct {
	ID               int64  `json:"id"`
	IDEmpleado       int64  `json:"id_empleado"`
	IDPuntoMarcacion int64  `json:"id_punto_marcacion"`
	FechaHora        string `json:"fecha_hora"`
	TipoEvento       string `json:"tipo_evento"`
	Mensaje          string `json:"mensaje"`
}

var _ = time.Now
