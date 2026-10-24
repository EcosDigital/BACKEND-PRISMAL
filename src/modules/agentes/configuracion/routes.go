package configuracion

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_configuracion(r fiber.Router) {

	protected := r.Group("/configuracion",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&AgenteRequest{}), CreateAgenteController)

	protected.Get("/",
		FindAgentesController)

	protected.Get("/modulos",
		FindModulosInstaladosController)

	protected.Get("/:id",
		FindAgenteByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&AgenteRequest{}), ChangeAgenteController)

}
