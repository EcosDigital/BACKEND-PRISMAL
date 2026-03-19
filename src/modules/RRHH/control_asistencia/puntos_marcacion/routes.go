package control_asistencia

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_puntos_marcacion(r fiber.Router) {

	protected := r.Group("/puntos-marcacion",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// ─── CRUD ─────────────────────────────────────────────────────────────────
	// POST   /control-asistencia/puntos-marcacion/
	protected.Post("/",
		middlewares.VallidateBody(&PuntoMarcacionRequest{}), CreatePuntoMarcacionController)

	// GET  BY STATUS
	protected.Get("/",
		ListPuntosMarcacionController)

	// GET  by ID
	protected.Get("/:id",
		ListPuntoMarcacionByIDController)

	// PUT
	protected.Put("/:id",
		middlewares.VallidateBody(&UpdatePuntoMarcacionRequest{}), ChangePuntoMarcacionController)

	protected.Patch("/:id/estado",
		middlewares.VallidateBody(&CambioEstadoRequest{}), ChangeEstadoPuntoMarcacionController)

	protected.Get("/ref/tipos",
		ListTiposPuntoMarcacionController)

	protected.Post("/marcar",
		middlewares.VallidateBody(&MarcacionRequest{}), RegistrarMarcacionController)

}
