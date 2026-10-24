package configuracion

import (
	"errors"

	"gorm.io/gorm"
)

// defaultModelo es el modelo de Gemini que usa el agente cuando no se
// elige otro al configurarlo.
const defaultModelo = "gemini-3.7-flash"

func RegisterAgente(db *gorm.DB, req *AgenteRequest) (int64, error) {

	if req.Modelo == "" {
		req.Modelo = defaultModelo
	}

	if err := validarEncargado(db, req, 0); err != nil {
		return 0, err
	}

	return CreateAgente(db, req)

}

func FilterAgentes(db *gorm.DB) ([]AgenteResponse, error) {
	return ListAgentes(db)
}

func FilterAgenteByID(db *gorm.DB, id int64) ([]AgenteResponse, error) {
	return ListAgenteByID(db, id)
}

func FilterModulosInstalados(db *gorm.DB) ([]ModuloResponse, error) {
	return ListModulosInstalados(db)
}

func EditAgente(db *gorm.DB, id int64, req *AgenteRequest) (int64, error) {

	//verificar si existe el registro por ID
	exists, err := ListAgenteByID(db, id)
	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		return 0, errors.New("No se encontraron registros para actualizar")
	}

	if req.Modelo == "" {
		req.Modelo = defaultModelo
	}

	if err := validarEncargado(db, req, id); err != nil {
		return 0, err
	}

	return UpdateAgente(db, id, req)

}

// solo puede haber un agente encendido encargado de cada módulo, para que
// el asistente sepa siempre a quién consultar
func validarEncargado(db *gorm.DB, req *AgenteRequest, id int64) error {

	if req.IDModuloRef == nil || req.IsActive == nil || !*req.IsActive {
		return nil
	}

	ocupado, err := ExistsEncargadoActivo(db, *req.IDModuloRef, id)
	if err != nil {
		return err
	}

	if ocupado {
		return errors.New("Ya hay otro agente encendido encargado de este módulo")
	}

	return nil

}
