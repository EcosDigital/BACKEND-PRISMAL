package middlewares

import (
	"context"
	"strings"

	"github.com/ecosistema/core/src/modules/_core/usuarios/tokens"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// TokenConexionProtected valida un token de conexión (prm_...) enviado como
// "Authorization: Bearer <token>". El token deja de servir si fue revocado o
// si el usuario dueño está desactivado.
func TokenConexionProtected(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodOptions {
		return c.Next()
	}

	token := strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if token == "" || !strings.HasPrefix(token, "prm_") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token de conexión requerido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := tokens.ValidateToken(db, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo validar el token"})
	}
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token inválido"})
	}

	authUser := AuthUser{ID: user.ID, Email: user.Email, RolID: user.IDRol, Rol: user.Rol}
	c.SetUserContext(context.WithValue(c.Context(), userCtxKey, authUser))
	c.Locals("authUser", authUser)
	c.Locals("userID", user.ID)
	c.Locals("RolID", user.IDRol)
	if user.IDEmpresa != nil {
		c.Locals("empresaID", *user.IDEmpresa)
	}

	return c.Next()
}
