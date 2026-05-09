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

	// Exportación Excel — debe registrarse ANTES de /:id para que Fiber
	// no interprete "export" como un parámetro dinámico de tipo ID.
	protected.Get("/export", ExportCuentasController)

	// CRUD
	protected.Get("/", FindCuentasController)
	protected.Post("/", middlewares.VallidateBody(&CuentaContableRequest{}), CreateCuentaController)
	protected.Get("/:id", FindCuentaByIDController)
	protected.Put("/:id", middlewares.VallidateBody(&CuentaContableUpdateRequest{}), UpdateCuentaController)
}
