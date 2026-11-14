package tokens

import (
	"net/http"

	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateTokenController(c *fiber.Ctx) error {
	var req TokenCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Json invalido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := CreateToken(db, middlewares.GetUserID(c), req.Nombre)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(resp)
}

func ListTokensController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := ListTokensByUser(db, middlewares.GetUserID(c))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func RevokeTokenController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ok, err := RevokeToken(db, c.Params("id"), middlewares.GetUserID(c))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if !ok {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "token no encontrado"})
	}
	return c.JSON(fiber.Map{"message": "token revocado"})
}
