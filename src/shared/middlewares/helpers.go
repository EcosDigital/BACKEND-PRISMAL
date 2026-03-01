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
