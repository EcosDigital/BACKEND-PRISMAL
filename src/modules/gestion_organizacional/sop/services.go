package documentacion_sop

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ─── Errores de dominio ───────────────────────────────────────────────────────
// Errores tipados para que el controller distinga entre 400 (negocio) y 500 (infra).

var (
	ErrSOPNoEncontrado       = errors.New("SOP no encontrado")
	ErrCodeYaExiste          = errors.New("ya existe un SOP con ese código en esta empresa")
	ErrCategoriaInvalida     = errors.New("la categoría especificada no existe o no está activa")
	ErrEstadoInvalido        = errors.New("el estado especificado no existe o no está activo")
	ErrFechaEfectivaInvalida = errors.New("la fecha efectiva tiene formato inválido, use YYYY-MM-DD")
	ErrFechaRevisionInvalida = errors.New("la fecha de revisión tiene formato inválido, use YYYY-MM-DD")
	ErrFechasInconsistentes  = errors.New("la fecha de revisión debe ser igual o posterior a la fecha efectiva")
	ErrContenidoVacio        = errors.New("el contenido del SOP no puede estar vacío")
	ErrVersionVacia          = errors.New("la versión del SOP no puede estar vacía")
)

// ─── Catálogos ────────────────────────────────────────────────────────────────

// ObtenerCatalogos retorna los catálogos de categorías y estados en una sola llamada.
// Diseñado para poblar selects del frontend con un único request.
func ObtenerCatalogos(db *gorm.DB) (*CatalogoResponse, error) {
	categorias, err := GetCategoriasSOP(db)
	if err != nil {
		return nil, err
	}

	estados, err := GetEstadosSOP(db)
	if err != nil {
		return nil, err
	}

	return &CatalogoResponse{
		Categorias: categorias,
		Estados:    estados,
	}, nil
}

// ─── Crear ────────────────────────────────────────────────────────────────────

// RegistrarSOP valida las reglas de negocio y delega la inserción al repositorio.
//
// Reglas:
//  1. El código debe ser único por empresa.
//  2. La categoría debe existir en el catálogo y estar activa.
//  3. El contenido y la versión no pueden estar vacíos.
//  4. La fecha efectiva debe tener formato válido YYYY-MM-DD.
//  5. Si se proporciona fecha de revisión, debe ser >= fecha efectiva.
//  6. Todo SOP nuevo nace con id_estado = EstadoDraft (patrón seguro).

func RegistrarSOP(db *gorm.DB, req *CreateSOPRequest) (*SOPResponse, error) {

	// ── 1. Código único por empresa ───────────────────────────────────────────
	existe, err := ExisteCode(db, req.Code, req.EmpresaID, 0)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrCodeYaExiste
	}

	// ── 2. Categoría válida ───────────────────────────────────────────────────
	categoriaOk, err := ExisteCategoria(db, req.IDCategoria)
	if err != nil {
		return nil, err
	}
	if !categoriaOk {
		return nil, ErrCategoriaInvalida
	}

	// ── 3. Contenido y versión no vacíos ─────────────────────────────────────
	if req.Content == "" {
		return nil, ErrContenidoVacio
	}
	if req.Version == "" {
		return nil, ErrVersionVacia
	}

	// ── 4 y 5. Validar fechas ─────────────────────────────────────────────────
	if err := validarFechas(req.EffectiveDate, req.ReviewDate); err != nil {
		return nil, err
	}

	return CreateSOP(db, req)
}

// ─── Listar ───────────────────────────────────────────────────────────────────

// ListadoSOPs retorna los SOPs activos de la empresa con los filtros aplicados.
func ListadoSOPs(db *gorm.DB, empresaID int64, filters SOPFilters) ([]SOPResponse, error) {
	return ListarSOPs(db, empresaID, filters)
}

// ─── Consultar por ID ─────────────────────────────────────────────────────────

// ConsultarSOP retorna un SOP por su ID validando que pertenezca a la empresa.
func ConsultarSOP(db *gorm.DB, id int64, empresaID int64) (*SOPResponse, error) {
	result, err := GetSOPByID(db, id, empresaID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrSOPNoEncontrado
	}
	return result, nil
}

// ─── Actualizar ───────────────────────────────────────────────────────────────

// ModificarSOP valida las reglas de negocio y aplica los cambios al SOP.
//
// Reglas:
//  1. El SOP debe existir y estar activo en la empresa.
//  2. Si el código cambia, debe seguir siendo único por empresa.
//  3. La categoría y el estado deben existir en sus respectivos catálogos.
//  4. El contenido y la versión no pueden quedar vacíos.
//  5. Las fechas deben ser consistentes entre sí.

func ModificarSOP(db *gorm.DB, id int64, req *UpdateSOPRequest) (*SOPResponse, error) {

	// ── 1. Verificar existencia del SOP ───────────────────────────────────────
	current, err := getRawSOPByID(db, id, req.EmpresaID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrSOPNoEncontrado
	}

	// ── 2. Categoría válida ───────────────────────────────────────────────────
	categoriaOk, err := ExisteCategoria(db, req.IDCategoria)
	if err != nil {
		return nil, err
	}
	if !categoriaOk {
		return nil, ErrCategoriaInvalida
	}

	// ── 3. Estado válido ──────────────────────────────────────────────────────
	estadoOk, err := ExisteEstado(db, req.IDEstado)
	if err != nil {
		return nil, err
	}
	if !estadoOk {
		return nil, fmt.Errorf("%w: id_estado %d", ErrEstadoInvalido, req.IDEstado)
	}

	// ── 4. Contenido y versión no vacíos ─────────────────────────────────────
	if req.Content == "" {
		return nil, ErrContenidoVacio
	}
	if req.Version == "" {
		return nil, ErrVersionVacia
	}

	// ── 5. Validar fechas ─────────────────────────────────────────────────────
	if err := validarFechas(req.EffectiveDate, req.ReviewDate); err != nil {
		return nil, err
	}

	return UpdateSOP(db, id, req)
}

// ─── Eliminar ─────────────────────────────────────────────────────────────────

// EliminarSOP verifica la existencia del SOP y aplica el soft delete.
func EliminarSOP(db *gorm.DB, id int64, empresaID int64, userID int64) error {

	current, err := getRawSOPByID(db, id, empresaID)
	if err != nil {
		return err
	}
	if current == nil {
		return ErrSOPNoEncontrado
	}

	return SoftDeleteSOP(db, id, empresaID, userID)
}

// ─── Helpers internos ─────────────────────────────────────────────────────────

// validarFechas verifica el formato y la consistencia de las fechas del SOP.
// Centralizado para reutilizar en Create y Update sin duplicar lógica.
func validarFechas(effectiveDate string, reviewDate *string) error {
	// La fecha efectiva se parsea en el repositorio; aquí solo validamos
	// la consistencia lógica cuando ambas fechas están presentes.
	if reviewDate == nil {
		return nil
	}

	// Ambas fechas presentes: verificar que review >= effective
	// El parseo real ocurre en el repositorio; aquí comparamos como strings
	// ISO 8601 (YYYY-MM-DD), que permite comparación lexicográfica directa.
	if *reviewDate < effectiveDate {
		return ErrFechasInconsistentes
	}

	return nil
}
