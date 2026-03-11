package modular

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ecosistema/core/src/database"
	"gorm.io/gorm"
)

/** START OPERACIONES CRUD DE PRODUCTS WEB*/

func ListProductByCode(db *gorm.DB, codigo string) ([]ProductsResponse, error) {

	var prodcts []ProductsResponse

	err := db.
		Table("configuracion.cfg_productos_software").
		Select(`
			id,
			codigo,
			nombre,
			descripcion,
			logo,
			verssion,
			is_active
		`).
		Where("codigo = ?", codigo).
		Scan(&prodcts).Error

	if err != nil {
		return nil, err
	}

	return prodcts, err

}

func CreateProductSoftware(db *gorm.DB, req *ProductRequest) (int64, error) {

	data := map[string]interface{}{
		"codigo":      req.Codigo,
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
		"logo":        req.Logo,
		"verssion":    req.Verssion,
		"is_active":   req.IsActive,
		"created_at":  time.Now(),
		"created_by":  req.UserID,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
	}

	tx := db.
		Table("configuracion.cfg_productos_software").
		Create(&data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil // el insert fue exitoso, pero no se requiere el ID
	}

	return id, nil

}

func ListProductLast(db *gorm.DB) ([]ProductsResponse, error) {

	var Products []ProductsResponse

	err := db.
		Table("configuracion.cfg_productos_software").
		Select(`
			id,
			codigo,
			nombre,
			descripcion,
			logo,
			verssion,
			is_active
		`).
		Order("id ASC").
		Limit(10).Scan(&Products).Error

	if err != nil {
		return nil, err
	}

	if Products == nil {
		Products = []ProductsResponse{}
	}

	return Products, nil

}

func ListProductByID(db *gorm.DB, id int64) ([]ProductsResponse, error) {

	var prodcts []ProductsResponse

	err := db.
		Table("configuracion.cfg_productos_software").
		Select(`
			id,
			codigo,
			nombre,
			descripcion,
			logo,
			verssion,
			is_active
		`).
		Where("id = ?", id).
		Scan(&prodcts).Error

	if err != nil {
		return nil, err
	}

	return prodcts, err

}

func UpdateProduct(db *gorm.DB, req *ProductUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
		"logo":        req.Logo,
		"verssion":    req.Verssion,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
		"update_at":   time.Now(),
		"update_by":   req.UserID,
	}

	if req.IsActive != nil {
		data["is_active"] = *req.IsActive
	}

	tx := db.
		Table("configuracion.cfg_productos_software").
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

