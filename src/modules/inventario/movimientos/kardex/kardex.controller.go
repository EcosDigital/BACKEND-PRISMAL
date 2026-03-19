package kardex

import (
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// GetKardexController godoc
// GET /inventario/kardex?id_articulo=1&id_bodega=1&lote=L001&fecha_inicio=2025-01-01&fecha_fin=2025-03-31
func GetKardexController(c *fiber.Ctx) error {

	var req KardexRequest

	// Parsear query params manualmente (Fiber QueryParser)
	idArt, err := strconv.ParseInt(c.Query("id_articulo", "0"), 10, 64)
	if err != nil || idArt == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_articulo requerido"})
	}

	idBodega, err := strconv.ParseInt(c.Query("id_bodega", "0"), 10, 64)
	if err != nil || idBodega == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_bodega requerido"})
	}

	fechaInicio := c.Query("fecha_inicio", "")
	fechaFin := c.Query("fecha_fin", "")

	if fechaInicio == "" || fechaFin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_inicio y fecha_fin son requeridos"})
	}

	lote := c.Query("lote", "")

	req = KardexRequest{
		IDArticulo:  idArt,
		IDBodega:    idBodega,
		Lote:        lote,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	reporte, err := FilterKardex(db, &req, empresaID)
	if err != nil {
		logging.Error.Printf("Error generando kardex artículo=%d bodega=%d: %v", idArt, idBodega, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reporte)
}

// GetLotesController godoc
// GET /inventario/kardex/ref/lotes?id_articulo=1&id_bodega=1
func GetLotesController(c *fiber.Ctx) error {

	idArt, err := strconv.ParseInt(c.Query("id_articulo", "0"), 10, 64)
	if err != nil || idArt == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_articulo requerido"})
	}

	idBodega, err := strconv.ParseInt(c.Query("id_bodega", "0"), 10, 64)
	if err != nil || idBodega == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id_bodega requerido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	lotes, err := FilterLotes(db, idArt, idBodega, empresaID)
	if err != nil {
		logging.Error.Printf("Error obteniendo lotes artículo=%d bodega=%d: %v", idArt, idBodega, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(lotes)
}
