package pedidos

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_pedidos registra las rutas del módulo de pedidos de domicilio
// (Mi Llave → Prismar). Se monta sobre el grupo /pedidos desde ventas.routes.go.
func Rutas_pedidos(r fiber.Router) {
	protected := r.Group("pedidos",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// GET /pedidos             lista de pedidos del tenant actual
	// GET /pedidos/stream      canal SSE de avisos en tiempo real (antes de /:id
	//                          para que "stream" no se interprete como un ID)
	// GET /pedidos/:id         detalle (cabecera + artículos + entrega)
	// PUT /pedidos/:id/estado  cambiar estado (valida contra ref_estado_pedido)
	protected.Get("/", FindPedidosController)
	protected.Get("/stream", StreamPedidosController)
	protected.Get("/:id", FindPedidoByIDController)
	protected.Put("/:id/estado",
		middlewares.VallidateBody(&CambiarEstadoRequest{}),
		ChangeEstadoPedidoController,
	)
}
