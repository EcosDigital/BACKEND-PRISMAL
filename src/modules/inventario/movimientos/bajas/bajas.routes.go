package bajas

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_bajas(r fiber.Router) {

	protected := r.Group("/bajas",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// ─── CRUD ─────────────────────────────────────────────────────
	protected.Post("/",
		middlewares.VallidateBody(&BajaRequest{}), CreateBajaController)
	protected.Get("/", FindBajasController)
	protected.Get("/:id", FindBajaByIDController)

	// ─── Reporte ──────────────────────────────────────────────────
	protected.Get("/:id/reporte", GetBajaReporteController)

	// ─── Anulación ────────────────────────────────────────────────
	protected.Put("/:id/anular",
		middlewares.VallidateBody(&AnulacionBajaRequest{}), AnularBajaController)

	// ─── Referencias ──────────────────────────────────────────────
	protected.Get("/ref/subtipos", FindSubtiposBajaController)
	protected.Get("/ref/comprobantes", FindComprobantesBajaController)
	// GET /inventario/bajas/ref/lotes/:id_articulo?id_bodega=X
	protected.Get("/ref/lotes/:id_articulo", FindLotesBajaController)

}
