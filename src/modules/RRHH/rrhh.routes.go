package RRHH

import (
	historial_asistencia "github.com/ecosistema/core/src/modules/RRHH/control_asistencia/historial"
	"github.com/ecosistema/core/src/modules/RRHH/control_asistencia/horarios"
	control_asistencia "github.com/ecosistema/core/src/modules/RRHH/control_asistencia/puntos_marcacion"
	reportes_asistencia "github.com/ecosistema/core/src/modules/RRHH/control_asistencia/reportes"
	"github.com/gofiber/fiber/v2"
)

func RegisterRRHHRoutes(r fiber.Router) {

	RRHH := r.Group("/RRHH")

	//puntos marcacion
	control_asistencia.Rutas_puntos_marcacion(RRHH)
	//horarios
	horarios.Rutas_horarios(RRHH)
	//historial
	historial_asistencia.Rutas_historial_asistencia(RRHH)
	//reportes
	reportes_asistencia.Rutas_reportes(RRHH)
}
