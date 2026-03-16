package reportes_inv

// ─── service ──────────────────────────────────────────────────

import (
	"errors"

	"gorm.io/gorm"
)

func FilterEntradaReporte(db *gorm.DB, idMovimiento int64, empresaID int64) (*EntradaReporte, error) {
	reporte, err := GetEntradaReporte(db, idMovimiento, empresaID)
	if err != nil {
		return nil, err
	}
	if reporte == nil {
		return nil, errors.New("entrada no encontrada")
	}
	return reporte, nil
}

func FilterTrasladoReporte(db *gorm.DB, idMovimiento int64, empresaID int64) (*TrasladoReporte, error) {
	reporte, err := GetTrasladoReporte(db, idMovimiento, empresaID)
	if err != nil {
		return nil, err
	}
	if reporte == nil {
		return nil, errors.New("traslado no encontrado")
	}
	return reporte, nil
}
