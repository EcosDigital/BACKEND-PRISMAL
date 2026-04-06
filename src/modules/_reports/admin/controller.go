package reports

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// userIDStr returns the authenticated user ID as string (used as created_by).
func userIDStr(c *fiber.Ctx) string {
	return strconv.FormatInt(int64(middlewares.GetUserID(c)), 10)
}

// POST /report-server-connections
func CreateConnectionController(c *fiber.Ctx) error {

	var req ReportServerConnectionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	newID, err := RegisterConnection(db, &req, userIDStr(c))
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Conexión registrada exitosamente",
		"id":      newID,
	})
}

// GET /report-server-connections
func FindConnectionsController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterConnections(db, userIDStr(c))
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// GET /report-server-connections/:id
func FindConnectionByIDController(c *fiber.Ctx) error {

	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := FilterConnectionByID(db, id, userIDStr(c))
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result == nil || result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conexión no encontrada"})
	}

	return c.JSON(result)
}

// PUT /report-server-connections/:id
func ChangeConnectionController(c *fiber.Ctx) error {

	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req ReportServerConnectionUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	uptID, err := EditConnection(db, id, &req, userIDStr(c))
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Conexión actualizada exitosamente",
		"id":      uptID,
	})
}

// DELETE /report-server-connections/:id
func DeleteConnectionController(c *fiber.Ctx) error {

	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	delID, err := RemoveConnection(db, id, userIDStr(c))
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Conexión eliminada exitosamente",
		"id":      delID,
	})
}

// PATCH /report-server-connections/:id/activate
func ActivateConnectionController(c *fiber.Ctx) error {

	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	actID, err := ActivateConnection(db, id, userIDStr(c))
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Conexión activada exitosamente",
		"id":      actID,
	})
}

// ── helpers ───────────────────────────────────────────────────────────────────

func parseID(c *fiber.Ctx) (int64, error) {
	return strconv.ParseInt(c.Params("id"), 10, 64)
}
