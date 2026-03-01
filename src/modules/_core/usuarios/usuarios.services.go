package usuarios

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ecosistema/core/src/modules/_core/terceros"
	"gorm.io/gorm"
)

// createUser valida existencia por tercero y crea registro
func RegisterUser(db *gorm.DB, req *UserRequest) (int, error) {
	//verificar si existe el tercero
	tercero, err := terceros.ListTerceroById(db, int64(req.IdTercero))
	if err != nil {
		return 0, err
	}
	if tercero == nil {
		return 0, errors.New("no existe el tercero")
	}

	//verificar si existe usuario asociado al tercero
	existencias, err := ListUserByIDTercero(db, req.IdTercero)
	if err != nil {
		return 0, err
	}

	if len(existencias) > 0 {
		return 0, errors.New("existe un usuario con este tercero")
	}

	//verificar existencia por email
	email, err := ListUserByEmail(db, req.Email)
	if err != nil {
		return 0, err
	}

	if email != nil {
		return 0, errors.New("existe un usuario con este email")
	}

	//insertar registro
	newID, err := CreateUser(db, req)
	if err != nil {
		return 0, err
	}

	return int(newID), nil
}

func FilterUserLast(db *gorm.DB) ([]UserResp, error) {
	return ListUserLast(db)
}

func FilterUserById(db *gorm.DB, id int64) ([]UserResponseFull, error) {
	return ListUserById(db, id)
}

func FIlterUserByEmail(db *gorm.DB, email string) (*UserResponse, error) {
	return ListUserByEmail(db, email)
}

func EditUser(db *gorm.DB, id int64, req *UserUpdateRequest) (int64, error) {
	//validar existencia
	user, err := ListUserById(db, id)
	if err != nil {
		return 0, err
	}
	if len(user) == 0 {
		return 0, fmt.Errorf("el usuario con ID %d no existe", id)
	}

	//verificar si existe el tercero
	tercero, err := terceros.ListTerceroById(db, int64(req.IdTercero))
	if err != nil {
		return 0, err
	}
	if tercero == nil {
		return 0, errors.New("no existe el tercero")
	}

	regID, err := UpdateUserData(db, id, req)
	if err != nil {
		return 0, err
	}

	//actualizar registro
	return regID, nil // UpdateUserData(db, id, req)
}

func ChangePasswordUser(ctx context.Context, db *sql.DB, id int, req *NewPasswordUser) (int, error) {
	/*validar existencia
	user, err := FindUserById(db, id)
	if err != nil {
		return 0, err
	}
	if len(user) == 0 {
		return 0, fmt.Errorf("el usuario con ID %d no existe", id)
	}

	//actualizar password
	return UpdatePasswordUser(db, id, req)*/
	return 0, nil
}
