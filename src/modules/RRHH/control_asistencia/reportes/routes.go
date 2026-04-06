package reportes_asistencia

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_reportes(r fiber.Router) {

	protected := r.Group("/reportes",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// GET /control-asistencia/reporte-diario?fecha_inicio=&fecha_fin=
	protected.Get("/diario", GetReporteDiarioController)
}
