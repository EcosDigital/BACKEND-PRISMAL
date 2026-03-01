package terceros

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Terceros(r fiber.Router) {

	protected := r.Group("/third-parties",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	//crear nuevo registro
	protected.Post("/",
		middlewares.JWTProtected,
		middlewares.VallidateBody(&TerceroRequest{}), CreateTerceroController)

	//busca los ultimos 20 registros
	protected.Get("/",
		middlewares.JWTProtected, FilterRecentTercerosController)

	//buscar registro por ID
	protected.Get("/:id",
		middlewares.JWTProtected, FilterTerceroByIDController)

	//actualizar registro
	protected.Put("/:id",
		middlewares.JWTProtected,
		middlewares.VallidateBody(&TerceroUpdateRequest{}), UpdateTerceroController)

	//actualizar identificacion
	protected.Put("/identity/:id",
		middlewares.JWTProtected,
		middlewares.VallidateBody(&TerceroUpdateIndenty{}), ChangeTerceroIdentityController)

}
