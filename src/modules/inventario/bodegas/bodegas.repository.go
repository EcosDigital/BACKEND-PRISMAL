package bodega

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func CreateBodega(db *gorm.DB, req *BodegaRequest) (int64, error) {

	data := map[string]interface{}{
		"id_tipo_bodega":      req.IDTipoBodega,
		"codigo":              req.Codigo,
		"nombre":              req.Nombre,
		"descripcion":         req.Descripcion,
		"is_central":          req.IsCentral,
		"permite_ventas":      req.PermiteVentas,
		"stock_minimo_global": req.StockMinimoGlobal,
		"stock_maximo_global": req.StockMaximoGlobal,
		"alerta_stock":        req.AlertaStock,
		"aplica_mov_contable": req.AplicaMovContable,
		"is_active":           req.IsActive,
		"created_at":          time.Now(),
		"created_by":          req.UserID,
		"id_empresa":          req.EmpresaID,
		"id_sede":             req.IDSede,
	}

	tx := db.Table("inventario.cfg_bodegas").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil

}

func ListBodegas(db *gorm.DB, empresaID int64) ([]BodegaResponse, error) {

	var results []BodegaResponse

	err := db.
		Table("inventario.cfg_bodegas b").
		Select(`
			b.id,
			b.codigo,
			b.nombre,
			t.nombre AS tipo_bodega,
			b.is_central,
			b.permite_ventas,
			b.is_active`).
		Joins("INNER JOIN inventario.ref_tipo_bodega t ON t.id = b.id_tipo_bodega").
		Where("b.id_empresa = ?", empresaID).
		Order("b.id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []BodegaResponse{}
	}

	return results, nil
}

func ListBodegaByID(db *gorm.DB, id int64) (*BodegaResponseFull, error) {

	var result BodegaResponseFull

	err := db.
		Table("inventario.cfg_bodegas b").
		Select(`
			b.id,
			b.id_tipo_bodega,
			t.nombre              AS tipo_bodega,
			b.codigo,
			b.nombre,
			b.descripcion,
			b.is_central,
			b.permite_ventas,
			b.stock_minimo_global,
			b.stock_maximo_global,
			b.alerta_stock,
			b.aplica_mov_contable,
			b.is_active`).
		Joins("INNER JOIN inventario.ref_tipo_bodega t ON t.id = b.id_tipo_bodega").
		Where("b.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateBodega(db *gorm.DB, req *BodegaUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"id_tipo_bodega":      req.IDTipoBodega,
		"nombre":              req.Nombre,
		"descripcion":         req.Descripcion,
		"is_central":          req.IsCentral,
		"permite_ventas":      req.PermiteVentas,
		"stock_minimo_global": req.StockMinimoGlobal,
		"stock_maximo_global": req.StockMaximoGlobal,
		"alerta_stock":        req.AlertaStock,
		"aplica_mov_contable": req.AplicaMovContable,
		"is_active":           req.IsActive,
		"updated_at":          time.Now(),
		"updated_by":          req.UserID,
		"id_empresa":          req.EmpresaID,
		"id_sede":             req.SedeID,
	}

	tx := db.
		Table("inventario.cfg_bodegas").
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
