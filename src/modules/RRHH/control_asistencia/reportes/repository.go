package reportes_asistencia

import (
	"gorm.io/gorm"
)

// GetReporteDiario invoca fn_reporte_diario_asistencia y construye el response.
// La función PostgreSQL encapsula toda la lógica de cruce y cálculo de estado.
func GetReporteDiario(
	db *gorm.DB,
	req *ReporteDiarioRequest,
	empresaID int64,
) (*ReporteDiarioResponse, error) {

	type filaRaw struct {
		IDTercero       int64  `gorm:"column:id_tercero"`
		NombreEmpleado  string `gorm:"column:nombre_empleado"`
		NumeroDocumento string `gorm:"column:numero_documento"`
		HoraEntrada     string `gorm:"column:hora_entrada"`
		HoraSalida      string `gorm:"column:hora_salida"`
		Estado          string `gorm:"column:estado"`
		HorarioAsignado string `gorm:"column:horario_asignado"`
		HoraProgramada  string `gorm:"column:hora_programada"`
		Observaciones   string `gorm:"column:observaciones"`
	}

	var rawFilas []filaRaw

	err := db.
		Table("control_asistencia.fn_reporte_diario_asistencia(?, ?, ?)",
			req.FechaInicio,
			req.FechaFin,
			empresaID,
		).
		Scan(&rawFilas).Error
	if err != nil {
		return nil, err
	}

	filas := make([]ReporteDiarioFila, len(rawFilas))
	for i, r := range rawFilas {
		filas[i] = ReporteDiarioFila{
			IDTercero:       r.IDTercero,
			NombreEmpleado:  r.NombreEmpleado,
			NumeroDocumento: r.NumeroDocumento,
			HoraEntrada:     r.HoraEntrada,
			HoraSalida:      r.HoraSalida,
			Estado:          r.Estado,
			HorarioAsignado: r.HorarioAsignado,
			HoraProgramada:  r.HoraProgramada,
			Observaciones:   r.Observaciones,
		}
	}

	return &ReporteDiarioResponse{
		FechaInicio:    req.FechaInicio,
		FechaFin:       req.FechaFin,
		TotalRegistros: len(filas),
		Filas:          filas,
	}, nil
}