/** START OPERACIONES CRUD ACTEGORIAS FOR PRODUCTS WEB*/
func ListCategoryByCode(db *gorm.DB, codigo string) ([]CategoryResponse, error) {

	var results []CategoryResponse

	err := db.
		Table("configuracion.cfg_categorias").
		Select(`
			id,
			id_producto,
			'' as producto,
			codigo,
			nombre, 
			descripcion,
			is_active`).
		Where("codigo = ?", codigo).
		Scan(results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func ListCategoryByIdProduct(db *gorm.DB, id int64) ([]CategoryResponse, error) {
	var results []CategoryResponse

	err := db.
		Table("configuracion.cfg_categorias").
		Select(`
			id,
			id_producto,
			'' as producto,
			codigo,
			nombre, 
			descripcion,
			is_active`).
		Where("id_producto = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err
}

func CreateCategoryForProduct(db *gorm.DB, req *CategoryRequest) (int64, error) {

	data := map[string]interface{}{
		"id_producto": req.IDProduct,
		"codigo":      req.Codigo,
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
		"is_active":   req.IsActive,
		"created_at":  time.Now(),
		"created_by":  req.UserID,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
	}

	tx := db.
		Table("configuracion.cfg_categorias").
		Create(&data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil
	}

	return id, nil

}

func ListCategoryLast(db *gorm.DB) ([]CategoryResponse, error) {

	var results []CategoryResponse

	err := db.
		Table("configuracion.cfg_categorias ct").
		Select(`
			ct.id as id,
			ct.id_producto as id_producto,
			p.nombre as producto,
			ct.codigo as codigo,
			ct.nombre as nombre,
			ct.descripcion as descripcion,
			ct.is_active as is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON ct.id_producto = p.id").
		Order("ct.id DESC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []CategoryResponse{}
	}

	return results, nil

}

func ListCategoryById(db *gorm.DB, id int64) ([]CategoryResponse, error) {

	var results []CategoryResponse

	err := db.
		Table("configuracion.cfg_categorias ct").
		Select(`
			ct.id as id,
			ct.id_producto as id_producto,
			p.nombre as producto,
			ct.codigo as codigo,
			ct.nombre as nombre,
			ct.descripcion as descripcion,
			ct.is_active as is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON ct.id_producto = p.id").
		Where("ct.id = ?", id).
		Order("ct.id DESC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

func UpdateCategory(db *gorm.DB, req *CategoryUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"id_producto": req.IDProduct,
		"nombre":      req.Nombre,
		"descripcion": req.Descripcion,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
		"update_at":   time.Now(),
		"update_by":   req.UserID,
	}

	if req.IDProduct != 0 {
		data["id_producto"] = req.IDProduct
	}

	if req.IsActive != nil {
		data["is_active"] = *req.IsActive
	}

	tx := db.
		Table("configuracion.cfg_categorias").
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

/** START OPERACIONES CRUD MODULOS */
func ListModuleByCode(db *gorm.DB, codigo string) ([]ModuleResponse, error) {
	var results []ModuleResponse

	err := db.
		Table("configuracion.cfg_modulos").
		Select(`
			id,
			id_producto,
			'' as producto,
			id_categoria,
			'' as categoria,
			codigo,
			nombre, 
			descripcion,
			color,
			bg_color,
			border_color,
			is_active`).
		Where("codigo = ?", codigo).
		Scan(results).Error

	if err != nil {
		return nil, err
	}

	return results, err
}

func CreateModule(db *gorm.DB, req *ModuleRequest) (int64, error) {

	data := map[string]interface{}{
		"id_producto":    req.IDProduct,
		"id_categoria":   req.IDCategory,
		"id_estado":      req.IdEstado,
		"codigo":         req.Codigo,
		"nombre":         req.Nombre,
		"descripcion":    req.Descripcion,
		"orden_lista":    req.Ordenista,
		"es_interno":     req.EsInterno,
		"color":          req.Color,
		"bg_color":       req.BgColor,
		"border_color":   req.BorderColor,
		"icono":          req.Icono,
		"migration_path": req.MigratioPath,
		"is_active":      req.IsActive,
		"created_at":     time.Now(),
		"created_by":     req.UserID,
		"id_empresa":     req.EmpresaID,
		"id_sede":        req.SedeID,
	}

	tx := db.
		Table("configuracion.cfg_modulos").
		Create(&data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil

}

func ListModuleLast(db *gorm.DB) ([]ModuleResponse, error) {

	var results []ModuleResponse

	err := db.
		Table("configuracion.cfg_modulos m").
		Select(`
			m.id,
			m.id_producto,
			p.nombre as producto,
			m.id_categoria,
			ct.nombre as categoria,
			m.codigo,
			m.nombre,
			m.descripcion,
			m.is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON m.id_producto = p.id").
		Joins("INNER JOIN configuracion.cfg_categorias ct ON ct.id = m.id_categoria").
		Order("m.id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ModuleResponse{}
	}

	return results, nil

}

func ListModuleAll(db *gorm.DB) ([]ModuleResponse, error) {

	var results []ModuleResponse

	err := db.
		Table("configuracion.cfg_modulos m").
		Select(`
			m.id,
			m.id_producto,
			p.nombre as producto,
			m.id_categoria,
			ct.nombre as categoria,
			m.codigo,
			m.nombre,
			m.descripcion,
			m.is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON m.id_producto = p.id").
		Joins("INNER JOIN configuracion.cfg_categorias ct ON ct.id = m.id_categoria").
		Order("m.orden_lista ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ModuleResponse{}
	}

	return results, nil

}

func GetDependenciasByModulo(db *gorm.DB, idModulo int) ([]DependenciaItem, error) {
	var results []DependenciaItem

	err := db.
		Table("configuracion.cfg_modulos_dependencias d").
		Select(`
            m.id,
            m.codigo,
            m.nombre
        `).
		Joins("INNER JOIN configuracion.cfg_modulos m ON m.id = d.id_modulo_dependencia").
		Where("d.id_modulo = ?", idModulo).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		return []DependenciaItem{}, nil
	}

	return results, nil
}

func ListModuleById(db *gorm.DB, id int64) ([]ModuleResponse, error) {

	var results []ModuleResponse

	err := db.
		Table("configuracion.cfg_modulos m").
		Select(`
			m.id,
			m.id_producto,
			p.nombre as producto,
			m.id_categoria,
			ct.nombre as categoria,
			m.codigo,
			m.id_estado,
			m.nombre,
			m.descripcion,
			m.migration_path,
			m.color,
			m.bg_color,
			m.border_color,
			m.icono,
			m.es_interno,
			m.is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON m.id_producto = p.id").
		Joins("INNER JOIN configuracion.cfg_categorias ct ON m.id_categoria = ct.id").
		Where("m.id = ?", id).
		Order("m.id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Cargar dependencias por cada módulo
	for i := range results {
		deps, err := GetDependenciasByModulo(db, results[i].ID)
		if err != nil {
			return nil, err
		}
		results[i].Dependencias = deps
	}

	return results, nil

}

func ListModuleByIdProduct(db *gorm.DB, id int64) ([]ModuleResponse, error) {

	var results []ModuleResponse

	err := db.
		Table("configuracion.cfg_modulos m").
		Select(`
			m.id,
			m.id_producto,
			p.nombre as producto,
			m.id_categoria,
			ct.nombre as categoria,
			m.codigo,
			m.id_estado,
			m.nombre,
			m.descripcion,
			m.color,
			m.bg_color,
			m.border_color,
			m.icono,
			m.es_interno,
			m.is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON m.id_producto = p.id").
		Joins("INNER JOIN configuracion.cfg_categorias ct ON m.id_categoria = ct.id").
		Where("m.id_producto = ?", id).
		Order("m.id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

func ListModuleByIdRol(db *gorm.DB, id int64) (*AccessReponse, error) {
	var results *AccessReponse

	err := db.
		Table("seguridad.cfg_sedes_roles sr").
		Select(`
			sr.id,
			sr.id_rol,
			sr.json_modules as json_access`).
		Where("sr.id_rol = ?", id).Limit(1).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func AddDependencias(db *gorm.DB, idModulo int64, dependencias []int, id int64) error {

	for _, idDep := range dependencias {

		data := map[string]interface{}{
			"id_modulo":             idModulo,
			"id_modulo_dependencia": idDep,
			"created_by":            id,
			"created_at":            time.Now(),
		}

		err := db.Table("configuracion.cfg_modulos_dependencias").Create(&data).Error
		if err != nil {
			return fmt.Errorf("error insertando dependencia %d: %v", idDep, err)
		}

	}
	return nil
}

func UpdateModule(db *gorm.DB, req *ModuleUpdateRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"id_producto":    req.IDProduct,
		"id_estado":      req.IdEstado,
		"nombre":         req.Nombre,
		"descripcion":    req.Descripcion,
		"orden_lista":    req.Ordenista,
		"es_interno":     req.EsInterno,
		"color":          req.Color,
		"bg_color":       req.BgColor,
		"border_color":   req.BorderColor,
		"icono":          req.Icono,
		"migration_path": req.MigratioPath,
		"id_empresa":     req.EmpresaID,
		"id_sede":        req.SedeID,
		"update_at":      time.Now(),
		"update_by":      req.UserID,
	}

	if req.IDCategory != 0 {
		data["id_categoria"] = req.IDCategory
	}

	if req.IsActive != nil {
		data["is_active"] = *req.IsActive
	}

	tx := db.
		Table("configuracion.cfg_modulos").
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

func SyncDependencias(db *gorm.DB, idModulo int64, dependencias []int, id int64) error {

	//Borrar dependencias actuales
	err := db.Table("configuracion.cfg_modulos_dependencias").
		Where("id_modulo = ?", idModulo).
		Delete(nil).Error
	if err != nil {
		return fmt.Errorf("error eliminando dependencias anteriores: %v", err)
	}

	// Insertar las nuevas si vienen
	if len(dependencias) > 0 {
		if err := AddDependencias(db, idModulo, dependencias, id); err != nil {
			return err
		}
	}

	return nil
}

/** START OPERACIONES CRUD FUNCION */

func CreateFuncion(db *gorm.DB, req *FuncionRequest) (int64, error) {

	data := map[string]interface{}{
		"id_producto": req.IDProduct,
		"id_modulo":   req.IDModulo,
		"nombre":      req.Nombre,
		"orden_lista": req.OrdenLista,
		"ruta_acceso": req.RutaAcceso,
		"is_active":   req.IsActive,
		"created_at":  time.Now(),
		"created_by":  req.UserID,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
	}

	tx := db.
		Table("configuracion.cfg_funciones").
		Create(&data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil
	}

	return id, nil

}

func ListFuncionLast(db *gorm.DB) ([]FuncionResponse, error) {

	var results []FuncionResponse

	err := db.
		Table("configuracion.cfg_funciones f").
		Select(`
			f.id,
			f.id_producto,
			p.nombre as producto,
			f.id_modulo,
			m.nombre as modulo,
			f.nombre,
			f.ruta_acceso,
			f.is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON f.id_producto = p.id").
		Joins("INNER JOIN configuracion.cfg_modulos m ON m.id = f.id_modulo").
		Order("f.id ASC").Limit(30).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []FuncionResponse{}
	}

	return results, nil
}

func ListFuncionByID(db *gorm.DB, id int64) ([]FuncionResponse, error) {

	var results []FuncionResponse

	err := db.
		Table("configuracion.cfg_funciones f").
		Select(`
			f.id,
			f.id_producto,
			p.nombre as producto,
			f.id_modulo,
			m.nombre as modulo,
			f.nombre,
			f.ruta_acceso,
			f.is_active`).
		Joins("INNER JOIN configuracion.cfg_productos_software p ON f.id_producto = p.id").
		Joins("INNER JOIN configuracion.cfg_modulos m ON m.id = f.id_modulo").
		Where("f.id = ?", id).
		Order("f.id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func UpdateFuncion(db *gorm.DB, req *FuncionRequest, id int64) (int64, error) {

	data := map[string]interface{}{
		"id_producto": req.IDProduct,
		"id_modulo":   req.IDModulo,
		"nombre":      req.Nombre,
		"orden_lista": req.OrdenLista,
		"ruta_acceso": req.RutaAcceso,
		"is_active":   req.IsActive,
		"created_at":  time.Now(),
		"created_by":  req.UserID,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
	}

	if req.IsActive != nil {
		data["is_active"] = *req.IsActive
	}

	tx := db.
		Table("configuracion.cfg_funciones").
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

func ListFuncionByIDModule(db *gorm.DB, id int64) ([]FuncionesResponse, error) {

	var funciones []FuncionesResponse

	err := db.
		Table("configuracion.cfg_funciones").
		Select(`
			id,
			id_modulo,
			nombre as title,
			ruta_acceso as path
		`).
		Where("id_modulo = ?", id).
		Order("orden_lista ASC").
		Scan(&funciones).
		Error

	if err != nil {
		return nil, err
	}

	for i := range funciones {
		var subfunciones []FuncionesResponse

		err := database.GormDB.
			Table("configuracion.cfg_subfunciones").
			Select(`
				id,
				nombre as title,
				ruta_acceso as path
			`).
			Where("id_funcion = ?", funciones[i].ID).
			Order("id ASC").
			Scan(&subfunciones).
			Error

		if err != nil {
			return nil, err
		}

		// Solo agregar children si hay subfunciones
		if len(subfunciones) > 0 {
			funciones[i].Children = subfunciones
		}
	}

	return funciones, nil

}

/** START OPERACIONES CRUD SUBFUNCION */

func CreateSubFuncion(db *gorm.DB, req *SubFuncionesRequest) (int64, error) {

	data := map[string]interface{}{
		"id_funcion":  req.IDFuncion,
		"nombre":      req.Nombre,
		"ruta_acceso": req.RutaAcceso,
		"created_at":  time.Now(),
		"created_by":  req.UserID,
		"id_empresa":  req.EmpresaID,
		"id_sede":     req.SedeID,
	}

	tx := db.
		Table("configuracion.cfg_subfunciones").
		Create(&data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil
	}

	return id, nil

}

func ListSubByIDFuncion(db *gorm.DB, id int64) ([]SubFuncionesRequest, error) {

	var results []SubFuncionesRequest

	err := db.
		Table("configuracion.cfg_subfunciones f").
		Select(`
			f.id,
			f.id_funcion,
			f.nombre,
			f.ruta_acceso`).
		Where("f.id_funcion = ?", id).
		Order("f.id ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []SubFuncionesRequest{}, nil
	}

	return results, err

}

/** OPCIONES GENERALES **/
func ListModulesByCodeLicence(codigo string) ([]ModuleResponse, error) {

	var results []ModuleResponse

	err := database.GormDB.
		Table("configuracion.mov_licencias_modulos lm").
		Select(`
			m.id,
			m.id_producto,
			p.nombre as producto,
			m.id_categoria as id_category,
			ct.nombre as categoria,
			m.codigo,
			m.id_estado,
			m.nombre,
			m.descripcion,
			m.color,
			m.bg_color,
			m.border_color,
			m.icono,
			m.is_active,
			m.es_interno`).
		Joins("INNER JOIN configuracion.cfg_licencias l ON l.id = lm.id_licencia").
		Joins("INNER JOIN configuracion.cfg_modulos m ON m.id = lm.id_modulo").
		Joins("INNER JOIN configuracion.cfg_productos_software p ON p.id = m.id_producto").
		Joins("INNER JOIN configuracion.cfg_categorias ct ON ct.id = m.id_categoria").
		Where("l.codigo_licencia = ?", codigo).
		Order("m.orden_lista ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []ModuleResponse{}, nil
	}

	return results, err
}
