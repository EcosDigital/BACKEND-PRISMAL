package onboarding

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/database"
	"gorm.io/gorm"
)

func CreateTerceroTenant(req *EmpresaPayload) (int64, error) {

	data := map[string]interface{}{
		"id_tipo_persona":     2,
		"id_tipo_documento":   5,
		"numero_documento":    req.Nit,
		"dv":                  req.Dv,
		"razon_social":        req.RazonSocial,
		"representante_legal": req.RepresentanteLegal,
		"id_genero":           3,
		"telefono":            req.Telefono,
		"email":               req.EmailContacto,
		"direccion":           " ",
		"id_pais":             1,
		"id_departamento":     req.IdDepartamento,
		"id_ciudad":           req.IdCiudad,
		"id_zona":             1,
		"created_at":          time.Now(),
	}

	err := database.GormDB.
		Table("configuracion.cfg_terceros").
		Create(&data).Error

	if err != nil {
		return 0, err
	}

	var id_reg int64
	database.GormDB.Raw("SELECT lastval()").Scan(&id_reg)

	return id_reg, nil

}

func CreateTenant(idTecero int64, slug string, dominio string) (int64, error) {

	data := map[string]interface{}{
		"id_tercero":   idTecero,
		"nombre":       slug,
		"dominio":      dominio,
		"is_active":    true,
		"delivery_app": false,
		"created_at":   time.Now(),
	}

	err := database.GormDB.
		Table("configuracion.cfg_tenants").
		Create(&data).Error

	if err != nil {
		return 0, err
	}

	var id_reg int64
	database.GormDB.Raw("SELECT lastval()").Scan(&id_reg)

	return id_reg, nil

}

func FindLicenceByCode(codigo string) (bool, error) {
	var count int64

	err := database.GormDB.
		Table("configuracion.cfg_licencias").
		Where("codigo_licencia = ?", codigo).
		Count(&count).Error

	return count > 0, err

}

func CreateLicencia(idTenant int64, idCliente int64) (int64, string, error) {

	const prefix = "prismar"
	const randomLem = 8

	var codigoLicencia string

	//generar codigo unico
	for {
		randomPart, err := randomString(randomLem)
		if err != nil {
			return 0, "", err
		}

		codigoLicencia = prefix + randomPart

		exists, err := FindLicenceByCode(codigoLicencia)
		if err != nil {
			return 0, "", err
		}

		if !exists {
			break
		}

	}

	data := map[string]interface{}{
		"codigo_licencia":  codigoLicencia,
		"id_tenant":        idTenant,
		"id_cliente":       idCliente,
		"id_producto":      1,
		"id_tipo_licencia": 1,
		"fecha_inicio":     time.Now(),
		"fecha_fin":        time.Now().AddDate(0, 0, 15),
		"id_estado":        1,
		"created_at":       time.Now(),
	}

	err := database.GormDB.
		Table("configuracion.cfg_licencias").
		Create(&data).Error

	if err != nil {
		return 0, "", err
	}

	var id_reg int64
	database.GormDB.Raw("SELECT lastval()").Scan(&id_reg)

	return id_reg, codigoLicencia, nil

}

func AddModule(licenceID int64, moduleCode string) error {

	moduleInfo, err := GetModuleInfoByCode(moduleCode)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"id_licencia": licenceID,
		"id_modulo":   moduleInfo.ID,
		"instalado":   true,
	}

	tx := database.GormDB.
		Table("configuracion.mov_licencias_modulos").
		Create(&data).Error

	return tx

}

