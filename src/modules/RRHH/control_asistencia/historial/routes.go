package historial_asistencia

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_historial_asistencia(r fiber.Router) {

	protected := r.Group("/historial-asistencia",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// GET /control-asistencia/historial?id_tercero=&fecha_inicio=&fecha_fin=
	protected.Get("/", GetHistorialAsistenciaController)
}
