package contabilidad

import (
	"github.com/ecosistema/core/src/modules/contabilidad/comprobantes"
	"github.com/ecosistema/core/src/modules/contabilidad/costos"
	"github.com/gofiber/fiber/v2"
)

func RegisterContabilidadRoutes(app *fiber.App) {

	api := app.Group("/api")
	contabilidad := api.Group("/contabilidad")

	costos.Rutas_Costo(contabilidad)
	comprobantes.Rutas_comprobante(contabilidad)

}
