package gestiones

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Gestiones(r fiber.Router) {

	protected := r.Group("/gestiones",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/ticket",
		middlewares.VallidateBody(&TicketRequest{}), CreateTicketController)

	protected.Get("/tickets", FindTicketsController)
	protected.Get("/ticket/:id", FindTicketByIDController)

}
