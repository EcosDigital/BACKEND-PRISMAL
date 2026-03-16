package traslados

import (
	"database/sql"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// ─── Registrar traslado ───────────────────────────────────────

func RegisterTraslado(db *gorm.DB, req *TrasladoRequest) (*TrasladoResponse, error) {

	detalleSQL := "ARRAY["
	for i, d := range req.Detalle {
		if i > 0 {
			detalleSQL += ","
		}
		lote := "NULL"
		if d.Lote != "" {
			lote = fmt.Sprintf("'%s'", escapePG(d.Lote))
		}
		obs := "NULL"
		if d.Observacion != "" {
			obs = fmt.Sprintf("'%s'", escapePG(d.Observacion))
		}
		detalleSQL += fmt.Sprintf(
			"ROW(%d,%s,%s,%s)::inventario.t_traslado_detalle",
			d.IDArticulo, lote, formatFloat(d.Cantidad), obs,
		)
	}
	detalleSQL += "]::inventario.t_traslado_detalle[]"

	refExterna := "NULL"
	if req.ReferenciaExterna != "" {
		refExterna = fmt.Sprintf("'%s'", escapePG(req.ReferenciaExterna))
	}
	observaciones := "NULL"
	if req.Observaciones != "" {
		observaciones = fmt.Sprintf("'%s'", escapePG(req.Observaciones))
	}

	query := fmt.Sprintf(`
		SELECT id_movimiento, id_mov_comprobante, consecutivo_generado, prefijo_comprobante
		FROM inventario.fn_registrar_traslado(%d,%d,%d,%s,%s,'%s'::date,%d,%d,%d,%s)`,
		req.IDComprobante,
		req.IDBodegaOrigen,
		req.IDBodegaDestino,
		refExterna,
		observaciones,
		escapePG(req.FechaMovimiento),
		req.UserID,
		req.EmpresaID,
		req.SedeID,
		detalleSQL,
	)

	var result TrasladoResponse
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func escapePG(s string) string {
	r := ""
	for _, c := range s {
		if c == '\'' {
			r += "''"
		} else {
			r += string(c)
		}
	}
	return r
}

// ─── Listado ──────────────────────────────────────────────────

func ListTraslados(db *gorm.DB, empresaID int64) ([]TrasladoListResponse, error) {

	var results []TrasladoListResponse

	err := db.
		Table("inventario.mov_movimientos m").
		Select(`
			m.id,
			mg.consecutivo_comprobante  AS consecutivo,
			mg.prefijo_comprobante      AS prefijo,
			m.fecha_movimiento::text    AS fecha_movimiento,
			bo.nombre                   AS bodega_origen,
			bd.nombre                   AS bodega_destino,
			COALESCE(mg.documento_soporte, '') AS referencia_externa,
			e.nombre                    AS estado`).
		Joins("INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante").
		Joins("INNER JOIN inventario.cfg_bodegas bo ON bo.id = m.id_bodega_origen").
		Joins("INNER JOIN inventario.cfg_bodegas bd ON bd.id = m.id_bodega_destino").
		Joins("INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado").
		Where(`m.id_empresa = ? AND m.id_tipo_movimiento = (
			SELECT id FROM inventario.ref_tipo_movimiento WHERE nombre = 'Traslados'
		)`, empresaID).
		Order("m.id DESC").
		Limit(50).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []TrasladoListResponse{}
	}
	return results, nil
}

// ─── Detalle por ID ───────────────────────────────────────────

func ListTrasladoByID(db *gorm.DB, id int64) (*TrasladoFullResponse, error) {

	var header TrasladoFullResponse

	err := db.
		Table("inventario.mov_movimientos m").
		Select(`
			m.id,
			mg.consecutivo_comprobante          AS consecutivo,
			mg.prefijo_comprobante              AS prefijo,
			c.nombre_comprobante                AS comprobante,
			m.fecha_movimiento::text            AS fecha_movimiento,
			bo.nombre                           AS bodega_origen,
			bd.nombre                           AS bodega_destino,
			COALESCE(mg.documento_soporte, '')  AS referencia_externa,
			COALESCE(mg.observaciones, '')      AS observaciones,
			e.nombre                            AS estado`).
		Joins("INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante").
		Joins("INNER JOIN comprobantes.cfg_comprobante c ON c.id = mg.id_comprobante").
		Joins("INNER JOIN inventario.cfg_bodegas bo ON bo.id = m.id_bodega_origen").
		Joins("INNER JOIN inventario.cfg_bodegas bd ON bd.id = m.id_bodega_destino").
		Joins("INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado").
		Where("m.id = ?", id).
		Scan(&header).Error

	if err != nil {
		return nil, err
	}
	if header.ID == 0 {
		return nil, nil
	}

	var detalle []TrasladoDetalleResponse
	err = db.
		Table("inventario.mov_movimientos_detalle d").
		Select(`
			d.id_articulo,
			a.codigo,
			a.nombre,
			u.nombre              AS unidad_medida,
			COALESCE(d.lote, '') AS lote,
			d.cantidad,
			COALESCE(d.observacion, '') AS observacion`).
		Joins("INNER JOIN inventario.cfg_articulos a ON a.id = d.id_articulo").
		Joins("INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida").
		Where("d.id_movimiento = ?", id).
		Scan(&detalle).Error

	if err != nil {
		return nil, err
	}
	header.Detalle = detalle
	return &header, nil
}

// ─── Anular traslado ──────────────────────────────────────────

func AnularTraslado(db *gorm.DB, idMovimiento int64, req *AnulacionTrasladoRequest) error {
	return db.Exec(
		`SELECT inventario.fn_anular_traslado($1, $2, $3, $4)`,
		idMovimiento, req.Motivo, req.UserID, req.EmpresaID,
	).Error
}

// ─── Refs ─────────────────────────────────────────────────────

func ListBodegas(db *gorm.DB, empresaID int64) ([]RefBodega, error) {
	var results []RefBodega
	err := db.
		Table("inventario.cfg_bodegas").
		Select("id, codigo, nombre").
		Where("id_empresa = ? AND is_active = true", empresaID).
		Order("nombre ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func ListLotesByArticuloBodega(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) ([]RefLote, error) {
	var results []RefLote
	err := db.
		Table("inventario.inv_existencias_articulos").
		Select("COALESCE(lote, '') AS lote, stock_actual").
		Where(`id_articulo = ? AND id_bodega = ? AND id_empresa = ? AND stock_actual > 0`,
			idArticulo, idBodega, empresaID).
		Order("lote ASC NULLS FIRST").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []RefLote{}
	}
	return results, nil
}

func ListComprobantesTraslado(db *gorm.DB, empresaID int64) ([]RefComprobanteTraslado, error) {
	var results []RefComprobanteTraslado
	err := db.
		Table("comprobantes.cfg_comprobante c").
		Select("c.id, c.nombre_comprobante AS nombre, c.prefijo_comprobante AS prefijo").
		Joins("INNER JOIN comprobantes.ref_tipo_operacion o ON o.id = c.id_tipo_operacion").
		Where("c.id_empresa = ? AND c.is_active = true AND o.nombre = 'Traslado'", empresaID).
		Order("c.id ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GetStockArticuloBodegaLote verifica stock actual para validación en frontend
func GetStockArticuloBodegaLote(db *gorm.DB, idArticulo int64, idBodega int64, lote string, empresaID int64) (float64, error) {
	var stock float64
	err := db.
		Table("inventario.inv_existencias_articulos").
		Select("COALESCE(stock_actual, 0)").
		Where(`id_articulo = ? AND id_bodega = ? AND id_empresa = ? AND COALESCE(lote,'') = ?`,
			idArticulo, idBodega, empresaID, lote).
		Scan(&stock).Error
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return stock, nil
}