package empresa

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterEmpresa(db *gorm.DB, req *EmpresaRequest) (int64, error) {

	//verificar existencia
	exists, err := ListEmpresaByNit(db, req.Nit)

	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateEmpresa(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func FilterLastEmpresa(db *gorm.DB) ([]EmpresaResponse, error) {
	return ListEmpresaLast(db)
}

func FilterEmpresaByID(db *gorm.DB, id int64) ([]EmpresaResponseFull, error) {
	return ListEmpresaByID(db, id)
}

func EditEmpresa(db *gorm.DB, id int64, req *EmpresaUpdateRequest) (int64, error) {

	//verificar registro por ID
	exist, err := ListEmpresaByID(db, id)
	if err != nil {
		return 0, err
	}

	if len(exist) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No Se encontraron resultados de este registro")
	}

	//actualizar registro
	regID, err := UpdateEmpresa(db, id, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

// SEDES
func FilterLsatSede(db *gorm.DB) ([]SedeResponse, error) {
	return ListSedeLast(db)
}

func FilterSedeByID(db *gorm.DB, id int64) ([]SedeResponseFull, error) {
	return ListSedeByID(db, id)
}

func RegisterSede(db *gorm.DB, req *SedeRequest) (int64, error) {

	//verificar existencia
	exists, err := ListSedeByCode(db, req.Codigo)

	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateSede(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func EditSede(db *gorm.DB, id int64, req *SedeUpdateRequest) (int64, error) {

	//verificar registro por ID
	exist, err := ListSedeByID(db, id)
	if err != nil {
		return 0, err
	}

	if len(exist) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No Se encontraron resultados de este registro")
	}

	//actualizar registro
	regID, err := UpdateSede(db, id, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}
