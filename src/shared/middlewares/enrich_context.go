package middlewares

import (
	"context"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func EnrichContext(c *fiber.Ctx) error {

	claimsInterface := c.Locals("user")
	if claimsInterface == nil {
		logging.Error.Printf("Usuario no valido")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "usuario no autenticado",
		})
	}

	//convertir claims
	claims, ok := claimsInterface.(*utils.JWTClaims)
	if !ok {
		logging.Error.Printf("usuario no valido")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "usuario no autenticado",
		})
	}

	//definicion
	AuthUser := AuthUser{
		ID:     claims.UserID,
		Email:  claims.Email,
		Nombre: claims.Nombre,
		RolID:  0,
		Rol:    claims.Rol,
	}

	if claims.RolID != nil {
		AuthUser.RolID = *claims.RolID
	}

	ctx := context.WithValue(c.Context(), userCtxKey, AuthUser)
	c.SetUserContext(ctx)

	c.Locals("authUser", AuthUser)
	c.Locals("userID", claims.UserID)

	if claims.EmpresaID != nil {
		c.Locals("empresaID", *claims.EmpresaID)
	}

	if claims.EmpresaID != nil {
		c.Locals("sedeID", *claims.SedeID)
	}

	if claims.RolID != nil {
		c.Locals("RolID", *claims.RolID)
	}

	return c.Next()

}
