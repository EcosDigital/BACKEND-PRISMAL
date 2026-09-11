package ordenes

import "gorm.io/gorm"

func RegisterOrden(db *gorm.DB, req *OrdenRequest) (int64, error) {
	return CreateOrden(db, req)
}

func FilterOrdenes(db *gorm.DB) ([]OrdenListItem, error) {
	return ListOrdenes(db)
}

func FilterOrdenDetalle(db *gorm.DB, id int64) (*OrdenDetalle, error) {
	return GetOrdenDetalle(db, id)
}

func CambiarEstadoOrden(db *gorm.DB, id int64, req *CambiarEstadoOrdenRequest) error {
	return UpdateEstadoOrden(db, id, req.CodigoEstado, req.MotivoRechazo)
}
