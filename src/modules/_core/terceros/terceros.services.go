package terceros

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	PersonaNatural  = 1
	PersonaJuridica = 2
)

// CreateTercero valida registros por documento y crea un nuevo tercero
func RegisterTercero(db *gorm.DB, req *TerceroRequest) (int64, error) {
	// Verifica si existe registro
	existentes, err := FindTerceroByDocument(db, req.NumeroDocumento)
	if err != nil {
		return 0, err
	}

	if len(existentes) > 0 {
		// Existe registro con número de documento

	}

	// Validaciones según tipo de persona
	switch req.IdTipoPersona {
	case PersonaNatural:
		if req.PrimerNombre == "" || req.PrimerApellido == "" {
			return 0, errors.New("para persona natural se requieren nombre y apellido")
		}
	case PersonaJuridica:
		if req.RazonSocial == "" || req.Dv == "" {
			return 0, errors.New("para persona jurídica se requiere razón social y NIT con dígito de verificación")
		}
	default:
		return 0, errors.New("tipo de persona no válido")
	}

	// Insertar registro
	newID, err := CreateTercero(db, req)
	if err != nil {
		return 0, err
	}

	//insertar clase de tercero (si existen)
	if len(req.ClaseTercero) > 0 {
		if err := AddClasesTercero(db, newID, req.ClaseTercero, req.UserID); err != nil {
			return 0, err
		}
	} //

	return newID, nil
}

func ListRecenttTerceros(db *gorm.DB) ([]TerceroResponse, error) {
	return FindLastTercero(db)
}

func FilterTerceroById(db *gorm.DB, id int64) (*TerceroResponseFull, error) {
	return ListTerceroById(db, id)
}

func EditTercero(db *gorm.DB, id int64, id_user int64, req *TerceroUpdateRequest) (int64, error) {

	results, err := ListTerceroById(db, id)
	if err != nil {
		return 0, err
	}
	if results == nil {
		return 0, fmt.Errorf("error, no existe registro")
	}

	//actualizar registro
	regID, err := UpdateTerceroData(db, id, req)
	if err != nil {
		return 0, err
	}

	//gestionar clase de tercero (si enviaron)
	if req.ClaseTercero != nil {
		//eliminar clases
		if err := RemoveClaseTerceros(db, id); err != nil {
			return 0, fmt.Errorf("error al eliminar clases anteriores: %w", err)
		}

		//insertar nuevas clases (si hay)
		if len(req.ClaseTercero) > 0 {
			if err := AddClasesTercero(db, id, req.ClaseTercero, id_user); err != nil {
				return 0, fmt.Errorf("Error al insertar nuevas clases: %w", err)
			}
		}

	}
	return regID, nil
}

// actualizar identificacion
func ChangeTerceroIdentity(db *gorm.DB, id int64, req *TerceroUpdateIndenty) (int, error) {
	//validar existencia
	tercero, err := FindClasesByTerceroId(db, id)
	if err != nil {
		return 0, err
	}
	if len(tercero) == 0 {
		return 0, fmt.Errorf("el tercero con ID %d no existe", id)
	}

	//actualizar identificacion
	return 0, nil
	//return updateTerceroIdentity(db, id, req)
}
