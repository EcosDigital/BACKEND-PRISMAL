package middlewares

import (
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/gofiber/fiber/v2"
)

func GetUserID(c *fiber.Ctx) int {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		logging.Error.Printf("Error suario no encontrado")
	}
	return userID
}

func GetEmpresaID(c *fiber.Ctx) int {
	empresaID, ok := c.Locals("empresaID").(int)
	if !ok {
		logging.Error.Printf("Error suario no encontrado")
	}
	return empresaID
}

func GetSedeID(c *fiber.Ctx) int {
	sedeID, ok := c.Locals("sedeID").(int)
	if !ok {
		logging.Error.Printf("Error suario no encontrado")
	}
	return sedeID
}

func GetRolID(c *fiber.Ctx) int {
	sedeID, ok := c.Locals("RolID").(int)
	if !ok {
		logging.Error.Printf("Error suario no encontrado")
	}
	return sedeID
}

// GetTenantSlug retorna el slug del tenant actual (ej. "core"), seteado por
// TenancyMiddleware. Vacío en modo single-tenant — no es un error, ese modo
// no tiene concepto de tenant.
func GetTenantSlug(c *fiber.Ctx) string {
	tenant, ok := c.Locals("tenant").(string)
	if !ok {
		return ""
	}
	return tenant
}
