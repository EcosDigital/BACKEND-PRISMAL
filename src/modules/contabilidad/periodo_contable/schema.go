package periodo_contable

import "time"

// ─── Constantes de dominio ────────────────────────────────────────────────────

// Estados de período (ref_estado_periodo)
const (
	EstadoPeriodoAbierto   = 1 // 'Abierto'
	EstadoPeriodoCerrado   = 2 // 'Cerrado'
	EstadoPeriodoAjuste    = 3 // 'Ajuste'
	EstadoPeriodoBloqueado = 4 // 'Bloqueado'
)

// Estados de módulo dentro de un período (ref_estado_modulo)
const (
	EstadoModuloAbierto     = 1 // 'Abierto'
	EstadoModuloCerrado     = 2 // 'Cerrado'
	EstadoModuloSoloLectura = 3 // 'Solo Lectura'
	EstadoModuloRestringido = 4 // 'Restringido'
)

// Módulos del sistema (ref_modulo_contable)
const (
	ModuloFacturacion  = 1
	ModuloContabilidad = 2
	ModuloInventario   = 3
	ModuloCompras      = 4
	ModuloTesoreria    = 5
	ModuloNomina       = 6
)

// NombreMes mapea número de mes a nombre en español.
var NombreMes = map[int]string{
	1:  "Enero",
	2:  "Febrero",
	3:  "Marzo",
	4:  "Abril",
	5:  "Mayo",
	6:  "Junio",
	7:  "Julio",
	8:  "Agosto",
	9:  "Septiembre",
	10: "Octubre",
	11: "Noviembre",
	12: "Diciembre",
}

// ─── Modelos de dominio (catálogos ref_*) ─────────────────────────────────────

// RefEstadoPeriodo representa una fila del catálogo de estados de período.
type RefEstadoPeriodo struct {
	ID          int    `json:"id"          gorm:"column:id"`
	Nombre      string `json:"nombre"      gorm:"column:nombre"`
	Descripcion string `json:"descripcion" gorm:"column:descripcion"`
	Color       string `json:"color"       gorm:"column:color"`
}

func (RefEstadoPeriodo) TableName() string {
	return "contabilidad.ref_estado_periodo"
}

// RefModuloContable representa un módulo del sistema que puede controlarse por período.
type RefModuloContable struct {
	ID          int    `json:"id"          gorm:"column:id"`
	Codigo      string `json:"codigo"      gorm:"column:codigo"`
	Nombre      string `json:"nombre"      gorm:"column:nombre"`
	Descripcion string `json:"descripcion" gorm:"column:descripcion"`
	IsActive    bool   `json:"is_active"   gorm:"column:is_active"`
}

func (RefModuloContable) TableName() string {
	return "contabilidad.ref_modulo_contable"
}

// RefEstadoModulo representa el estado operacional de un módulo en un período.
type RefEstadoModulo struct {
	ID          int    `json:"id"           gorm:"column:id"`
	Nombre      string `json:"nombre"       gorm:"column:nombre"`
	Descripcion string `json:"descripcion"  gorm:"column:descripcion"`
	PermiteOp   bool   `json:"permite_op"   gorm:"column:permite_op"`
	Color       string `json:"color"        gorm:"column:color"`
}

func (RefEstadoModulo) TableName() string {
	return "contabilidad.ref_estado_modulo"
}

// ─── Modelos de dominio (tablas cfg_*) ────────────────────────────────────────

// PeriodoContable es el modelo principal de un período mensual.
type PeriodoContable struct {
	ID          int64     `json:"id"            gorm:"column:id;primaryKey"`
	IDAñoFiscal int64     `json:"id_año_fiscal" gorm:"column:id_año_fiscal"`
	EmpresaID   int64     `json:"empresa_id"    gorm:"column:empresa_id"`
	NumeroMes   int       `json:"numero_mes"    gorm:"column:numero_mes"`
	NombreMes   string    `json:"nombre_mes"    gorm:"column:nombre_mes"`
	FechaInicio time.Time `json:"fecha_inicio"  gorm:"column:fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"     gorm:"column:fecha_fin"`
	IDEstado    int       `json:"id_estado"     gorm:"column:id_estado"`
	Notas       string    `json:"notas"         gorm:"column:notas"`
	CreatedBy   int64     `json:"created_by"    gorm:"column:created_by"`
	UpdatedBy   int64     `json:"updated_by"    gorm:"column:updated_by"`
	CreatedAt   time.Time `json:"created_at"    gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at"    gorm:"column:updated_at"`
}

func (PeriodoContable) TableName() string {
	return "contabilidad.cfg_periodos_contables"
}

