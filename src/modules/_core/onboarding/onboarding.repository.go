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
		"id_tercero": idTecero,
		"nombre":     slug,
		"dominio":    dominio,
		"is_active":  true,
		"created_at": time.Now(),
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

func GetModuleMigrationPath(codigoModulo string) (string, error) {

	//mapeo de rutas
	modulePaths := map[string]string{
		"MD-003": "src/modules/contabilidad/_migrations",
	}

	path, exists := modulePaths[codigoModulo]
	if !exists {
		return "", fmt.Errorf("código de módulo '%s' no tiene ruta configurada", codigoModulo)
	}

	return path, nil

}

func GetModuleInfoByCode(code string) (*ModuleInfo, error) {

	var module ModuleInfo

	err := database.GormDB.
		Table("configuracion.cfg_modulos").
		Select("id, codigo as code, nombre as name").
		Where("codigo = ? AND is_active = ?", code, true).
		First(&module).Error

	if err != nil {
		return nil, fmt.Errorf("módulo con código '%s' no encontrado: %v", code, err)
	}

	return &module, nil

}

func ExecuteModuleMigrations(dbName string, moduleCode string) error {

	//obtener modulo data
	moduleInfo, err := GetModuleInfoByCode(moduleCode)
	if err != nil {
		return err
	}

	//obtener ruta de las migratciones by code
	migrationPath, err := GetModuleMigrationPath(moduleCode)
	if err != nil {
		return fmt.Errorf("error obteniendo ruta: %v", err)
	}

	//  Conectar a la BD del tenant
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error conectando a BD: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error de ping: %v", err)
	}

	//verificar existencia de la carpeta

	fullPath := filepath.Join(".", migrationPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("❌ Directorio no encontrado: %s", fullPath)
	}

	//leer archivos sql
	files, err := filepath.Glob(filepath.Join(fullPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("error leyendo archivos: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("⚠️  No hay archivos SQL en: %s\n", fullPath)
		return nil
	}

	//ordenas y ejecutar
	sort.Strings(files)
	fmt.Printf("📄 Archivos encontrados: %d\n", len(files))

	for _, file := range files {
		fileName := filepath.Base(file)
		fmt.Printf("  ▶️  %s\n", fileName)

		sqlContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %v", fileName, err)
		}

		statements := strings.Split(string(sqlContent), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("❌ Error en %s: %v", fileName, err)
			}
		}

		fmt.Printf("Completado\n")
	}

	fmt.Printf("Migraciones de '%s' ejecutadas\n", moduleInfo.Name)
	return nil

}

func InstallModule(dbName string, licenceID int64, code string) error {

	if err := AddModule(licenceID, code); err != nil {
		return fmt.Errorf("error agregando a licencia: %v", err)
	}

	//ejecutar migraciones BDTENANT
	if err := ExecuteModuleMigrations(dbName, code); err != nil {
		return err
	}

	return nil

}

func InstallModulesByCode(dbName string, licenceID int64, modules []string) error {

	if len(modules) == 0 {
		fmt.Printf("⚠️  No hay módulos para instalar\n")
		return nil
	}

	for i, moduleCode := range modules {
		fmt.Printf("\n" + strings.Repeat("=", 60))
		fmt.Printf("\n[%d/%d] ", i+1, len(modules))

		if err := InstallModule(dbName, licenceID, moduleCode); err != nil {
			return fmt.Errorf("❌ Error instalando '%s': %v", moduleCode, err)
		}

		fmt.Printf(strings.Repeat("=", 60) + "\n")
	}

	return nil

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

func AddBaseModulesToLicence(licenceID int64) error {
	baseCodes := []string{"MD-001", "MD-002"}
	return AddModulesLicence(licenceID, baseCodes)
}
