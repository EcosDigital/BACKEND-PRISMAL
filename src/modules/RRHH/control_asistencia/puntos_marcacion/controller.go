package control_asistencia

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreatePuntoMarcacionController(c *fiber.Ctx) error {

	var req PuntoMarcacionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := RegisterPuntoMarcacion(db, &req)
	if err != nil {
		logging.Error.Printf("Error creando punto de marcación: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Punto de marcación creado exitosamente",
		"data":    result,
	})

}

func ListPuntosMarcacionController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterPuntosMarcacion(db, int64(middlewares.GetEmpresaID(c)))
	if err != nil {
		logging.Error.Printf("Error listando puntos de marcación: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)

}

func ListPuntoMarcacionByIDController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	result, err := FilterPuntoMarcacionByID(db, id)
	if err != nil {
		logging.Error.Printf("Error obteniendo punto de marcación %d: %v", id, err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func ChangePuntoMarcacionController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req UpdatePuntoMarcacionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := UpdatePuntoMarcacion(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error actualizando punto de marcación %d: %v", id, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Punto de marcación actualizado exitosamente",
		"data":    result,
	})

}

func ChangeEstadoPuntoMarcacionController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CambioEstadoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ProcessCambioEstado(db, id, &req); err != nil {
		logging.Error.Printf("Error cambiando estado del punto de marcación %d: %v", id, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	estado := "inactivado"
	if req.Activo {
		estado = "activado"
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Punto de marcación " + estado + " exitosamente",
	})

}

// ─── Refs ─────────────────────────────────────────────────────────────────────

func ListTiposPuntoMarcacionController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterTiposPuntoMarcacion(db)
	if err != nil {
		logging.Error.Printf("Error obteniendo tipos de punto de marcación: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// ─── MARCACIONES (MOV) ─────────────────────────────────────────────────────────────────────

// RegistrarMarcacionController maneja POST /api/asistencia/marcar
// Requiere JWT válido. El id_empleado, IP y user_agent se capturan
// desde el contexto de Fiber — no se aceptan del body.

func RegistrarMarcacionController(c *fiber.Ctx) error {

	var req MarcacionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	// Datos de trazabilidad capturados en backend — no del cliente
	req.IDEmpleado = int64(middlewares.GetUserID(c))
	req.DireccionIP = c.IP()
	req.UserAgent = c.Get("User-Agent")

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := ProcessMarcacion(db, &req)
	if err != nil {
		logging.Error.Printf("Error registrando marcación [empleado=%d token=%s ip=%s]: %v",
			req.IDEmpleado, req.Token, req.DireccionIP, err)

		// Errores de negocio → 400/422; cualquier otro → 500
		status := http.StatusBadRequest
		if isErrorInterno(err) {
			status = http.StatusInternalServerError
		}

		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": result.Mensaje,
		"data":    result,
	})

}

// isErrorInterno distingue errores de negocio (esperados) de errores
// de infraestructura (DB caída, etc.) para asignar el HTTP status correcto.
func isErrorInterno(err error) bool {
	erroresNegocio := []error{
		ErrTokenInvalido,
		ErrFueraDeRadio,
		ErrMarcacionDuplicada,
		ErrOrigenNoConfigurado,
		ErrCoordenadasInvalidas,
	}
	for _, e := range erroresNegocio {
		if errors.Is(err, e) {
			return false
		}
	}
	return true
}
