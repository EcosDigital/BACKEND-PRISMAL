package onboarding

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_onBoarding(r fiber.Router) {
	api := r.Group("/onboarding")

	//crear nuevo tenat + DB
	api.Post("/tenant",
		middlewares.VallidateBody(&TenantRequest{}), CreateTenantController)

	//buscar los modulos habiles para clientes que estan en produccion
	api.Get("/modules", FindModuloByProduccionController)
	api.Get("/check-domain", CheckDomainController)

}
