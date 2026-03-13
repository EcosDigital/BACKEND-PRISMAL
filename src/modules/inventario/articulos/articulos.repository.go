package articulos

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func CreateArticulo(db *gorm.DB, req *ArticuloRequest) (int64, error) {

	// factor_unidad_base por defecto es 1 si no se especifica
	factor := req.FactorUnidadBase
	if factor == 0 {
		factor = 1
	}

	data := map[string]interface{}{
		"codigo":                req.Codigo,
		"nombre":                req.Nombre,
		"descripcion":           req.Descripcion,
		"id_tipo_articulo":      req.IDTipoArticulo,
		"id_unidad_medida":      req.IDUnidadMedida,
		"factor_unidad_base":    factor,
		"id_forma_farmaceutica": req.IDFormaFarmaceutica,
		"id_presentacion":       req.IDPresentacion,
		"is_active":             req.IsActive,
		"created_at":            time.Now(),
		"created_by":            req.UserID,
		"id_empresa":            req.EmpresaID,
		"id_sede":               req.SedeID,
	}

	tx := db.Table("inventario.cfg_articulos").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

func ListArticulos(db *gorm.DB, empresaID int64) ([]ArticuloResponse, error) {

	var results []ArticuloResponse

	err := db.
		Table("inventario.cfg_articulos a").
		Select(`
			a.id,
			a.codigo,
			a.nombre,
			g.nombre  AS grupo_articulo,
			u.nombre  AS unidad_medida,
			a.factor_unidad_base,
			a.is_active`).
		Joins("INNER JOIN inventario.cfg_grupo_articulos g ON g.id = a.id_tipo_articulo").
		Joins("INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida").
		Where("a.id_empresa = ?", empresaID).
		Order("a.id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ArticuloResponse{}
	}

	return results, nil
}

func ListArticuloByID(db *gorm.DB, id int64) (*ArticuloResponseFull, error) {

	var result ArticuloResponseFull

	err := db.
		Table("inventario.cfg_articulos a").
		Select(`
			a.id,
			a.codigo,
			a.nombre,
			a.descripcion,
			a.id_tipo_articulo,
			g.nombre  AS grupo_articulo,
			a.id_unidad_medida,
			u.nombre  AS unidad_medida,
			a.factor_unidad_base,
			a.id_forma_farmaceutica,
			a.id_presentacion,
			a.is_active`).
		Joins("INNER JOIN inventario.cfg_grupo_articulos g ON g.id = a.id_tipo_articulo").
		Joins("INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida").
		Where("a.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateArticulo(db *gorm.DB, req *ArticuloUpdateRequest, id int64) (int64, error) {

	factor := req.FactorUnidadBase
	if factor == 0 {
		factor = 1
	}

	data := map[string]interface{}{
		"nombre":                req.Nombre,
		"descripcion":           req.Descripcion,
		"id_tipo_articulo":      req.IDTipoArticulo,
		"id_unidad_medida":      req.IDUnidadMedida,
		"factor_unidad_base":    factor,
		"id_forma_farmaceutica": req.IDFormaFarmaceutica,
		"id_presentacion":       req.IDPresentacion,
		"is_active":             req.IsActive,
		"updated_at":            time.Now(),
		"updated_by":            req.UserID,
	}

	tx := db.
		Table("inventario.cfg_articulos").
		Where("id = ?", id).
		Updates(data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil
}
