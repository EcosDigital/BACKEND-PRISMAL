package empresa

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Sede(r fiber.Router) {

	protected := r.Group("/sedes",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&SedeRequest{}), CreateSedeController)

	protected.Get("/", FindLastSedeController)

	protected.Get("/:id",
		FindSedeByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&SedeUpdateRequest{}), ChangeSedeController)

}
