package contabilidad

import (
	"github.com/ecosistema/core/src/modules/contabilidad/costos"
	"github.com/ecosistema/core/src/modules/contabilidad/plan_cuentas"
	"github.com/gofiber/fiber/v2"
)

func RegisterContabilidadRoutes(r fiber.Router) {

	cont := r.Group("/contabilidad")

	costos.Rutas_Costo(cont)
	plan_cuentas.Rutas_cuentas(cont)

}
