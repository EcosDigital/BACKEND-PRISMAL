package fiados

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// ─── Cuentas ──────────────────────────────────────────────────────────────────

// FindCuentasController GET /fiados/cuentas
// Soporta ?q=nombre_o_cedula para el buscador en tiempo real del frontend.
func FindCuentasController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	q := c.Query("q", "")

	results, err := GetCuentas(db, empresaID, q)
	if err != nil {
		logging.Error.Printf("FindCuentas error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// FindCuentaByIDController GET /fiados/cuentas/:id
func FindCuentaByIDController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	result, err := GetCuenta(db, id, empresaID)
	if err != nil {
		logging.Error.Printf("FindCuentaByID error: %v", err)
		if err.Error() == "cuenta no encontrada" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// ChangeCuentaController PUT /fiados/cuentas/:id
// Actualiza límite de crédito, observaciones o QR de cobro.
func ChangeCuentaController(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req ActualizarCuentaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := EditCuenta(db, id, &req); err != nil {
		logging.Error.Printf("ChangeCuenta error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Cuenta actualizada exitosamente",
		"id":      id,
	})
}

// ─── Movimientos ──────────────────────────────────────────────────────────────

// CreateFiadoController POST /fiados/fiado
// Registra una nueva deuda. Crea tercero y cuenta si no existen.
func CreateFiadoController(c *fiber.Ctx) error {
	var req RegistrarFiadoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := RegistrarFiado(db, &req)
	if err != nil {
		logging.Error.Printf("CreateFiado error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Fiado registrado exitosamente",
		"data":    result,
	})
}

// CreateAbonoController POST /fiados/abono
// Registra un pago parcial o total sobre una cuenta existente.

func CreateAbonoController(c *fiber.Ctx) error {
	var req RegistrarAbonoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := RegistrarAbono(db, &req)
	if err != nil {
		logging.Error.Printf("CreateAbono error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	msg := "Abono registrado exitosamente"
	if result.CuentaSaldada {
		msg = "Abono registrado — cuenta saldada y cerrada"
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": msg,
		"data":    result,
	})
}

// AnularMovimientoController PUT /fiados/movimientos/:id/anular
func AnularMovimientoController(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req AnularMovimientoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := AnularMovimiento(db, id, &req); err != nil {
		logging.Error.Printf("AnularMovimiento error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Movimiento anulado exitosamente",
		"id":      id,
	})
}

// FindHistorialController GET /fiados/cuentas/:id/historial
// Soporta ?fecha_inicio=YYYY-MM-DD&fecha_fin=YYYY-MM-DD (opcionales).
func FindHistorialController(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	fechaInicio := c.Query("fecha_inicio", "")
	fechaFin := c.Query("fecha_fin", "")

	results, err := GetHistorial(db, id, empresaID, fechaInicio, fechaFin)
	if err != nil {
		logging.Error.Printf("FindHistorial error: %v", err)
		if err.Error() == "cuenta no encontrada" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// FindResumenController GET /fiados/resumen
// Retorna los stats del dashboard: cartera total, cuentas por estado, etc.
func FindResumenController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	// id_sede es opcional en el query string
	var sedeID *int64
	if sedeParam := c.Query("id_sede", ""); sedeParam != "" {
		if v, err := strconv.ParseInt(sedeParam, 10, 64); err == nil {
			sedeID = &v
		}
	}

	result, err := GetResumen(db, empresaID, sedeID)
	if err != nil {
		logging.Error.Printf("FindResumen error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
