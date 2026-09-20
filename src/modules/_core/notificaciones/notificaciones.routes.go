package notificaciones

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Notificaciones(r fiber.Router) {

	protected := r.Group("/notificaciones",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Get("/tipos",
		FindTiposNotificacionController)

	protected.Post("/dispositivos",
		middlewares.VallidateBody(&DispositivoPushRequest{}), CreateDispositivoPushController)

	protected.Put("/dispositivos/baja",
		middlewares.VallidateBody(&BajaDispositivoRequest{}), BajaDispositivoPushController)

	protected.Get("/origenes",
		FindOrigenesDestinatarioController)

	protected.Post("/",
		middlewares.VallidateBody(&NotificacionRequest{}), CreateNotificacionController)

	protected.Get("/mias",
		FindNotificacionesUsuarioController)

	protected.Put("/leer-todas",
		MarcarLeidasTodasController)

	protected.Put("/:id/leido",
		MarcarLeidaController)

	protected.Get("/destinos/usuarios",
		FindUsuariosDestinoController)

	protected.Get("/destinos/roles",
		FindRolesDestinoController)

}
