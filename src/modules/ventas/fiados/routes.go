package fiados

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_fiados registra todas las rutas del módulo de fiados.
// Se monta sobre el grupo /fiados desde el router principal.
//
// Ejemplo en main.go:
//
//	api := app.Group("/api/v1")
//	fiados.Rutas_fiados(api.Group("/fiados"))
func Rutas_fiados(r fiber.Router) {
	protected := r.Group("fiados",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// ── Dashboard ─────────────────────────────────────────────
	// GET /fiados/resumen
	// Stats de cartera para el header del módulo.
	protected.Get("/resumen", FindResumenController)

	// ── Cuentas ───────────────────────────────────────────────
	// GET  /fiados/cuentas          lista + buscador (?q=)
	// GET  /fiados/cuentas/:id      detalle de una cuenta
	// PUT  /fiados/cuentas/:id      editar límite / QR / observaciones
	protected.Get("/cuentas", FindCuentasController)
	protected.Get("/cuentas/:id", FindCuentaByIDController)
	protected.Put("/cuentas/:id",
		middlewares.VallidateBody(&ActualizarCuentaRequest{}),
		ChangeCuentaController,
	)

	// ── Historial de movimientos de una cuenta ─────────────────
	// GET /fiados/cuentas/:id/historial
	//     ?fecha_inicio=YYYY-MM-DD&fecha_fin=YYYY-MM-DD (opcionales)
	protected.Get("/cuentas/:id/historial", FindHistorialController)

	// ── Registrar fiado (nueva deuda) ─────────────────────────
	// POST /fiados/fiado
	// Crea tercero y cuenta si no existen.
	protected.Post("/credito",
		middlewares.VallidateBody(&RegistrarFiadoRequest{}),
		CreateFiadoController,
	)

	// ── Registrar abono (pago) ────────────────────────────────
	// POST /fiados/abono
	protected.Post("/abono",
		middlewares.VallidateBody(&RegistrarAbonoRequest{}),
		CreateAbonoController,
	)

	// ── Anular movimiento ─────────────────────────────────────
	// PUT /fiados/movimientos/:id/anular
	protected.Put("/movimientos/:id/anular",
		middlewares.VallidateBody(&AnularMovimientoRequest{}),
		AnularMovimientoController,
	)
}
