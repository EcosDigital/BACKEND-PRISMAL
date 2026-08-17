package proyectos

import "time"

// ─── Constantes de dominio ────────────────────────────────────────────────────

const (
	// Estados del ciclo de vida de un proyecto (ref_estado_proyecto.id)
	EstadoPlaneacion = 1
	EstadoEnProgreso = 2
	EstadoEnPausa    = 3
	EstadoFinalizado = 4
	EstadoCancelado  = 5

	// Prioridades de un proyecto (ref_prioridad_proyecto.id)
	PrioridadBaja    = 1
	PrioridadMedia   = 2
	PrioridadAlta    = 3
	PrioridadCritica = 4

	// Tipos de evidencia (ref_tipo_evidencia.id)
	TipoEvidenciaImagen         = 1
	TipoEvidenciaAudio          = 2
	TipoEvidenciaVideo          = 3
	TipoEvidenciaDocumento      = 4
	TipoEvidenciaPDF            = 5
	TipoEvidenciaArchivoGeneral = 6
)

// ─── Modelos de catálogo ──────────────────────────────────────────────────────

type RefEstadoProyecto struct {
	ID       int    `json:"id"        gorm:"column:id"`
	Nombre   string `json:"nombre"    gorm:"column:nombre"`
	Color    string `json:"color"     gorm:"column:color"`
	EsActivo bool   `json:"es_activo" gorm:"column:es_activo"`
}

func (RefEstadoProyecto) TableName() string {
	return "organizacion.ref_estado_proyecto"
}

type RefPrioridadProyecto struct {
	ID       int    `json:"id"        gorm:"column:id"`
	Nombre   string `json:"nombre"    gorm:"column:nombre"`
	Color    string `json:"color"     gorm:"column:color"`
	EsActivo bool   `json:"es_activo" gorm:"column:es_activo"`
}

func (RefPrioridadProyecto) TableName() string {
	return "organizacion.ref_prioridad_proyecto"
}

type RefTipoEvidencia struct {
	ID       int    `json:"id"        gorm:"column:id"`
	Nombre   string `json:"nombre"    gorm:"column:nombre"`
	EsActivo bool   `json:"es_activo" gorm:"column:es_activo"`
}

func (RefTipoEvidencia) TableName() string {
	return "organizacion.ref_tipo_evidencia"
}

// ─── Modelos de dominio ───────────────────────────────────────────────────────

// Proyecto es el modelo principal del módulo.
// Los nombres de campo en Go son descriptivos; gorm:"column:..." apunta
// al nombre exacto de la columna en español definido en la migración.

type Proyecto struct {
	ID             int64      `json:"id"              gorm:"column:id;primaryKey"`
	Codigo         string     `json:"codigo"          gorm:"column:codigo"`
	Nombre         string     `json:"nombre"          gorm:"column:nombre"`
	Descripcion    string     `json:"descripcion"     gorm:"column:descripcion"`
	FechaInicio    time.Time  `json:"fecha_inicio"    gorm:"column:fecha_inicio"`
	FechaFin       *time.Time `json:"fecha_fin"       gorm:"column:fecha_fin"`
	IDEstado       int        `json:"id_estado"       gorm:"column:id_estado"`
	IDPrioridad    int        `json:"id_prioridad"    gorm:"column:id_prioridad"`
	IDResponsable  int64      `json:"id_responsable"  gorm:"column:id_responsable"`
	EsActivo       bool       `json:"es_activo"       gorm:"column:es_activo"`
	IDEmpresa      int64      `json:"id_empresa"      gorm:"column:id_empresa"`
	CreadoPor      int64      `json:"created_by"      gorm:"column:created_by"`
	ActualizadoPor int64      `json:"updated_by" gorm:"column:updated_by"`
	CreadoEn       time.Time  `json:"created_at"       gorm:"column:created_at"`
	ActualizadoEn  time.Time  `json:"updated_at"  gorm:"column:updated_at"`
}

func (Proyecto) TableName() string {
	return "organizacion.cfg_proyectos"
}

