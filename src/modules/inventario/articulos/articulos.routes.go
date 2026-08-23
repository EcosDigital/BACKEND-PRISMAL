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

	// Rutas estáticas antes de "/:id" — Fiber matchea por orden de registro,
	// no prioriza rutas estáticas automáticamente sobre parámetros.
	protected.Get("/search", SearchArticulosController)

	protected.Post("/import", ImportArticulosController)

	protected.Get("/:id",
		FindArticuloByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&ArticuloUpdateRequest{}), ChangeArticuloController)

	protected.Post("/:id/imagen", SetArticuloImagenController)

}
