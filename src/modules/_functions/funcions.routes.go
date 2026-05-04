package internal

import (
	"github.com/ecosistema/core/src/modules/_functions/prospeccion/scrapping"
	"github.com/gofiber/fiber/v2"
)

func RegisterFunctionsRoutes(r fiber.Router) {

	fun := r.Group("/internal")

	scrapping.Rutas_prospeccion(fun)

}
