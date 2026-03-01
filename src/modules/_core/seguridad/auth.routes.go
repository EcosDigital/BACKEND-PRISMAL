package seguridad

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_seguridad(r fiber.Router) {
	api := r.Group("/auth")

	//signin
	api.Post("/signin",
		middlewares.VallidateBody(&SigninRequest{}), SigninController)

	api.Get("/access", middlewares.JWTProtectedBasic, AccessGetController)

	api.Post("/access", middlewares.JWTProtectedBasic, AccessController)

	api.Post("/signout",
		middlewares.JWTProtectedBasic, SignoutController)

	api.Post("/signout-all",
		middlewares.JWTProtectedBasic, SignoutAllController)

	api.Get("/verifyToken", VerifyTokenController)

	/* TEMPORAL: Endpoint de debugging
	api.Get("/auth/debug/sessions",
		controller_seguridad.DebugSessionsController) */

}
