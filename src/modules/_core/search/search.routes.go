package search

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SearchDynamicsCore(r fiber.Router) {

	api := r.Group("/search-dynamics")

	protected := api.Group("/",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&SearchDynamics{}), FindDataDynamicsController)

}
