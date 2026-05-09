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

func CreateGestionController(c *fiber.Ctx) error {

	ticketID, err := parseID(c, "id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID de ticket inválido"})
	}

	var req GestionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	newID, err := RegisterGestion(db, ticketID, &req)
	if err != nil {
		logging.Error.Printf("CreateGestionController error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Gestión registrada exitosamente",
		"id":      newID,
	})
}

func FindGestionesController(c *fiber.Ctx) error {

	ticketID, err := parseID(c, "id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID de ticket inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterGestiones(db, ticketID)
	if err != nil {
		logging.Error.Printf("FindGestionesController error: %v", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// ─── Asignación ───────────────────────────────────────────────────────────────
// AssignCollaboratorsController  POST /ticket/:id/asignar
func AssignCollaboratorsController(c *fiber.Ctx) error {

	ticketID, err := parseID(c, "id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID de ticket inválido"})
	}

	var req AsignacionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := AssignCollaborators(db, ticketID, &req); err != nil {
		logging.Error.Printf("AssignCollaboratorsController error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Colaboradores asignados exitosamente",
	})

}

// RemoveCollaboratorController  DELETE /ticket/:id/asignar/:colab_id
func RemoveCollaboratorController(c *fiber.Ctx) error {

	ticketID, err := parseID(c, "id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID de ticket inválido"})
	}

	colabID, err := parseID(c, "colab_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID de colaborador inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := RemoveCollaborator(db, ticketID, colabID); err != nil {
		logging.Error.Printf("RemoveCollaboratorController error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Colaborador removido exitosamente"})
}

// FindAsignadosController  GET /ticket/:id/asignados
func FindAsignadosController(c *fiber.Ctx) error {

	ticketID, err := parseID(c, "id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID de ticket inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterAsignados(db, ticketID)
	if err != nil {
		logging.Error.Printf("FindAsignadosController error: %v", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// ─── Estadísticas globales ────────────────────────────────────────────────────

// GetStatsController  GET /stats
// Devuelve contadores globales por estado, independiente de cualquier filtro.
func GetStatsController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	stats, err := GetGlobalStats(db)
	if err != nil {
		logging.Error.Printf("GetStatsController error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stats)
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func parseID(c *fiber.Ctx, param string) (int64, error) {
	return strconv.ParseInt(c.Params(param), 10, 64)
}
