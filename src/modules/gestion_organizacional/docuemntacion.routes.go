package docuemntacion

import (
	"github.com/ecosistema/core/src/modules/gestion_organizacional/proyectos"
	documentacion_sop "github.com/ecosistema/core/src/modules/gestion_organizacional/sop"
	"github.com/gofiber/fiber/v2"
)

func RegisterDocumentacionRoutes(r fiber.Router) {

	org := r.Group("/organizacion")

	//SOP
	documentacion_sop.Rutas_SOP(org)
	//Proyectos
	proyectos.Rutas_Proyectos(org)

}
