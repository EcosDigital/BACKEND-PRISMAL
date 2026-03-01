package seguridad

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/modules/_core/usuarios"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func SigninController(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Obtener
	user, err := usuarios.FIlterUserByEmail(db, req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	info := &UserInfo{
		UserID: user.ID,
		Email:  user.Email,
		RolID:  &user.IdRol,
	}

	// Generar token inicial
	token, err := GenerateBasicToken(info)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	// Guardar en cookie segura
	c.Cookie(&fiber.Cookie{
		Name:     "ecosis_auth",
		Value:    token,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Path:     "/",
		SameSite: "",
		Secure:   false,
		Domain:   strings.Split(c.Hostname(), ":")[0],
	})

	return c.JSON(fiber.Map{"message": "login ok", "token": token})
}

func AccessGetController(c *fiber.Ctx) error {

	claims := c.Locals("user").(*utils.JWTClaims)

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	access, err := UserAccessAuthorized(db, claims)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(access)

}

func AccessController(c *fiber.Ctx) error {
	//  Parsear payload
	var payload struct {
		IdEmpresa int  `json:"id_empresa"`
		IdSede    *int `json:"id_sede"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//  Obtener claims desde middleware o fallback a cookie
	var claims *utils.JWTClaims
	if v := c.Locals("user"); v != nil {
		claims = v.(*utils.JWTClaims)
	} else {
		tokenCookie := c.Cookies("ecosis_auth")
		if tokenCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "auth cookie missing",
			})
		}
		var err error
		claims, err = utils.ValidateJWT(tokenCookie)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}
	}

	//  Obtener datos completos del usuario
	user, err := usuarios.ListUserByEmail(db, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Validar existencia de empresa/sede (seguridad adicional)
	info := UserInfo{
		UserID:    claims.UserID,
		Email:     user.Email,
		Nombre:    user.Nombre,
		Rol:       user.Rol,
		RolID:     &user.IdRol,
		EmpresaID: &payload.IdEmpresa,
		SedeID:    payload.IdSede,
	}

	// Generar token final (full_authenticated = true)
	tokenFinal, err := GenerateFinalToken(&info)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not generate final token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "ecosis_auth",
		Value:    tokenFinal,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
		SameSite: "",    // importante para cross-site
		Secure:   false, // true en produccion
		Domain:   strings.Split(c.Hostname(), ":")[0],
	})

	return c.JSON(fiber.Map{
		"message": "acceso exitoso",
		"token":   tokenFinal,
	})
}

func VerifyTokenController(c *fiber.Ctx) error {

	var tokenString string

	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		tokenString = strings.TrimPrefix(auth, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
	}

	// Si no está en Authorization, intentar cookie
	if tokenString == "" {
		tokenString = c.Cookies("ecosis_auth")
		tokenString = strings.TrimSpace(tokenString)
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "auth token missing",
		})
	}

	// Validar JWT
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":                 claims.UserID,
			"email":              claims.Email,
			"nombre":             claims.Nombre,
			"rol":                claims.Rol,
			"id_empresa":         claims.EmpresaID,
			"id_sede":            claims.SedeID,
			"full_authenticated": claims.FullAuthenticated,
		},
		"token": tokenString,
	})

}

func SignoutController(c *fiber.Ctx) error {

	// Eliminar cookie
	c.Cookie(&fiber.Cookie{
		Name:     "ecosis_auth",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return c.JSON(fiber.Map{"message": "sesión cerrada exitosamente"})
}

func SignoutAllController(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:     "ecosis_auth",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})

	return c.JSON(fiber.Map{"message": "todas las sesiones cerradas"})
}

func DebugSessionsController(c *fiber.Ctx) error {
	ctx := context.Background()
	pattern := "session:*"

	var sessions []map[string]interface{}
	var cursor uint64

	for {
		keys, nextCursor, err := core.RedisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		for _, key := range keys {
			data, err := core.RedisClient.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var session utils.SessionData
			if err := json.Unmarshal([]byte(data), &session); err != nil {
				continue
			}

			sessions = append(sessions, map[string]interface{}{
				"key":     key,
				"user_id": session.UserID,
				"email":   session.Email,
				"expires": session.ExpiresAt,
			})
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return c.JSON(fiber.Map{
		"total":    len(sessions),
		"sessions": sessions,
	})
}
