package historial_asistencia

import (
	"errors"

	"gorm.io/gorm"
)

// FilterHistorialAsistencia valida el request y delega al repository.
func FilterHistorialAsistencia(
	db *gorm.DB,
	req *HistorialAsistenciaRequest,
	empresaID int64,
) (*HistorialAsistenciaReporte, error) {

	if req.FechaInicio > req.FechaFin {
		return nil, errors.New("la fecha de inicio no puede ser mayor a la fecha fin")
	}

	return GetHistorialAsistencia(db, req, empresaID)
}
