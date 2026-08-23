package catalogo

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_catalogo(r fiber.Router) {
	protected := r.Group("/catalogo",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// ── Cabecera ──────────────────────────────────────────────
	protected.Post("/", CreateCatalogoController)
	protected.Get("/", FindCatalogosController)
	protected.Get("/:id", FindCatalogoByIDController)
	protected.Put("/:id", ChangeCatalogoController)

	// ── Detalle: artículos publicados ────────────────────────
	protected.Post("/:id_catalogo/articulos", CreateCatalogoArticuloController)
	protected.Get("/:id_catalogo/articulos", FindCatalogoArticulosController)
	protected.Put("/articulos/:id", ChangeCatalogoArticuloController)
	protected.Delete("/articulos/:id", DeleteCatalogoArticuloController)
}
