package bodega

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterBodega(db *gorm.DB, req *BodegaRequest) (int64, error) {
	regID, err := CreateBodega(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func FilterBodegas(db *gorm.DB, empresaID int64) ([]BodegaResponse, error) {
	return ListBodegas(db, empresaID)
}

func FilterBodegaByID(db *gorm.DB, id int64) (*BodegaResponseFull, error) {
	return ListBodegaByID(db, id)
}

func EditBodega(db *gorm.DB, id int64, req *BodegaUpdateRequest) (int64, error) {

	// Verificar que el registro existe
	exists, err := ListBodegaByID(db, id)
	if err != nil {
		return 0, err
	}
	if exists == nil || exists.ID == 0 {
		return 0, errors.New("no se encontró bodega con este ID")
	}

	return UpdateBodega(db, req, id)
}
