package ano_fiscal

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// ─── Helpers de respuesta ─────────────────────────────────────────────────────

func ok(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(APIResponse{Success: true, Message: message, Data: data})
}

func fail(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(APIResponse{Success: false, Message: err.Error()})
}

// httpStatus mapea errores de dominio a su código HTTP correspondiente.
func httpStatus(err error) int {
	switch {
	case errors.Is(err, ErrYearYaExiste),
		errors.Is(err, ErrEstadoInvalido),
		errors.Is(err, ErrYaEnEseEstado),
		errors.Is(err, ErrAsientosPendientes):
		return http.StatusBadRequest
	case errors.Is(err, ErrAñoFiscalNoEncontrado):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

// GetEstadosAñoFiscalController maneja GET /contabilidad/anos-fiscales/estados
// Devuelve el catálogo de estados para poblar selects en el frontend.
func GetEstadosAñoFiscalController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	estados, err := GetEstadosAñoFiscal(db)
	if err != nil {
		logging.Error.Printf("Error cargando estados año fiscal: %v", err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", estados)
}

// CreateAñoFiscalController maneja POST /contabilidad/anos-fiscales
func CreateAñoFiscalController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req CreateAñoFiscalRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	result, err := RegistrarAñoFiscal(db, &req)
	if err != nil {
		logging.Error.Printf("Error creando año fiscal: %v", err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusCreated, "Año fiscal creado exitosamente", result)
}

// GetAñosFiscalesController maneja GET /contabilidad/anos-fiscales
func GetAñosFiscalesController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	list, err := ListadoAñosFiscales(db, empresaID)
	if err != nil {
		logging.Error.Printf("Error listando años fiscales: %v", err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", list)
}

// UpdateEstadoAñoFiscalController maneja PATCH /contabilidad/anos-fiscales/:id/estado
func UpdateEstadoAñoFiscalController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de año fiscal inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req UpdateEstadoAñoFiscalRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))

	result, err := CambiarEstadoAñoFiscal(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error cambiando estado año fiscal %d: %v", id, err)
		return fail(c, httpStatus(err), err)
	}

	msg := "Año fiscal cerrado exitosamente"
	if req.IDEstado == EstadoAbierto {
		msg = "Año fiscal abierto exitosamente"
	}

	return ok(c, http.StatusOK, msg, result)
}
