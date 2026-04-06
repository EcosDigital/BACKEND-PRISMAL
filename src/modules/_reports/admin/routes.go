package reports

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RutasReport(r fiber.Router) {

	protected := r.Group("/connections",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// CRUD base
	protected.Post("/",
		middlewares.VallidateBody(&ReportServerConnectionRequest{}),
		CreateConnectionController)

	protected.Get("/",
		FindConnectionsController)

	protected.Get("/:id",
		FindConnectionByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&ReportServerConnectionUpdateRequest{}),
		ChangeConnectionController)

	protected.Delete("/:id",
		DeleteConnectionController)

	// Activar conexión — solo una activa por usuario a la vez
	protected.Patch("/:id/activate",
		ActivateConnectionController)
}
