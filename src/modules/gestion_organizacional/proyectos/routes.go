package proyectos

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Nota: /catalogos se registra ANTES de /:id para que Fiber no interprete
// la palabra "catalogos" como un parámetro de ruta numérico.

func Rutas_Proyectos(r fiber.Router) {
	protegido := r.Group("/proyectos",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// GET  /organizacion/proyectos/catalogos
	//      Devuelve estados y prioridades para poblar selects en el frontend.

	protegido.Get("/catalogos", GetCatalogosController)

	// GET    /organizacion/proyectos
	//        Filtros opcionales: id_estado, id_prioridad, id_responsable, busqueda
	protegido.Get("/", GetProyectosController)

	// POST   /organizacion/proyectos
	//        Crea un proyecto nuevo. Nace con estado = Planeación.
	protegido.Post("/", CrearProyectoController)

	// GET    /organizacion/proyectos/:id
	protegido.Get("/:id", GetProyectoByIDController)

	// PUT    /organizacion/proyectos/:id
	//        Actualiza todos los campos editables del proyecto.
	protegido.Put("/:id", ActualizarProyectoController)

	// DELETE /organizacion/proyectos/:id
	//        Eliminación lógica (es_activo = false).
	protegido.Delete("/:id", EliminarProyectoController)

	// ── Partes interesadas ────────────────────────────────────────────────────
	// GET    /organizacion/proyectos/:id/partes-interesadas
	protegido.Get("/:id/partes-interesadas", GetPartesInteresadasController)

	// POST   /organizacion/proyectos/:id/partes-interesadas
	protegido.Post("/:id/partes-interesadas", CrearParteInteresadaController)

	// PUT    /organizacion/proyectos/:id/partes-interesadas/:pid
	protegido.Put("/:id/partes-interesadas/:pid", ActualizarParteInteresadaController)

	// DELETE /organizacion/proyectos/:id/partes-interesadas/:pid
	//        Eliminación lógica (es_activo = false).
	protegido.Delete("/:id/partes-interesadas/:pid", EliminarParteInteresadaController)

	// ── Bitácora de comentarios ───────────────────────────────────────────────
	// GET    /organizacion/proyectos/:id/bitacora
	//        Retorna todos los comentarios con sus evidencias embebidas.
	protegido.Get("/:id/bitacora", GetBitacoraController)

	// POST   /organizacion/proyectos/:id/bitacora
	//        Registra un nuevo comentario. Sin soft delete: historial inmutable.
	protegido.Post("/:id/bitacora", CrearComentarioController)

	// ── Evidencias (sub-recurso de comentario) ────────────────────────────────
	// POST   /organizacion/proyectos/:id/bitacora/:cid/evidencias
	protegido.Post("/:id/bitacora/:cid/evidencias", CrearEvidenciaController)

	// DELETE /organizacion/proyectos/:id/bitacora/:cid/evidencias/:eid
	//        Borrado físico: las evidencias no tienen soft delete.
	protegido.Delete("/:id/bitacora/:cid/evidencias/:eid", EliminarEvidenciaController)
}
