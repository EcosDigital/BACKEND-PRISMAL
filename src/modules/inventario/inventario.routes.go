package inventario

import (
	"github.com/ecosistema/core/src/modules/inventario/articulos"
	bodega "github.com/ecosistema/core/src/modules/inventario/bodegas"
	"github.com/ecosistema/core/src/modules/inventario/movimientos/bajas"
	"github.com/ecosistema/core/src/modules/inventario/movimientos/entradas"
	"github.com/ecosistema/core/src/modules/inventario/movimientos/traslados"
	"github.com/gofiber/fiber/v2"
)

func RegisterInvetarioRoutes(r fiber.Router) {

	inv := r.Group("/inventario")

	bodega.Rutas_bodegas(inv)      //bodegas
	articulos.Rutas_articulos(inv) //articulos

	entradas.Rutas_entradas(inv)   //entradas
	traslados.Rutas_traslados(inv) //traslados
	bajas.Rutas_bajas(inv)

}
