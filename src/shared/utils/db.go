package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetDB(c *fiber.Ctx) (*gorm.DB, error) {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return nil, errors.New("no se pudo obtener conexión a base de datos")
	}
	return db, nil
}
