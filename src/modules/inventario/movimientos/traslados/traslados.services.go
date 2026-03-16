package traslados

import (
	"errors"

	"gorm.io/gorm"
)

func CreateTraslado(db *gorm.DB, req *TrasladoRequest) (*TrasladoResponse, error) {
	if len(req.Detalle) == 0 {
		return nil, errors.New("el traslado debe tener al menos un artículo")
	}
	if req.IDBodegaOrigen == req.IDBodegaDestino {
		return nil, errors.New("la bodega origen y destino no pueden ser la misma")
	}
	return RegisterTraslado(db, req)
}

func FilterTraslados(db *gorm.DB, empresaID int64) ([]TrasladoListResponse, error) {
	return ListTraslados(db, empresaID)
}

func FilterTrasladoByID(db *gorm.DB, id int64) (*TrasladoFullResponse, error) {
	result, err := ListTrasladoByID(db, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("traslado no encontrado")
	}
	return result, nil
}

func ProcessAnulacionTraslado(db *gorm.DB, id int64, req *AnulacionTrasladoRequest) error {
	return AnularTraslado(db, id, req)
}

func FilterBodegasTraslado(db *gorm.DB, empresaID int64) ([]RefBodega, error) {
	return ListBodegas(db, empresaID)
}

func FilterLotes(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) ([]RefLote, error) {
	return ListLotesByArticuloBodega(db, idArticulo, idBodega, empresaID)
}

func FilterComprobantesTraslado(db *gorm.DB, empresaID int64) ([]RefComprobanteTraslado, error) {
	return ListComprobantesTraslado(db, empresaID)
}

func FilterStockArticulo(db *gorm.DB, idArticulo int64, idBodega int64, lote string, empresaID int64) (float64, error) {
	return GetStockArticuloBodegaLote(db, idArticulo, idBodega, lote, empresaID)
}
