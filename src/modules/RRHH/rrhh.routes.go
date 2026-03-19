package RRHH

import (
	control_asistencia "github.com/ecosistema/core/src/modules/RRHH/control_asistencia/puntos_marcacion"
	"github.com/gofiber/fiber/v2"
)

func RegisterRRHHRoutes(r fiber.Router) {

	RRHH := r.Group("/RRHH")

	//puntos marcacion
	control_asistencia.Rutas_puntos_marcacion(RRHH)

}
