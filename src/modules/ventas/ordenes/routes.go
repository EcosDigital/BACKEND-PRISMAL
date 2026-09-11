package ordenes

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_ordenes registra las rutas del módulo de órdenes de mesa (el pedido
// que toma el mesero, tenant-scoped). Se monta sobre el grupo /ordenes
// desde ventas.routes.go.
func Rutas_ordenes(r fiber.Router) {
	protected := r.Group("ordenes",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// GET /ordenes             lista de órdenes del tenant actual
	// GET /ordenes/stream      canal SSE de avisos en tiempo real (antes de
	//                          /:id para que "stream" no se interprete como un ID)
	// GET /ordenes/:id         detalle (cabecera + artículos)
	// POST /ordenes            registrar una orden nueva
	// PUT /ordenes/:id/estado  cambiar estado (valida contra ref_estado_orden)
	protected.Get("/", FindOrdenesController)
	protected.Get("/stream", StreamOrdenesController)
	protected.Get("/:id", FindOrdenByIDController)
	protected.Post("/",
		middlewares.VallidateBody(&OrdenRequest{}),
		CreateOrdenController,
	)
	protected.Put("/:id/estado",
		middlewares.VallidateBody(&CambiarEstadoOrdenRequest{}),
		ChangeEstadoOrdenController,
	)
}
