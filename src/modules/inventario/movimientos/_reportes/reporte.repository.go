package reportes_inv

import (
	"gorm.io/gorm"
)

func GetEntradaReporte(db *gorm.DB, idMovimiento int64, empresaID int64) (*EntradaReporte, error) {

	// ── Cabecera del reporte ──────────────────────────────────────
	var reporte EntradaReporte

	err := db.Raw(`
		SELECT
			m.id,
			mg.consecutivo_comprobante                                  AS consecutivo,
			mg.prefijo_comprobante                                      AS prefijo,
			c.nombre_comprobante                                        AS comprobante,
			m.fecha_movimiento::text                                    AS fecha_movimiento,
			e.nombre                                                    AS estado,
			COALESCE(mg.documento_soporte, '')                          AS referencia_externa,
			COALESCE(mg.observaciones, '')                              AS observaciones,
			-- Proveedor
			COALESCE(t.razon_social, t.primer_nombre || ' ' || t.primer_apellido, '') AS proveedor,
			COALESCE(t.numero_documento, '')                            AS nit_proveedor,
			COALESCE(t.direccion, '')                                   AS direccion_proveedor,
			COALESCE(t.telefono, '')                                    AS telefono_proveedor,
			-- Bodega
			b.nombre                                                    AS bodega,
			-- Usuario que elaboró
			COALESCE(
				u_tercero.primer_nombre || ' ' || u_tercero.primer_apellido,
				m.created_by::text
			)                                                           AS elaboro,
			-- Totales del comprobante
			COALESCE(mg.valor_total_comprobante - mg.valor_impuesto, 0) AS subtotal,
			COALESCE(mg.valor_impuesto, 0)                              AS total_iva,
			COALESCE(mg.valor_descuento, 0)                             AS total_descuento,
			COALESCE(mg.valor_total_comprobante, 0)                     AS total_neto
		FROM inventario.mov_movimientos m
		INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante
		INNER JOIN comprobantes.cfg_comprobante c ON c.id = mg.id_comprobante
		INNER JOIN inventario.cfg_bodegas b ON b.id = m.id_bodega_destino
		INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado
		LEFT  JOIN configuracion.cfg_terceros t ON t.id = mg.id_tercero
		LEFT  JOIN seguridad.cfg_usuarios u ON u.id = m.created_by
		LEFT  JOIN configuracion.cfg_terceros u_tercero ON u_tercero.id = u.id_tercero
		WHERE m.id = ? AND m.id_empresa = ?`,
		idMovimiento, empresaID,
	).Scan(&reporte).Error

	if err != nil {
		return nil, err
	}

	if reporte.ID == 0 {
		return nil, nil
	}

	// ── Datos de la empresa ───────────────────────────────────────
	var empresa EmpresaReporte

	err = db.Raw(`
		SELECT
			emp.razon_social,
			emp.nit,
			emp.dv,
			emp.direccion,
			mun.nombre_municipio   AS ciudad,
			dep.nombre   AS departamento,
			emp.telefono,
			emp.email,
			COALESCE(emp.logo_url, '') AS logo_url
		FROM configuracion.cfg_empresas emp
		INNER JOIN configuracion.ref_municipios mun ON mun.id = emp.id_ciudad
		INNER JOIN configuracion.ref_departamentos dep ON dep.id = emp.id_departamento
		WHERE emp.id = ?`,
		empresaID,
	).Scan(&empresa).Error

	if err != nil {
		return nil, err
	}

	reporte.Empresa = empresa

	// ── Detalle de artículos ──────────────────────────────────────
	var detalle []EntradaReporteDetalle

	err = db.Raw(`
		SELECT
			a.codigo,
			a.nombre,
			COALESCE(p.nombre, u.nombre) AS presentacion,
			u.nombre                     AS unidad_medida,
			d.cantidad,
			d.valor_unitario,
			COALESCE(
				(SELECT i.porcentaje
				 FROM contabilidad.cfg_impuestos i
				 WHERE i.id = (
					SELECT imp_ref.id FROM contabilidad.cfg_impuestos imp_ref
					WHERE ROUND(d.cantidad * d.valor_unitario * imp_ref.porcentaje / 100, 2)
					      = d.valor_impuesto
					LIMIT 1
				 )), 0
			)                            AS porcentaje_iva,
			d.valor_impuesto             AS valor_iva,
			ROUND(d.cantidad * d.valor_unitario, 2) AS subtotal
		FROM inventario.mov_movimientos_detalle d
		INNER JOIN inventario.cfg_articulos a ON a.id = d.id_articulo
		INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida
		LEFT  JOIN inventario.ref_presentacion_articulo p ON p.id = a.id_presentacion
		WHERE d.id_movimiento = ?
		ORDER BY d.id ASC`,
		idMovimiento,
	).Scan(&detalle).Error

	if err != nil {
		return nil, err
	}

	reporte.Detalle = detalle

	return &reporte, nil
}

