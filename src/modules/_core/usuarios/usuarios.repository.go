package usuarios

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/ecosistema/core/src/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ListUserLast(db *gorm.DB) ([]UserResp, error) {
	var results []UserResp

	err := database.GormDB.
		Table("seguridad.cfg_usuarios u").
		Joins("INNER JOIN seguridad.cfg_roles_usuario r on u.id_rol = r.id").
		Joins("INNER JOIN configuracion.cfg_terceros t on t.id = u.id_tercero").
		Select(`
			u.id,
			u.email,
			r.nombre as rol,
			CASE
				WHEN t.id_tipo_persona = 1 THEN
					TRIM(CONCAT_WS(' ', t.primer_nombre, t.segundo_nombre, t.primer_apellido, t.segundo_apellido))
				WHEN t.id_tipo_persona = 2 THEN
					t.razon_social
				ELSE ''	
			END as tercero,
			u.activo as is_active`).
		Limit(10).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []UserResp{}
	}

	return results, nil

}

// valida existencia de usuario por ID
func ListUserByIDTercero(db *gorm.DB, id int) ([]UserResponse, error) {
	var results []UserResponse

	err := db.
		Table("seguridad.cfg_usuarios u").
		Joins("INNER JOIN seguridad.cfg_roles_usuario r on u.id_rol = r.id").
		Joins("INNER JOIN configuracion.cfg_terceros t on t.id = u.id_tercero").
		Select(`
			u.id,
			u.email,
			u.password,
			u.id_rol,
			r.nombre as rol,
			t.primer_nombre as nombre,
			u.activo`).
		Where("u.id_tercero = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []UserResponse{}
	}

	return results, nil

}

// crear registro
func CreateUser(db *gorm.DB, req *UserRequest) (int64, error) {
	// hashear password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("error al hashear contraseña")
	}

	data := map[string]interface{}{
		"email":         req.Email,
		"password":      hashedPassword,
		"id_rol":        req.IdRol,
		"id_tercero":    req.IdTercero,
		"image_profile": req.ImageProfile,
		"activo":        true,
		"created_at":    time.Now(),
		"created_by":    req.UserID,
	}

	tx := db.
		Table("seguridad.cfg_usuarios").
		Create(&data).Error

	if tx != nil {
		return 0, err
	}

	var id_reg int64
	database.GormDB.Raw("SELECT lastval()").Scan(&id_reg)

	return id_reg, nil

}

func ListUserById(db *gorm.DB, id int64) ([]UserResponseFull, error) {
	var results []UserResponseFull

	err := database.GormDB.
		Table("seguridad.cfg_usuarios u").
		Joins("INNER JOIN seguridad.cfg_roles_usuario r on u.id_rol = r.id").
		Joins("INNER JOIN configuracion.cfg_terceros t on t.id = u.id_tercero").
		Select(`
			u.id,
			u.id_tercero,
			u.email,
			u.id_rol,
			u.activo as is_active`).
		Where("u.id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []UserResponseFull{}
	}

	return results, nil

}

func ListUserByEmail(db *gorm.DB, email string) (*UserResponse, error) {

	var results *UserResponse

	err := db.
		Table("seguridad.cfg_usuarios u").
		Joins("INNER JOIN seguridad.cfg_roles_usuario r on u.id_rol = r.id").
		Joins("INNER JOIN configuracion.cfg_terceros t on t.id = u.id_tercero").
		Select(`
			u.id,
			u.email,
			u.password,
			u.id_rol,
			r.nombre as rol,
			t.primer_nombre as nombre,
			u.activo`).
		Where("u.email = ?", email).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

// update registro
func UpdateUserData(db *gorm.DB, id int64, req *UserUpdateRequest) (int64, error) {

	data := map[string]interface{}{
		"id_tercero":    req.IdTercero,
		"id_rol":        req.IdRol,
		"image_profile": req.ImageProfile,
		"activo":        req.Activo,
		"update_at":     time.Now(),
		"update_by":     req.UserID,
	}

	err := db.
		Table("seguridad.cfg_usuarios").
		Where("id = ?", id).
		Updates(data)

	if err.Error != nil {
		return 0, err.Error
	}

	if err.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil

}

// update passw
func UpdatePasswordUser(db *sql.DB, id int, req *NewPasswordUser) (int, error) {
	var jsonResult []byte
	var id_usuario int = 1

	query := `
		SELECT seguridad.qry_usuarios(
			operacion => $1,
			password_p => $2,
			id_usuario_p => $3
			id_registro_p => $4)`

	err := db.QueryRow(query, 4, id_usuario, id).Scan(&jsonResult)

	if err != nil {
		return 0, err
	}

	var result struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(jsonResult, &result); err != nil {
		return 0, err
	}

	return result.ID, nil

}
