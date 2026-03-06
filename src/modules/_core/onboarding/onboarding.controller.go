package onboarding

import (
	"net/http"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateTenantController(c *fiber.Ctx) error {

	var req TenantRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//invocar services (logica de negocio)
	newID, err := RegisterTenantOnboarding(db, &req)
	if err != nil {

		logging.Error.Printf("Hubo un error: %v", err)

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Registro exitoso...",
		"id":      newID,
	})

}

func FindModuloByProduccionController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	modules, err := FilterModulosProduccion(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(modules)
}

func CheckDomainController(c *fiber.Ctx) error {
	dominio := c.Query("dominio")
	if dominio == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "dominio requerido",
		})
	}

	available, err := CheckDomain(dominio)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "error verificando dominio",
		})
	}

	return c.JSON(fiber.Map{
		"dominio":   dominio,
		"available": available,
	})
}
