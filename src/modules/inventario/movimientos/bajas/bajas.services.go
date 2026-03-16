package bajas

import (
	"errors"

	"gorm.io/gorm"
)

func CreateBaja(db *gorm.DB, req *BajaRequest) (*BajaResponse, error) {
	if len(req.Detalle) == 0 {
		return nil, errors.New("la baja debe tener al menos un artículo")
	}
	return RegisterBaja(db, req)
}

func FilterBajas(db *gorm.DB, empresaID int64) ([]BajaListResponse, error) {
	return ListBajas(db, empresaID)
}

func FilterBajaByID(db *gorm.DB, id int64) (*BajaFullResponse, error) {
	result, err := ListBajaByID(db, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("baja no encontrada")
	}
	return result, nil
}

func ProcessAnulacionBaja(db *gorm.DB, id int64, req *AnulacionBajaRequest) error {
	return AnularBaja(db, id, req)
}

func FilterSubtiposMovimientoBaja(db *gorm.DB) ([]RefSubtipoMovimiento, error) {
	return ListSubtiposMovimientoBaja(db)
}

func FilterComprobantesBaja(db *gorm.DB, empresaID int64) ([]RefComprobanteBaja, error) {
	return ListComprobantesBaja(db, empresaID)
}

func FilterLotesBaja(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) ([]RefLoteBaja, error) {
	return ListLotesBaja(db, idArticulo, idBodega, empresaID)
}

func FilterBajaReporte(db *gorm.DB, idMovimiento int64, empresaID int64) (*BajaReporte, error) {
	reporte, err := GetBajaReporte(db, idMovimiento, empresaID)
	if err != nil {
		return nil, err
	}
	if reporte == nil {
		return nil, errors.New("baja no encontrada")
	}
	return reporte, nil
}