// ParteInteresada representa un stakeholder asociado a un proyecto.
type ParteInteresada struct {
	ID             int64     `json:"id"              gorm:"column:id;primaryKey"`
	IDProyecto     int64     `json:"id_proyecto"     gorm:"column:id_proyecto"`
	Nombre         string    `json:"nombre"          gorm:"column:nombre"`
	Rol            string    `json:"rol"             gorm:"column:rol"`
	Organizacion   string    `json:"organizacion"    gorm:"column:organizacion"`
	Correo         string    `json:"correo"          gorm:"column:correo"`
	Telefono       string    `json:"telefono"        gorm:"column:telefono"`
	Observaciones  string    `json:"observaciones"   gorm:"column:observaciones"`
	EsActivo       bool      `json:"es_activo"       gorm:"column:es_activo"`
	IDEmpresa      int64     `json:"id_empresa"      gorm:"column:id_empresa"`
	CreadoPor      int64     `json:"created_by"      gorm:"column:created_by"`
	ActualizadoPor int64     `json:"updated_by" gorm:"column:updated_by"`
	CreadoEn       time.Time `json:"created_at"       gorm:"column:created_at"`
	ActualizadoEn  time.Time `json:"updated_at"  gorm:"column:updated_at"`
}

func (ParteInteresada) TableName() string {
	return "organizacion.cfg_partes_interesadas"
}

// Comentario representa una entrada de la bitácora cronológica.
// Sin soft delete: el historial es inmutable.
type Comentario struct {
	ID             int64     `json:"id"              gorm:"column:id;primaryKey"`
	IDProyecto     int64     `json:"id_proyecto"     gorm:"column:id_proyecto"`
	FechaActividad time.Time `json:"fecha_actividad" gorm:"column:fecha_actividad"`
	Comentario     string    `json:"comentario"      gorm:"column:comentario"`
	IDEmpresa      int64     `json:"id_empresa"      gorm:"column:id_empresa"`
	CreadoPor      int64     `json:"created_by"      gorm:"column:created_by"`
	CreadoEn       time.Time `json:"created_at"       gorm:"column:created_at"`
}

func (Comentario) TableName() string {
	return "organizacion.bit_comentarios"
}

// Evidencia representa un archivo adjunto a un comentario de la bitácora.
type Evidencia struct {
	ID            int64     `json:"id"             gorm:"column:id;primaryKey"`
	IDComentario  int64     `json:"id_comentario"  gorm:"column:id_comentario"`
	NombreArchivo string    `json:"nombre_archivo" gorm:"column:nombre_archivo"`
	RutaArchivo   string    `json:"ruta_archivo"   gorm:"column:ruta_archivo"`
	IDTipo        int       `json:"id_tipo"        gorm:"column:id_tipo"`
	IDEmpresa     int64     `json:"id_empresa"     gorm:"column:id_empresa"`
	SubidoPor     int64     `json:"created_by"     gorm:"column:subido_por"`
	SubidoEn      time.Time `json:"created_at"      gorm:"column:subido_en"`
}

func (Evidencia) TableName() string {
	return "organizacion.bit_evidencias"
}

// ─── Requests: Proyectos ──────────────────────────────────────────────────────

type CrearProyectoRequest struct {
	Codigo        string  `json:"codigo"        validate:"required,max=50"`
	Nombre        string  `json:"nombre"        validate:"required,max=500"`
	Descripcion   string  `json:"descripcion"`
	FechaInicio   string  `json:"fecha_inicio"  validate:"required"`
	FechaFin      *string `json:"fecha_fin"`
	IDPrioridad   int     `json:"id_prioridad"  validate:"required,min=1"`
	IDResponsable int64   `json:"id_responsable" validate:"required,min=1"`

	// Inyectados desde JWT — no vienen del body
	IDEmpresa int64 `json:"-"`
	IDUsuario int64 `json:"-"`
}

type ActualizarProyectoRequest struct {
	Nombre        string  `json:"nombre"        validate:"required,max=500"`
	Descripcion   string  `json:"descripcion"`
	FechaInicio   string  `json:"fecha_inicio"  validate:"required"` // YYYY-MM-DD
	FechaFin      *string `json:"fecha_fin"`                         // YYYY-MM-DD, opcional
	IDEstado      int     `json:"id_estado"     validate:"required,min=1"`
	IDPrioridad   int     `json:"id_prioridad"  validate:"required,min=1"`
	IDResponsable int64   `json:"id_responsable" validate:"required,min=1"`

	// Inyectados desde JWT
	IDEmpresa int64 `json:"-"`
	IDUsuario int64 `json:"-"`
}

type FiltrosProyecto struct {
	IDEstado      int    // 0 = sin filtro
	IDPrioridad   int    // 0 = sin filtro
	IDResponsable int64  // 0 = sin filtro
	Busqueda      string // "" = sin filtro (nombre o código)
}

// ─── Requests: Partes interesadas ─────────────────────────────────────────────

