package horarios

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_horarios(r fiber.Router) {

	protected := r.Group("/horarios",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	// POST  /control-asistencia/horarios/
	// Crea un horario con sus bloques (transacción atómica)
	protected.Post("/",
		middlewares.VallidateBody(&HorarioRequest{}), CreateHorarioController)

	// GET   /control-asistencia/horarios/
	// Lista horarios recientes (sin bloques)
	protected.Get("/",
		ListHorariosController)

	// GET   /control-asistencia/horarios/:id
	// Detalle de horario con sus bloques activos
	protected.Get("/:id",
		GetHorarioByIDController)

	// PUT   /control-asistencia/horarios/:id
	// Actualiza datos básicos del horario (bloques son inmutables)
	protected.Put("/:id",
		middlewares.VallidateBody(&HorarioUpdateRequest{}), UpdateHorarioController)

	// PATCH /control-asistencia/horarios/:id/estado
	// Activa o inactiva un horario (soft delete)
	protected.Patch("/:id/estado",
		middlewares.VallidateBody(&CambioEstadoHorarioRequest{}), ChangeEstadoHorarioController)

	// GET   /control-asistencia/horarios/:id/asignaciones
	// Lista empleados actualmente asignados a un horario
	protected.Get("/:id/asignaciones", ListAsignacionesHorarioController)

	// DELETE /control-asistencia/horarios/asignaciones/:id
	// Desactiva (soft delete) una asignación específica
	protected.Delete("/asignaciones/:id", DeleteAsignacionController)

	// POST  /control-asistencia/horarios/asignar
	// Asigna un horario a uno o más empleados (id_tercero)
	protected.Post("/asignar",
		middlewares.VallidateBody(&AsignacionHorarioRequest{}), AsignarHorarioController)

}
