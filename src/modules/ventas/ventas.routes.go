package ventas

import (
	"github.com/ecosistema/core/src/modules/ventas/fiados"
	"github.com/ecosistema/core/src/modules/ventas/ordenes"
	"github.com/ecosistema/core/src/modules/ventas/pedidos"
	"github.com/gofiber/fiber/v2"
)

func RegisterVentasRoutes(r fiber.Router) {

	vent := r.Group("/ventas")

	fiados.Rutas_fiados(vent)
	pedidos.Rutas_pedidos(vent)
	ordenes.Rutas_ordenes(vent)

}
