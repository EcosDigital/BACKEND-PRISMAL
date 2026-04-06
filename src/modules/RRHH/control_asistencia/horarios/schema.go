package horarios

// ─── Requests ─────────────────────────────────────────────────────────────────

type BloqueHorarioRequest struct {
	HoraEntrada string `json:"hora_entrada" validate:"required"`
	HoraSalida  string `json:"hora_salida"  validate:"required"`
	Orden       int    `json:"orden"        validate:"required,gt=0"`
}

type HorarioRequest struct {
	Codigo            string                 `json:"codigo"             validate:"required,min=1,max=20"`
	Nombre            string                 `json:"nombre"             validate:"required,min=3,max=150"`
	Descripcion       string                 `json:"descripcion"        validate:"omitempty"`
	MinutosTolerancia int                    `json:"minutos_tolerancia" validate:"omitempty,min=0"`
	Bloques           []BloqueHorarioRequest `json:"bloques"            validate:"required,min=1,dive"`
}

type HorarioUpdateRequest struct {
	Nombre            string `json:"nombre"             validate:"required,min=3,max=150"`
	Descripcion       string `json:"descripcion"        validate:"omitempty"`
	MinutosTolerancia int    `json:"minutos_tolerancia" validate:"omitempty,min=0"`
}

type AsignacionHorarioRequest struct {
	IDHorario     int64  `json:"id_horario"    validate:"required,gt=0"`
	IDTercero     int64  `json:"id_tercero"    validate:"required,gt=0"`
	FechaInicio   string `json:"fecha_inicio"  validate:"required"`
	Observaciones string `json:"observaciones" validate:"omitempty"`
	// Inyectado desde contexto JWT
	RegistradoPor int64 `json:"-"`
}

type CambioEstadoHorarioRequest struct {
	Activo bool `json:"activo"`
}

type BloqueHorarioResponse struct {
	ID          int64  `json:"id"`
	IDHorario   int64  `json:"id_horario"`
	HoraEntrada string `json:"hora_entrada"`
	HoraSalida  string `json:"hora_salida"`
	Orden       int    `json:"orden"`
	Activo      bool   `json:"activo"`
}

type HorarioResponse struct {
	ID                int64                   `json:"id"`
	Codigo            string                  `json:"codigo"`
	Nombre            string                  `json:"nombre"`
	Descripcion       string                  `json:"descripcion"`
	MinutosTolerancia int                     `json:"minutos_tolerancia"`
	Activo            bool                    `json:"activo"`
	FechaCreacion     string                  `json:"fecha_creacion"`
	Bloques           []BloqueHorarioResponse `json:"bloques"`
}

type HorarioListResponse struct {
	ID                int64  `json:"id"`
	Codigo            string `json:"codigo"`
	Nombre            string `json:"nombre"`
	MinutosTolerancia int    `json:"minutos_tolerancia"`
	Activo            bool   `json:"activo"`
	FechaCreacion     string `json:"fecha_creacion"`
}

type AsignacionResponse struct {
	IDHorario   int64  `json:"id_horario"`
	IDTercero   int64  `json:"id_tercero"`
	FechaInicio string `json:"fecha_inicio"`
	Total       int    `json:"total_asignados"`
}

type AsignacionListResponse struct {
	ID              int64  `json:"id"`
	IDTercero       int64  `json:"id_tercero"`
	NombreTercero   string `json:"nombre_tercero"`
	NumeroDocumento string `json:"numero_documento"`
	FechaInicio     string `json:"fecha_inicio"`
	Observaciones   string `json:"observaciones"`
}

type horarioRow struct {
	ID                int64  `gorm:"column:id"`
	Codigo            string `gorm:"column:codigo"`
	Nombre            string `gorm:"column:nombre"`
	Descripcion       string `gorm:"column:descripcion"`
	MinutosTolerancia int    `gorm:"column:minutos_tolerancia"`
	Activo            bool   `gorm:"column:activo"`
	FechaCreacion     string `gorm:"column:fecha_creacion"`
}

type bloqueRow struct {
	ID          int64  `gorm:"column:id"`
	IDHorario   int64  `gorm:"column:id_horario"`
	HoraEntrada string `gorm:"column:hora_entrada"`
	HoraSalida  string `gorm:"column:hora_salida"`
	Orden       int    `gorm:"column:orden"`
	Activo      bool   `gorm:"column:activo"`
}
