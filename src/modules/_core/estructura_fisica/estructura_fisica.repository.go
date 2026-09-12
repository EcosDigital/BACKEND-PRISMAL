package estructura_fisica

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// registrar mesa
func CreateMesa(db *gorm.DB, req *MesaRequest) (int64, error) {

	data := map[string]interface{}{
		"codigo":     req.Codigo,
		"nombre":     req.Nombre,
		"id_estado":  req.IDEstado,
		"is_active":  req.IsActive,
		"created_at": time.Now(),
		"created_by": req.UserID,
		"id_empresa": req.EmpresaID,
		"id_sede":    req.SedeID,
	}

	err := db.
		Table("configuracion.cfg_mesas").
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

// consultar todos los registros de mesa
func ListMesas(db *gorm.DB) ([]MesaResponse, error) {

	var results []MesaResponse

	err := db.
		Table("configuracion.cfg_mesas m").
		Select(`
			m.id,
			m.codigo,
			m.nombre,
			m.id_estado,
			e.nombre as estado,
			m.is_active`).
		Joins("INNER JOIN configuracion.ref_estado_mesa e ON e.id = m.id_estado").
		Order("m.codigo ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []MesaResponse{}
	}

	return results, nil

}

// consultar las mesas con estado is_active en true
func ListMesasActivas(db *gorm.DB) ([]MesaResponse, error) {

	var results []MesaResponse

	err := db.
		Table("configuracion.cfg_mesas m").
		Select(`
			m.id,
			m.codigo,
			m.nombre,
			m.id_estado,
			e.nombre as estado,
			m.is_active`).
		Joins("INNER JOIN configuracion.ref_estado_mesa e ON e.id = m.id_estado").
		Where("m.is_active = true").
		Order("m.codigo ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []MesaResponse{}
	}

	return results, nil

}

// consultar un registro
func ListMesaByID(db *gorm.DB, id int64) ([]MesaResponse, error) {

	var results []MesaResponse

	err := db.
		Table("configuracion.cfg_mesas m").
		Select(`
			m.id,
			m.codigo,
			m.nombre,
			m.id_estado,
			e.nombre as estado,
			m.is_active`).
		Joins("INNER JOIN configuracion.ref_estado_mesa e ON e.id = m.id_estado").
		Where("m.id = ?", id).
		Order("m.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []MesaResponse{}
	}

	return results, nil

}

// update a un registro mesa
func UpdateMesa(db *gorm.DB, id int64, req *MesaRequest) (int64, error) {

	data := map[string]interface{}{
		"codigo":    req.Codigo,
		"nombre":    req.Nombre,
		"id_estado": req.IDEstado,
		"is_active": req.IsActive,
		"update_at": time.Now(),
		"update_by": req.UserID,
	}

	err := db.Table("configuracion.cfg_mesas").
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

// cambiar solo el estado de una mesa (por id_estado, ya resuelto por el
// frontend contra configuracion.ref_estado_mesa).
func UpdateEstadoMesa(db *gorm.DB, id int64, idEstado int) error {

	result := db.Table("configuracion.cfg_mesas").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"id_estado": idEstado,
			"update_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}
