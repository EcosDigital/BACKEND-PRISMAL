package reportes_asistencia

import (
	"errors"

	"gorm.io/gorm"
)

// FilterReporteDiario valida el request y delega al repository.
func FilterReporteDiario(
	db *gorm.DB,
	req *ReporteDiarioRequest,
	empresaID int64,
) (*ReporteDiarioResponse, error) {

	if req.FechaInicio == "" || req.FechaFin == "" {
		return nil, errors.New("fecha_inicio y fecha_fin son requeridos")
	}
	if req.FechaInicio > req.FechaFin {
		return nil, errors.New("la fecha de inicio no puede ser mayor a la fecha fin")
	}

	return GetReporteDiario(db, req, empresaID)
}
