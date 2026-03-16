package traslados

import (
	reportes_inv "github.com/ecosistema/core/src/modules/inventario/movimientos/_reportes"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_traslados(r fiber.Router) {

	protected := r.Group("/traslados",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// ─── CRUD ────────────────────────────────────────────────────
	protected.Post("/",
		middlewares.VallidateBody(&TrasladoRequest{}), CreateTrasladoController)

	protected.Get("/",
		FindTrasladosController)

	protected.Get("/:id",
		FindTrasladoByIDController)

	protected.Get("/:id/reporte", reportes_inv.GetTrasladoReporteController)

	// ─── Anulación ────────────────────────────────────────────────
	// PUT /inventario/traslados/:id/anular
	protected.Put("/:id/anular",
		middlewares.VallidateBody(&AnulacionTrasladoRequest{}), AnularTrasladoController)

	// ─── Referencias para el formulario ──────────────────────────
	protected.Get("/ref/bodegas",
		FindBodegasTrasladoController)

	protected.Get("/ref/comprobantes",
		FindComprobantesTrasladorController)

	// GET /inventario/traslados/ref/lotes/:id_articulo?id_bodega=X
	protected.Get("/ref/lotes/:id_articulo",
		FindLotesController)
}
