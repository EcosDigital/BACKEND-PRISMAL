package empresa

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func ListEmpresaByNit(db *gorm.DB, nit string) ([]EmpresaResponse, error) {

	var results []EmpresaResponse

	err := db.
		Table("configuracion.cfg_empresas").
		Select(`
			id,
			nit,
			dv,
			razon_social,
			usa_sedes,
			is_active as estado`).
		Where("nit = ?", nit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

func CreateEmpresa(db *gorm.DB, req *EmpresaRequest) (int64, error) {

	data := map[string]interface{}{
		"id_tipo_empresa":        req.IdTipoEmpresa,
		"nit":                    req.Nit,
		"dv":                     req.Dv,
		"razon_social":           req.RazonSocial,
		"descripcion":            req.Descripcion,
		"direccion":              req.Direccion,
		"id_pais":                req.IdPais,
		"id_departamento":        req.IdDepartamento,
		"id_ciudad":              req.Id_ciudad,
		"id_zona":                req.IdZona,
		"telefono":               req.Telefono,
		"telefono_2":             req.Telefono_2,
		"email":                  req.Email,
		"fax":                    req.Fax,
		"pagina_web":             req.PaginaWeb,
		"logo_url":               req.LogoUrl,
		"usa_sedes":              req.UsaSedes,
		"numero_sedes":           req.NumeroSedes,
		"id_naturaleza":          req.IdNaturaleza,
		"id_actividad_economica": req.IdActividadEconomica,
		"id_tipo_contribuyente":  req.IdTipoContribuyente,
		"id_regimen_tributario":  req.IdRegimenTributario,
		"tipo_documento_rep":     req.IdTipoDocRepresentante,
		"numero_documento":       req.NumeroDocumento,
		"representante_legal":    req.RepresentanteLegal,
		"codigo_licencia":        req.CodigoLicencia,
		"is_active":              req.Estado,
		"created_by":             req.UserID,
		"created_at":             time.Now(),
	}

	err := db.
		Table("configuracion.cfg_empresas").
		Create(&data)

	if err.Error != nil {
		return 0, err.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil // el insert fue exitoso, pero no se requiere el ID
	}

	return id, nil

}

func ListEmpresaLast(db *gorm.DB) ([]EmpresaResponse, error) {

	var results []EmpresaResponse

	err := db.
		Table("configuracion.cfg_empresas e").
		Select(`
			e.id,
			nit,
			dv,
			razon_social,
			usa_sedes,
			te.nombre as tipo,
			is_active`).
		Joins("INNER JOIN configuracion.ref_tipo_empresa te ON te.id = e.id_tipo_empresa").
		Order("id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func ListEmpresaByID(db *gorm.DB, id int64) ([]EmpresaResponseFull, error) {

	var results []EmpresaResponseFull

	err := db.
		Table("configuracion.cfg_empresas").
		Select(`
			id,
			id_tipo_empresa,
			nit,
			dv,
			razon_social,
			descripcion,
			direccion,
			id_pais,
			id_departamento,
			id_ciudad,
			id_zona,
			telefono,
			telefono_2,
			email,
			fax,
			pagina_web,
			logo_url,
			usa_sedes,
			numero_sedes,
			id_naturaleza,
			id_actividad_economica,
			id_tipo_contribuyente,
			id_regimen_tributario,
			tipo_documento_rep as id_tipo_doc_representante,
			numero_documento,
			representante_legal,
			codigo_licencia,
			is_active as is_active`).
		Where("id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func UpdateEmpresa(db *gorm.DB, id int64, req *EmpresaUpdateRequest) (int64, error) {

	data := map[string]interface{}{
		"id_tipo_empresa":        req.IdTipoEmpresa,
		"nit":                    req.Nit,
		"dv":                     req.Dv,
		"razon_social":           req.RazonSocial,
		"descripcion":            req.Descripcion,
		"direccion":              req.Direccion,
		"id_pais":                req.IdPais,
		"id_departamento":        req.IdDepartamento,
		"id_ciudad":              req.Id_ciudad,
		"id_zona":                req.IdZona,
		"telefono":               req.Telefono,
		"telefono_2":             req.Telefono_2,
		"email":                  req.Email,
		"fax":                    req.Fax,
		"pagina_web":             req.PaginaWeb,
		"logo_url":               req.LogoUrl,
		"usa_sedes":              req.UsaSedes,
		"numero_sedes":           req.NumeroSedes,
		"id_naturaleza":          req.IdNaturaleza,
		"id_actividad_economica": req.IdActividadEconomica,
		"id_tipo_contribuyente":  req.IdTipoContribuyente,
		"id_regimen_tributario":  req.IdRegimenTributario,
		"tipo_documento_rep":     req.IdTipoDocRepresentante,
		"numero_documento":       req.NumeroDocumento,
		"representante_legal":    req.RepresentanteLegal,
		"codigo_licencia":        req.CodigoLicencia,
		"is_active":              req.Estado,
		"update_by":              req.UserID,
		"update_at":              time.Now(),
	}

	err := db.
		Table("configuracion.cfg_empresas").
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

// Funciones Sedes
func ListSedeByCode(db *gorm.DB, codigo string) ([]SedeResponse, error) {

	var results []SedeResponse

	err := db.
		Table("configuracion.cfg_sedes s").
		Select(`
			s.id,
			s.codigo,
			e.razon_social as  empresa,
			t.nombre as tipo,
			s.estado as is_active`).
		Joins("INNER JOIN configuracion.cfg_empresas e ON e.id = s.id_empresa").
		Joins("INNER JOIN configuracion.ref_tipo_sedes t ON t.id = s.id_tipo_sede").
		Where("s.codigo = ?", codigo).
		Order("id ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func ListSedeLast(db *gorm.DB) ([]SedeResponse, error) {
	var results []SedeResponse

	err := db.
		Table("configuracion.cfg_sedes s").
		Select(`
			s.id,
			s.codigo as codigo,
			s.nombre,
			e.razon_social as  empresa,
			t.nombre as tipo,
			s.estado is_active`).
		Joins("INNER JOIN configuracion.cfg_empresas e ON e.id = s.id_empresa").
		Joins("INNER JOIN configuracion.ref_tipo_sedes t ON t.id = s.id_tipo_sede").
		Order("id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err
}

func ListSedeByID(db *gorm.DB, id int64) ([]SedeResponseFull, error) {

	var results []SedeResponseFull

	err := db.
		Table("configuracion.cfg_sedes").
		Where("id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err
}

func CreateSede(db *gorm.DB, req *SedeRequest) (int64, error) {
	data := map[string]interface{}{
		"id_empresa":        req.IdEmpresa,
		"nombre":            req.Nombre,
		"codigo":            req.Codigo,
		"id_tipo_sede":      req.IdTipoSede,
		"direccion":         req.Direccion,
		"id_pais":           req.IdPais,
		"id_departamento":   req.IdDepartamento,
		"id_ciudad":         req.Id_ciudad,
		"id_zona":           req.IdZona,
		"telefono":          req.Telefono,
		"telefono_2":        req.Telefono_2,
		"email":             req.Email,
		"fax":               req.Fax,
		"pagina_web":        req.PaginaWeb,
		"responsable_sede":  req.ResponsableSede,
		"cargo_responsable": req.CargoResponsable,
		"estado":            req.Estado,
		"created_by":        req.UserID,
		"created_at":        time.Now(),
	}

	err := db.
		Table("configuracion.cfg_sedes").
		Create(&data)

	if err.Error != nil {
		return 0, err.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil // el insert fue exitoso, pero no se requiere el ID
	}

	return id, nil
}

func UpdateSede(db *gorm.DB, id int64, req *SedeUpdateRequest) (int64, error) {

	data := map[string]interface{}{
		"id_empresa":        req.IdEmpresa,
		"nombre":            req.Nombre,
		"codigo":            req.Codigo,
		"id_tipo_sede":      req.IdTipoSede,
		"direccion":         req.Direccion,
		"id_pais":           req.IdPais,
		"id_departamento":   req.IdDepartamento,
		"id_ciudad":         req.Id_ciudad,
		"id_zona":           req.IdZona,
		"telefono":          req.Telefono,
		"telefono_2":        req.Telefono_2,
		"email":             req.Email,
		"fax":               req.Fax,
		"pagina_web":        req.PaginaWeb,
		"responsable_sede":  req.ResponsableSede,
		"cargo_responsable": req.CargoResponsable,
		"estado":            req.Estado,
		"update_by":         req.UserID,
		"update_at":         time.Now(),
	}

	err := db.
		Table("configuracion.cfg_sedes").
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
