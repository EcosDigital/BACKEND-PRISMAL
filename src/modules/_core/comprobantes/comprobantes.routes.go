package comprobantes

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_comprobante(r fiber.Router) {

	protected := r.Group("/comprobantes",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&ComprobanteRequest{}), CreateComprobanteController)

	protected.Get("/",
		FindLastComprobantesController)

	protected.Get("/:id",
		FindComprobanteByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&ComprobanteUpdateRequest{}), ChangeComprobanteController)

}
