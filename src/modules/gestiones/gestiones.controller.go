package gestiones

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateTicketController  POST /ticket
func CreateTicketController(c *fiber.Ctx) error {

	var req TicketRequest

	// Parsear body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	// Inyectar usuario desde el JWT
	req.UserID = int64(middlewares.GetUserID(c))

	// Obtener conexión a la base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Invocar servicio
	newID, err := RegisterTicket(db, &req)
	if err != nil {
		logging.Error.Printf("CreateTicketController error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Ticket creado exitosamente",
		"id":      newID,
	})
}

func FindTicketsController(c *fiber.Ctx) error {

	// Leer filtros desde query params
	var filtros TicketFiltros
	if err := c.QueryParser(&filtros); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Parámetros de filtro inválidos",
		})
	}

	// UserID del JWT para el filtro solo_asignados
	userID := int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := FilterTickets(db, &filtros, userID)
	if err != nil {
		logging.Error.Printf("FindTicketsController error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}

// FindTicketByIDController  GET /ticket/:id
func FindTicketByIDController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido",
		})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := FilterTicketByID(db, id)
	if err != nil {
		logging.Error.Printf("FindTicketByIDController error: %v", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}
