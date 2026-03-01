package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ecosistema/core/src/shared/logging"
)

var validate = validator.New()

// validateBody recibe un puntero a un struct y valida el body de la request
func VallidateBody(schema interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//parsear el body a schema
		if err := c.BodyParser(schema); err != nil {
			logging.Error.Printf("❌ [VALIDATOR] Error en BodyParser: %v", err)
			logging.Error.Printf("❌ [VALIDATOR] Schema esperado: %T", schema)
			
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "JSON inválidado",
				"detalle": err.Error(),
			})
		}


		//validar el schema con go-playground
		if err := validate.Struct(schema); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Datos inválidos",
				"detalle": err.Error(),
			})
		}

		//Guardar el schema validado en el contexto
		c.Locals("body", schema)

		return c.Next()

	}
}
