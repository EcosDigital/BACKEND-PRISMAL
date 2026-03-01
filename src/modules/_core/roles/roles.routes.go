package roles

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Roles(r fiber.Router) {

	protected := r.Group("/roles",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&RolesRequest{}), CreateRolUserController)

	protected.Get("/",
		FindLastRolUserController)

	protected.Get("/:id",
		FindRolUserByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&RolesResponse{}), ChangeRolUserController)

	protected.Post("/config",
		middlewares.VallidateBody(&ConfigRolRequest{}), CreateConfigRolController)

}
