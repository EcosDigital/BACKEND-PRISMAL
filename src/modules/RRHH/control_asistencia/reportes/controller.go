package reportes_asistencia

import (
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// GetReporteDiarioController godoc
// GET /control-asistencia/reporte-diario?fecha_inicio=2026-01-01&fecha_fin=2026-04-05
func GetReporteDiarioController(c *fiber.Ctx) error {

	fechaInicio := c.Query("fecha_inicio", "")
	fechaFin := c.Query("fecha_fin", "")

	if fechaInicio == "" || fechaFin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fecha_inicio y fecha_fin son requeridos",
		})
	}

	req := ReporteDiarioRequest{
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	reporte, err := FilterReporteDiario(db, &req, empresaID)
	if err != nil {
		logging.Error.Printf("Error generando reporte diario [%s → %s]: %v", fechaInicio, fechaFin, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reporte)
}
