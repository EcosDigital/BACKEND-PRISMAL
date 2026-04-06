package horarios

import (
	"errors"

	"github.com/ecosistema/core/src/modules/_core/terceros"
	"gorm.io/gorm"
)

// ─── Errores esperados (negocio) ──────────────────────────────────────────────

var (
	ErrCodigoHorarioDuplicado = errors.New("ya existe un horario con ese código")
	ErrHorarioNoEncontrado    = errors.New("horario no encontrado")
	ErrTerceroNoEncontrado    = errors.New("el id_tercero no existe en el sistema")
	ErrHorarioInactivo        = errors.New("no se puede asignar un horario inactivo")
	ErrBloquesSinOrdenUnico   = errors.New("los bloques deben tener un orden único por horario")
)

// ─── Horarios ─────────────────────────────────────────────────────────────────

// nuevo registro
func RegisterHorario(db *gorm.DB, req *HorarioRequest) (int64, error) {

	// Validar unicidad del código
	existe, err := FindHorarioByCodigo(db, req.Codigo, 0)
	if err != nil {
		return 0, err
	}
	if existe {
		return 0, ErrCodigoHorarioDuplicado
	}

	// Validar que los órdenes de los bloques no se repitan
	if err := validarOrdenBloques(req.Bloques); err != nil {
		return 0, err
	}

	var idHorario int64

	err = db.Transaction(func(db *gorm.DB) error {

		id, err := CreateHorario(db, req)
		if err != nil {
			return err
		}

		if err := CreateBloques(db, id, req.Bloques); err != nil {
			return err
		}

		idHorario = id
		return nil
	})

	if err != nil {
		return 0, err
	}

	return idHorario, nil

}

// FilterHorarios devuelve el listado de horarios recientes con límite (10).
func FilterHorarios(db *gorm.DB) ([]HorarioListResponse, error) {
	return FindLastHorarios(db)
}

func FilterHorarioByID(db *gorm.DB, id int64) (*HorarioResponse, error) {

	row, err := FindHorarioByID(db, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrHorarioNoEncontrado
	}

	bloques, err := FindBloquesByHorario(db, id)
	if err != nil {
		return nil, err
	}

	return buildHorarioResponse(row, bloques), nil
}

func EditEstadoHorario(db *gorm.DB, id int64, req *CambioEstadoHorarioRequest) error {

	// Verificar que el horario existe antes de actualizar
	row, err := FindHorarioByID(db, id)
	if err != nil {
		return err
	}
	if row == nil {
		return ErrHorarioNoEncontrado
	}

	return ChangeHorarioEstado(db, id, req.Activo)
}

// UpdateHorario actualiza los datos básicos de un horario existente.
// No modifica los bloques — son inmutables para preservar integridad histórica.
func UpdateHorario(db *gorm.DB, id int64, req *HorarioUpdateRequest) (*HorarioResponse, error) {

	// Verificar que el horario existe
	row, err := FindHorarioByID(db, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrHorarioNoEncontrado
	}

	if err := UpdateHorarioData(db, id, req); err != nil {
		return nil, err
	}

	return FilterHorarioByID(db, id)
}

// ─── Asignaciones ─────────────────────────────────────────────────────────────
func RegisterHorarioTercero(db *gorm.DB, req *AsignacionHorarioRequest) (*AsignacionResponse, error) {

	// Validar que el horario existe y está activo
	horario, err := FindHorarioByID(db, req.IDHorario)
	if err != nil {
		return nil, err
	}
	if horario == nil {
		return nil, ErrHorarioNoEncontrado
	}
	if !horario.Activo {
		return nil, ErrHorarioInactivo
	}

	// Validar que el tercero existe en cfg_terceros
	tercero, err := terceros.FilterTerceroById(db, int64(req.IDTercero))
	if err != nil {
		return nil, err
	}
	if tercero == nil {
		return nil, ErrTerceroNoEncontrado
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return InsertAsignacion(tx, req)
	})
	if err != nil {
		return nil, err
	}

	return &AsignacionResponse{
		IDHorario:   req.IDHorario,
		IDTercero:   req.IDTercero,
		FechaInicio: req.FechaInicio,
	}, nil
}

// FilterAsignacionesByHorario devuelve las asignaciones activas de un horario.
func FilterAsignacionesByHorario(db *gorm.DB, idHorario int64) ([]AsignacionListResponse, error) {
	return FindAsignacionesByHorario(db, idHorario)
}

// RemoveAsignacion desactiva una asignación existente.
func RemoveAsignacion(db *gorm.DB, idAsignacion int64) error {
	return DeleteAsignacion(db, idAsignacion)
}

// ─── Helpers internos ─────────────────────────────────────────────────────────

// validarOrdenBloques verifica que no haya órdenes duplicados en el slice de bloques.
func validarOrdenBloques(bloques []BloqueHorarioRequest) error {
	seen := make(map[int]struct{}, len(bloques))
	for _, b := range bloques {
		if _, existe := seen[b.Orden]; existe {
			return ErrBloquesSinOrdenUnico
		}
		seen[b.Orden] = struct{}{}
	}
	return nil
}

// buildHorarioResponse construye el DTO de respuesta completo a partir de los rows de BD.
func buildHorarioResponse(row *horarioRow, bloques []bloqueRow) *HorarioResponse {

	bloquesResp := make([]BloqueHorarioResponse, 0, len(bloques))
	for _, b := range bloques {
		bloquesResp = append(bloquesResp, BloqueHorarioResponse{
			ID:          b.ID,
			IDHorario:   b.IDHorario,
			HoraEntrada: b.HoraEntrada,
			HoraSalida:  b.HoraSalida,
			Orden:       b.Orden,
			Activo:      b.Activo,
		})
	}

	return &HorarioResponse{
		ID:                row.ID,
		Codigo:            row.Codigo,
		Nombre:            row.Nombre,
		Descripcion:       row.Descripcion,
		MinutosTolerancia: row.MinutosTolerancia,
		Activo:            row.Activo,
		FechaCreacion:     row.FechaCreacion,
		Bloques:           bloquesResp,
	}
}