type CrearParteInteresadaRequest struct {
	IDProyecto    int64  `json:"id_proyecto"   validate:"required,min=1"`
	Nombre        string `json:"nombre"        validate:"required,max=300"`
	Rol           string `json:"rol"`
	Organizacion  string `json:"organizacion"`
	Correo        string `json:"correo"`
	Telefono      string `json:"telefono"`
	Observaciones string `json:"observaciones"`

	// Inyectados desde JWT
	IDEmpresa int64 `json:"-"`
	IDUsuario int64 `json:"-"`
}

type ActualizarParteInteresadaRequest struct {
	Nombre        string `json:"nombre"        validate:"required,max=300"`
	Rol           string `json:"rol"`
	Organizacion  string `json:"organizacion"`
	Correo        string `json:"correo"`
	Telefono      string `json:"telefono"`
	Observaciones string `json:"observaciones"`

	// Inyectados desde JWT
	IDEmpresa int64 `json:"-"`
	IDUsuario int64 `json:"-"`
}

// ─── Requests: Comentarios ────────────────────────────────────────────────────

type CrearComentarioRequest struct {
	IDProyecto     int64  `json:"id_proyecto"     validate:"required,min=1"`
	FechaActividad string `json:"fecha_actividad" validate:"required"`
	Comentario     string `json:"comentario"      validate:"required"`

	// Inyectados desde JWT
	IDEmpresa int64 `json:"-"`
	IDUsuario int64 `json:"-"`
}

// ─── Requests: Evidencias ─────────────────────────────────────────────────────
type CrearEvidenciaRequest struct {
	IDComentario  int64  `json:"id_comentario"  validate:"required,min=1"`
	NombreArchivo string `json:"nombre_archivo" validate:"required,max=500"`
	RutaArchivo   string `json:"ruta_archivo"   validate:"required"`
	IDTipo        int    `json:"id_tipo"        validate:"required,min=1"`

	// Inyectados desde JWT
	IDEmpresa int64 `json:"-"`
	IDUsuario int64 `json:"-"`
}

// ─── Responses ────────────────────────────────────────────────────────────────

type ProyectoResponse struct {
	ID             int64   `json:"id"`
	Codigo         string  `json:"codigo"`
	Nombre         string  `json:"nombre"`
	Descripcion    string  `json:"descripcion"`
	FechaInicio    string  `json:"fecha_inicio"`
	FechaFin       *string `json:"fecha_fin"`
	IDEstado       int     `json:"id_estado"`
	Estado         string  `json:"estado"`
	ColorEstado    string  `json:"color_estado"`
	IDPrioridad    int     `json:"id_prioridad"`
	Prioridad      string  `json:"prioridad"`
	ColorPrioridad string  `json:"color_prioridad"`
	IDResponsable  int64   `json:"id_responsable"`
	CreadoEn       string  `json:"created_at"`
	ActualizadoEn  string  `json:"updated_at"`
}

type ParteInteresadaResponse struct {
	ID            int64  `json:"id"`
	IDProyecto    int64  `json:"id_proyecto"`
	Nombre        string `json:"nombre"`
	Rol           string `json:"rol"`
	Organizacion  string `json:"organizacion"`
	Correo        string `json:"correo"`
	Telefono      string `json:"telefono"`
	Observaciones string `json:"observaciones"`
	CreadoEn      string `json:"created_at"`
	ActualizadoEn string `json:"updated_at"`
}

type ComentarioResponse struct {
	ID             int64               `json:"id"`
	IDProyecto     int64               `json:"id_proyecto"`
	FechaActividad string              `json:"fecha_actividad"`
	Comentario     string              `json:"comentario"`
	CreadoPor      int64               `json:"created_by"`
	CreadoEn       string              `json:"created_at"`
	Evidencias     []EvidenciaResponse `json:"evidencias"`
}

type EvidenciaResponse struct {
	ID            int64  `json:"id"`
	IDComentario  int64  `json:"id_comentario"`
	NombreArchivo string `json:"nombre_archivo"`
	RutaArchivo   string `json:"ruta_archivo"`
	IDTipo        int    `json:"id_tipo"`
	TipoEvidencia string `json:"tipo_evidencia"`
	SubidoEn      string `json:"subido_en"`
}

// CatalogoResponse agrupa estados y prioridades en una sola llamada al frontend.
type CatalogoResponse struct {
	Estados     []RefEstadoProyecto    `json:"estados"`
	Prioridades []RefPrioridadProyecto `json:"prioridades"`
}

// ─── Envoltorio estándar de respuesta HTTP ────────────────────────────────────

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
