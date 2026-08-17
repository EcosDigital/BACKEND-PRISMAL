package documentacion_sop

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
	case errors.Is(err, ErrCodeYaExiste),
		errors.Is(err, ErrCategoriaInvalida),
		errors.Is(err, ErrEstadoInvalido),
		errors.Is(err, ErrFechaEfectivaInvalida),
		errors.Is(err, ErrFechaRevisionInvalida),
		errors.Is(err, ErrFechasInconsistentes),
		errors.Is(err, ErrContenidoVacio),
		errors.Is(err, ErrVersionVacia):
		return http.StatusBadRequest
	case errors.Is(err, ErrSOPNoEncontrado):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

// GetCatalogosSOPController maneja GET /documentacion/sop/catalogos
// Devuelve categorías y estados en una sola llamada para poblar selects del frontend.
func GetCatalogosSOPController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	catalogos, err := ObtenerCatalogos(db)
	if err != nil {
		logging.Error.Printf("Error cargando catálogos SOP: %v", err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", catalogos)
}

// CreateSOPController maneja POST /documentacion/sop
func CreateSOPController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req CreateSOPRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))

	result, err := RegistrarSOP(db, &req)
	if err != nil {
		logging.Error.Printf("Error creando SOP: %v", err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusCreated, "SOP creado exitosamente", result)
}

// GetSOPsController maneja GET /documentacion/sop
// Acepta query params opcionales: id_categoria, id_estado, area, search
func GetSOPsController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	// Construir filtros desde query params — valores inválidos se ignoran (0/"")
	idCategoria, _ := strconv.Atoi(c.Query("id_categoria"))
	idEstado, _ := strconv.Atoi(c.Query("id_estado"))

	filters := SOPFilters{
		IDCategoria: idCategoria,
		IDEstado:    idEstado,
		Area:        c.Query("area"),
		Search:      c.Query("search"),
	}

	list, err := ListadoSOPs(db, empresaID, filters)
	if err != nil {
		logging.Error.Printf("Error listando SOPs: %v", err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", list)
}

// GetSOPByIDController maneja GET /documentacion/sop/:id
func GetSOPByIDController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de SOP inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	result, err := ConsultarSOP(db, id, empresaID)
	if err != nil {
		logging.Error.Printf("Error consultando SOP %d: %v", id, err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusOK, "", result)
}

// UpdateSOPController maneja PUT /documentacion/sop/:id
func UpdateSOPController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de SOP inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req UpdateSOPRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))

	result, err := ModificarSOP(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error actualizando SOP %d: %v", id, err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusOK, "SOP actualizado exitosamente", result)
}

// DeleteSOPController maneja DELETE /documentacion/sop/:id
// Ejecuta soft delete (is_active = false); el registro permanece en la BD.
func DeleteSOPController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de SOP inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	userID := int64(middlewares.GetUserID(c))

	if err := EliminarSOP(db, id, empresaID, userID); err != nil {
		logging.Error.Printf("Error eliminando SOP %d: %v", id, err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusOK, "SOP eliminado exitosamente", nil)
}
