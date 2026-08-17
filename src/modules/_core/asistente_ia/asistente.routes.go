package asistente

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Asistente(r fiber.Router) {

	protected := r.Group("/ia",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/chat",
		middlewares.JWTProtected,
		middlewares.VallidateBody(&ChatRequest{}),
		ChatController)

}
