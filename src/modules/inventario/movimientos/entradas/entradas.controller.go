package entradas

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateEntradaController(c *fiber.Ctx) error {

	var req EntradaRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	logging.Info.Printf("IDTercero recibido: %v", req.IDTercero)

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := CreateEntrada(db, &req)
	if err != nil {
		logging.Error.Printf("Error registrando entrada: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":            "Entrada registrada exitosamente",
		"id_movimiento":      result.IDMovimiento,
		"id_mov_comprobante": result.IDMovComprobante,
		"consecutivo":        result.Consecutivo,
		"prefijo":            result.Prefijo,
	})
}

func FindEntradasController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := FilterEntradas(db, empresaID)
	if err != nil {
		logging.Error.Printf("Error listando entradas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

func FindEntradaByIDController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	result, err := FilterEntradaByID(db, id)
	if err != nil {
		logging.Error.Printf("Error obteniendo entrada: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// ─── Refs ─────────────────────────────────────────────────────────────────────

func FindComprobantesEntradaController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := FilterComprobantesEntrada(db, empresaID)
	if err != nil {
		logging.Error.Printf("Error obteniendo comprobantes: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

func FindImpuestosController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := FilterImpuestos(db, empresaID)
	if err != nil {
		logging.Error.Printf("Error obteniendo impuestos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

func FindExistenciaArticuloController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	idArticulo, err := strconv.ParseInt(c.Params("id_articulo"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_articulo inválido"})
	}

	idBodega, err := strconv.ParseInt(c.Query("id_bodega", "0"), 10, 64)
	if err != nil || idBodega == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_bodega requerido"})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	result, err := FilterExistenciaArticulo(db, idArticulo, idBodega, empresaID)
	if err != nil {
		logging.Error.Printf("Error obteniendo existencia: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

func AnularEntradaController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req AnulacionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ProcessAnulacion(db, id, &req); err != nil {
		logging.Error.Printf("Error anulando entrada %d: %v", id, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Entrada anulada exitosamente",
	})
}
