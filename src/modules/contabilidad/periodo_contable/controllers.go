package periodo_contable

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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
	case errors.Is(err, ErrPeriodosYaExisten),
		errors.Is(err, ErrAñoFiscalCerrado),
		errors.Is(err, ErrPeriodoEstadoInvalido),
		errors.Is(err, ErrPeriodoYaEnEseEstado),
		errors.Is(err, ErrEstadoModuloInvalido),
		errors.Is(err, ErrModuloNoEncontrado),
		errors.Is(err, ErrOperacionNoPermitida):
		return http.StatusBadRequest
	case errors.Is(err, ErrPeriodoNoEncontrado),
		errors.Is(err, ErrAñoFiscalNoEncontrado):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

// GetCatalogosPeriodosController maneja GET /contabilidad/periodos/catalogos
// Retorna todos los catálogos de referencia en una sola llamada para optimizar el frontend.
func GetCatalogosPeriodosController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	catalogos, err := GetCatalogosPeriodos(db)
	if err != nil {
		logging.Error.Printf("Error cargando catálogos de períodos: %v", err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", catalogos)
}

// GenerarPeriodosController maneja POST /contabilidad/periodos/generar
// Genera los 12 períodos mensuales para un año fiscal.
func GenerarPeriodosController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req GenerarPeriodosRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))

	result, err := GenerarPeriodosAñoFiscal(db, &req)
	if err != nil {
		logging.Error.Printf("Error generando períodos: %v", err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusCreated, "Períodos contables generados exitosamente", result)
}

// GetPeriodosPorAñoController maneja GET /contabilidad/periodos/:id_año_fiscal
// Retorna los 12 períodos de un año fiscal con sus módulos.
func GetPeriodosPorAñoController(c *fiber.Ctx) error {
	idAñoFiscal, err := strconv.ParseInt(c.Params("id_año_fiscal"), 10, 64)
	if err != nil || idAñoFiscal <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de año fiscal inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	list, err := ObtenerPeriodosPorAño(db, idAñoFiscal, empresaID)
	if err != nil {
		logging.Error.Printf("Error listando períodos del año %d: %v", idAñoFiscal, err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", list)
}

// UpdateEstadoPeriodoController maneja PATCH /contabilidad/periodos/:id/estado
// Cambia el estado de un período contable (Abierto, Cerrado, Ajuste, Bloqueado).
func UpdateEstadoPeriodoController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de período inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req UpdateEstadoPeriodoRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))

	result, err := CambiarEstadoPeriodo(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error cambiando estado período %d: %v", id, err)
		return fail(c, httpStatus(err), err)
	}

	mensajes := map[int]string{
		EstadoPeriodoAbierto:   "Período contable abierto exitosamente",
		EstadoPeriodoCerrado:   "Período contable cerrado exitosamente",
		EstadoPeriodoAjuste:    "Período contable en modo ajuste",
		EstadoPeriodoBloqueado: "Período contable bloqueado exitosamente",
	}
	msg := mensajes[req.IDEstado]
	if msg == "" {
		msg = "Estado de período actualizado"
	}

	return ok(c, http.StatusOK, msg, result)
}

// UpdateEstadoModuloPeriodoController maneja PATCH /contabilidad/periodos/:id/modulos/:id_modulo/estado
// Cambia el estado operacional de un módulo específico dentro de un período.
func UpdateEstadoModuloPeriodoController(c *fiber.Ctx) error {
	idPeriodo, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || idPeriodo <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de período inválido"))
	}

	idModulo, err := strconv.Atoi(c.Params("id_modulo"))
	if err != nil || idModulo <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("ID de módulo inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	var req UpdateEstadoModuloPeriodoRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.UserID = int64(middlewares.GetUserID(c))

	result, err := CambiarEstadoModuloPeriodo(db, idPeriodo, idModulo, &req)
	if err != nil {
		logging.Error.Printf("Error cambiando estado módulo %d en período %d: %v", idModulo, idPeriodo, err)
		return fail(c, httpStatus(err), err)
	}

	return ok(c, http.StatusOK, "Estado de módulo actualizado exitosamente", result)
}

// ValidarOperacionController maneja GET /contabilidad/periodos/validar
// Endpoint de validación centralizada: consulta si una operación puede ejecutarse.
// Query params: empresa_id (desde JWT), modulo (int), fecha (YYYY-MM-DD)
func ValidarOperacionController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fail(c, http.StatusInternalServerError, err)
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	idModulo, err := strconv.Atoi(c.Query("modulo"))
	if err != nil || idModulo <= 0 {
		return fail(c, http.StatusBadRequest, errors.New("parámetro 'modulo' inválido"))
	}

	fechaStr := c.Query("fecha")
	if fechaStr == "" {
		fechaStr = time.Now().Format("2006-01-02")
	}
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		return fail(c, http.StatusBadRequest, errors.New("formato de fecha inválido, use YYYY-MM-DD"))
	}

	result, err := validarOperacion(db, empresaID, fecha, idModulo)
	if err != nil {
		logging.Error.Printf("Error validando operación módulo %d fecha %s: %v", idModulo, fechaStr, err)
		return fail(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", result)
}
