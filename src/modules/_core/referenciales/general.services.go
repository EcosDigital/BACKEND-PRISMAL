package referenciales

import (
	"gorm.io/gorm"
)

func GetTipoPersona(db *gorm.DB) ([]TipoPersonaDB, error) {
	return ListTipoPersona(db)
}

func GetTipoDocumento(db *gorm.DB) ([]TipoDocumentoDB, error) {
	return ListTipoDocument(db)
}

func GetGeneros(db *gorm.DB) ([]GeneroDB, error) {
	return ListGenero(db)
}

func GetDepartamentos(db *gorm.DB) ([]DepartamentoDB, error) {
	return ListDepartament(db)
}

func GetMunicipios(db *gorm.DB, id_dep int) ([]MunicipiosDB, error) {
	return ListMunicipio(db, id_dep)
}

func GetZonaRurales(db *gorm.DB) ([]ZonaRuralesDB, error) {
	return ListZonaRurales(db)
}

func GetActividadesEco(db *gorm.DB) ([]ActividadesEconomicas, error) {
	return ListActividades(db)
}

func GetAmbitosTerceros(db *gorm.DB) ([]AmbitosTerceros, error) {
	return ListAmbitosTerceros(db)
}

func GetCentralizaciones(db *gorm.DB) ([]Centralizaciones, error) {
	return ListCentralizaciones(db)
}

func GetResponsabilidadDian(db *gorm.DB) ([]ResponsabilidadDian, error) {
	return ListResponsabilidadesDian(db)
}

func GetTipoContribuyente(db *gorm.DB) ([]TipoContribuyente, error) {
	return ListTipoContribuyente(db)
}

func GetRegimenIva(db *gorm.DB) ([]RegimenIva, error) {
	return ListRegimenIva(db)
}

func GetRegimenDian(db *gorm.DB) ([]RegimenDian, error) {
	return ListRegimenDian(db)
}

func GetClaseTercero(db *gorm.DB) ([]ClasePersonas, error) {
	return ListClaseTercero(db)
}

func GetEstadoModule() ([]EstadoModulo, error) {
	return ListEstadoModule()
}

func GetTipoEmpresa(db *gorm.DB) ([]TipoEmpresa, error) {
	return ListTipoEmpresa(db)
}

func GetPaises(db *gorm.DB) ([]Paises, error) {
	return ListPaises(db)
}

func GetNaturalezaEmpresa(db *gorm.DB) ([]NaturalezaEmpresa, error) {
	return ListNaturalezaEmpresa(db)
}

func GetTipoSede(db *gorm.DB) ([]TipoSede, error) {
	return ListTipoSedes(db)
}

func GetTipoRol(db *gorm.DB) ([]TipoRol, error) {
	return ListTipoRoles(db)
}

// === comprobantes ==== ///
func GetModulosComprobantes(db *gorm.DB) ([]ModuleComprobantes, error) {
	return ListModuleComprobantes(db)
}

func GetTipoOperacion(db *gorm.DB, id int64) ([]TipoOperacion, error) {
	return ListTipoOperaciones(db, id)
}

// === GESTIONES === //
func GetNivelesCaso(db *gorm.DB) ([]PrioridadTicket, error) {
	return ListPrioridadTicket(db)
}

func GetEstadosCaso(db *gorm.DB) ([]EstadoTicket, error) {
	return ListEstadosTicket(db)
}

func GetColaboradores(db *gorm.DB) ([]Colaboradores, error) {
	return ListColaboradores(db)
}

// ==== BODEGAS === ///
func GetTipoBodega(db *gorm.DB) ([]TipoBodega, error) {
	return ListTipoBodega(db)
}

func GetUnidadesMedida(db *gorm.DB) ([]UnidadeMedida, error) {
	return ListUnidadesMedida(db)
}

func GetGrupoArticulos(db *gorm.DB) ([]GrupoArticulos, error) {
	return ListGrupoArticulos(db)
}

func GetPresentacionArticulos(db *gorm.DB) ([]Presentacionrticulos, error) {
	return ListPresentacionArticulos(db)
}

// ===== CENTROS DE COSTO ======= //
func GetAreaCosto(db *gorm.DB) ([]ResultGeneral, error) {
	return ListAreaCosto(db)
}

func GetUnidadFuncional(db *gorm.DB) ([]ResultGeneral, error) {
	return ListUnidadFuncional(db)
}