// ControlModuloPeriodo representa el estado de un módulo dentro de un período.
type ControlModuloPeriodo struct {
	ID        int64     `json:"id"          gorm:"column:id;primaryKey"`
	IDPeriodo int64     `json:"id_periodo"  gorm:"column:id_periodo"`
	IDModulo  int       `json:"id_modulo"   gorm:"column:id_modulo"`
	IDEstado  int       `json:"id_estado"   gorm:"column:id_estado"`
	EmpresaID int64     `json:"empresa_id"  gorm:"column:empresa_id"`
	Notas     string    `json:"notas"       gorm:"column:notas"`
	UpdatedBy int64     `json:"updated_by"  gorm:"column:updated_by"`
	UpdatedAt time.Time `json:"updated_at"  gorm:"column:updated_at"`
	CreatedAt time.Time `json:"created_at"  gorm:"column:created_at"`
}

func (ControlModuloPeriodo) TableName() string {
	return "contabilidad.cfg_control_modulos_periodo"
}

// ─── Requests ─────────────────────────────────────────────────────────────────

// GenerarPeriodosRequest es el DTO para generar los 12 períodos de un año fiscal.
type GenerarPeriodosRequest struct {
	IDAñoFiscal int64 `json:"id_año_fiscal" validate:"required"`

	// Inyectados desde JWT
	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
}

// UpdateEstadoPeriodoRequest cambia el estado de un período contable.
type UpdateEstadoPeriodoRequest struct {
	IDEstado int    `json:"id_estado" validate:"required,min=1"`
	Notas    string `json:"notas"`

	// Inyectados desde JWT
	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
}

// UpdateEstadoModuloPeriodoRequest cambia el estado de un módulo en un período.
type UpdateEstadoModuloPeriodoRequest struct {
	IDEstadoModulo int    `json:"id_estado_modulo" validate:"required,min=1"`
	Notas          string `json:"notas"`

	// Inyectados desde JWT
	EmpresaID int64 `json:"-"`
	UserID    int64 `json:"-"`
}

// ValidarOperacionRequest es el DTO para la validación centralizada de operaciones.
type ValidarOperacionRequest struct {
	EmpresaID int64
	Fecha     time.Time
	IDModulo  int
}

// ─── Responses ────────────────────────────────────────────────────────────────

// ControlModuloPeriodoResponse es el DTO público del control de módulos.
type ControlModuloPeriodoResponse struct {
	IDModulo     int    `json:"id_modulo"`
	CodigoModulo string `json:"codigo_modulo"`
	NombreModulo string `json:"nombre_modulo"`
	IDEstado     int    `json:"id_estado"`
	Estado       string `json:"estado"`
	PermiteOp    bool   `json:"permite_op"`
	Color        string `json:"color"`
	Notas        string `json:"notas"`
}

// PeriodoContableResponse es el DTO público con estados y módulos resueltos.
type PeriodoContableResponse struct {
	ID          int64                          `json:"id"`
	IDAñoFiscal int64                          `json:"id_año_fiscal"`
	AñoFiscal   int                            `json:"año_fiscal"`
	EmpresaID   int64                          `json:"empresa_id"`
	NumeroMes   int                            `json:"numero_mes"`
	NombreMes   string                         `json:"nombre_mes"`
	FechaInicio string                         `json:"fecha_inicio"` // "DD/MM/YYYY"
	FechaFin    string                         `json:"fecha_fin"`    // "DD/MM/YYYY"
	IDEstado    int                            `json:"id_estado"`
	Estado      string                         `json:"estado"`
	ColorEstado string                         `json:"color_estado"`
	Notas       string                         `json:"notas"`
	Modulos     []ControlModuloPeriodoResponse `json:"modulos"`
	CreatedAt   string                         `json:"created_at"`
	UpdatedAt   string                         `json:"updated_at"`
}

// ValidacionOperacionResponse indica si una operación puede ejecutarse.
type ValidacionOperacionResponse struct {
	Puede     bool   `json:"puede"`
	Motivo    string `json:"motivo,omitempty"`
	IDPeriodo int64  `json:"id_periodo,omitempty"`
	IDEstado  int    `json:"id_estado,omitempty"`
}

// CatalogosResponse agrupa todos los catálogos de referencia en un solo endpoint.
type CatalogosResponse struct {
	EstadosPeriodo []RefEstadoPeriodo  `json:"estados_periodo"`
	Modulos        []RefModuloContable `json:"modulos"`
	EstadosModulo  []RefEstadoModulo   `json:"estados_modulo"`
}

// ─── Envoltorio estándar de respuesta HTTP ────────────────────────────────────
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