func AddModulesLicence(licenceID int64, modulos []string) error {
	for _, moduleCode := range modulos {
		err := AddModule(licenceID, moduleCode)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDependenciasCodes(moduleID int) ([]string, error) {
	var codes []string

	err := database.GormDB.
		Table("configuracion.cfg_modulos_dependencias d").
		Select("m.codigo").
		Joins("INNER JOIN configuracion.cfg_modulos m ON m.id = d.id_modulo_dependencia").
		Where("d.id_modulo = ? AND m.is_active = ?", moduleID, true).
		Scan(&codes).Error

	if err != nil {
		return nil, fmt.Errorf("error obteniendo dependencias: %v", err)
	}

	return codes, nil
}

func GetModuleInfoByCode(code string) (*ModuleInfo, error) {

	var module ModuleInfo

	err := database.GormDB.
		Table("configuracion.cfg_modulos").
		Select("id, codigo as code, nombre as name, migration_path").
		Where("codigo = ? AND is_active = ?", code, true).
		First(&module).Error

	if err != nil {
		return nil, fmt.Errorf("módulo con código '%s' no encontrado: %v", code, err)
	}

	return &module, nil

}

func IsModuleInstalledInLicence(licenceID int64, moduleID int) (bool, error) {
	var count int64

	err := database.GormDB.
		Table("configuracion.mov_licencias_modulos").
		Where("id_licencia = ? AND id_modulo = ?", licenceID, moduleID).
		Count(&count).Error

	return count > 0, err
}

func ResolveInstallOrder(requestedCodes []string) ([]string, error) {
	visited := map[string]bool{}
	ordered := []string{}

	var resolve func(code string) error
	resolve = func(code string) error {
		if visited[code] {
			return nil
		}
		visited[code] = true

		module, err := GetModuleInfoByCode(code)
		if err != nil {
			return err
		}

		depCodes, err := GetDependenciasCodes(module.ID)
		if err != nil {
			return err
		}

		for _, depCode := range depCodes {
			if err := resolve(depCode); err != nil {
				return err
			}
		}

		ordered = append(ordered, code)
		return nil
	}

	for _, code := range requestedCodes {
		if err := resolve(code); err != nil {
			return nil, err
		}
	}

	return ordered, nil
}

func ExecuteModuleMigrations(dbName string, moduleInfo *ModuleInfo) error {

	if moduleInfo.MigrationPath == "" {
		return fmt.Errorf("módulo '%s' no tiene migration_path configurado", moduleInfo.Name)
	}

	//  Conectar a la BD del tenant
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
		dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error conectando a BD: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error de ping: %v", err)
	}

	// Resolver ruta absoluta igual que main.go hace con uploads
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio de trabajo: %v", err)
	}

	fullPath := filepath.Join(wd, moduleInfo.MigrationPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("directorio no encontrado: %s", fullPath)
	}

	files, err := filepath.Glob(filepath.Join(fullPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("error leyendo archivos: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("⚠️  No hay archivos SQL en: %s\n", fullPath)
		return nil
	}

	sort.Strings(files)
	fmt.Printf("📄 %d archivos en %s\n", len(files), moduleInfo.MigrationPath)

	for _, file := range files {
		fileName := filepath.Base(file)
		sqlContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %v", fileName, err)
		}

		statements := splitSQLStatements(string(sqlContent))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err = db.Exec(stmt); err != nil {
				return fmt.Errorf("error en %s: %v", fileName, err)
			}
		}
		fmt.Printf("  ✅ %s\n", fileName)
	}

	fmt.Printf("✅ Migraciones de '%s' completadas\n", moduleInfo.Name)
	return nil

}

func InstallModule(dbName string, licenceID int64, code string) error {

	moduleInfo, err := GetModuleInfoByCode(code)
	if err != nil {
		return err
	}

	//validar si esta instalado
	already, err := IsModuleInstalledInLicence(licenceID, moduleInfo.ID)
	if err != nil {
		return fmt.Errorf("error verificando módulo en licencia: %v", err)
	}
	if already {
		fmt.Printf("⚠️  Módulo '%s' ya instalado, omitiendo...\n", code)
		return nil
	}

	// Registrar en licencia
	if err := AddModule(licenceID, code); err != nil {
		return fmt.Errorf("error agregando a licencia: %v", err)
	}

	// Ejecutar migraciones solo si tiene ruta configurada
	if moduleInfo.MigrationPath != "" {
		if err := ExecuteModuleMigrations(dbName, moduleInfo); err != nil {
			return fmt.Errorf("error en migraciones de '%s': %v", code, err)
		}
	} else {
		fmt.Printf("ℹ️  Módulo '%s' sin migraciones configuradas\n", code)
	}

	//registra referencia en local tenant
	tenantDB, err := ConnectToTenantDB(dbName)
	if err != nil {
		return fmt.Errorf("error conectando a tenant para registrar ref_modulo: %v", err)
	}

	if err := RegisterModuloRefTenant(tenantDB, moduleInfo); err != nil {
		return fmt.Errorf("error registrando ref_modulo en tenant: %v", err)
	}

	fmt.Printf("✅ Módulo '%s' instalado\n", code)
	return nil

}

