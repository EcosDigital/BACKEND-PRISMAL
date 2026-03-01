package costos

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterCentroCosto(db *gorm.DB, req *CostosRequest) (int64, error) {

	//verificar existencia por codigo
	exists, err := ListCentroCostoByCode(db, req.Codigo)
	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateCentroCosto(req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func FilterLastCentrosCosto() ([]CostosResponse, error) {
	return ListCentrosCostosLast()
}

func FilterCentroCostoByID(id int64) ([]CostosResponse, error) {
	return ListCentrosCostosByID(id)
}

func EditCentroCosto(req *CostosUpdateRequest, id int64) (int64, error) {

	//verificar registro por ID
	exists, err := ListCentrosCostosByID(id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := UpdateCentroCosto(req, id)
	if err != nil {
		return 0, err
	}

	return regID, err

}
