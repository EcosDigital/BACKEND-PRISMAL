package kardex

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_kardex(r fiber.Router) {

	protected := r.Group("/kardex",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// GET /inventario/kardex?id_articulo=&id_bodega=&lote=&fecha_inicio=&fecha_fin=
	protected.Get("/", GetKardexController)

	// GET /inventario/kardex/ref/lotes?id_articulo=&id_bodega=
	protected.Get("/ref/lotes", GetLotesController)
}
