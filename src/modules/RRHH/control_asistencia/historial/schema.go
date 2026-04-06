package historial_asistencia

// ─── Request ──────────────────────────────────────────────────

type HistorialAsistenciaRequest struct {
	IDTercero   int64  `query:"id_tercero"   validate:"gt=0"`
	FechaInicio string `query:"fecha_inicio" validate:"required"`
	FechaFin    string `query:"fecha_fin"    validate:"required"`
}

// ─── Fila del Historial ──────────────────────────────────────────

type HistorialAsistenciaFila struct {
	FilaOrden      int    `json:"fila_orden"`
	Fecha          string `json:"fecha"`
	Hora           string `json:"hora"`
	DiaSemana      string `json:"dia_semana"`
	TipoMarcacion  string `json:"tipo_marcacion"`
	PuntoMarcacion string `json:"punto_marcacion"`
	// Estado calculado en backend mediante cruce con horarios y tolerancias.
	// Valores: PUNTUAL | TOLERANCIA | TARDE | SIN_HORARIO
	Estado            string `json:"estado"`
	HorarioAsignado   string `json:"horario_asignado"`
	HoraProgramada    string `json:"hora_programada"`
	MinutosDiferencia int    `json:"minutos_diferencia"`
	Observaciones     string `json:"observaciones"`
}

// ─── Cabecera del reporte ─────────────────────────────────────

type HistorialAsistenciaCabecera struct {
	IDTercero       int64  `json:"id_tercero"`
	NombreEmpleado  string `json:"nombre_empleado"`
	NumeroDocumento string `json:"numero_documento"`
	FechaInicio     string `json:"fecha_inicio"`
	FechaFin        string `json:"fecha_fin"`
	TotalRegistros  int    `json:"total_registros"`
}

// ─── Empresa ──────────────────────────────────────────────────

type HistorialAsistenciaEmpresa struct {
	RazonSocial  string `json:"razon_social"`
	Nit          string `json:"nit"`
	Dv           string `json:"dv"`
	Direccion    string `json:"direccion"`
	Ciudad       string `json:"ciudad"`
	Departamento string `json:"departamento"`
	Telefono     string `json:"telefono"`
	LogoURL      string `json:"logo_url"`
}

// ─── Response completo ────────────────────────────────────────

type HistorialAsistenciaReporte struct {
	Empresa  HistorialAsistenciaEmpresa  `json:"empresa"`
	Cabecera HistorialAsistenciaCabecera `json:"cabecera"`
	Filas    []HistorialAsistenciaFila   `json:"filas"`
}
