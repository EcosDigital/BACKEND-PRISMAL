package referenciales

import (
	"github.com/ecosistema/core/src/database"
	"gorm.io/gorm"
)

func ListTipoPersona(db *gorm.DB) ([]TipoPersonaDB, error) {
	var results []TipoPersonaDB

	err := db.
		Table("configuracion.ref_tipo_persona").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoPersonaDB{}
	}

	return results, nil
}

func ListTipoDocument(db *gorm.DB) ([]TipoDocumentoDB, error) {

	var results []TipoDocumentoDB

	err := db.
		Table("configuracion.ref_tipo_documento").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoDocumentoDB{}
	}

	return results, nil
}

func ListGenero(db *gorm.DB) ([]GeneroDB, error) {
	var results []GeneroDB

	err := db.
		Table("configuracion.ref_genero").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []GeneroDB{}
	}

	return results, nil
}

func ListDepartament(db *gorm.DB) ([]DepartamentoDB, error) {
	var results []DepartamentoDB

	err := db.
		Table("configuracion.ref_departamentos").
		Select(`
			id,
			codigo,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []DepartamentoDB{}
	}

	return results, nil
}

func ListMunicipio(db *gorm.DB, id_dep int) ([]MunicipiosDB, error) {

	var results []MunicipiosDB

	err := db.
		Table("configuracion.ref_municipios").
		Select(`
			id,
			codigo_departamento,
			codigo_municipio,
			nombre_municipio`).
		Where(
			"codigo_departamento = (?)",
			db.Table("configuracion.ref_departamentos").
				Select("codigo").
				Where("id = ?", id_dep),
		).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []MunicipiosDB{}
	}

	return results, nil

}

func ListZonaRurales(db *gorm.DB) ([]ZonaRuralesDB, error) {

	var results []ZonaRuralesDB

	err := db.
		Table("configuracion.ref_zona_residencial").
		Select(`
			id,
			codigo,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ZonaRuralesDB{}
	}

	return results, nil

}

func ListActividades(db *gorm.DB) ([]ActividadesEconomicas, error) {

	var results []ActividadesEconomicas

	err := db.
		Table("configuracion.ref_actividades_economicas").
		Select(`
			id,
			codigo,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ActividadesEconomicas{}
	}

	return results, nil
}

func ListAmbitosTerceros(db *gorm.DB) ([]AmbitosTerceros, error) {

	var results []AmbitosTerceros

	err := db.
		Table("configuracion.ref_ambitos_terceros").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []AmbitosTerceros{}
	}

	return results, nil
}

func ListCentralizaciones(db *gorm.DB) ([]Centralizaciones, error) {

	var results []Centralizaciones

	err := db.
		Table("configuracion.ref_centralizacion").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []Centralizaciones{}
	}

	return results, nil

}

func ListResponsabilidadesDian(db *gorm.DB) ([]ResponsabilidadDian, error) {

	var results []ResponsabilidadDian

	err := db.
		Table("configuracion.ref_responsabilidades_dian").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ResponsabilidadDian{}
	}

	return results, nil
}

func ListTipoContribuyente(db *gorm.DB) ([]TipoContribuyente, error) {

	var results []TipoContribuyente

	err := db.
		Table("configuracion.ref_tipo_contribuyente").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoContribuyente{}
	}

	return results, nil
}

func ListRegimenIva(db *gorm.DB) ([]RegimenIva, error) {

	var results []RegimenIva

	err := db.
		Table("configuracion.ref_regimen_iva").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []RegimenIva{}
	}

	return results, nil

}

func ListRegimenDian(db *gorm.DB) ([]RegimenDian, error) {
	var results []RegimenDian

	err := db.
		Table("configuracion.ref_regimen_dian").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []RegimenDian{}
	}

	return results, nil
}

func ListClaseTercero(db *gorm.DB) ([]ClasePersonas, error) {
	var results []ClasePersonas

	err := db.
		Table("configuracion.ref_clase_personas").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ClasePersonas{}
	}

	return results, nil
}

func ListPaises(db *gorm.DB) ([]Paises, error) {

	var results []Paises

	err := db.
		Table("configuracion.ref_paises").
		Select(`
			id,
			nombre`).
		Order("id DESC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []Paises{}
	}

	return results, nil
}

func ListTipoEmpresa(db *gorm.DB) ([]TipoEmpresa, error) {
	var results []TipoEmpresa

	err := db.
		Table("configuracion.ref_tipo_empresa").
		Select(`
			id,
			nombre`).
		Order("id DESC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoEmpresa{}
	}

	return results, nil
}

func ListNaturalezaEmpresa(db *gorm.DB) ([]NaturalezaEmpresa, error) {
	var results []NaturalezaEmpresa

	err := db.
		Table("configuracion.ref_naturaleza_empresa").
		Select(`
			id,
			nombre`).
		Order("id DESC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []NaturalezaEmpresa{}
	}

	return results, nil
}

func ListTipoSedes(db *gorm.DB) ([]TipoSede, error) {
	var results []TipoSede

	err := db.
		Table("configuracion.ref_tipo_sedes").
		Order("id DESC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoSede{}
	}

	return results, nil

}

func ListTipoRoles(db *gorm.DB) ([]TipoRol, error) {

	var results []TipoRol

	err := db.
		Table("seguridad.ref_tipo_rol").
		Select(`
			id,
			nombre`).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoRol{}
	}

	return results, nil

}

func ListEstadoModule() ([]EstadoModulo, error) {

	var results []EstadoModulo

	err := database.GormDB.
		Table("configuracion.ref_estado_modulo").
		Select("id,nombre").
		Order("id ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []EstadoModulo{}
	}

	return results, nil

}

// ===== COMPROBANTES ======= //
func ListModuleComprobantes(db *gorm.DB) ([]ModuleComprobantes, error) {
	var results []ModuleComprobantes

	err := database.GormDB.
		Table("configuracion.ref_modulos_tenant").
		Order("id ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ModuleComprobantes{}
	}

	return results, nil
}

// ===== GESTIONES ======== //
func ListPrioridadTicket(db *gorm.DB) ([]PrioridadTicket, error) {
	var results []PrioridadTicket

	err := db.
		Table("gestiones.cfg_niveles_caso").
		Order("id ASC").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []PrioridadTicket{}
	}

	return results, nil
}
