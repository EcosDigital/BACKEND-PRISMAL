package comprobantes

import "errors"

func RegisterComprobante(req *ComprobanteRequest) (int64, error) {
	regID, err := CreateComprobante(req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func FilterLastComprobantes() ([]ComprobanteResponse, error) {
	return ListComprobantesLast()
}

func FilterComprobanteByID(id int64) ([]ComprobanteResponse, error) {
	return ListComprobanteById(id)
}

func EditComprobante(id int64, req *ComprobanteUpdateRequest) (int64, error) {

	//verificar registro por ID
	exists, err := ListComprobanteById(id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := EditComprobante(id, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}
