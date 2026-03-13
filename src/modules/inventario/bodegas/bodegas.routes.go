package bodega

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_bodegas(r fiber.Router) {

	protected := r.Group("/bodegas",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&BodegaRequest{}), CreateBodegaController)

	protected.Get("/",
		FindBodegasController)

	protected.Get("/:id",
		FindBodegaByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&BodegaUpdateRequest{}), ChangeBodegaController)

}
