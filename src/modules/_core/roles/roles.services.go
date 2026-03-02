package roles

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterRolUsuario(db *gorm.DB, req *RolesRequest) (int64, error) {
	//insertar registro
	regID, err := CreateRolUsuario(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func FilterLastRolUser(db *gorm.DB) ([]RolesResponse, error) {
	return ListRolUsuarioLast(db)
}

func FilterRolUserByID(db *gorm.DB, id int64) ([]RolesResponse, error) {
	return ListRolUsuarioByID(db, id)
}

func EditRolUser(db *gorm.DB, id int64, req *RolesRequest) (int64, error) {

	//varificar existencias por ID
	exists, err := ListRolUsuarioByID(db, id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := UpdateRolUsuario(db, id, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

//* CONFIG  ROLES * //

func RegisterConfigRoles(db *gorm.DB, req *ConfigRolRequest) (int64, error) {

	//verificar existencia del rol
	exists, err := ListRolUsuarioByID(db, int64(req.IDRol))
	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		return 0, errors.New("no existe registro asociado")
	}

	//insertar sedes
	if len(req.Sedes) > 0 {
		if err := AddSedesRolUser(db, req.IDRol, req.Sedes, req.UserID, req.Json); err != nil {
			return 0, err
		}
	}

	return 0, nil

}
