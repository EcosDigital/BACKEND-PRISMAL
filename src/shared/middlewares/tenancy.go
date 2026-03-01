package middlewares

import (
	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/database"
	"github.com/gofiber/fiber/v2"
)

func TenancyMiddleware(c *fiber.Ctx) error {

	if core.Cfg.TenancyMode != "multi" {
		// single tenant → usa conexión global
		c.Locals("db", database.GormDB)
		return c.Next()
	}

	tenant := c.Get("X-Tenant-ID")
	if tenant == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Tenant inválido",
		})
	}

	db, err := database.GetGormDBByTenant(tenant)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "No se pudo conectar a la base de datos del tenant",
		})
	}

	c.Locals("tenant", tenant)
	c.Locals("db", db)

	return c.Next()
}
