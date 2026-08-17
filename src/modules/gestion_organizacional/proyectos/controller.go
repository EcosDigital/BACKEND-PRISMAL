package proyectos

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// ─── Helpers de respuesta ─────────────────────────────────────────────────────

func ok(c *fiber.Ctx, estado int, mensaje string, datos interface{}) error {
	return c.Status(estado).JSON(APIResponse{Success: true, Message: mensaje, Data: datos})
}

func fallo(c *fiber.Ctx, estado int, err error) error {
	return c.Status(estado).JSON(APIResponse{Success: false, Message: err.Error()})
}

// codigoHTTP mapea errores de dominio a su código HTTP correspondiente.
func codigoHTTP(err error) int {
	switch {
	case errors.Is(err, ErrCodigoProyectoYaExiste),
		errors.Is(err, ErrEstadoProyectoInvalido),
		errors.Is(err, ErrPrioridadInvalida),
		errors.Is(err, ErrFechaInicioInvalida),
		errors.Is(err, ErrFechaFinInvalida),
		errors.Is(err, ErrFechasProyectoInvalidas),
		errors.Is(err, ErrNombreProyectoVacio),
		errors.Is(err, ErrNombreParteInteresadaVacio),
		errors.Is(err, ErrComentarioVacio),
		errors.Is(err, ErrFechaActividadInvalida),
		errors.Is(err, ErrTipoEvidenciaInvalido),
		errors.Is(err, ErrNombreArchivoVacio),
		errors.Is(err, ErrRutaArchivoVacia):
		return http.StatusBadRequest
	case errors.Is(err, ErrProyectoNoEncontrado),
		errors.Is(err, ErrParteInteresadaNoEncontrada),
		errors.Is(err, ErrComentarioNoEncontrado),
		errors.Is(err, ErrEvidenciaNoEncontrada):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// CATÁLOGOS
// ═══════════════════════════════════════════════════════════════════════════════

// GetCatalogosController maneja GET /organizacion/proyectos/catalogos
func GetCatalogosController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	catalogos, err := ObtenerCatalogos(db)
	if err != nil {
		logging.Error.Printf("Error cargando catálogos de proyectos: %v", err)
		return fallo(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", catalogos)
}

// ═══════════════════════════════════════════════════════════════════════════════
// PROYECTOS
// ═══════════════════════════════════════════════════════════════════════════════

// CrearProyectoController maneja POST /organizacion/proyectos
func CrearProyectoController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	var req CrearProyectoRequest
	if err := c.BodyParser(&req); err != nil {
		return fallo(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.IDEmpresa = int64(middlewares.GetEmpresaID(c))
	req.IDUsuario = int64(middlewares.GetUserID(c))

	resultado, err := RegistrarProyecto(db, &req)
	if err != nil {
		logging.Error.Printf("Error creando proyecto: %v", err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusCreated, "Proyecto creado exitosamente", resultado)
}

// GetProyectosController maneja GET /organizacion/proyectos
// Query params opcionales: id_estado, id_prioridad, id_responsable, busqueda
func GetProyectosController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))

	idEstado, _ := strconv.Atoi(c.Query("id_estado"))
	idPrioridad, _ := strconv.Atoi(c.Query("id_prioridad"))
	idResponsable, _ := strconv.ParseInt(c.Query("id_responsable"), 10, 64)

	filtros := FiltrosProyecto{
		IDEstado:      idEstado,
		IDPrioridad:   idPrioridad,
		IDResponsable: idResponsable,
		Busqueda:      c.Query("busqueda"),
	}

	lista, err := ConsultarProyectos(db, idEmpresa, filtros)
	if err != nil {
		logging.Error.Printf("Error listando proyectos: %v", err)
		return fallo(c, http.StatusInternalServerError, err)
	}

	return ok(c, http.StatusOK, "", lista)
}

// GetProyectoByIDController maneja GET /organizacion/proyectos/:id
func GetProyectoByIDController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))

	resultado, err := ConsultarProyecto(db, id, idEmpresa)
	if err != nil {
		logging.Error.Printf("Error consultando proyecto %d: %v", id, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "", resultado)
}

// ActualizarProyectoController maneja PUT /organizacion/proyectos/:id
func ActualizarProyectoController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	var req ActualizarProyectoRequest
	if err := c.BodyParser(&req); err != nil {
		return fallo(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.IDEmpresa = int64(middlewares.GetEmpresaID(c))
	req.IDUsuario = int64(middlewares.GetUserID(c))

	resultado, err := ModificarProyecto(db, id, &req)
	if err != nil {
		logging.Error.Printf("Error actualizando proyecto %d: %v", id, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "Proyecto actualizado exitosamente", resultado)
}

// EliminarProyectoController maneja DELETE /organizacion/proyectos/:id
func EliminarProyectoController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))
	idUsuario := int64(middlewares.GetUserID(c))

	if err := EliminarProyecto(db, id, idEmpresa, idUsuario); err != nil {
		logging.Error.Printf("Error eliminando proyecto %d: %v", id, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "Proyecto eliminado exitosamente", nil)
}

// ═══════════════════════════════════════════════════════════════════════════════
// PARTES INTERESADAS
// ═══════════════════════════════════════════════════════════════════════════════