func GetTrasladoReporte(db *gorm.DB, idMovimiento int64, empresaID int64) (*TrasladoReporte, error) {

	// ── Cabecera del reporte ──────────────────────────────────────
	var reporte TrasladoReporte

	err := db.Raw(`
		SELECT
			m.id,
			mg.consecutivo_comprobante                   AS consecutivo,
			mg.prefijo_comprobante                       AS prefijo,
			c.nombre_comprobante                         AS comprobante,
			m.fecha_movimiento::text                     AS fecha_movimiento,
			e.nombre                                     AS estado,
			COALESCE(mg.observaciones, '')               AS observaciones,
			bo.nombre                                    AS bodega_origen,
			bd.nombre                                    AS bodega_destino,
			COALESCE(
				u_tercero.primer_nombre || ' ' || u_tercero.primer_apellido,
				m.created_by::text
			)                                            AS elaboro
		FROM inventario.mov_movimientos m
		INNER JOIN comprobantes.mov_gestion_comprobantes mg ON mg.id = m.id_mov_comprobante
		INNER JOIN comprobantes.cfg_comprobante c ON c.id = mg.id_comprobante
		INNER JOIN inventario.cfg_bodegas bo ON bo.id = m.id_bodega_origen
		INNER JOIN inventario.cfg_bodegas bd ON bd.id = m.id_bodega_destino
		INNER JOIN inventario.ref_estado_movimiento e ON e.id = m.id_estado
		LEFT  JOIN seguridad.cfg_usuarios u ON u.id = m.created_by
		LEFT  JOIN configuracion.cfg_terceros u_tercero ON u_tercero.id = u.id_tercero
		WHERE m.id = ? AND m.id_empresa = ?`,
		idMovimiento, empresaID,
	).Scan(&reporte).Error

	if err != nil {
		return nil, err
	}
	if reporte.ID == 0 {
		return nil, nil
	}

	// ── Datos empresa ─────────────────────────────────────────────
	var empresa EmpresaReporteTraslado

	err = db.Raw(`
		SELECT
			emp.razon_social,
			emp.nit,
			emp.dv,
			emp.direccion,
			mun.nombre_municipio   AS ciudad,
			dep.nombre   AS departamento,
			emp.telefono,
			COALESCE(emp.logo_url, '') AS logo_url
		FROM configuracion.cfg_empresas emp
		INNER JOIN configuracion.ref_municipios mun ON mun.id = emp.id_ciudad
		INNER JOIN configuracion.ref_departamentos dep ON dep.id = emp.id_departamento
		WHERE emp.id = ?`,
		empresaID,
	).Scan(&empresa).Error

	if err != nil {
		return nil, err
	}
	reporte.Empresa = empresa

	// ── Detalle artículos ─────────────────────────────────────────
	var detalle []TrasladoReporteDetalle

	err = db.Raw(`
		SELECT
			a.codigo,
			a.nombre,
			u.nombre                     AS unidad_medida,
			COALESCE(d.lote, '')         AS lote,
			d.cantidad,
			d.valor_unitario             AS costo_unitario,
			d.valor_total                AS costo_total
		FROM inventario.mov_movimientos_detalle d
		INNER JOIN inventario.cfg_articulos a ON a.id = d.id_articulo
		INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida
		WHERE d.id_movimiento = ?
		ORDER BY d.id ASC`,
		idMovimiento,
	).Scan(&detalle).Error

	if err != nil {
		return nil, err
	}
	reporte.Detalle = detalle

	// ── Calcular totales ──────────────────────────────────────────
	for _, d := range detalle {
		reporte.TotalUnidades += d.Cantidad
		reporte.TotalCosto += d.CostoTotal
	}

	return &reporte, nil
}
