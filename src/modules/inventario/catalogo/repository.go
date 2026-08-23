package catalogo

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// NextCatalogoCodigo genera el siguiente código consecutivo (CAT-0001, CAT-0002...).
func NextCatalogoCodigo(db *gorm.DB) (string, error) {

	var maxNum int

	err := db.Raw(`
		SELECT COALESCE(MAX(CAST(SUBSTRING(codigo FROM 5) AS INT)), 0)
		FROM inventario.cfg_catalogos
		WHERE codigo ~ '^CAT-[0-9]+$'
	`).Scan(&maxNum).Error

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CAT-%04d", maxNum+1), nil
}

// ExistsCatalogoBySede indica si la sede ya tiene un catálogo creado.
func ExistsCatalogoBySede(db *gorm.DB, sedeID int64) (bool, error) {
	var count int64
	err := db.
		Table("inventario.cfg_catalogos").
		Where("id_sede = ?", sedeID).
		Count(&count).Error
	return count > 0, err
}

func CreateCatalogo(db *gorm.DB, req *CatalogoRequest, codigo string) (int64, error) {

	data := map[string]interface{}{
		"codigo":      codigo,
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.IDSede,
		"is_active":   true,
		"created_at":  time.Now(),
		"created_by":  req.UserID,
	}

	tx := db.Table("inventario.cfg_catalogos").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

func ListCatalogos(db *gorm.DB, empresaID int64) ([]CatalogoResponse, error) {

	var results []CatalogoResponse

	err := db.
		Table("inventario.cfg_catalogos c").
		Select(`
			c.id,
			c.codigo,
			c.nombre,
			s.nombre AS sede,
			c.is_active`).
		Joins("INNER JOIN configuracion.cfg_sedes s ON s.id = c.id_sede").
		Where("c.id_empresa = ?", empresaID).
		Order("c.id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []CatalogoResponse{}
	}

	return results, nil
}

func ListCatalogoByID(db *gorm.DB, id int64) (*CatalogoResponseFull, error) {

	var result CatalogoResponseFull

	err := db.
		Table("inventario.cfg_catalogos c").
		Select(`
			c.id,
			c.codigo,
			c.nombre,
			c.descripcion,
			c.id_sede,
			s.nombre AS sede,
			c.is_active`).
		Joins("INNER JOIN configuracion.cfg_sedes s ON s.id = c.id_sede").
		Where("c.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateCatalogo(db *gorm.DB, req *CatalogoUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
		"is_active":   req.IsActive,
		"updated_at":  time.Now(),
		"updated_by":  req.UserID,
	}

	tx := db.
		Table("inventario.cfg_catalogos").
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

// ─── Detalle: artículos publicados dentro del catálogo ──────────────────────

// ExistsCatalogoArticulo indica si el artículo ya está publicado en ese catálogo.
func ExistsCatalogoArticulo(db *gorm.DB, catalogoID, articuloID int64) (bool, error) {
	var count int64
	err := db.
		Table("inventario.cfg_catalogo_articulos").
		Where("id_catalogo = ? AND id_articulo = ?", catalogoID, articuloID).
		Count(&count).Error
	return count > 0, err
}

func CreateCatalogoArticulo(db *gorm.DB, catalogoID int64, req *CatalogoArticuloRequest) (int64, error) {

	aplicaDomicilio := false
	if req.AplicaDomicilio != nil {
		aplicaDomicilio = *req.AplicaDomicilio
	}

	data := map[string]interface{}{
		"id_catalogo":      catalogoID,
		"id_articulo":      req.IDArticulo,
		"precio_publico":   req.PrecioPublico,
		"aplica_domicilio": aplicaDomicilio,
		"created_at":       time.Now(),
		"created_by":       req.UserID,
	}

	tx := db.Table("inventario.cfg_catalogo_articulos").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

func ListCatalogoArticulos(db *gorm.DB, catalogoID int64) ([]CatalogoArticuloResponse, error) {

	var results []CatalogoArticuloResponse

	err := db.
		Table("inventario.cfg_catalogo_articulos ca").
		Select(`
			ca.id,
			ca.id_articulo,
			a.codigo   AS codigo_articulo,
			a.nombre   AS nombre_articulo,
			COALESCE(a.imagen_url, '') AS imagen_url,
			ca.precio_publico,
			ca.aplica_domicilio,
			ca.sync_ok`).
		Joins("INNER JOIN inventario.cfg_articulos a ON a.id = ca.id_articulo").
		Where("ca.id_catalogo = ?", catalogoID).
		Order("a.nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []CatalogoArticuloResponse{}
	}

	return results, nil
}

func ListCatalogoArticuloByID(db *gorm.DB, id int64) (*CatalogoArticuloResponse, error) {

	var result CatalogoArticuloResponse

	err := db.
		Table("inventario.cfg_catalogo_articulos ca").
		Select(`
			ca.id,
			ca.id_articulo,
			a.codigo   AS codigo_articulo,
			a.nombre   AS nombre_articulo,
			COALESCE(a.imagen_url, '') AS imagen_url,
			ca.precio_publico,
			ca.aplica_domicilio,
			ca.sync_ok`).
		Joins("INNER JOIN inventario.cfg_articulos a ON a.id = ca.id_articulo").
		Where("ca.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateCatalogoArticulo(db *gorm.DB, req *CatalogoArticuloUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"precio_publico":   req.PrecioPublico,
		"aplica_domicilio": req.AplicaDomicilio,
		"updated_at":       time.Now(),
		"updated_by":       req.UserID,
	}

	tx := db.
		Table("inventario.cfg_catalogo_articulos").
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

func DeleteCatalogoArticulo(db *gorm.DB, id int64) error {

	tx := db.Exec("DELETE FROM inventario.cfg_catalogo_articulos WHERE id = ?", id)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetIDCatalogoByArticuloID resuelve a qué catálogo pertenece una fila de detalle —
// se necesita para disparar la sincronización con la base central tras editar/borrar.
func GetIDCatalogoByArticuloID(db *gorm.DB, id int64) (int64, error) {

	var idCatalogo int64

	err := db.
		Table("inventario.cfg_catalogo_articulos").
		Select("id_catalogo").
		Where("id = ?", id).
		Row().
		Scan(&idCatalogo)

	if err != nil {
		return 0, err
	}

	return idCatalogo, nil
}

// SetSyncOkForCatalogo marca el estado de sincronización con la base central
// para todos los artículos de un catálogo (la sincronización recalcula el
// catálogo completo, así que el estado aplica a todas sus filas por igual).
func SetSyncOkForCatalogo(db *gorm.DB, idCatalogo int64, ok bool) error {
	return db.
		Table("inventario.cfg_catalogo_articulos").
		Where("id_catalogo = ?", idCatalogo).
		Update("sync_ok", ok).Error
}
