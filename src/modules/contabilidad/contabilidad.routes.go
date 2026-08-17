package contabilidad

import (
	"github.com/ecosistema/core/src/modules/contabilidad/ano_fiscal"
	"github.com/ecosistema/core/src/modules/contabilidad/costos"
	"github.com/ecosistema/core/src/modules/contabilidad/periodo_contable"
	"github.com/ecosistema/core/src/modules/contabilidad/plan_cuentas"
	"github.com/gofiber/fiber/v2"
)

func RegisterContabilidadRoutes(r fiber.Router) {

	cont := r.Group("/contabilidad")

	costos.Rutas_Costo(cont)
	plan_cuentas.Rutas_cuentas(cont)
	ano_fiscal.Rutas_añosFiscales(cont)
	periodo_contable.Rutas_periodosContables(cont)

}
