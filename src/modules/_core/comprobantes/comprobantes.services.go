package comprobantes

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterComprobante(db *gorm.DB, req *ComprobanteRequest) (int64, error) {
	regID, err := CreateComprobante(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func FilterLastComprobantes(db *gorm.DB) ([]ComprobanteResponse, error) {
	return ListComprobantesLast(db)
}

func FilterComprobanteByID(db *gorm.DB, id int64) ([]ComprobanteResponseFull, error) {
	return ListComprobanteById(db, id)
}

func EditComprobante(db *gorm.DB, id int64, req *ComprobanteUpdateRequest) (int64, error) {

	//verificar registro por ID
	exists, err := ListComprobanteById(db, id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := EditComprobante(db, id, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}
