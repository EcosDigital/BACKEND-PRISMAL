package reportes_asistencia

// ─── Request ──────────────────────────────────────────────────

type ReporteDiarioRequest struct {
	FechaInicio string `query:"fecha_inicio" validate:"required"`
	FechaFin    string `query:"fecha_fin"    validate:"required"`
}

// ─── Fila del reporte ─────────────────────────────────────────

type ReporteDiarioFila struct {
	IDTercero       int64  `json:"id_tercero"`
	NombreEmpleado  string `json:"nombre_empleado"`
	NumeroDocumento string `json:"numero_documento"`
	HoraEntrada     string `json:"hora_entrada"` // HH:MM:SS — vacío si ausente
	HoraSalida      string `json:"hora_salida"`  // HH:MM:SS — vacío si ausente
	// Estado calculado en backend: PUNTUAL | TOLERANCIA | TARDE | AUSENTE
	Estado          string `json:"estado"`
	HorarioAsignado string `json:"horario_asignado"`
	HoraProgramada  string `json:"hora_programada"`
	Observaciones   string `json:"observaciones"`
}

// ─── Response completo ────────────────────────────────────────

type ReporteDiarioResponse struct {
	FechaInicio    string              `json:"fecha_inicio"`
	FechaFin       string              `json:"fecha_fin"`
	TotalRegistros int                 `json:"total_registros"`
	Filas          []ReporteDiarioFila `json:"filas"`
}
