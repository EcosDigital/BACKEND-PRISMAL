package plan_cuentas

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_cuentas(r fiber.Router) {
	protected := r.Group("/cuentas",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// Referencias (naturalezas, tipos, niveles, cuentas padre) — una sola llamada
	protected.Get("/referencias", GetReferenciasController)

	// Carga masiva PUC
	protected.Post("/import", ImportCuentasController)

	// CRUD
	protected.Get("/", FindCuentasController)
	protected.Post("/", middlewares.VallidateBody(&CuentaContableRequest{}), CreateCuentaController)
	protected.Get("/:id", FindCuentaByIDController)
	protected.Put("/:id", middlewares.VallidateBody(&CuentaContableUpdateRequest{}), UpdateCuentaController)
}
