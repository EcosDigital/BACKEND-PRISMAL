package configuracion

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// registrar agente
func CreateAgente(db *gorm.DB, req *AgenteRequest) (int64, error) {

	data := map[string]interface{}{
		"nombre":        req.Nombre,
		"instrucciones": req.Instrucciones,
		"modelo":        req.Modelo,
		"is_active":     req.IsActive,
		"id_modulo_ref": req.IDModuloRef,
		"created_at":    time.Now(),
		"created_by":    req.UserID,
		"id_empresa":    req.EmpresaID,
		"id_sede":       req.SedeID,
	}

	err := db.
		Table("agentes.cfg_agentes").
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

// consultar todos los agentes
func ListAgentes(db *gorm.DB) ([]AgenteResponse, error) {

	var results []AgenteResponse

	err := db.
		Table("agentes.cfg_agentes a").
		Select(`
			a.id,
			a.nombre,
			a.instrucciones,
			a.modelo,
			a.is_active,
			a.id_modulo_ref,
			COALESCE(m.nombre, '') as modulo`).
		Joins("LEFT JOIN configuracion.ref_modulos_tenant m ON m.id = a.id_modulo_ref").
		Order("a.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []AgenteResponse{}
	}

	return results, nil

}

// consultar agente por id
func ListAgenteByID(db *gorm.DB, id int64) ([]AgenteResponse, error) {

	var results []AgenteResponse

	err := db.
		Table("agentes.cfg_agentes a").
		Select(`
			a.id,
			a.nombre,
			a.instrucciones,
			a.modelo,
			a.is_active,
			a.id_modulo_ref,
			COALESCE(m.nombre, '') as modulo`).
		Joins("LEFT JOIN configuracion.ref_modulos_tenant m ON m.id = a.id_modulo_ref").
		Where("a.id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []AgenteResponse{}
	}

	return results, nil

}

// actualizar agente
func UpdateAgente(db *gorm.DB, id int64, req *AgenteRequest) (int64, error) {

	data := map[string]interface{}{
		"nombre":        req.Nombre,
		"instrucciones": req.Instrucciones,
		"modelo":        req.Modelo,
		"is_active":     req.IsActive,
		"id_modulo_ref": req.IDModuloRef,
		"update_at":     time.Now(),
		"update_by":     req.UserID,
	}

	err := db.Table("agentes.cfg_agentes").
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

// consultar los módulos instalados en el negocio
func ListModulosInstalados(db *gorm.DB) ([]ModuloResponse, error) {

	var results []ModuloResponse

	err := db.
		Table("configuracion.ref_modulos_tenant").
		Select("id, nombre").
		Order("nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ModuloResponse{}
	}

	return results, nil

}

// verificar si ya hay otro agente encendido encargado de ese módulo
func ExistsEncargadoActivo(db *gorm.DB, idModuloRef int64, excluirID int64) (bool, error) {

	var count int64

	err := db.
		Table("agentes.cfg_agentes").
		Where("id_modulo_ref = ? AND is_active = true AND id <> ?", idModuloRef, excluirID).
		Count(&count).Error

	return count > 0, err

}
