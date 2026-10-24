package herramientas

import (
	"gorm.io/gorm"
)

type ArticuloExistencia struct {
	Codigo          string   `json:"codigo"`
	Nombre          string   `json:"nombre"`
	Existencias     float64  `json:"existencias"`
	PrecioPublico   *float64 `json:"precio_publico"`
	AplicaDomicilio *bool    `json:"aplica_domicilio"`
}

type MovimientoResumen struct {
	Tipo        string  `json:"tipo"`
	Veces       int     `json:"veces"`
	Cantidad    float64 `json:"cantidad"`
	UltimaFecha string  `json:"ultima_fecha"`
}

type ArticuloBajo struct {
	Codigo      string  `json:"codigo"`
	Nombre      string  `json:"nombre"`
	Bodega      string  `json:"bodega"`
	Existencias float64 `json:"existencias"`
	Minimo      float64 `json:"minimo"`
}

// BuscarArticulos busca artículos activos por nombre o código y devuelve
// sus existencias en las bodegas de la sede, y su precio público y
// domicilio según el catálogo de la sede (si el artículo está publicado).
func BuscarArticulos(db *gorm.DB, sedeID int64, texto string) ([]ArticuloExistencia, error) {

	var results []ArticuloExistencia

	err := db.Raw(`
		SELECT a.codigo, a.nombre,
		       COALESCE(SUM(e.stock_actual), 0) AS existencias,
		       ca.precio_publico, ca.aplica_domicilio
		FROM inventario.cfg_articulos a
		LEFT JOIN inventario.inv_existencias_articulos e ON e.id_articulo = a.id
		     AND e.id_bodega IN (SELECT id FROM inventario.cfg_bodegas WHERE id_sede = ? AND is_active)
		LEFT JOIN inventario.cfg_catalogos c ON c.id_sede = ? AND c.is_active
		LEFT JOIN inventario.cfg_catalogo_articulos ca ON ca.id_catalogo = c.id AND ca.id_articulo = a.id
		WHERE a.is_active AND (a.nombre ILIKE ? OR a.codigo ILIKE ?)
		GROUP BY a.id, a.codigo, a.nombre, ca.precio_publico, ca.aplica_domicilio
		ORDER BY a.nombre
		LIMIT 15`,
		sedeID, sedeID, "%"+texto+"%", "%"+texto+"%").Scan(&results).Error

	if results == nil {
		results = []ArticuloExistencia{}
	}
	return results, err
}

// ListarCatalogo devuelve los artículos publicados en el catálogo de la
// sede con su precio, domicilio y existencias, primero los que tienen
// existencias. Sirve para recomendar u ofrecer productos sin un nombre
// concreto.
func ListarCatalogo(db *gorm.DB, sedeID int64) ([]ArticuloExistencia, error) {

	var results []ArticuloExistencia

	err := db.Raw(`
		SELECT a.codigo, a.nombre,
		       COALESCE(SUM(e.stock_actual), 0) AS existencias,
		       ca.precio_publico, ca.aplica_domicilio
		FROM inventario.cfg_catalogos c
		JOIN inventario.cfg_catalogo_articulos ca ON ca.id_catalogo = c.id
		JOIN inventario.cfg_articulos a ON a.id = ca.id_articulo AND a.is_active
		LEFT JOIN inventario.inv_existencias_articulos e ON e.id_articulo = a.id
		     AND e.id_bodega IN (SELECT id FROM inventario.cfg_bodegas WHERE id_sede = ? AND is_active)
		WHERE c.id_sede = ? AND c.is_active
		GROUP BY a.id, a.codigo, a.nombre, ca.precio_publico, ca.aplica_domicilio
		ORDER BY COALESCE(SUM(e.stock_actual), 0) > 0 DESC, a.nombre
		LIMIT 40`,
		sedeID, sedeID).Scan(&results).Error

	if results == nil {
		results = []ArticuloExistencia{}
	}
	return results, err
}

// MovimientosArticulo resume, por tipo de movimiento, cuántas veces y qué
// cantidad se movió un artículo en la sede entre dos fechas (sin anulados).
func MovimientosArticulo(db *gorm.DB, sedeID int64, texto, desde, hasta string) ([]MovimientoResumen, error) {

	var results []MovimientoResumen

	err := db.Raw(`
		SELECT t.nombre AS tipo, COUNT(DISTINCT m.id) AS veces,
		       COALESCE(SUM(d.cantidad), 0) AS cantidad,
		       TO_CHAR(MAX(m.fecha_movimiento), 'YYYY-MM-DD') AS ultima_fecha
		FROM inventario.mov_movimientos m
		JOIN inventario.mov_movimientos_detalle d ON d.id_movimiento = m.id
		JOIN inventario.cfg_articulos a ON a.id = d.id_articulo
		JOIN inventario.ref_tipo_movimiento t ON t.id = m.id_tipo_movimiento
		WHERE m.id_sede = ? AND m.anulado_at IS NULL
		  AND m.fecha_movimiento BETWEEN ?::date AND ?::date
		  AND (a.nombre ILIKE ? OR a.codigo ILIKE ?)
		GROUP BY t.nombre
		ORDER BY t.nombre`,
		sedeID, desde, hasta, "%"+texto+"%", "%"+texto+"%").Scan(&results).Error

	if results == nil {
		results = []MovimientoResumen{}
	}
	return results, err
}

// ArticulosBajoMinimo lista los artículos cuyas existencias en una bodega
// de la sede están en o por debajo del mínimo configurado en esa bodega.
func ArticulosBajoMinimo(db *gorm.DB, sedeID int64) ([]ArticuloBajo, error) {

	var results []ArticuloBajo

	err := db.Raw(`
		SELECT a.codigo, a.nombre, b.nombre AS bodega,
		       e.stock_actual AS existencias, b.stock_minimo_global AS minimo
		FROM inventario.inv_existencias_articulos e
		JOIN inventario.cfg_articulos a ON a.id = e.id_articulo AND a.is_active
		JOIN inventario.cfg_bodegas b ON b.id = e.id_bodega AND b.is_active
		WHERE b.id_sede = ? AND b.stock_minimo_global IS NOT NULL
		  AND e.stock_actual <= b.stock_minimo_global
		ORDER BY e.stock_actual ASC
		LIMIT 30`,
		sedeID).Scan(&results).Error

	if results == nil {
		results = []ArticuloBajo{}
	}
	return results, err
}
