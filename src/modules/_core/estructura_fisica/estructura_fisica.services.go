package estructura_fisica

import (
	"errors"

	"gorm.io/gorm"
)

// registrar mesa
func RegisterMesa(db *gorm.DB, req *MesaRequest) (int64, error) {
	return CreateMesa(db, req)
}

// consultar todos los registros de mesa
func FilterMesas(db *gorm.DB) ([]MesaResponse, error) {
	return ListMesas(db)
}

// consultar las mesas con estado is_active en true
func FilterMesasActivas(db *gorm.DB) ([]MesaResponse, error) {
	return ListMesasActivas(db)
}

// consultar un registro
func FilterMesaByID(db *gorm.DB, id int64) ([]MesaResponse, error) {
	return ListMesaByID(db, id)
}

// update a un registro mesa
func EditMesa(db *gorm.DB, id int64, req *MesaRequest) (int64, error) {

	//verificar existencia por ID
	exists, err := ListMesaByID(db, id)
	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		return 0, errors.New("no se encontró una mesa con ese ID")
	}

	//actualizar registro
	return UpdateMesa(db, id, req)

}

// cambiar solo el estado de una mesa
func EditEstadoMesa(db *gorm.DB, id int64, req *CambiarEstadoMesaRequest) error {

	exists, err := ListMesaByID(db, id)
	if err != nil {
		return err
	}

	if len(exists) <= 0 {
		return errors.New("no se encontró una mesa con ese ID")
	}

	return UpdateEstadoMesa(db, id, req.IDEstado)

}
