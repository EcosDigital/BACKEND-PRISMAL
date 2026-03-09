package comprobantes

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func CreateComprobante(db *gorm.DB, req *ComprobanteRequest) (int64, error) {

	data := map[string]interface{}{
		"id_modulo":           req.IDModulo,
		"id_tipo_operacion":   req.IDTipoOperacion,
		"nombre_comprobante":  req.Nombre,
		"pregijo_comprobante": req.Prefijo,
		"consecutivo_inicial": req.ConsecutivoFinal,
		"consecutivo_final":   req.ConsecutivoFinal,
		"consecutivo_actual":  0,
		"consecutivo_fin":     req.ConsecutivoFinal,
		"permite_anulacion":   req.PermiteAnular,
		"fecha_inicio":        req.FechaInicio,
		"fecha_final":         req.FechaFinal,
		"token":               req.Token,
		"resolucion":          req.Resolucion,
		"is_active":           req.IsActive,
		"created_at":          time.Now(),
		"created_by":          req.UserID,
		"id_empresa":          req.EmpresaID,
		"id_sede":             req.SedeID,
	}

	err := db.
		Table("comprobantes.cfg_comprobante").
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

func ListComprobantesLast(db *gorm.DB) ([]ComprobanteResponse, error) {

	var results []ComprobanteResponse

	err := db.
		Table("comprobantes.cfg_comprobante c").
		Select(`
			c.id,
			C.id_modulo,
			m.nombre as modulo,
			o.nombre as operacion,
			c.nombre_comprobante || ' (' || c.prefijo_comprobante || ')' as nombre,
			c.consecutivo_actual as consecutivo,
			c.is_active `).
		Joins("INNER JOIN comprobantes.ref_tipo_operacion o ON o.id = c.id_tipo_operacion").
		Joins("INNER JOIN configuracion.ref_modulos_tenant m ON m.id_ref = c.id_modulo").
		Order("c.id ASC").Limit(20).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ComprobanteResponse{}
	}

	return results, nil

}

func ListComprobanteById(db *gorm.DB, id int64) ([]ComprobanteResponse, error) {
	var results []ComprobanteResponse

	err := db.
		Table("comprobantes.cfg_comprobante c").
		Select(`
			c.id as id,
			c.id_modulo,
			m.nombre as modulo,
			c.id_tipo_operacion,
			o.nombre as operacion,
			c.nombre_comprobante,
			c.prefijo_comprobante,
			c.consecutivo_inicial,
			c.consecutivo_actual,
			c.permite_anulacion,
			c.fecha_inicio,
			c.fecha_final,
			c.token,
			c.resolucion,
			c.is_active`).
		Joins("INNER JOIN configuracion.cfg_modulos m ON m.id = c.id_modulo").
		Joins("INNER JOIN comprobantes.ref_tipo_operacion o ON o.id = c.id_tipo_operacion").
		Where("c.id = ?", id).
		Order("c.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ComprobanteResponse{}
	}

	return results, nil
}

func UpdateComprobante(db *gorm.DB, req *ComprobanteUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"nombre_comprobante":  req.Nombre,
		"pregijo_comprobante": req.Prefijo,
		"consecutivo_inicial": req.ConsecutivoFinal,
		"consecutivo_final":   req.ConsecutivoFinal,
		"consecutivo_actual":  0,
		"consecutivo_fin":     req.ConsecutivoFinal,
		"permite_anulacion":   req.PermiteAnular,
		"fecha_inicio":        req.FechaInicio,
		"fecha_final":         req.FechaFinal,
		"token":               req.Token,
		"resolucion":          req.Resolucion,
		"is_active":           req.IsActive,
		"update_at":           time.Now(),
		"update_by":           req.UserID,
		"id_empresa":          req.EmpresaID,
		"id_sede":             req.SedeID,
	}

	err := db.
		Table("comprobantes.cfg_comprobante").
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
