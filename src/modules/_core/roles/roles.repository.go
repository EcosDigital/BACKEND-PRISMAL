package roles

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	empresa "github.com/ecosistema/core/src/modules/_core/gestion_empresa"
	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)

func CreateRolUsuario(db *gorm.DB, req *RolesRequest) (int64, error) {

	data := map[string]interface{}{
		"nombre":          req.Nombre,
		"id_tipo":         req.IDTipo,
		"intentos_login":  req.IntentosLogin,
		"is_multisession": req.IsMultisession,
		"is_externo":      req.IsExterno,
		"created_at":      time.Now(),
		"created_by":      req.UserID,
		"id_empresa":      req.EmpresaID,
		"id_sede":         req.SedeID,
	}

	err := db.
		Table("seguridad.cfg_roles_usuario").
		Create(&data)

	if err.Error != nil {
		return 0, err.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil
	}

	return id, nil

}

func ListRolUsuarioLast(db *gorm.DB) ([]RolesResponse, error) {

	var results []RolesResponse

	err := db.
		Table("seguridad.cfg_roles_usuario r").
		Select(`
			r.id,
			r.nombre,
			r.id_tipo,
			tr.nombre as tipo,
			r.intentos_login,
			r.is_multisession,
			r.is_externo`).
		Joins("INNER JOIN seguridad.ref_tipo_rol tr ON tr.id = r.id_tipo").
		Order("r.id ASC").Limit(20).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []RolesResponse{}
	}

	return results, nil

}

func ListRolUsuarioByID(db *gorm.DB, id int64) ([]RolesResponse, error) {

	var results []RolesResponse

	err := db.
		Table("seguridad.cfg_roles_usuario r").
		Select(`
			r.id,
			r.nombre,
			r.id_tipo,
			tr.nombre as tipo,
			r.intentos_login,
			r.is_multisession,
			r.is_externo`).
		Joins("INNER JOIN seguridad.ref_tipo_rol tr ON tr.id = r.id_tipo").
		Where("r.id = ?", id).
		Order("r.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []RolesResponse{}
	}

	return results, nil

}

func UpdateRolUsuario(db *gorm.DB, id int64, req *RolesRequest) (int64, error) {

	data := map[string]interface{}{
		"nombre":          req.Nombre,
		"id_tipo":         req.IDTipo,
		"intentos_login":  req.IntentosLogin,
		"is_multisession": req.IsMultisession,
		"is_externo":      req.IsExterno,
		"created_by":      time.Now(),
		"updated_by":      req.UserID,
		"id_empresa":      req.EmpresaID,
		"id_sede":         req.SedeID,
	}

	err := db.Table("seguridad.cfg_roles_usuario").
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

//** CONFIG ROLES * //

func AddSedeRol(db *gorm.DB, idRol int, idSede int, idUser int64, jsonModules []Item) error {

	sedes, err := empresa.ListSedeByID(db, int64(idSede))
	if err != nil {
		return err
	}

	if len(sedes) == 0 {
		return fmt.Errorf("no se encontró la sede")
	}

	sede := sedes[0]

	//serializar el json
	jsonBytes, err := json.Marshal(jsonModules)
	if err != nil {
		logging.Error.Printf("error al serializar json_modules: %w", err)
	}

	data := map[string]interface{}{
		"id_empresa":   sede.IdEmpresa,
		"id_sede":      idSede,
		"id_rol":       idRol,
		"json_modules": string(jsonBytes),
		"created_by":   idUser,
		"created_at":   time.Now(),
	}

	tx := db.Table("seguridad.cfg_sedes_roles").
		Create(&data).Error

	return tx

}

func AddSedesRolUser(db *gorm.DB, idRol int, sedes []int, idUser int64, jsonModules []Item) error {
	for _, IdSede := range sedes {
		err := AddSedeRol(db, idRol, IdSede, idUser, jsonModules)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveSedesRol(db *gorm.DB, idRol int64) error {

	err := db.
		Table("seguridad.cfg_sedes_usuario").
		Where("id_rol", idRol).
		Delete(nil).Error

	if err != nil {
		return err
	}

	return nil

}

func GetConfigRolByRolID(db *gorm.DB, idRol int64) (*ConfigRolResponse, error) {

	type rawRow struct {
		ID          int64  `gorm:"column:id"`
		IDRol       int64  `gorm:"column:id_rol"`
		IDEmpresa   int64  `gorm:"column:id_empresa"`
		IDSede      int64  `gorm:"column:id_sede"`
		JsonModules string `gorm:"column:json_modules"`
	}

	var row rawRow

	err := db.
		Table("seguridad.cfg_sedes_roles").
		Select("id, id_rol, id_empresa, id_sede, json_modules").
		Where("id_rol = ?", idRol).
		Order("id DESC").
		First(&row).Error

	// Sin registro → devolvemos exists=false sin error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &ConfigRolResponse{Exists: false}, nil
		}
		return nil, err
	}

	// Deserializar el JSON de módulos
	var items []Item
	if row.JsonModules != "" {
		if jsonErr := json.Unmarshal([]byte(row.JsonModules), &items); jsonErr != nil {
			logging.Error.Printf("error al deserializar json_modules: %v", jsonErr)
			items = []Item{}
		}
	}

	return &ConfigRolResponse{
		ID:          row.ID,
		IDRol:       row.IDRol,
		IDEmpresa:   row.IDEmpresa,
		IDSede:      row.IDSede,
		JsonModules: items,
		Exists:      true,
	}, nil

}

func UpdateConfigRol(db *gorm.DB, id int64, req *UpdateConfigRolRequest) (int64, error) {

	jsonBytes, err := json.Marshal(req.Json)
	if err != nil {
		return 0, fmt.Errorf("error al serializar json_modules: %w", err)
	}

	data := map[string]interface{}{
		"json_modules": string(jsonBytes),
		"updated_by":   req.UserID,
		"updated_at":   time.Now(),
	}

	result := db.Table("seguridad.cfg_sedes_roles").
		Where("id = ?", id).
		Updates(data)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil
}
