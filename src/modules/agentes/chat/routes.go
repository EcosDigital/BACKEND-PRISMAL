package chat

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_chat(r fiber.Router) {

	protected := r.Group("/chat",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&ChatRequest{}), EnviarMensajeController)

}
