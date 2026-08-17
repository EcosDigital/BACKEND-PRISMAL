package periodo_contable

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_periodosContables registra todos los endpoints del módulo bajo el
// prefijo /periodos dentro del grupo de contabilidad.
//
// Estructura de URLs:
//
//	GET    /contabilidad/periodos/catalogos                              → catálogos de referencia
//	GET    /contabilidad/periodos/validar?modulo=1&fecha=2025-06-15      → validación centralizada
//	POST   /contabilidad/periodos/generar                               → generar 12 períodos de un año fiscal
//	GET    /contabilidad/periodos/:id_año_fiscal                        → listar períodos de un año
//	PATCH  /contabilidad/periodos/:id/estado                            → cambiar estado de período
//	PATCH  /contabilidad/periodos/:id/modulos/:id_modulo/estado         → cambiar estado de módulo en período
func Rutas_periodosContables(r fiber.Router) {
	protected := r.Group("/periodos",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// Rutas estáticas primero (antes de /:id) para evitar conflictos de routing
	protected.Get("/catalogos", GetCatalogosPeriodosController)
	protected.Get("/validar", ValidarOperacionController)
	protected.Post("/generar", GenerarPeriodosController)

	// Rutas con parámetros
	protected.Get("/:id_año_fiscal", GetPeriodosPorAñoController)
	protected.Patch("/:id/estado", UpdateEstadoPeriodoController)
	protected.Patch("/:id/modulos/:id_modulo/estado", UpdateEstadoModuloPeriodoController)
}
