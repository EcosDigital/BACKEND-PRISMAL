package proyectos

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ─── Errores de dominio ───────────────────────────────────────────────────────
// Errores tipados para que el controller distinga entre 400 (negocio) y 500 (infra).

var (
	// Proyectos
	ErrProyectoNoEncontrado    = errors.New("proyecto no encontrado")
	ErrCodigoProyectoYaExiste  = errors.New("ya existe un proyecto con ese código en esta empresa")
	ErrEstadoProyectoInvalido  = errors.New("el estado especificado no existe o no está activo")
	ErrPrioridadInvalida       = errors.New("la prioridad especificada no existe o no está activa")
	ErrFechaInicioInvalida     = errors.New("la fecha de inicio tiene formato inválido, use YYYY-MM-DD")
	ErrFechaFinInvalida        = errors.New("la fecha de fin tiene formato inválido, use YYYY-MM-DD")
	ErrFechasProyectoInvalidas = errors.New("la fecha de fin debe ser igual o posterior a la fecha de inicio")
	ErrNombreProyectoVacio     = errors.New("el nombre del proyecto no puede estar vacío")

	// Partes interesadas
	ErrParteInteresadaNoEncontrada = errors.New("parte interesada no encontrada")
	ErrNombreParteInteresadaVacio  = errors.New("el nombre de la parte interesada no puede estar vacío")

	// Comentarios
	ErrComentarioNoEncontrado = errors.New("comentario no encontrado")
	ErrComentarioVacio        = errors.New("el comentario no puede estar vacío")
	ErrFechaActividadInvalida = errors.New("la fecha de actividad tiene formato inválido, use YYYY-MM-DD")

	// Evidencias
	ErrEvidenciaNoEncontrada = errors.New("evidencia no encontrada")
	ErrTipoEvidenciaInvalido = errors.New("el tipo de evidencia especificado no existe o no está activo")
	ErrNombreArchivoVacio    = errors.New("el nombre del archivo no puede estar vacío")
	ErrRutaArchivoVacia      = errors.New("la ruta del archivo no puede estar vacía")
)

// ═══════════════════════════════════════════════════════════════════════════════
// CATÁLOGOS
// ═══════════════════════════════════════════════════════════════════════════════

// ObtenerCatalogos retorna estados y prioridades en una sola llamada al frontend.
func ObtenerCatalogos(db *gorm.DB) (*CatalogoResponse, error) {
	estados, err := ObtenerEstadosProyecto(db)
	if err != nil {
		return nil, err
	}

	prioridades, err := ObtenerPrioridadesProyecto(db)
	if err != nil {
		return nil, err
	}

	return &CatalogoResponse{
		Estados:     estados,
		Prioridades: prioridades,
	}, nil
}

// ═══════════════════════════════════════════════════════════════════════════════
// PROYECTOS
// ═══════════════════════════════════════════════════════════════════════════════

// RegistrarProyecto valida las reglas de negocio y delega la inserción al repositorio.
//
// Reglas:
//  1. El código debe ser único por empresa.
//  2. La prioridad debe existir en el catálogo.
//  3. El nombre no puede estar vacío.
//  4. La fecha de inicio debe tener formato válido.
//  5. Si se proporciona fecha de fin, debe ser >= fecha de inicio.
//  6. Todo proyecto nace con id_estado = EstadoPlaneacion (patrón seguro).

func RegistrarProyecto(db *gorm.DB, req *CrearProyectoRequest) (*ProyectoResponse, error) {

	// ── 1. Código único por empresa ───────────────────────────────────────────
	existe, err := ExisteCodigoProyecto(db, req.Codigo, req.IDEmpresa, 0)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrCodigoProyectoYaExiste
	}

	// ── 2. Prioridad válida ───────────────────────────────────────────────────
	prioridadOk, err := ExistePrioridad(db, req.IDPrioridad)
	if err != nil {
		return nil, err
	}
	if !prioridadOk {
		return nil, ErrPrioridadInvalida
	}

	// ── 3. Nombre no vacío ────────────────────────────────────────────────────
	if req.Nombre == "" {
		return nil, ErrNombreProyectoVacio
	}

	// ── 4 y 5. Fechas consistentes ────────────────────────────────────────────
	if err := validarFechasProyecto(req.FechaInicio, req.FechaFin); err != nil {
		return nil, err
	}

	return InsertarProyecto(db, req)
}

