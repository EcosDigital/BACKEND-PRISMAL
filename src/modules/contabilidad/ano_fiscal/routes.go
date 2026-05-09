package ano_fiscal

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_añosFiscales registra todos los endpoints del módulo bajo el
// prefijo /anos-fiscales dentro del grupo de contabilidad.
//
// Nota: se usa "anos" (sin tilde) en la URL para compatibilidad HTTP estricta.
func Rutas_añosFiscales(r fiber.Router) {
	protected := r.Group("/anos-fiscales",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// GET    /contabilidad/anos-fiscales/estados      → catálogo de estados (ref)
	// Registrado ANTES de /:id para evitar que Fiber interprete "estados" como ID
	protected.Get("/estados", GetEstadosAñoFiscalController)

	// GET    /contabilidad/anos-fiscales              → listar años fiscales de la empresa
	protected.Get("/", GetAñosFiscalesController)

	// POST   /contabilidad/anos-fiscales              → crear nuevo año fiscal
	protected.Post("/", CreateAñoFiscalController)

	// PATCH  /contabilidad/anos-fiscales/:id/estado   → cambiar estado (Abierto|Cerrado)
	// PATCH porque solo se modifica un campo parcial del recurso.
	protected.Patch("/:id/estado", UpdateEstadoAñoFiscalController)
}
