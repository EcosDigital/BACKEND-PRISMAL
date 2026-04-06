package horarios

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// Crea un horario con sus bloques en una sola transacción.
func CreateHorarioController(c *fiber.Ctx) error {
	var req HorarioRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := RegisterHorario(db, &req)
	if err != nil {
		logging.Error.Printf("Error creando horario: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Horario creado exitosamente",
		"data":    fiber.Map{"id": id},
	})
}

// Devuelve los horarios más recientes (sin bloques, para listados).
func ListHorariosController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterHorarios(db)
	if err != nil {
		logging.Error.Printf("Error listando horarios: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// Devuelve un horario con todos sus bloques activos.
func GetHorarioByIDController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := FilterHorarioByID(db, id)
	if err != nil {
		logging.Error.Printf("Error obteniendo horario %d: %v", id, err)
		status := http.StatusInternalServerError
		if errors.Is(err, ErrHorarioNoEncontrado) {
			status = http.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// Actualiza solo nombre, descripcion y minutos_tolerancia.
// Los bloques son inmutables y se ignoran aunque vengan en el body.
func UpdateHorarioController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req HorarioUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := UpdateHorario(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error actualizando horario %d: %v", id, err)
		status := http.StatusBadRequest
		if errors.Is(err, ErrHorarioNoEncontrado) {
			status = http.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Horario actualizado exitosamente",
		"data":    result,
	})
}

// Activa o inactiva un horario (soft delete). No elimina físicamente.
func ChangeEstadoHorarioController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CambioEstadoHorarioRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := EditEstadoHorario(db, id, &req); err != nil {
		logging.Error.Printf("Error cambiando estado del horario %d: %v", id, err)
		status := http.StatusBadRequest
		if errors.Is(err, ErrHorarioNoEncontrado) {
			status = http.StatusNotFound
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	estado := "inactivado"
	if req.Activo {
		estado = "activado"
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Horario " + estado + " exitosamente",
	})
}

// ListAsignacionesHorarioController maneja GET /api/control-asistencia/horarios/:id/asignaciones
// Devuelve los empleados actualmente asignados al horario.
func ListAsignacionesHorarioController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterAsignacionesByHorario(db, id)
	if err != nil {
		logging.Error.Printf("Error listando asignaciones del horario %d: %v", id, err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// DeleteAsignacionController maneja DELETE /api/control-asistencia/horarios/asignaciones/:id
// Desactiva (soft delete) una asignación de empleado a horario.
func DeleteAsignacionController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := RemoveAsignacion(db, id); err != nil {
		logging.Error.Printf("Error eliminando asignación %d: %v", id, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Asignación eliminada exitosamente",
	})
}

// Asigna un horario a uno o más empleados (id_tercero).
// Cierra asignaciones activas previas y registra las nuevas.
func AsignarHorarioController(c *fiber.Ctx) error {

	var req AsignacionHorarioRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	// Inyectar usuario desde contexto JWT
	req.RegistradoPor = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := RegisterHorarioTercero(db, &req)
	if err != nil {
		logging.Error.Printf("Error asignando horario [id_horario=%d]: %v", req.IDHorario, err)
		status := http.StatusBadRequest
		if isErrorInternoHorario(err) {
			status = http.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Horario asignado exitosamente",
		"data":    result,
	})
}

// isErrorInternoHorario clasifica si el error es de negocio (esperado) o de infraestructura.
// Los errores de negocio retornan 400; los de infraestructura retornan 500.
func isErrorInternoHorario(err error) bool {
	erroresNegocio := []error{
		ErrCodigoHorarioDuplicado,
		ErrHorarioNoEncontrado,
		ErrTerceroNoEncontrado,
		ErrHorarioInactivo,
		ErrBloquesSinOrdenUnico,
	}
	for _, e := range erroresNegocio {
		if errors.Is(err, e) {
			return false
		}
	}
	return true
}
