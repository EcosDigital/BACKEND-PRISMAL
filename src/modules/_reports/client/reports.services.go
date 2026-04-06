package reports

import (
	"errors"

	"gorm.io/gorm"
)

func GenerateReport(
	db *gorm.DB,
	reportID int64,
	filters map[string]interface{},
) ([]map[string]interface{}, error) {

	// validar que el reporte exista
	reporte, err := ListReportByID(db, int(reportID))
	if err != nil {
		return nil, err
	}

	if !reporte.IsActive {
		return nil, errors.New("el informe no está activo")
	}

	// ejecutar funcion sql
	result, err := ExecuteReport(
		db,
		reporte.FuncionSQL,
		filters
	)

	if err != nil {
		return nil, err
	}

	return result, nil

}
