package kardex

import (
	"errors"

	"gorm.io/gorm"
)

func FilterKardex(db *gorm.DB, req *KardexRequest, empresaID int64) (*KardexReporte, error) {
	if req.FechaInicio > req.FechaFin {
		return nil, errors.New("la fecha de inicio no puede ser mayor a la fecha fin")
	}

	reporte, err := GetKardex(db, req, empresaID)
	if err != nil {
		return nil, err
	}

	return reporte, nil
}

func FilterLotes(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) ([]KardexRefLote, error) {
	return ListLotesByArticuloBodega(db, idArticulo, idBodega, empresaID)
}
