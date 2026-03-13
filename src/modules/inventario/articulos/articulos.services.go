package articulos

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterArticulo(db *gorm.DB, req *ArticuloRequest) (int64, error) {
	return CreateArticulo(db, req)
}

func FilterArticulos(db *gorm.DB, empresaID int64) ([]ArticuloResponse, error) {
	return ListArticulos(db, empresaID)
}

func FilterArticuloByID(db *gorm.DB, id int64) (*ArticuloResponseFull, error) {
	return ListArticuloByID(db, id)
}

func EditArticulo(db *gorm.DB, id int64, req *ArticuloUpdateRequest) (int64, error) {

	exists, err := ListArticuloByID(db, id)
	if err != nil {
		return 0, err
	}
	if exists == nil || exists.ID == 0 {
		return 0, errors.New("no se encontró artículo con este ID")
	}

	return UpdateArticulo(db, req, id)
}
