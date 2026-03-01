package loader

import (
	core "github.com/ecosistema/core/src/modules/_core"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	core.RegisterCoreRoutes(v1)
	//contabilidad.RegisterContabilidadRoutes(app)
}
