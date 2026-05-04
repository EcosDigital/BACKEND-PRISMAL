package articulos

import (
	"database/sql"
	"fmt"
	"strings"
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

func SearchArticulos(db *gorm.DB, nombreCodigo string, empresaID int64) ([]ArticuloResponse, error) {

	var results []ArticuloResponse

	q := db.
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
		Where("a.id_empresa = ? AND a.is_active = true", empresaID)

	if nombreCodigo != "" {
		like := "%" + nombreCodigo + "%"
		q = q.Where("a.nombre ILIKE ? OR a.codigo ILIKE ?", like, like)
	}

	err := q.Order("a.nombre ASC").Limit(30).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ArticuloResponse{}
	}

	return results, nil
}

// ─── Carga masiva ─────────────────────────────────────────────────────────────

// LoadRefCache carga de una sola vez todos los mapas nombre→id necesarios
// para resolver las columnas de referencia del plano sin N+1 queries.
func LoadRefCache(db *gorm.DB, empresaID int64) (*refCache, error) {

	cache := &refCache{
		Grupos:         make(map[string]int),
		Unidades:       make(map[string]int),
		Presentaciones: make(map[string]int),
		FormasFarma:    make(map[string]int),
	}

	// Grupos de artículos
	type row struct {
		ID     int
		Nombre string
	}

	var grupos []row
	if err := db.Raw(`SELECT id, nombre FROM inventario.cfg_grupo_articulos WHERE is_active = true`).
		Scan(&grupos).Error; err != nil {
		return nil, err
	}
	for _, g := range grupos {
		cache.Grupos[strings.ToLower(strings.TrimSpace(g.Nombre))] = g.ID
	}

	// Unidades de medida
	var unidades []row
	if err := db.Raw(`SELECT id, nombre FROM inventario.ref_unidad_medidas`).
		Scan(&unidades).Error; err != nil {
		return nil, err
	}
	for _, u := range unidades {
		cache.Unidades[strings.ToLower(strings.TrimSpace(u.Nombre))] = u.ID
	}

	// Presentaciones
	var presentaciones []row
	if err := db.Raw(`SELECT id, nombre FROM inventario.ref_presentacion_articulo`).
		Scan(&presentaciones).Error; err != nil {
		return nil, err
	}
	for _, p := range presentaciones {
		cache.Presentaciones[strings.ToLower(strings.TrimSpace(p.Nombre))] = p.ID
	}

	// Formas farmacéuticas (opcional)
	var formas []row
	if err := db.Raw(`SELECT id, nombre FROM inventario.ref_forma_farmaceutica`).
		Scan(&formas).Error; err != nil {
		return nil, err
	}
	for _, f := range formas {
		cache.FormasFarma[strings.ToLower(strings.TrimSpace(f.Nombre))] = f.ID
	}

	return cache, nil
}

// FindArticuloByCodigoEmpresa busca un artículo por código dentro de una empresa.
// Retorna (id, encontrado, error).
func FindArticuloByCodigoEmpresa(db *gorm.DB, codigo string, empresaID int64) (int64, bool, error) {
	var id int64
	err := db.Raw(
		`SELECT id FROM inventario.cfg_articulos WHERE codigo = ? AND id_empresa = ? LIMIT 1`,
		codigo, empresaID,
	).Scan(&id).Error

	if err != nil {
		return 0, false, err
	}
	return id, id > 0, nil
}

// BulkUpsertArticulo inserta o actualiza un artículo a partir de una fila del plano.
// Retorna la acción ejecutada ("creado" | "actualizado") y el error si ocurre.
func BulkUpsertArticulo(
	db *gorm.DB,
	fila *PlanoArticuloRow,
	cache *refCache,
	userID int64,
	empresaID int64,
	sedeID int64,
) (accion string, err error) {

	// ── 1. Resolver referencias nombre → ID ───────────────────────────────────

	grupoID, ok := cache.Grupos[strings.ToLower(strings.TrimSpace(fila.IDTipoArticuloNombre))]
	if !ok {
		return "error", fmt.Errorf("grupo '%s' no encontrado", fila.IDTipoArticuloNombre)
	}

	unidadID, ok := cache.Unidades[strings.ToLower(strings.TrimSpace(fila.IDUnidadMediaNombre))]
	if !ok {
		return "error", fmt.Errorf("unidad de medida '%s' no encontrada", fila.IDUnidadMediaNombre)
	}

	var presID *int
	if fila.IDPresentacionNombre != "" {
		id, ok := cache.Presentaciones[strings.ToLower(strings.TrimSpace(fila.IDPresentacionNombre))]
		if !ok {
			return "error", fmt.Errorf("presentación '%s' no encontrada", fila.IDPresentacionNombre)
		}
		presID = &id
	}

	var formaID *int
	if fila.IDFormaNombre != "" {
		id, ok := cache.FormasFarma[strings.ToLower(strings.TrimSpace(fila.IDFormaNombre))]
		if !ok {
			return "error", fmt.Errorf("forma farmacéutica '%s' no encontrada", fila.IDFormaNombre)
		}
		formaID = &id
	}

	// ── 2. Resolver estado ────────────────────────────────────────────────────

	isActive := true
	estadoNorm := strings.ToLower(strings.TrimSpace(fila.Estado))
	if estadoNorm == "inactivo" || estadoNorm == "false" || estadoNorm == "0" {
		isActive = false
	}

	// ── 3. Factor por defecto ─────────────────────────────────────────────────

	factor := fila.FactorBase
	if factor == 0 {
		factor = 1
	}

	// ── 4. ¿Existe el código en esta empresa? ─────────────────────────────────

	existingID, found, err := FindArticuloByCodigoEmpresa(db, strings.TrimSpace(fila.Codigo), empresaID)
	if err != nil {
		return "error", err
	}

	now := time.Now()

	if found {
		// ── UPDATE ────────────────────────────────────────────────────────────
		data := map[string]interface{}{
			"nombre":                strings.TrimSpace(fila.Nombre),
			"descripcion":           strings.TrimSpace(fila.Descripcion),
			"id_tipo_articulo":      grupoID,
			"id_unidad_medida":      unidadID,
			"factor_unidad_base":    factor,
			"id_forma_farmaceutica": formaID,
			"id_presentacion":       presID,
			"is_active":             isActive,
			"updated_at":            now,
			"updated_by":            userID,
		}

		tx := db.Table("inventario.cfg_articulos").
			Where("id = ?", existingID).
			Updates(data)
		if tx.Error != nil {
			return "error", tx.Error
		}
		return "actualizado", nil
	}

	// ── INSERT ────────────────────────────────────────────────────────────────
	data := map[string]interface{}{
		"codigo":                strings.TrimSpace(fila.Codigo),
		"nombre":                strings.TrimSpace(fila.Nombre),
		"descripcion":           strings.TrimSpace(fila.Descripcion),
		"id_tipo_articulo":      grupoID,
		"id_unidad_medida":      unidadID,
		"factor_unidad_base":    factor,
		"id_forma_farmaceutica": formaID,
		"id_presentacion":       presID,
		"is_active":             isActive,
		"created_at":            now,
		"created_by":            userID,
		"id_empresa":            empresaID,
		"id_sede":               sedeID,
	}

	tx := db.Table("inventario.cfg_articulos").Create(&data)
	if tx.Error != nil {
		return "error", tx.Error
	}

	return "creado", nil
}
