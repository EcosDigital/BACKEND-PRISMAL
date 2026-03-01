package middlewares

import (
	"strings"

	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func JWTProtectedGeneral(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodOptions {
		return c.Next()
	}

	if err := JWTProtectedBasic(c); err == nil {
		return c.Next()
	}

	if err := JWTProtected(c); err == nil {
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "token inválido",
	})

}

func JWTProtectedBasic(c *fiber.Ctx) error {
	// Preflight: permitimos pasar
	if c.Method() == fiber.MethodOptions {
		return c.Next()
	}

	var token string

	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		token = strings.TrimPrefix(auth, "Bearer ")
		token = strings.TrimSpace(token)
	}

	if token == "" {
		token = c.Cookies("ecosis_auth")
		if token != "" {
			token = strings.TrimSpace(token)
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "auth token missing",
		})
	}

	//validar estructura del jwt
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	// Guardar claims para los controladores
	c.Locals("user", claims)
	return c.Next()
}

func JWTProtected(c *fiber.Ctx) error {

	if c.Method() == fiber.MethodOptions {
		return c.Next()
	}

	var tokenString string

	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		tokenString = strings.TrimPrefix(auth, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
	}

	if tokenString == "" {
		tokenString = c.Cookies("ecosis_auth")
		if tokenString != "" {
			tokenString = strings.TrimSpace(tokenString)
		}
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "token requerired",
		})
	}

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	if !claims.FullAuthenticated {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	// Verificar que tenga empresa seleccionada
	if claims.EmpresaID == nil || *claims.EmpresaID == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "empresa no seleccionada",
		})
	}

	c.Locals("user", claims)

	return c.Next()
}
