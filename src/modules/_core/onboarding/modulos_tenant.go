package onboarding

// Instalación de módulos en un tenant YA EXISTENTE (fuera del flujo de
// onboarding). Reutiliza la misma maquinaria que usa el registro de un
// tenant nuevo (InstallModule) — la única pieza nueva es resolver la
// licencia del tenant actual a partir del slug de la request.

import (
	"fmt"

	"github.com/ecosistema/core/src/database"
	modular "github.com/ecosistema/core/src/modules/_core/gestion_modular"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// GetModulosDisponibles retorna los módulos del catálogo central que el
// tenant (identificado por su slug) todavía NO tiene instalados en su licencia.
func GetModulosDisponibles(tenantSlug string) ([]modular.ModuleResponse, error) {

	_, licenciaID, err := FindLicenciaActivaByTenantSlug(tenantSlug)
	if err != nil {
		return nil, err
	}

	var codigoLicencia string
	if err := database.GormDB.
		Table("configuracion.cfg_licencias").
		Select("codigo_licencia").
		Where("id = ?", licenciaID).
		Scan(&codigoLicencia).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo código de licencia: %v", err)
	}

	todos, err := modular.FilterAllModule(database.GormDB)
	if err != nil {
		return nil, err
	}

	instalados, err := modular.FilterModulesByCodeLicence(codigoLicencia)
	if err != nil {
		return nil, err
	}

	instaladosSet := make(map[string]bool, len(instalados))
	for _, m := range instalados {
		instaladosSet[m.Codigo] = true
	}

	disponibles := make([]modular.ModuleResponse, 0)
	for _, m := range todos {
		if !instaladosSet[m.Codigo] {
			disponibles = append(disponibles, m)
		}
	}

	return disponibles, nil
}

// InstallModuloEnTenant instala un módulo (por código) en el tenant actual:
// lo agrega a su licencia y corre las migraciones del módulo contra su BD.
func InstallModuloEnTenant(tenantSlug, codigo string) error {

	_, licenciaID, err := FindLicenciaActivaByTenantSlug(tenantSlug)
	if err != nil {
		return err
	}

	dbName := fmt.Sprintf("prismar_%s", tenantSlug)

	return InstallModule(dbName, licenciaID, codigo)
}

// ─── Controllers ────────────────────────────────────────────────────────────

// FindModulosDisponiblesController GET /core/modulos/disponibles
func FindModulosDisponiblesController(c *fiber.Ctx) error {

	tenantSlug := middlewares.GetTenantSlug(c)
	if tenantSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tenant no identificado"})
	}

	results, err := GetModulosDisponibles(tenantSlug)
	if err != nil {
		logging.Error.Printf("FindModulosDisponibles error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// InstallModuloTenantController POST /core/modulos/:codigo/instalar
func InstallModuloTenantController(c *fiber.Ctx) error {

	codigo := c.Params("codigo")
	if codigo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "código inválido"})
	}

	tenantSlug := middlewares.GetTenantSlug(c)
	if tenantSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tenant no identificado"})
	}

	if err := InstallModuloEnTenant(tenantSlug, codigo); err != nil {
		logging.Error.Printf("InstallModuloTenant error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Módulo instalado exitosamente",
		"codigo":  codigo,
	})
}
