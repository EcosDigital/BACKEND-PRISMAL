package onboarding

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ecosistema/core/src/database"
	empresa "github.com/ecosistema/core/src/modules/_core/gestion_empresa"
	"github.com/ecosistema/core/src/modules/_core/terceros"
	"github.com/ecosistema/core/src/modules/_core/usuarios"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func BoolPtr(b bool) *bool {
	return &b
}

func RegisterTenantOnboarding(db *gorm.DB, req *TenantRequest) (int64, error) {
	// 1 validar existencia por NIT
	exists, err := terceros.FindTerceroByDocument(db, req.Empresa.Nit)
	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		return 0, errors.New("ya existe una entidad con este documento")
	}

	// 2 crear tercero cliente DBADMIN
	newID, err := CreateTerceroTenant(&req.Empresa)
	if err != nil {
		return 0, err
	}

	ClaseTercero := []int{2}
	if err := terceros.AddClasesTercero(db, newID, ClaseTercero, 1); err != nil {
		return 0, err
	}

	// 3 Crear Slug y TENAT
	baseSlug := slug.Make(strings.ToLower(req.Empresa.RazonSocial))
	slugTenant := fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix())

	tnantsID, err := CreateTenant(newID, slugTenant, req.Empresa.Subdominio)
	if err != nil {
		return 0, err
	}

	// 4 crear la licencia
	licenceID, codigo, err := CreateLicencia(tnantsID, newID)
	if err != nil {
		return 0, err
	}

	// 5 CREAR BASE DE DATOS TENANT
	words := strings.Fields(req.Empresa.RazonSocial)
	dbSlug := strings.ToLower(words[0])
	dbName := fmt.Sprintf("prismar_%s", dbSlug)

	if err := CreateTenantDB(dbSlug); err != nil {
		return 0, err
	}

	// 6 AGREGAR MODULOS BASE A LA LICENCIA
	if err := AddBaseModulesToLicence(licenceID); err != nil {
		return 0, fmt.Errorf("error agregando módulos base a licencia: %v", err)
	}

	// 6.1 INSTALAR MODULOS SELECCIONADOS POR EL CLIENTE
	if len(req.Modulos) > 0 {
		if err := InstallModulesByCode(dbName, licenceID, req.Modulos); err != nil {
			return 0, fmt.Errorf("error instalando módulos: %v", err)
		}
	}

	// 6.2 CAMBIAR LA CONEXION A BD DEL TENANT
	tenantDB, err := ConnectToTenantDB(dbName)
	if err != nil {
		return 0, fmt.Errorf("error conectando a BD tenant: %v", err)
	}

	// Guardar la conexión original
	originalDB := database.GormDB

	// Cambiar temporalmente a BD del tenant
	database.GormDB = tenantDB

	defer func() {
		database.GormDB = originalDB
	}()

	// 7 CREAR EMPRESA EN BASE DE DATOS DEL TENANT

	EmpresaRequest := empresa.EmpresaRequest{
		IdTipoEmpresa:          1,
		Nit:                    req.Empresa.Nit,
		Dv:                     req.Empresa.Dv,
		RazonSocial:            req.Empresa.RazonSocial,
		Descripcion:            "",
		Direccion:              "",
		IdPais:                 1,
		IdDepartamento:         req.Empresa.IdDepartamento,
		Id_ciudad:              req.Empresa.IdCiudad,
		IdZona:                 1,
		Telefono:               req.Empresa.Telefono,
		Telefono_2:             "",
		Email:                  req.Admin.Email,
		Fax:                    "",
		PaginaWeb:              "",
		LogoUrl:                "",
		UsaSedes:               BoolPtr(false),
		NumeroSedes:            0,
		IdNaturaleza:           1,
		IdActividadEconomica:   1,
		IdTipoContribuyente:    1,
		IdRegimenTributario:    1,
		IdTipoDocRepresentante: 2,
		NumeroDocumento:        "00001",
		RepresentanteLegal:     req.Empresa.RepresentanteLegal,
		CodigoLicencia:         codigo,
		Estado:                 BoolPtr(true),
	}

	empresaID, err := empresa.CreateEmpresa(database.GormDB, &EmpresaRequest)
	if err != nil {
		return 0, fmt.Errorf("error creando empresa tenant: %v", err)
	}

	// 8 CREAR TERCERO ADMIN EN BD DEL TENANT

	ClaseTerceroTenant := []int{12}

	TerceroRequest := terceros.TerceroRequest{
		IdTipoPersona:      1,
		IdTipoDocumento:    2,
		NumeroDocumento:    req.Admin.NumeroDocumento,
		Dv:                 "",
		PrimerNombre:       req.Admin.PrimerNombre,
		SegundoNombre:      req.Admin.SegundoNombre,
		PrimerApellido:     req.Admin.PrimerApellido,
		SegundoApellido:    req.Admin.SegundoApellido,
		RazonSocial:        "",
		RepresentanteLegal: "",
		IdGenero:           3,
		Telefono:           req.Empresa.Telefono,
		Telefono_2:         "",
		Email:              req.Admin.Email,
		PaginaWeb:          "",
		Direccion:          "",
		IdPais:             1,
		IdDepartamento:     req.Empresa.IdDepartamento,
		IdCiudad:           req.Empresa.IdCiudad,
		IdZona:             1,
		Estado:             BoolPtr(true),
		ClaseTercero:       ClaseTerceroTenant,
		EmpresaID:          empresaID,
	}

	terceroID, err := terceros.RegisterTercero(database.GormDB, &TerceroRequest)
	if err != nil {
		return 0, fmt.Errorf("error creando tercero admin: %v", err)
	}

	// 9 CREAR USUARIO ADMIN EN BD DEL TENANT
	UserRequest := usuarios.UserRequest{
		Email:        req.Admin.Email,
		Password:     req.Admin.Password,
		IdTercero:    int(terceroID),
		IdRol:        1,
		ImageProfile: "",
		Activo:       BoolPtr(true),
	}

	userID, err := usuarios.CreateUser(database.GormDB, &UserRequest)
	if err != nil {
		return 0, fmt.Errorf("error creando usuario admin: %v", err)
	}

	return userID, err

}

func FilterModulosProduccion(db *gorm.DB) ([]ModuleCatalogCategory, error) {
	return ListModulosActiveProduccion(db)
}
