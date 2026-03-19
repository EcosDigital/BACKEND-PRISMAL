package kardex

import (
	"strings"

	"gorm.io/gorm"
)

// derivarClase obtiene la clase visual a partir del tipo_movimiento.
// Se calcula en Go para no depender de columnas extra en la función SQL.
func derivarClase(tipoMovimiento string) string {
	t := strings.ToLower(tipoMovimiento)
	switch {
	case strings.Contains(t, "saldo inicial"):
		return "saldo_inicial"
	case strings.Contains(t, "entrada"):
		return "entrada"
	case strings.Contains(t, "baja"):
		return "baja"
	case strings.Contains(t, "traslado") && strings.Contains(t, "salida"):
		return "traslado_salida"
	case strings.Contains(t, "traslado") && strings.Contains(t, "entrada"):
		return "traslado_entrada"
	case strings.Contains(t, "traslado"):
		// Sin sufijo: se determina en el frontend según entrada/salida > 0
		return "traslado_entrada"
	case strings.Contains(t, "despacho"):
		return "despacho"
	default:
		return "saldo_inicial"
	}
}

// GetKardex ejecuta fn_kardex_articulo y arma el reporte completo.
// Lote es OPCIONAL: si req.Lote == "" se retornan todos los lotes.
func GetKardex(
	db *gorm.DB,
	req *KardexRequest,
	empresaID int64,
) (*KardexReporte, error) {

	// ── 1. Datos del artículo ─────────────────────────────────────
	type ArticuloInfo struct {
		Codigo       string `gorm:"column:codigo"`
		Nombre       string `gorm:"column:nombre"`
		UnidadMedida string `gorm:"column:unidad_medida"`
	}
	var artInfo ArticuloInfo

	err := db.Raw(`
		SELECT a.codigo, a.nombre, u.nombre AS unidad_medida
		FROM inventario.cfg_articulos a
		INNER JOIN inventario.ref_unidad_medidas u ON u.id = a.id_unidad_medida
		WHERE a.id = ?`, req.IDArticulo,
	).Scan(&artInfo).Error
	if err != nil {
		return nil, err
	}

	// ── 2. Datos de la bodega ─────────────────────────────────────
	type BodegaInfo struct {
		Nombre string `gorm:"column:nombre"`
	}
	var bodInfo BodegaInfo

	err = db.Raw(`
		SELECT nombre FROM inventario.cfg_bodegas WHERE id = ?`, req.IDBodega,
	).Scan(&bodInfo).Error
	if err != nil {
		return nil, err
	}

	// ── 3. Datos de la empresa ────────────────────────────────────
	var empresa KardexEmpresa

	err = db.Raw(`
		SELECT
			emp.razon_social,
			emp.nit,
			emp.dv,
			emp.direccion,
			mun.nombre_municipio AS ciudad,
			dep.nombre           AS departamento,
			emp.telefono,
			COALESCE(emp.logo_url, '') AS logo_url
		FROM configuracion.cfg_empresas emp
		INNER JOIN configuracion.ref_municipios mun ON mun.id = emp.id_ciudad
		INNER JOIN configuracion.ref_departamentos dep ON dep.id = emp.id_departamento
		WHERE emp.id = ?`, empresaID,
	).Scan(&empresa).Error
	if err != nil {
		return nil, err
	}

	// ── 4. Ejecutar fn_kardex_articulo ────────────────────────────
	// Lote: si viene vacío se pasa NULL para que la función retorne
	// todos los lotes (la función evalúa p_lote IS NOT NULL AND p_lote <> '').
	var loteParam interface{}
	if req.Lote == "" {
		loteParam = nil
	} else {
		loteParam = req.Lote
	}

	// Seleccionamos exactamente las columnas que retorna la función —
	// sin "clase", que no existe en el RETURNS TABLE de fn_kardex_articulo.
	// La clase se deriva en Go con derivarClase().
	type filaRaw struct {
		FilaOrden         int     `gorm:"column:fila_orden"`
		Fecha             string  `gorm:"column:fecha"`
		TipoMovimiento    string  `gorm:"column:tipo_movimiento"`
		Documento         string  `gorm:"column:documento"`
		ReferenciaExterna string  `gorm:"column:referencia_externa"`
		Observaciones     string  `gorm:"column:observaciones"`
		CantidadEntrada   float64 `gorm:"column:cantidad_entrada"`
		CantidadSalida    float64 `gorm:"column:cantidad_salida"`
		CostoUnitario     float64 `gorm:"column:costo_unitario"`
		SaldoCantidad     float64 `gorm:"column:saldo_cantidad"`
		EsAnulado         bool    `gorm:"column:es_anulado"`
		MotivoAnulacion   string  `gorm:"column:motivo_anulacion"`
		IDMovimiento      int     `gorm:"column:id_movimiento"`
	}

	var rawFilas []filaRaw

	err = db.Raw(`
		SELECT
			fila_orden,
			fecha,
			tipo_movimiento,
			documento,
			referencia_externa,
			observaciones,
			cantidad_entrada,
			cantidad_salida,
			costo_unitario,
			saldo_cantidad,
			es_anulado,
			COALESCE(motivo_anulacion, '') AS motivo_anulacion,
			id_movimiento
		FROM inventario.fn_kardex_articulo($1, $2, $3, $4::date, $5::date, $6)
		ORDER BY fila_orden ASC`,
		req.IDArticulo,
		req.IDBodega,
		loteParam,
		req.FechaInicio,
		req.FechaFin,
		empresaID,
	).Scan(&rawFilas).Error

	if err != nil {
		return nil, err
	}

	// ── 5. Mapear a KardexFila enriqueciendo con Clase ────────────
	filas := make([]KardexFila, len(rawFilas))
	for i, r := range rawFilas {
		filas[i] = KardexFila{
			FilaOrden:         r.FilaOrden,
			Fecha:             r.Fecha,
			TipoMovimiento:    r.TipoMovimiento,
			Documento:         r.Documento,
			ReferenciaExterna: r.ReferenciaExterna,
			Observaciones:     r.Observaciones,
			CantidadEntrada:   r.CantidadEntrada,
			CantidadSalida:    r.CantidadSalida,
			CostoUnitario:     r.CostoUnitario,
			SaldoCantidad:     r.SaldoCantidad,
			EsAnulado:         r.EsAnulado,
			MotivoAnulacion:   r.MotivoAnulacion,
			IDMovimiento:      r.IDMovimiento,
			// Clase derivada en Go — no viene de la función SQL
			Clase: derivarClase(r.TipoMovimiento),
		}
	}

	// ── 6. Stock final = saldo de la última fila ──────────────────
	stockFinal := 0.0
	if len(filas) > 0 {
		stockFinal = filas[len(filas)-1].SaldoCantidad
	}

	// ── 7. Construir reporte ──────────────────────────────────────
	loteLabel := req.Lote
	if loteLabel == "" {
		loteLabel = "Todos los lotes"
	}

	reporte := &KardexReporte{
		Empresa: empresa,
		Cabecera: KardexCabecera{
			IDArticulo:   req.IDArticulo,
			CodigoArt:    artInfo.Codigo,
			NombreArt:    artInfo.Nombre,
			UnidadMedida: artInfo.UnidadMedida,
			IDBodega:     req.IDBodega,
			NombreBodega: bodInfo.Nombre,
			Lote:         loteLabel,
			FechaInicio:  req.FechaInicio,
			FechaFin:     req.FechaFin,
			StockFinal:   stockFinal,
		},
		Filas: filas,
	}

	return reporte, nil
}

// ListLotesByArticuloBodega — lotes disponibles para el selector del formulario.
func ListLotesByArticuloBodega(
	db *gorm.DB,
	idArticulo int64,
	idBodega int64,
	empresaID int64,
) ([]KardexRefLote, error) {

	var results []KardexRefLote

	err := db.Raw(`
		SELECT DISTINCT COALESCE(lote, '') AS lote
		FROM inventario.inv_existencias_articulos
		WHERE id_articulo = ?
		  AND id_bodega   = ?
		  AND id_empresa  = ?
		ORDER BY lote ASC`,
		idArticulo, idBodega, empresaID,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
