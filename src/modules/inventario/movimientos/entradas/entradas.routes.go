package entradas

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_entradas(r fiber.Router) {

	protected := r.Group("/entradas",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// ─── CRUD entradas ────────────────────────────────────────────────────────
	protected.Post("/",
		middlewares.VallidateBody(&EntradaRequest{}), CreateEntradaController)

	protected.Get("/",
		FindEntradasController)

	protected.Get("/:id",
		FindEntradaByIDController)

	// ─── Referencias para el formulario ──────────────────────────────────────
	protected.Get("/ref/comprobantes",
		FindComprobantesEntradaController)

	protected.Get("/ref/impuestos",
		FindImpuestosController)

	// Existencia de un artículo en una bodega específica
	// GET /inventario/entradas/ref/existencia/:id_articulo?id_bodega=1
	protected.Get("/ref/existencia/:id_articulo",
		FindExistenciaArticuloController)
}
