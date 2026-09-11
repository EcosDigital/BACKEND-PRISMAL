package estructura_fisica

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_EstructuraFisica(r fiber.Router) {

	protected := r.Group("/mesas",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&MesaRequest{}), CreateMesaController)

	protected.Get("/",
		FindMesasController)

	protected.Get("/activas",
		FindMesasActivasController)

	protected.Get("/:id",
		FindMesaByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&MesaRequest{}), ChangeMesaController)

}
