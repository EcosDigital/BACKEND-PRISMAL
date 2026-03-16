package bajas

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
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

// ─── Registrar baja ───────────────────────────────────────────

func RegisterBaja(db *gorm.DB, req *BajaRequest) (*BajaResponse, error) {

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
		// id_subtipo_movimiento por línea es opcional — NULL hereda el del encabezado
		subtipoLinea := "NULL"
		if d.IDSubtipoMovimiento > 0 {
			subtipoLinea = fmt.Sprintf("%d", d.IDSubtipoMovimiento)
		}
		detalleSQL += fmt.Sprintf(
			"ROW(%d,%s,%s,%s,%s)::inventario.t_baja_detalle",
			d.IDArticulo, lote, formatFloat(d.Cantidad), subtipoLinea, obs,
		)
	}
	detalleSQL += "]::inventario.t_baja_detalle[]"

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
		FROM inventario.fn_registrar_baja(%d,%d,%d,%s,%s,'%s'::date,%d,%d,%d,%s)`,
		req.IDComprobante, req.IDBodega,
		req.IDSubtipoMovimiento, // nuevo parámetro
		refExterna, observaciones,
		escapePG(req.FechaMovimiento),
		req.UserID, req.EmpresaID, req.SedeID,
		detalleSQL,
	)

	var result BajaResponse
	if err := db.Raw(query).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// ─── Listado ──────────────────────────────────────────────────

func ListBajas(db *gorm.DB, empresaID int64) ([]BajaListResponse, error) {
	var results []BajaListResponse
	err := db.
		Table("inventario.mov_movimientos m").
		Select(`
			m.id,
			mg.consecutivo_comprobante              AS consecutivo,
			mg.prefijo_comprobante                  AS prefijo,
			m.fecha_movimiento::text                AS fecha_movimiento,
			b.nombre                                AS bodega,
			COALESCE(mg.documento_soporte, '')      AS referencia_externa,
			e.nombre                                AS estado,
			COALESCE(mg.valor_total_comprobante, 0) AS total_costo`).
		Joins("INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante").
		Joins("INNER JOIN inventario.cfg_bodegas b ON b.id = m.id_bodega_origen").
		Joins("INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado").
		Where(`m.id_empresa = ? AND m.id_tipo_movimiento = (
			SELECT id FROM inventario.ref_tipo_movimiento WHERE nombre = 'Bajas'
		)`, empresaID).
		Order("m.id DESC").Limit(50).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []BajaListResponse{}
	}
	return results, nil
}

// ─── Detalle por ID ───────────────────────────────────────────

func ListBajaByID(db *gorm.DB, id int64) (*BajaFullResponse, error) {
	var header BajaFullResponse
	err := db.
		Table("inventario.mov_movimientos m").
		Select(`
			m.id,
			mg.consecutivo_comprobante                  AS consecutivo,
			mg.prefijo_comprobante                      AS prefijo,
			c.nombre_comprobante                        AS comprobante,
			m.fecha_movimiento::text                    AS fecha_movimiento,
			b.nombre                                    AS bodega,
			COALESCE(mg.documento_soporte, '')          AS referencia_externa,
			COALESCE(mg.observaciones, '')              AS observaciones,
			e.nombre                                    AS estado,
			COALESCE(mg.valor_total_comprobante, 0)     AS total_costo`).
		Joins("INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante").
		Joins("INNER JOIN comprobantes.cfg_comprobante c ON c.id = mg.id_comprobante").
		Joins("INNER JOIN inventario.cfg_bodegas b ON b.id = m.id_bodega_origen").
		Joins("INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado").
		Where("m.id = ?", id).
		Scan(&header).Error
	if err != nil {
		return nil, err
	}
	if header.ID == 0 {
		return nil, nil
	}

	var detalle []BajaDetalleResponse
	err = db.
		Table("inventario.mov_movimientos_detalle d").
		Select(`
			d.id_articulo,
			a.codigo,
			a.nombre,
			u.nombre                        AS unidad_medida,
			COALESCE(d.lote, '')            AS lote,
			d.cantidad,
			COALESCE(tb.nombre, '')         AS tipo_baja,
			d.valor_unitario                AS costo_unitario,
			d.valor_total                   AS costo_total,
			COALESCE(d.observacion, '')     AS observacion`).
		Joins("INNER JOIN inventario.cfg_articulos a ON a.id = d.id_articulo").
		Joins("INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida").
		Joins("LEFT JOIN inventario.ref_tipo_baja tb ON tb.id = d.id_tipo_baja").
		Where("d.id_movimiento = ?", id).
		Scan(&detalle).Error
	if err != nil {
		return nil, err
	}
	header.Detalle = detalle
	return &header, nil
}

// ─── Anular ───────────────────────────────────────────────────

func AnularBaja(db *gorm.DB, idMovimiento int64, req *AnulacionBajaRequest) error {
	return db.Exec(
		`SELECT inventario.fn_anular_baja($1, $2, $3, $4)`,
		idMovimiento, req.Motivo, req.UserID, req.EmpresaID,
	).Error
}

// ─── Refs ─────────────────────────────────────────────────────

func ListSubtiposMovimientoBaja(db *gorm.DB) ([]RefSubtipoMovimiento, error) {
	var results []RefSubtipoMovimiento
	err := db.
		Table("inventario.ref_subtipo_movimiento s").
		Select("s.id, s.id_tipo_movimiento, s.codigo, s.nombre").
		Joins("INNER JOIN inventario.ref_tipo_movimiento t ON t.id = s.id_tipo_movimiento").
		Where("t.nombre = 'Bajas' AND s.is_active = true").
		Order("s.nombre ASC").
		Scan(&results).Error
	return results, err
}

func ListComprobantesBaja(db *gorm.DB, empresaID int64) ([]RefComprobanteBaja, error) {
	var results []RefComprobanteBaja
	err := db.
		Table("comprobantes.cfg_comprobante c").
		Select("c.id, c.nombre_comprobante AS nombre, c.prefijo_comprobante AS prefijo").
		Joins("INNER JOIN comprobantes.ref_tipo_operacion o ON o.id = c.id_tipo_operacion").
		Where("c.id_empresa = ? AND c.is_active = true AND o.nombre = 'Baja'", empresaID).
		Order("c.id ASC").Scan(&results).Error
	return results, err
}

func ListLotesBaja(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) ([]RefLoteBaja, error) {
	var results []RefLoteBaja
	err := db.
		Table("inventario.inv_existencias_articulos").
		Select("COALESCE(lote, '') AS lote, stock_actual").
		Where("id_articulo = ? AND id_bodega = ? AND id_empresa = ? AND stock_actual > 0",
			idArticulo, idBodega, empresaID).
		Order("lote ASC NULLS FIRST").
		Scan(&results).Error
	if results == nil {
		results = []RefLoteBaja{}
	}
	return results, err
}

// ─── Reporte ──────────────────────────────────────────────────

func GetBajaReporte(db *gorm.DB, idMovimiento int64, empresaID int64) (*BajaReporte, error) {
	var reporte BajaReporte
	err := db.Raw(`
		SELECT
			m.id,
			mg.consecutivo_comprobante              AS consecutivo,
			mg.prefijo_comprobante                  AS prefijo,
			c.nombre_comprobante                    AS comprobante,
			m.fecha_movimiento::text                AS fecha_movimiento,
			e.nombre                                AS estado,
			b.nombre                                AS bodega,
			COALESCE(mg.documento_soporte, '')      AS referencia_externa,
			COALESCE(mg.observaciones, '')          AS observaciones,
			COALESCE(
				u_t.primer_nombre || ' ' || u_t.primer_apellido,
				m.created_by::text
			)                                       AS elaboro
		FROM inventario.mov_movimientos m
		INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante
		INNER JOIN comprobantes.cfg_comprobante c ON c.id = mg.id_comprobante
		INNER JOIN inventario.cfg_bodegas b ON b.id = m.id_bodega_origen
		INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado
		LEFT  JOIN seguridad.cfg_usuarios u ON u.id = m.created_by
		LEFT  JOIN configuracion.cfg_terceros u_t ON u_t.id = u.id_tercero
		WHERE m.id = ? AND m.id_empresa = ?`,
		idMovimiento, empresaID,
	).Scan(&reporte).Error
	if err != nil {
		return nil, err
	}
	if reporte.ID == 0 {
		return nil, nil
	}

	var empresa EmpresaReporteBaja
	err = db.Raw(`
		SELECT emp.razon_social, emp.nit, emp.dv, emp.direccion,
			mun.nombre_municipio AS ciudad, dep.nombre AS departamento,
			emp.telefono, COALESCE(emp.logo_url,'') AS logo_url
		FROM configuracion.cfg_empresas emp
		INNER JOIN configuracion.ref_municipios mun ON mun.id = emp.id_ciudad
		INNER JOIN configuracion.ref_departamentos dep ON dep.id = emp.id_departamento
		WHERE emp.id = ?`, empresaID,
	).Scan(&empresa).Error
	if err != nil {
		return nil, err
	}
	reporte.Empresa = empresa

	var detalle []BajaReporteDetalle
	err = db.Raw(`
		SELECT
			a.codigo,
			a.nombre,
			u.nombre                    AS unidad_medida,
			COALESCE(d.lote, '')        AS lote,
			COALESCE(tb.nombre, '')     AS tipo_baja,
			d.cantidad,
			d.valor_unitario            AS costo_unitario,
			d.valor_total               AS costo_total,
			COALESCE(d.observacion, '') AS observacion
		FROM inventario.mov_movimientos_detalle d
		INNER JOIN inventario.cfg_articulos a ON a.id = d.id_articulo
		INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida
		LEFT  JOIN inventario.ref_subtipo_movimiento tb ON tb.id = d.id_subtipo_movimiento
		WHERE d.id_movimiento = ?
		ORDER BY d.id ASC`, idMovimiento,
	).Scan(&detalle).Error
	if err != nil {
		return nil, err
	}
	reporte.Detalle = detalle
	for _, d := range detalle {
		reporte.TotalUnidades += d.Cantidad
		reporte.TotalCosto += d.CostoTotal
	}
	return &reporte, nil
}
