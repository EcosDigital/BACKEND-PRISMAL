package notificaciones

import (
	"net/http"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// consultar los eventos configurables del tenant, agrupados por modulo
func FindEventosConfigurablesController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := FilterEventosConfigurables(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// consultar los usuarios activos con su rol
func FindUsuariosDestinoDetalleController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := FilterUsuariosDestinoDetalle(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// consultar a quien le llega un evento en la sede del usuario conectado
func FindDestinatariosEventoController(c *fiber.Ctx) error {

	codigo := c.Params("codigo")
	idSede := int64(middlewares.GetSedeID(c))

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := FilterDestinatariosEvento(db, codigo, idSede)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// guardar a quien le llega un evento en la sede del usuario conectado
func SaveDestinatariosEventoController(c *fiber.Ctx) error {

	codigo := c.Params("codigo")

	var req DestinatariosEventoRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := SaveDestinatariosEvento(db, codigo, &req); err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Destinatarios guardados"})

}
