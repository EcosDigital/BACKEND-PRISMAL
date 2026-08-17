package documentacion_sop

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_SOP registra todos los endpoints del módulo bajo el
// prefijo /sop dentro del grupo de documentación.

// Nota: el endpoint /catalogos se registra ANTES de /:id para evitar
// que Fiber interprete "catalogos" como un parámetro de ruta.

func Rutas_SOP(r fiber.Router) {

	protected := r.Group("/sop",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// GET    /documentacion/sop/catalogos   → catálogos de categorías y estados (ref)
	// Registrado ANTES de /:id para evitar que Fiber interprete "catalogos" como ID
	protected.Get("/catalogos", GetCatalogosSOPController)

	// GET    /documentacion/sop             → listar SOPs con filtros opcionales
	//   Query params: id_categoria, id_estado, area, search
	protected.Get("/", GetSOPsController)

	// POST   /documentacion/sop             → crear nuevo SOP
	protected.Post("/", CreateSOPController)

	// GET    /documentacion/sop/:id         → consultar SOP por ID
	protected.Get("/:id", GetSOPByIDController)

	// PUT    /documentacion/sop/:id         → actualizar SOP completo
	// PUT porque se reemplazan todos los campos editables del recurso.
	protected.Put("/:id", UpdateSOPController)

	// DELETE /documentacion/sop/:id         → eliminar SOP (soft delete)
	protected.Delete("/:id", DeleteSOPController)
}
