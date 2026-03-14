package entradas

import (
	"fmt"

	"gorm.io/gorm"
)

// RegisterEntrada llama la función PostgreSQL fn_registrar_entrada
// pasando el arreglo de detalle como un array de tipo compuesto.
func RegisterEntrada(db *gorm.DB, req *EntradaRequest) (*EntradaResponse, error) {

	// Construir el array de detalle en formato PostgreSQL:
	// ROW(id_articulo, id_proveedor, marca, lote, cantidad, valor_unitario, id_impuesto, observacion)
	detalleSQL := "ARRAY["
	for i, d := range req.Detalle {
		if i > 0 {
			detalleSQL += ","
		}

		idProveedor := "NULL"
		if d.IDProveedor != nil {
			idProveedor = fmt.Sprintf("%d", *d.IDProveedor)
		}

		idImpuesto := "NULL"
		if d.IDImpuesto != nil {
			idImpuesto = fmt.Sprintf("%d", *d.IDImpuesto)
		}

		lote := "NULL"
		if d.Lote != "" {
			lote = fmt.Sprintf("'%s'", d.Lote)
		}

		marca := "NULL"
		if d.Marca != "" {
			marca = fmt.Sprintf("'%s'", d.Marca)
		}

		observacion := "NULL"
		if d.Observacion != "" {
			observacion = fmt.Sprintf("'%s'", d.Observacion)
		}

		detalleSQL += fmt.Sprintf(
			"ROW(%d,%s,%s,%s,%f,%f,%s,%s)::inventario.t_entrada_detalle",
			d.IDArticulo,
			idProveedor,
			marca,
			lote,
			d.Cantidad,
			d.ValorUnitario,
			idImpuesto,
			observacion,
		)
	}
	detalleSQL += "]::inventario.t_entrada_detalle[]"

	// Referencia externa y observaciones opcionales
	refExterna := "NULL"
	if req.ReferenciaExterna != "" {
		refExterna = fmt.Sprintf("'%s'", req.ReferenciaExterna)
	}

	observaciones := "NULL"
	if req.Observaciones != "" {
		observaciones = fmt.Sprintf("'%s'", req.Observaciones)
	}

	idTercero := "NULL"
	if req.IDTercero != nil {
		idTercero = fmt.Sprintf("%d", *req.IDTercero)
	}

	query := fmt.Sprintf(`
		SELECT
			id_movimiento,
			id_mov_comprobante,
			consecutivo_generado,
			prefijo_comprobante
		FROM inventario.fn_registrar_entrada(
			%d,         -- p_id_comprobante
			%d,         -- p_id_bodega
			%s,         -- p_id_tercero
			%s,         -- p_referencia_externa
			%s,         -- p_observaciones
			'%s'::date, -- p_fecha_movimiento
			%d,         -- p_created_by
			%d,         -- p_id_empresa
			%d,         -- p_id_sede
			%s          -- p_detalle
		)`,
		req.IDComprobante,
		req.IDBodega,
		idTercero,
		refExterna,
		observaciones,
		req.FechaMovimiento,
		req.UserID,
		req.EmpresaID,
		req.SedeID,
		detalleSQL,
	)

	var result EntradaResponse
	err := db.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListEntradas — listado de entradas por empresa con join a comprobante y bodega
func ListEntradas(db *gorm.DB, empresaID int64) ([]EntradaListResponse, error) {

	var results []EntradaListResponse

	err := db.
		Table("inventario.mov_movimientos m").
		Select(`
			m.id,
			mg.consecutivo_comprobante  AS consecutivo,
			mg.prefijo_comprobante      AS prefijo,
			m.fecha_movimiento::text    AS fecha_movimiento,
			COALESCE(t.razon_social, t.nombres || ' ' || t.apellidos, '') AS proveedor,
			COALESCE(mg.documento_soporte, '')  AS referencia_externa,
			b.nombre                    AS bodega,
			COALESCE(mg.valor_total_comprobante, 0) AS valor_total,
			e.nombre                    AS estado`).
		Joins("INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante").
		Joins("INNER JOIN inventario.cfg_bodegas b ON b.id = m.id_bodega_destino").
		Joins("INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado").
		Joins("LEFT  JOIN configuracion.cfg_terceros t ON t.id = mg.id_tercero").
		Where("m.id_empresa = ? AND m.id_tipo_movimiento = (SELECT id FROM inventario.ref_tipo_movimiento WHERE nombre = 'Entradas')", empresaID).
		Order("m.id DESC").
		Limit(50).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []EntradaListResponse{}
	}

	return results, nil
}

