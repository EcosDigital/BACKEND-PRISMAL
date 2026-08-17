package periodo_contable

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ─── Errores de dominio ───────────────────────────────────────────────────────

var (
	ErrPeriodosYaExisten     = errors.New("ya existen períodos generados para este año fiscal")
	ErrAñoFiscalCerrado      = errors.New("el año fiscal está cerrado, no se pueden generar períodos")
	ErrPeriodoNoEncontrado   = errors.New("período contable no encontrado")
	ErrPeriodoEstadoInvalido = errors.New("estado de período inválido")
	ErrPeriodoYaEnEseEstado  = errors.New("el período ya se encuentra en ese estado")
	ErrModuloNoEncontrado    = errors.New("módulo no encontrado en este período")
	ErrEstadoModuloInvalido  = errors.New("estado de módulo inválido")
	ErrOperacionNoPermitida  = errors.New("la operación no está permitida en el período/módulo actual")
	ErrAñoFiscalNoEncontrado = errors.New("año fiscal no encontrado")
)

// EstadosPeriodoValidos para validación rápida en el servicio.
var EstadosPeriodoValidos = map[int]bool{
	EstadoPeriodoAbierto:   true,
	EstadoPeriodoCerrado:   true,
	EstadoPeriodoAjuste:    true,
	EstadoPeriodoBloqueado: true,
}

// EstadosModuloValidos para validación rápida.
var EstadosModuloValidos = map[int]bool{
	EstadoModuloAbierto:     true,
	EstadoModuloCerrado:     true,
	EstadoModuloSoloLectura: true,
	EstadoModuloRestringido: true,
}

// ─── Generar períodos ─────────────────────────────────────────────────────────

// GenerarPeriodosAñoFiscal genera los 12 períodos mensuales para un año fiscal.
//
// Reglas:
//  1. El año fiscal debe existir y pertenecer a la empresa.
//  2. No se puede generar dos veces (idempotencia controlada).
//  3. Todos los períodos nacen en estado Cerrado.
//  4. Se crean los controles de módulos por defecto (todos Cerrados).
//  5. La generación es atómica: o se crean los 12 o ninguno.
func GenerarPeriodosAñoFiscal(db *gorm.DB, req *GenerarPeriodosRequest) ([]PeriodoContableResponse, error) {

	// 1. Verificar que el año fiscal existe y pertenece a la empresa
	var añoFiscal struct {
		ID       int64
		Year     int
		IDEstado int
	}
	err := db.Table("contabilidad.cfg_años_fiscales").
		Select("id, year, id_estado").
		Where("id = ? AND empresa_id = ?", req.IDAñoFiscal, req.EmpresaID).
		First(&añoFiscal).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAñoFiscalNoEncontrado
		}
		return nil, err
	}

	// 2. Verificar que no se hayan generado ya los períodos
	existen, err := ExistenPeriodosParaAño(db, req.IDAñoFiscal, req.EmpresaID)
	if err != nil {
		return nil, err
	}
	if existen {
		return nil, ErrPeriodosYaExisten
	}

	// 3. Generación atómica en transacción
	var result []PeriodoContableResponse
	txErr := db.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = CrearPeriodosDelAño(tx, req, añoFiscal.Year)
		return err
	})

	if txErr != nil {
		return nil, txErr
	}

	return result, nil
}

// ─── Listar períodos ──────────────────────────────────────────────────────────

// ObtenerPeriodosPorAño retorna todos los períodos de un año fiscal con módulos.
func ObtenerPeriodosPorAño(db *gorm.DB, idAñoFiscal int64, empresaID int64) ([]PeriodoContableResponse, error) {
	return ListarPeriodosPorAño(db, idAñoFiscal, empresaID)
}

// ─── Cambiar estado de período ────────────────────────────────────────────────

// CambiarEstadoPeriodo aplica reglas de negocio antes de cambiar el estado de un período.
//
// Reglas críticas:
//   - Solo un período ABIERTO por año fiscal a la vez.
//   - Al abrir, se cierra automáticamente cualquier otro abierto del mismo año (transacción atómica).
//   - Períodos Bloqueados no pueden cambiar de estado directamente (requiere flujo de desbloqueo).
//   - Los estados AJUSTE y BLOQUEADO son estados especiales con reglas propias.
func CambiarEstadoPeriodo(
	db *gorm.DB,
	id int64,
	req *UpdateEstadoPeriodoRequest,
) (*PeriodoContableResponse, error) {

	// 1. Validar que el estado solicitado es válido
	if !EstadosPeriodoValidos[req.IDEstado] {
		return nil, ErrPeriodoEstadoInvalido
	}

	// 2. Obtener el período actual
	current, err := getRawPeriodoByID(db, id, req.EmpresaID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrPeriodoNoEncontrado
	}

	// 3. Idempotencia: ya está en el estado solicitado
	if current.IDEstado == req.IDEstado {
		return nil, fmt.Errorf("%w: período %s ya tiene estado %d",
			ErrPeriodoYaEnEseEstado, current.NombreMes, req.IDEstado)
	}

	// 4. Períodos Bloqueados no pueden moverse directamente
	if current.IDEstado == EstadoPeriodoBloqueado && req.IDEstado != EstadoPeriodoCerrado {
		return nil, errors.New("un período Bloqueado solo puede pasar a Cerrado mediante proceso de desbloqueo autorizado")
	}

	// 5. APERTURA: transacción atómica
	if req.IDEstado == EstadoPeriodoAbierto {
		txErr := db.Transaction(func(tx *gorm.DB) error {
			// Cerrar todos los que estén abiertos en el mismo año fiscal
			if err := CerrarPeriodosAbiertosPorAño(tx, current.IDAñoFiscal, req.EmpresaID, id, req.UserID); err != nil {
				return fmt.Errorf("error cerrando períodos previos: %w", err)
			}
			// Abrir el período solicitado y sus módulos configurados
			if err := SetEstadoPeriodo(tx, id, EstadoPeriodoAbierto, req.Notas, req.UserID); err != nil {
				return fmt.Errorf("error abriendo período: %w", err)
			}
			// Al abrir el período, los módulos no se abren automáticamente
			// (control granular independiente: el usuario decide qué módulos habilita)
			return nil
		})

		if txErr != nil {
			return nil, txErr
		}

		return GetPeriodoByID(db, id, req.EmpresaID)
	}

	// 6. Otros cambios de estado (Cerrar, Ajuste, Bloquear)
	if err := SetEstadoPeriodo(db, id, req.IDEstado, req.Notas, req.UserID); err != nil {
		return nil, err
	}

	// Al cerrar/bloquear el período, cerrar también todos los módulos
	if req.IDEstado == EstadoPeriodoCerrado || req.IDEstado == EstadoPeriodoBloqueado {
		if err := cerrarTodosLosModulos(db, id, req.UserID); err != nil {
			return nil, err
		}
	}

	return GetPeriodoByID(db, id, req.EmpresaID)
}