// ConsultarProyectos retorna los proyectos activos de la empresa con filtros opcionales.
func ConsultarProyectos(db *gorm.DB, idEmpresa int64, filtros FiltrosProyecto) ([]ProyectoResponse, error) {
	return ListarProyectos(db, idEmpresa, filtros)
}

// ConsultarProyecto retorna un proyecto por ID validando que pertenezca a la empresa.
func ConsultarProyecto(db *gorm.DB, id int64, idEmpresa int64) (*ProyectoResponse, error) {
	resultado, err := BuscarProyectoPorID(db, id, idEmpresa)
	if err != nil {
		return nil, err
	}
	if resultado == nil {
		return nil, ErrProyectoNoEncontrado
	}
	return resultado, nil
}

// ModificarProyecto valida las reglas de negocio y aplica los cambios.
//
// Reglas:
//  1. El proyecto debe existir y estar activo en la empresa.
//  2. El estado y la prioridad deben ser válidos en sus catálogos.
//  3. El nombre no puede quedar vacío.
//  4. Las fechas deben ser consistentes entre sí.

func ModificarProyecto(db *gorm.DB, id int64, req *ActualizarProyectoRequest) (*ProyectoResponse, error) {

	// ── 1. Existencia del proyecto ────────────────────────────────────────────
	actual, err := buscarProyectoCrudo(db, id, req.IDEmpresa)
	if err != nil {
		return nil, err
	}
	if actual == nil {
		return nil, ErrProyectoNoEncontrado
	}

	// ── 2a. Estado válido ─────────────────────────────────────────────────────
	estadoOk, err := ExisteEstadoProyecto(db, req.IDEstado)
	if err != nil {
		return nil, err
	}
	if !estadoOk {
		return nil, fmt.Errorf("%w: id_estado %d", ErrEstadoProyectoInvalido, req.IDEstado)
	}

	// ── 2b. Prioridad válida ──────────────────────────────────────────────────
	prioridadOk, err := ExistePrioridad(db, req.IDPrioridad)
	if err != nil {
		return nil, err
	}
	if !prioridadOk {
		return nil, ErrPrioridadInvalida
	}

	// ── 3. Nombre no vacío ────────────────────────────────────────────────────
	if req.Nombre == "" {
		return nil, ErrNombreProyectoVacio
	}

	// ── 4. Fechas consistentes ────────────────────────────────────────────────
	if err := validarFechasProyecto(req.FechaInicio, req.FechaFin); err != nil {
		return nil, err
	}

	return ActualizarProyecto(db, id, req)
}

// EliminarProyecto verifica existencia y aplica el soft delete.
func EliminarProyecto(db *gorm.DB, id int64, idEmpresa int64, idUsuario int64) error {
	actual, err := buscarProyectoCrudo(db, id, idEmpresa)
	if err != nil {
		return err
	}
	if actual == nil {
		return ErrProyectoNoEncontrado
	}
	return EliminarProyectoLogico(db, id, idEmpresa, idUsuario)
}

// ═══════════════════════════════════════════════════════════════════════════════
// PARTES INTERESADAS
// ═══════════════════════════════════════════════════════════════════════════════

// RegistrarParteInteresada verifica que el proyecto exista y registra la parte interesada.

func RegistrarParteInteresada(db *gorm.DB, req *CrearParteInteresadaRequest) (*ParteInteresadaResponse, error) {

	// Verificar que el proyecto padre existe y pertenece a la empresa
	proyecto, err := buscarProyectoCrudo(db, req.IDProyecto, req.IDEmpresa)
	if err != nil {
		return nil, err
	}
	if proyecto == nil {
		return nil, ErrProyectoNoEncontrado
	}

	if req.Nombre == "" {
		return nil, ErrNombreParteInteresadaVacio
	}

	return InsertarParteInteresada(db, req)
}

// ConsultarPartesInteresadas retorna las partes interesadas activas de un proyecto.
func ConsultarPartesInteresadas(db *gorm.DB, idProyecto int64, idEmpresa int64) ([]ParteInteresadaResponse, error) {
	proyecto, err := buscarProyectoCrudo(db, idProyecto, idEmpresa)
	if err != nil {
		return nil, err
	}
	if proyecto == nil {
		return nil, ErrProyectoNoEncontrado
	}
	return ListarPartesInteresadas(db, idProyecto, idEmpresa)
}

// ConsultarParteInteresada retorna una parte interesada por ID.
func ConsultarParteInteresada(db *gorm.DB, id int64, idEmpresa int64) (*ParteInteresadaResponse, error) {
	resultado, err := BuscarParteInteresadaPorID(db, id, idEmpresa)
	if err != nil {
		return nil, err
	}
	if resultado == nil {
		return nil, ErrParteInteresadaNoEncontrada
	}
	return resultado, nil
}

