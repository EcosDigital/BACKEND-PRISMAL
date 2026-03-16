package traslados

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateTrasladoController(c *fiber.Ctx) error {
	var req TrasladoRequest
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

	result, err := CreateTraslado(db, &req)
	if err != nil {
		logging.Error.Printf("Error registrando traslado: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":            "Traslado registrado exitosamente",
		"id_movimiento":      result.IDMovimiento,
		"id_mov_comprobante": result.IDMovComprobante,
		"consecutivo":        result.Consecutivo,
		"prefijo":            result.Prefijo,
	})
}

func FindTrasladosController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	results, err := FilterTraslados(db, int64(middlewares.GetEmpresaID(c)))
	if err != nil {
		logging.Error.Printf("Error listando traslados: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func FindTrasladoByIDController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	result, err := FilterTrasladoByID(db, id)
	if err != nil {
		logging.Error.Printf("Error obteniendo traslado: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func AnularTrasladoController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}
	var req AnulacionTrasladoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}
	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := ProcessAnulacionTraslado(db, id, &req); err != nil {
		logging.Error.Printf("Error anulando traslado %d: %v", id, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Traslado anulado exitosamente"})
}

// ─── Refs ─────────────────────────────────────────────────────

func FindBodegasTrasladoController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	results, err := FilterBodegasTraslado(db, int64(middlewares.GetEmpresaID(c)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func FindLotesController(c *fiber.Ctx) error {
	idArticulo, err := strconv.ParseInt(c.Params("id_articulo"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_articulo inválido"})
	}
	idBodega, err := strconv.ParseInt(c.Query("id_bodega", "0"), 10, 64)
	if err != nil || idBodega == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_bodega requerido"})
	}
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	results, err := FilterLotes(db, idArticulo, idBodega, int64(middlewares.GetEmpresaID(c)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func FindComprobantesTrasladorController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	results, err := FilterComprobantesTraslado(db, int64(middlewares.GetEmpresaID(c)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}
