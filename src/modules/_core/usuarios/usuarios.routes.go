package usuarios

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Usuarios(r fiber.Router) {
	api := r.Group("/users")

	api.Get("/", middlewares.JWTProtected, ListUserLastController)

	//nuevo registro
	api.Post("/",
		middlewares.JWTProtected, middlewares.VallidateBody(&UserRequest{}), CreateUserController)

	//buscar resgistro por ID
	api.Get("/:id", middlewares.JWTProtected, ListUserByIdController)

	api.Put("/:id", middlewares.JWTProtected,
		middlewares.VallidateBody(&UserUpdateRequest{}), UpdateUserController)

	//update password
	api.Put("/users/passw/:id", middlewares.JWTProtected, middlewares.VallidateBody(&NewPasswordUser{}), UpdatePasswordController)

}
