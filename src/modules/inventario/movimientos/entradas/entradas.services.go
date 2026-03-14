package entradas

import (
	"errors"

	"gorm.io/gorm"
)

func CreateEntrada(db *gorm.DB, req *EntradaRequest) (*EntradaResponse, error) {

	if len(req.Detalle) == 0 {
		return nil, errors.New("la entrada debe tener al menos un artículo en el detalle")
	}

	return RegisterEntrada(db, req)
}

func FilterEntradas(db *gorm.DB, empresaID int64) ([]EntradaListResponse, error) {
	return ListEntradas(db, empresaID)
}

func FilterEntradaByID(db *gorm.DB, id int64) (*EntradaFullResponse, error) {
	result, err := ListEntradaByID(db, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("entrada no encontrada")
	}
	return result, nil
}

func FilterComprobantesEntrada(db *gorm.DB, empresaID int64) ([]RefComprobante, error) {
	return ListComprobantesEntrada(db, empresaID)
}

func FilterImpuestos(db *gorm.DB, empresaID int64) ([]RefImpuesto, error) {
	return ListImpuestos(db, empresaID)
}

func FilterExistenciaArticulo(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) (*ArticuloExistencia, error) {
	return GetExistenciaArticulo(db, idArticulo, idBodega, empresaID)
}
