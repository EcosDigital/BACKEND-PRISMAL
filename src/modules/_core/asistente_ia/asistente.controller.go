package asistente

import (
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func ChatController(c *fiber.Ctx) error {

	var req ChatRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	resultado, err := AskCopilot(db, empresaID, req.Pregunta)
	if err != nil {
		logging.Error.Printf("error en chat del asistente de IA: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resultado)
}