// ListEntradaByID — cabecera + detalle de una entrada específica
func ListEntradaByID(db *gorm.DB, id int64) (*EntradaFullResponse, error) {

	var header EntradaFullResponse

	err := db.
		Table("inventario.mov_movimientos m").
		Select(`
			m.id,
			mg.consecutivo_comprobante          AS consecutivo,
			mg.prefijo_comprobante              AS prefijo,
			c.nombre_comprobante                AS comprobante,
			m.fecha_movimiento::text            AS fecha_movimiento,
			b.nombre                            AS bodega,
			COALESCE(t.razon_social, t.nombres || ' ' || t.apellidos, '') AS proveedor,
			COALESCE(mg.documento_soporte, '')  AS referencia_externa,
			COALESCE(mg.observaciones, '')      AS observaciones,
			e.nombre                            AS estado,
			COALESCE(mg.valor_total_comprobante - mg.valor_impuesto, 0) AS valor_subtotal,
			COALESCE(mg.valor_impuesto, 0)      AS valor_impuesto,
			COALESCE(mg.valor_total_comprobante, 0) AS valor_total`).
		Joins("INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante").
		Joins("INNER JOIN comprobantes.cfg_comprobante c ON c.id = mg.id_comprobante").
		Joins("INNER JOIN inventario.cfg_bodegas b ON b.id = m.id_bodega_destino").
		Joins("INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado").
		Joins("LEFT  JOIN configuracion.cfg_terceros t ON t.id = mg.id_tercero").
		Where("m.id = ?", id).
		Scan(&header).Error

	if err != nil {
		return nil, err
	}

	if header.ID == 0 {
		return nil, nil
	}

	// Cargar detalle de líneas
	var detalle []EntradaDetalleResponse

	err = db.
		Table("inventario.mov_movimientos_detalle d").
		Select(`
			d.id_articulo,
			a.codigo,
			a.nombre,
			u.nombre                        AS unidad_medida,
			COALESCE(d.marca, '')           AS marca,
			COALESCE(d.lote, '')            AS lote,
			d.cantidad,
			0                               AS valor_unitario,
			0                               AS valor_impuesto,
			0                               AS valor_total,
			COALESCE(d.observacion, '')     AS observacion`).
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

func ListComprobantesEntrada(db *gorm.DB, empresaID int64) ([]RefComprobante, error) {

	var results []RefComprobante

	err := db.
		Table("comprobantes.cfg_comprobante c").
		Select("c.id, c.nombre_comprobante AS nombre, c.prefijo_comprobante AS prefijo").
		Joins("INNER JOIN comprobantes.ref_tipo_operacion o ON o.id = c.id_tipo_operacion").
		Where("c.id_empresa = ? AND c.is_active = true AND o.nombre = 'Entrada'", empresaID).
		Order("c.id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetExistenciaArticulo(db *gorm.DB, idArticulo int64, idBodega int64, empresaID int64) (*ArticuloExistencia, error) {

	var result ArticuloExistencia

	err := db.
		Table("inventario.cfg_articulos a").
		Select(`
			a.id                                AS id_articulo,
			a.codigo,
			a.nombre,
			u.nombre                            AS unidad_medida,
			COALESCE(e.stock_actual, 0)         AS stock_actual,
			COALESCE(e.costo_promedio, 0)       AS costo_promedio`).
		Joins("INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida").
		Joins(`LEFT JOIN inventario.inv_existencias_articulos e
			ON e.id_articulo = a.id
			AND e.id_bodega  = ?
			AND e.id_empresa = ?`, idBodega, empresaID).
		Where("a.id = ?", idArticulo).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func ListImpuestos(db *gorm.DB, empresaID int64) ([]RefImpuesto, error) {

	var results []RefImpuesto

	err := db.
		Table("contabilidad.cfg_impuestos").
		Select("id, nombre, porcentaje").
		Where("id_empresa = ? AND is_active = true", empresaID).
		Order("nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
