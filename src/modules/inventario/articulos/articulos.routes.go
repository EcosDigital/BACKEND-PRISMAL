package articulos

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_articulos(r fiber.Router) {
	protected := r.Group("/articulos",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&ArticuloRequest{}), CreateArticuloController)

	protected.Get("/",
		FindArticulosController)

	protected.Get("/articulos/:id",
		FindArticuloByIDController)

	protected.Put("/articulos/:id",
		middlewares.VallidateBody(&ArticuloUpdateRequest{}), ChangeArticuloController)

	protected.Get("/search", SearchArticulosController)

}