// CrearParteInteresadaController maneja POST /organizacion/proyectos/:id/partes-interesadas
func CrearParteInteresadaController(c *fiber.Ctx) error {
	idProyecto, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || idProyecto <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	var req CrearParteInteresadaRequest
	if err := c.BodyParser(&req); err != nil {
		return fallo(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.IDProyecto = idProyecto
	req.IDEmpresa = int64(middlewares.GetEmpresaID(c))
	req.IDUsuario = int64(middlewares.GetUserID(c))

	resultado, err := RegistrarParteInteresada(db, &req)
	if err != nil {
		logging.Error.Printf("Error registrando parte interesada en proyecto %d: %v", idProyecto, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusCreated, "Parte interesada registrada exitosamente", resultado)
}

// GetPartesInteresadasController maneja GET /organizacion/proyectos/:id/partes-interesadas
func GetPartesInteresadasController(c *fiber.Ctx) error {
	idProyecto, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || idProyecto <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))

	lista, err := ConsultarPartesInteresadas(db, idProyecto, idEmpresa)
	if err != nil {
		logging.Error.Printf("Error listando partes interesadas del proyecto %d: %v", idProyecto, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "", lista)
}

// ActualizarParteInteresadaController maneja PUT /organizacion/proyectos/:id/partes-interesadas/:pid
func ActualizarParteInteresadaController(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil || pid <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de parte interesada inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	var req ActualizarParteInteresadaRequest
	if err := c.BodyParser(&req); err != nil {
		return fallo(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.IDEmpresa = int64(middlewares.GetEmpresaID(c))
	req.IDUsuario = int64(middlewares.GetUserID(c))

	resultado, err := ModificarParteInteresada(db, pid, &req)
	if err != nil {
		logging.Error.Printf("Error actualizando parte interesada %d: %v", pid, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "Parte interesada actualizada exitosamente", resultado)
}

// EliminarParteInteresadaController maneja DELETE /organizacion/proyectos/:id/partes-interesadas/:pid
func EliminarParteInteresadaController(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil || pid <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de parte interesada inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))
	idUsuario := int64(middlewares.GetUserID(c))

	if err := EliminarParteInteresada(db, pid, idEmpresa, idUsuario); err != nil {
		logging.Error.Printf("Error eliminando parte interesada %d: %v", pid, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "Parte interesada eliminada exitosamente", nil)
}

// ═══════════════════════════════════════════════════════════════════════════════
// COMENTARIOS (BITÁCORA)
// ═══════════════════════════════════════════════════════════════════════════════

// CrearComentarioController maneja POST /organizacion/proyectos/:id/bitacora
func CrearComentarioController(c *fiber.Ctx) error {
	idProyecto, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || idProyecto <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	var req CrearComentarioRequest
	if err := c.BodyParser(&req); err != nil {
		return fallo(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.IDProyecto = idProyecto
	req.IDEmpresa = int64(middlewares.GetEmpresaID(c))
	req.IDUsuario = int64(middlewares.GetUserID(c))

	resultado, err := RegistrarComentario(db, &req)
	if err != nil {
		logging.Error.Printf("Error registrando comentario en proyecto %d: %v", idProyecto, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusCreated, "Comentario registrado exitosamente", resultado)
}

// GetBitacoraController maneja GET /organizacion/proyectos/:id/bitacora
func GetBitacoraController(c *fiber.Ctx) error {
	idProyecto, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || idProyecto <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de proyecto inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))

	lista, err := ConsultarBitacora(db, idProyecto, idEmpresa)
	if err != nil {
		logging.Error.Printf("Error consultando bitácora del proyecto %d: %v", idProyecto, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "", lista)
}

// ═══════════════════════════════════════════════════════════════════════════════
// EVIDENCIAS
// ═══════════════════════════════════════════════════════════════════════════════

// CrearEvidenciaController maneja POST /organizacion/proyectos/:id/bitacora/:cid/evidencias
func CrearEvidenciaController(c *fiber.Ctx) error {
	idComentario, err := strconv.ParseInt(c.Params("cid"), 10, 64)
	if err != nil || idComentario <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de comentario inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	var req CrearEvidenciaRequest
	if err := c.BodyParser(&req); err != nil {
		return fallo(c, http.StatusBadRequest, errors.New("JSON inválido"))
	}

	req.IDComentario = idComentario
	req.IDEmpresa = int64(middlewares.GetEmpresaID(c))
	req.IDUsuario = int64(middlewares.GetUserID(c))

	resultado, err := RegistrarEvidencia(db, &req)
	if err != nil {
		logging.Error.Printf("Error registrando evidencia en comentario %d: %v", idComentario, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusCreated, "Evidencia registrada exitosamente", resultado)
}

// EliminarEvidenciaController maneja DELETE /organizacion/proyectos/:id/bitacora/:cid/evidencias/:eid
func EliminarEvidenciaController(c *fiber.Ctx) error {
	eid, err := strconv.ParseInt(c.Params("eid"), 10, 64)
	if err != nil || eid <= 0 {
		return fallo(c, http.StatusBadRequest, errors.New("ID de evidencia inválido"))
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return fallo(c, http.StatusInternalServerError, err)
	}

	idEmpresa := int64(middlewares.GetEmpresaID(c))

	if err := EliminarEvidencia(db, eid, idEmpresa); err != nil {
		logging.Error.Printf("Error eliminando evidencia %d: %v", eid, err)
		return fallo(c, codigoHTTP(err), err)
	}

	return ok(c, http.StatusOK, "Evidencia eliminada exitosamente", nil)
}
