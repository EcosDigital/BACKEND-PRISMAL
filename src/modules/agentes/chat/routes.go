package chat

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_chat(r fiber.Router) {

	// Meta no inicia sesión en Prismar: estas dos rutas son públicas y el
	// aviso se valida con la firma de Meta
	r.Get("/chat/whatsapp", VerificarWhatsAppController)
	r.Post("/chat/whatsapp", RecibirWhatsAppController)

	protected := r.Group("/chat",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&ChatRequest{}), EnviarMensajeController)

}
