package costos

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Costo(r fiber.Router) {

	protected := r.Group("/costos",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&CostosRequest{}), CreateCentroCostoController)

	protected.Get("/",
		FindLastCentrosCostoController)

	protected.Get("/:id",
		FindCentroCostoByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&CostosUpdateRequest{}), ChangeCentroCostoController)

}