func InstallModulesByCode(dbName string, licenceID int64, modules []string) error {

	if len(modules) == 0 {
		fmt.Printf("⚠️  No hay módulos para instalar\n")
		return nil
	}

	orderedModules, err := ResolveInstallOrder(modules)
	if err != nil {
		return fmt.Errorf("error resolviendo dependencias: %v", err)
	}

	fmt.Printf("\n📦 Plan de instalación: %v\n", orderedModules)

	for i, moduleCode := range orderedModules {
		fmt.Printf("\n%s\n[%d/%d] Instalando: %s\n",
			strings.Repeat("=", 60), i+1, len(orderedModules), moduleCode)

		if err := InstallModule(dbName, licenceID, moduleCode); err != nil {
			return fmt.Errorf("error instalando '%s': %v", moduleCode, err)
		}

		fmt.Printf("%s\n", strings.Repeat("=", 60))
	}

	return nil

}

func RegisterModuloRefTenant(tenantDB *gorm.DB, moduleInfo *ModuleInfo) error {

	return tenantDB.Exec(`
        INSERT INTO configuracion.ref_modulos_tenant (id_ref, codigo, nombre)
        VALUES (?, ?, ?)
        ON CONFLICT (id_ref) DO NOTHING
    `, moduleInfo.ID, moduleInfo.Code, moduleInfo.Name).Error

}

func ListModulosActiveProduccion(db *gorm.DB) ([]ModuleCatalogCategory, error) {

	var categories []ModuleCatalogCategory

	err := db.
		Table("configuracion.cfg_categorias").
		Select(`
			id AS id_categoria,
			nombre AS nombre_categoria,
			descripcion as descripcion_categoria
		`).
		Scan(&categories).
		Error

	if err != nil {
		return nil, err
	}

	result := make([]ModuleCatalogCategory, 0)

	for _, category := range categories {

		var modules []ModuleCatalogItem

		err := database.GormDB.
			Table("configuracion.cfg_modulos").
			Select(`
				id,
				nombre,
				codigo,
				descripcion,
				icono,
				color,
				bg_color,
				border_color
			`).
			Where(
				"id_categoria = ? AND is_active = ? AND id_estado = ? AND es_interno = ?",
				category.CategoryID,
				true,
				3,
				false,
			).
			Scan(&modules).
			Error

		if err != nil {
			return nil, err
		}

		if len(modules) > 0 {
			category.Modules = modules
			result = append(result, category)
		}

	}

	return result, nil

}

func AddBaseModulesToLicence(licenceID int64, dbName string) error {
	baseCodes := []string{"MD-001", "MD-002"}

	tenantDB, err := ConnectToTenantDB(dbName)
	if err != nil {
		return fmt.Errorf("error conectando a tenant: %v", err)
	}

	for _, code := range baseCodes {
		if err := AddModule(licenceID, code); err != nil {
			return err
		}

		moduleInfo, err := GetModuleInfoByCode(code)
		if err != nil {
			return err
		}

		if err := RegisterModuloRefTenant(tenantDB, moduleInfo); err != nil {
			return fmt.Errorf("error registrando ref base '%s': %v", code, err)
		}
	}

	return nil
}

func CheckDomainAvailable(dominio string) (bool, error) {
	var count int64
	err := database.GormDB.
		Table("configuracion.cfg_tenants").
		Where("dominio = ?", strings.ToLower(strings.TrimSpace(dominio))).
		Count(&count).Error
	return count == 0, err
}

// FindLicenciaActivaByTenantSlug resuelve el tenant y su licencia más reciente
// a partir del slug usado en el header X-Tenant-ID (el mismo valor con el que
// se arma el nombre de la BD del tenant: "prismar_"+slug). El slug coincide
// directamente con cfg_tenants.dominio, tal como se guarda en el onboarding
// (CreateTenant recibe req.Empresa.Subdominio sin transformar).
func FindLicenciaActivaByTenantSlug(tenantSlug string) (int64, int64, error) {

	var tenant struct {
		ID int64
	}

	err := database.GormDB.
		Table("configuracion.cfg_tenants").
		Select("id").
		Where("LOWER(dominio) = ?", strings.ToLower(strings.TrimSpace(tenantSlug))).
		First(&tenant).Error
	if err != nil {
		return 0, 0, fmt.Errorf("tenant no encontrado para '%s': %v", tenantSlug, err)
	}

	var licencia struct {
		ID int64
	}

	err = database.GormDB.
		Table("configuracion.cfg_licencias").
		Select("id").
		Where("id_tenant = ?", tenant.ID).
		Order("created_at DESC").
		Limit(1).
		First(&licencia).Error
	if err != nil {
		return 0, 0, fmt.Errorf("licencia no encontrada para el tenant '%s': %v", tenantSlug, err)
	}

	return tenant.ID, licencia.ID, nil
}
