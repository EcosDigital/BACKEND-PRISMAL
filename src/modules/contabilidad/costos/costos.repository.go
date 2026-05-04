package costos

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func ListCentroCostoByCode(db *gorm.DB, codigo string) ([]CostosResponse, error) {

	var results []CostosResponse

	err := db.
		Table("costos.cfg_centros_costo").
		Select(`
			id,
			codigo,
			nombre,
			id_area,
			'' as area,
			id_unidad_funcional,
			'' as unidad,
			is_active
		`).
		Where("codigo = ?", codigo).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func CreateCentroCosto(db *gorm.DB, req *CostosRequest) (int64, error) {

	data := map[string]interface{}{
		"codigo":              req.Codigo,
		"nombre":              req.Nombre,
		"id_area":             req.IDArea,
		"id_unidad_funcional": req.IDUnidad,
		"is_active":           req.IsActive,
		"created_at":          time.Now(),
		"created_by":          req.UserID,
		"id_empresa":          req.EmpresaID,
		"id_sede":             req.SedeID,
	}

	err := db.
		Table("costos.cfg_centros_costo").
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

func ListCentrosCostosLast(db *gorm.DB) ([]CostosResponse, error) {

	var results []CostosResponse

	err := db.
		Table("costos.cfg_centros_costo cc").
		Select(`
			cc.id as id,
			cc.codigo as codigo,
			cc.nombre as nombre,
			cc.id_area as id_area,
			ac.nombre as area,
			cc.id_unidad_funcional as id_unidad,
			uf.nombre as unidad,
			cc.is_active`).
		Joins("LEFT JOIN costos.cfg_area_costo ac ON ac.id = cc.id_area").
		Joins("LEFT JOIN costos.cfg_unidad_funcional uf ON uf.id = cc.id_unidad_funcional").
		Order("cc.id DESC").Limit(20).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []CostosResponse{}
	}

	return results, nil

}

func ListCentrosCostosByID(db *gorm.DB, id int64) ([]CostosResponse, error) {

	var results []CostosResponse

	err := db.
		Table("costos.cfg_centros_costo cc").
		Select(`
			cc.id as id,
			cc.codigo as codigo,
			cc.nombre as nombre,
			cc.id_area as id_area,
			ac.nombre as area,
			cc.id_unidad_funcional as id_unidad,
			cc.is_active`).
		Joins("LEFT JOIN costos.cfg_area_costo ac ON ac.id = cc.id_area").
		Joins("LEFT JOIN costos.cfg_unidad_funcional uf ON uf.id = cc.id_unidad_funcional").
		Where("cc.id = ?", id).
		Order("cc.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []CostosResponse{}
	}

	return results, nil

}

func UpdateCentroCosto(db *gorm.DB, req *CostosUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"nombre":              req.Nombre,
		"id_area":             req.IDArea,
		"id_unidad_funcional": req.IDUnidad,
		"is_active":           req.IsActive,
		"updated_at":          time.Now(),
		"updated_by":          req.UserID,
		"id_empresa":          req.EmpresaID,
		"id_sede":             req.SedeID,
	}

	err := db.
		Table("costos.cfg_centros_costo").
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