// ModificarParteInteresada valida existencia y aplica los cambios.
func ModificarParteInteresada(db *gorm.DB, id int64, req *ActualizarParteInteresadaRequest) (*ParteInteresadaResponse, error) {
	actual, err := buscarParteInteresadaCruda(db, id, req.IDEmpresa)
	if err != nil {
		return nil, err
	}
	if actual == nil {
		return nil, ErrParteInteresadaNoEncontrada
	}

	if req.Nombre == "" {
		return nil, ErrNombreParteInteresadaVacio
	}

	return ActualizarParteInteresada(db, id, req)
}

// EliminarParteInteresada verifica existencia y aplica el soft delete.
func EliminarParteInteresada(db *gorm.DB, id int64, idEmpresa int64, idUsuario int64) error {
	actual, err := buscarParteInteresadaCruda(db, id, idEmpresa)
	if err != nil {
		return err
	}
	if actual == nil {
		return ErrParteInteresadaNoEncontrada
	}
	return EliminarParteInteresadaLogica(db, id, idEmpresa, idUsuario)
}

// ═══════════════════════════════════════════════════════════════════════════════
// COMENTARIOS (BITÁCORA)
// ═══════════════════════════════════════════════════════════════════════════════

// RegistrarComentario verifica que el proyecto exista y persiste el comentario.
func RegistrarComentario(db *gorm.DB, req *CrearComentarioRequest) (*ComentarioResponse, error) {

	proyecto, err := buscarProyectoCrudo(db, req.IDProyecto, req.IDEmpresa)
	if err != nil {
		return nil, err
	}
	if proyecto == nil {
		return nil, ErrProyectoNoEncontrado
	}

	if req.Comentario == "" {
		return nil, ErrComentarioVacio
	}

	return InsertarComentario(db, req)
}

// ConsultarBitacora retorna la bitácora completa de un proyecto con sus evidencias.
func ConsultarBitacora(db *gorm.DB, idProyecto int64, idEmpresa int64) ([]ComentarioResponse, error) {
	proyecto, err := buscarProyectoCrudo(db, idProyecto, idEmpresa)
	if err != nil {
		return nil, err
	}
	if proyecto == nil {
		return nil, ErrProyectoNoEncontrado
	}
	return ListarComentarios(db, idProyecto, idEmpresa)
}

// ═══════════════════════════════════════════════════════════════════════════════
// EVIDENCIAS
// ═══════════════════════════════════════════════════════════════════════════════

// RegistrarEvidencia verifica que el comentario exista y adjunta la evidencia.
func RegistrarEvidencia(db *gorm.DB, req *CrearEvidenciaRequest) (*EvidenciaResponse, error) {

	comentario, err := buscarComentarioCrudo(db, req.IDComentario, req.IDEmpresa)
	if err != nil {
		return nil, err
	}
	if comentario == nil {
		return nil, ErrComentarioNoEncontrado
	}

	tipoOk, err := ExisteTipoEvidencia(db, req.IDTipo)
	if err != nil {
		return nil, err
	}
	if !tipoOk {
		return nil, ErrTipoEvidenciaInvalido
	}

	if req.NombreArchivo == "" {
		return nil, ErrNombreArchivoVacio
	}
	if req.RutaArchivo == "" {
		return nil, ErrRutaArchivoVacia
	}

	return InsertarEvidencia(db, req)
}

// EliminarEvidencia verifica existencia y realiza borrado físico.
func EliminarEvidencia(db *gorm.DB, id int64, idEmpresa int64) error {
	evidencia, err := buscarEvidenciaCruda(db, id, idEmpresa)
	if err != nil {
		return err
	}
	if evidencia == nil {
		return ErrEvidenciaNoEncontrada
	}
	return EliminarEvidenciaFisica(db, id, idEmpresa)
}

// ─── Helpers internos ─────────────────────────────────────────────────────────

// validarFechasProyecto verifica la consistencia lógica entre fecha_inicio y fecha_fin.
// Comparación lexicográfica ISO 8601 (YYYY-MM-DD).
func validarFechasProyecto(fechaInicio string, fechaFin *string) error {
	if fechaFin == nil {
		return nil
	}
	if *fechaFin < fechaInicio {
		return ErrFechasProyectoInvalidas
	}
	return nil
}
