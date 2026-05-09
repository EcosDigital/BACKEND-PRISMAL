package ano_fiscal

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ─── Errores de dominio ───────────────────────────────────────────────────────
// Errores tipados para que el controller distinga entre 400 (negocio) y 500 (infra).

var (
	ErrYearYaExiste          = errors.New("ya existe un año fiscal para ese año en esta empresa")
	ErrAñoFiscalNoEncontrado = errors.New("año fiscal no encontrado")
	ErrEstadoInvalido        = errors.New("estado inválido: use id_estado 1 (Abierto) o 2 (Cerrado)")
	ErrYaEnEseEstado         = errors.New("el año fiscal ya se encuentra en ese estado")
	ErrAsientosPendientes    = errors.New("el año fiscal tiene asientos pendientes y no puede cerrarse")
)

// ─── Crear ────────────────────────────────────────────────────────────────────

// RegistrarAñoFiscal valida reglas de negocio y delega la inserción al repositorio.
//
// Reglas:
//  1. El año debe ser único por empresa.
//  2. Todo año nuevo nace con id_estado = EstadoCerrado (patrón seguro).
func RegistrarAñoFiscal(db *gorm.DB, req *CreateAñoFiscalRequest) (*AñoFiscalResponse, error) {

	existe, err := ExisteYear(db, req.Year, req.EmpresaID, 0)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrYearYaExiste
	}

	return CreateAñoFiscal(db, req)
}

// ─── Listar ───────────────────────────────────────────────────────────────────

// ListadoAñosFiscales retorna todos los años fiscales de la empresa.
func ListadoAñosFiscales(db *gorm.DB, empresaID int64) ([]AñoFiscalResponse, error) {
	return ListarAñosFiscales(db, empresaID)
}

// ─── Cambiar estado ───────────────────────────────────────────────────────────

// CambiarEstadoAñoFiscal aplica todas las reglas de negocio antes de cambiar el estado.
//
// Reglas críticas:
//
//   - Solo puede existir UN año fiscal con id_estado = Abierto por empresa.
//     Al abrir uno, el servicio cierra automáticamente cualquier otro abierto
//     dentro de una transacción atómica (no puede quedar la BD en estado parcial).
//
//   - Antes de cerrar se ejecutan validaciones contables (asientos pendientes, etc.)
//     mediante la función validatePreCierre, diseñada para crecer con el sistema.
//
//   - Se rechaza con error descriptivo si el año ya está en el estado solicitado
//     (idempotencia explícita — el frontend puede mostrar el mensaje al usuario).
func CambiarEstadoAñoFiscal(
	db *gorm.DB,
	id int64,
	req *UpdateEstadoAñoFiscalRequest,
) (*AñoFiscalResponse, error) {

	// ── 1. Validar que el id_estado sea un valor conocido ─────────────────────
	if req.IDEstado != EstadoAbierto && req.IDEstado != EstadoCerrado {
		return nil, ErrEstadoInvalido
	}

	// ── 2. Obtener el año fiscal actual (modelo crudo para comparar id_estado) ─
	current, err := getRawAñoFiscalByID(db, id, req.EmpresaID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrAñoFiscalNoEncontrado
	}

	// ── 3. Idempotencia: ya está en el estado solicitado ──────────────────────
	if current.IDEstado == req.IDEstado {
		return nil, fmt.Errorf("%w: año %d ya tiene id_estado %d",
			ErrYaEnEseEstado, current.Year, req.IDEstado)
	}

	// ── 4. Validaciones previas al CIERRE ─────────────────────────────────────
	if req.IDEstado == EstadoCerrado {
		if err := validatePreCierre(db, current); err != nil {
			return nil, err
		}
	}

	// ── 5. APERTURA: transacción atómica para garantizar año único abierto ────
	if req.IDEstado == EstadoAbierto {
		txErr := db.Transaction(func(tx *gorm.DB) error {

			// 5a. Cerrar cualquier año que esté actualmente abierto
			if err := CerrarTodosLosAbiertosExcepto(tx, req.EmpresaID, id, req.UserID); err != nil {
				return fmt.Errorf("error cerrando año fiscal previo: %w", err)
			}

			// 5b. Abrir el año solicitado
			if err := SetEstadoAñoFiscal(tx, id, EstadoAbierto, req.UserID); err != nil {
				return fmt.Errorf("error abriendo año fiscal: %w", err)
			}

			return nil
		})

		if txErr != nil {
			return nil, txErr
		}

		return GetAñoFiscalByID(db, id, req.EmpresaID)
	}

	// ── 6. CIERRE simple (ya pasó validatePreCierre) ──────────────────────────
	if err := SetEstadoAñoFiscal(db, id, EstadoCerrado, req.UserID); err != nil {
		return nil, err
	}

	return GetAñoFiscalByID(db, id, req.EmpresaID)
}

// ─── Validaciones previas al cierre ──────────────────────────────────────────

// validatePreCierre agrupa todas las validaciones que deben pasar antes de
// cerrar un año fiscal. Diseñado para crecer con cada módulo contable nuevo.
func validatePreCierre(db *gorm.DB, af *AñoFiscal) error {

	// Asientos pendientes (placeholder activo)
	pendientes, err := TieneAsientosPendientes(db, af.ID)
	if err != nil {
		return err
	}
	if pendientes {
		return ErrAsientosPendientes
	}

	// TODO: HasPendingBankReconciliations(db, af.ID) — módulo conciliación bancaria
	// TODO: HasUnclosedBalances(db, af.ID)          — módulo saldos de cierre

	return nil
}
