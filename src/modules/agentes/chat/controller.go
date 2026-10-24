package chat

import (
	"net/http"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// enviar mensaje al agente y devolver su respuesta
func EnviarMensajeController(c *fiber.Ctx) error {
	var req ChatRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))
	req.TenantSlug = middlewares.GetTenantSlug(c)

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//invocar services
	resultado, err := EnviarMensaje(db, &req)
	if err != nil {

		logging.Error.Printf("error en chat del agente: %v", err)

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resultado)
}
