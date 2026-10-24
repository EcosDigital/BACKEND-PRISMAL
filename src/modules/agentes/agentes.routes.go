package agentes

import (
	"github.com/ecosistema/core/src/modules/agentes/chat"
	"github.com/ecosistema/core/src/modules/agentes/configuracion"
	"github.com/gofiber/fiber/v2"
)

func RegisterAgentesRoutes(r fiber.Router) {

	agen := r.Group("/agentes")

	configuracion.Rutas_configuracion(agen)
	chat.Rutas_chat(agen)

}
