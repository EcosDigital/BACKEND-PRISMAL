package historial_asistencia

import (
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// GetHistorialAsistenciaController godoc
// GET /control-asistencia/historial?id_tercero=1&fecha_inicio=2026-01-01&fecha_fin=2026-04-05
func GetHistorialAsistenciaController(c *fiber.Ctx) error {

	idTercero, err := strconv.ParseInt(c.Query("id_tercero", "0"), 10, 64)
	if err != nil || idTercero == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_tercero requerido",
		})
	}

	fechaInicio := c.Query("fecha_inicio", "")
	fechaFin := c.Query("fecha_fin", "")

	if fechaInicio == "" || fechaFin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fecha_inicio y fecha_fin son requeridos",
		})
	}

	req := HistorialAsistenciaRequest{
		IDTercero:   idTercero,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	reporte, err := FilterHistorialAsistencia(db, &req, empresaID)
	if err != nil {
		logging.Error.Printf("Error generando historial asistencia tercero=%d: %v", idTercero, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reporte)
}
