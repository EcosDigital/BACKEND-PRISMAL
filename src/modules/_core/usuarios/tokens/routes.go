package tokens

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Tokens(r fiber.Router) {
	api := r.Group("/users/tokens", middlewares.JWTProtected, middlewares.EnrichContext)

	api.Post("/", middlewares.VallidateBody(&TokenCreateRequest{}), CreateTokenController)
	api.Get("/", ListTokensController)
	api.Delete("/:id", RevokeTokenController)
}