// cerrarTodosLosModulos pone todos los módulos de un período en estado Cerrado.
func cerrarTodosLosModulos(db *gorm.DB, idPeriodo int64, updatedBy int64) error {
	return db.Model(&ControlModuloPeriodo{}).
		Where("id_periodo = ?", idPeriodo).
		Updates(map[string]interface{}{
			"id_estado":  EstadoModuloCerrado,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

// ─── Cambiar estado de módulo ─────────────────────────────────────────────────

// CambiarEstadoModuloPeriodo cambia el estado operacional de un módulo en un período.
//
// Reglas:
//   - El período debe existir y estar en estado que permita modificaciones (Abierto o Ajuste).
//   - El módulo debe existir en el período.
func CambiarEstadoModuloPeriodo(
	db *gorm.DB,
	idPeriodo int64,
	idModulo int,
	req *UpdateEstadoModuloPeriodoRequest,
) (*PeriodoContableResponse, error) {

	// 1. Validar estado de módulo solicitado
	if !EstadosModuloValidos[req.IDEstadoModulo] {
		return nil, ErrEstadoModuloInvalido
	}

	// 2. Verificar que el período existe
	periodo, err := getRawPeriodoByID(db, idPeriodo, req.EmpresaID)
	if err != nil {
		return nil, err
	}
	if periodo == nil {
		return nil, ErrPeriodoNoEncontrado
	}

	// 3. Solo se pueden cambiar módulos en períodos Abiertos o en Ajuste
	if periodo.IDEstado != EstadoPeriodoAbierto && periodo.IDEstado != EstadoPeriodoAjuste {
		return nil, errors.New("solo se pueden modificar módulos en períodos con estado Abierto o Ajuste")
	}

	// 4. Verificar que el módulo existe en el período
	control, err := GetControlModulo(db, idPeriodo, idModulo)
	if err != nil {
		return nil, err
	}
	if control == nil {
		return nil, ErrModuloNoEncontrado
	}

	// 5. Actualizar el estado del módulo
	if err := SetEstadoModuloPeriodo(db, idPeriodo, idModulo, req.IDEstadoModulo, req.Notas, req.UserID); err != nil {
		return nil, err
	}

	return GetPeriodoByID(db, idPeriodo, req.EmpresaID)
}

// ─── Validación centralizada de operaciones ───────────────────────────────────
//
// Este es el servicio más importante del módulo: cualquier otro módulo del ERP
// debe llamar a estas funciones antes de registrar un movimiento.

// CanCreateInvoice valida si se puede crear una factura para una empresa en una fecha.
func CanCreateInvoice(db *gorm.DB, empresaID int64, fecha time.Time) (*ValidacionOperacionResponse, error) {
	return validarOperacion(db, empresaID, fecha, ModuloFacturacion)
}

// CanCreateJournalEntry valida si se puede crear un asiento contable.
func CanCreateJournalEntry(db *gorm.DB, empresaID int64, fecha time.Time) (*ValidacionOperacionResponse, error) {
	return validarOperacion(db, empresaID, fecha, ModuloContabilidad)
}

// CanCreateCreditNote valida si se puede crear una nota crédito.
func CanCreateCreditNote(db *gorm.DB, empresaID int64, fecha time.Time) (*ValidacionOperacionResponse, error) {
	return validarOperacion(db, empresaID, fecha, ModuloFacturacion)
}

// CanCreateInventoryMovement valida si se puede crear un movimiento de inventario.
func CanCreateInventoryMovement(db *gorm.DB, empresaID int64, fecha time.Time) (*ValidacionOperacionResponse, error) {
	return validarOperacion(db, empresaID, fecha, ModuloInventario)
}

// CanCreatePurchaseOrder valida si se puede crear una orden de compra.
func CanCreatePurchaseOrder(db *gorm.DB, empresaID int64, fecha time.Time) (*ValidacionOperacionResponse, error) {
	return validarOperacion(db, empresaID, fecha, ModuloCompras)
}

// validarOperacion es la función núcleo de validación centralizada.
// Agrupa toda la lógica de verificación de período + módulo para extensibilidad.
func validarOperacion(db *gorm.DB, empresaID int64, fecha time.Time, idModulo int) (*ValidacionOperacionResponse, error) {
	puede, motivo, idPeriodo, err := PuedeOperarModulo(db, empresaID, fecha, idModulo)
	if err != nil {
		return nil, err
	}

	resp := &ValidacionOperacionResponse{
		Puede:     puede,
		Motivo:    motivo,
		IDPeriodo: idPeriodo,
	}

	return resp, nil
}
