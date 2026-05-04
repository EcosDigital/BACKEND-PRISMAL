package plan_cuentas

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// GET /contabilidad/cuentas
func FindCuentasController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := FilterCuentas(db, empresaID)
	if err != nil {
		logging.Error.Printf("Error listando cuentas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// GET /contabilidad/cuentas/:id
func FindCuentaByIDController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	result, err := FilterCuentaByID(db, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// POST /contabilidad/cuentas
func CreateCuentaController(c *fiber.Ctx) error {
	var req CuentaContableRequest

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

	newID, err := RegisterCuenta(db, &req)
	if err != nil {
		logging.Error.Printf("Error creando cuenta: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Cuenta contable registrada exitosamente",
		"id":      newID,
	})
}

// PUT /contabilidad/cuentas/:id
func UpdateCuentaController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CuentaContableUpdateRequest
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

	updID, err := EditCuenta(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error actualizando cuenta: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Cuenta contable actualizada exitosamente",
		"id":      updID,
	})
}

// GET /contabilidad/cuentas/referencias?exclude_id=
// Devuelve naturalezas, tipos, niveles y cuentas padre en una sola llamada
func GetReferenciasController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	excludeID, _ := strconv.ParseInt(c.Query("exclude_id", "0"), 10, 64)

	refs, err := GetReferencias(db, empresaID, excludeID)
	if err != nil {
		logging.Error.Printf("Error cargando referencias: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(refs)
}

// ImportCuentasController recibe el JSON con las filas del plano PUC y hace upsert.
//
// POST /contabilidad/cuentas/import
// Body: { "filas": [ { "CODIGO_CUENTA": "1", "NOMBRE": "Activo", ... } ] }
// Respuesta: { "creadas": N, "actualizadas": N, "errores": N, "detalle": [...] }
//
// NOTA IMPORTANTE: el plano debe enviarse en ORDEN jerárquico (padres antes que hijos)
// ya que una cuenta hijo puede referenciar como padre a una cuenta nueva del mismo plano.
func ImportCuentasController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var req ImportCuentasRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}
	if len(req.Filas) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "El plano no contiene filas"})
	}

	result, err := ImportCuentas(
		db, &req,
		int64(middlewares.GetUserID(c)),
		int64(middlewares.GetEmpresaID(c)),
		int64(middlewares.GetSedeID(c)),
	)
	if err != nil {
		logging.Error.Printf("Error en carga masiva de cuentas: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(result)
}
