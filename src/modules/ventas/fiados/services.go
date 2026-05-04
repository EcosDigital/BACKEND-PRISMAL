package fiados

import (
	"errors"
	"strings"

	"github.com/ecosistema/core/src/modules/_core/terceros"
	"gorm.io/gorm"
)

// ─── Cuentas ──────────────────────────────────────────────────────────────────

// GetCuentas devuelve el listado de cuentas filtrado por empresa y búsqueda
// opcional (nombre, cédula o teléfono del cliente).
func GetCuentas(db *gorm.DB, empresaID int64, q string) ([]CuentaFiadoResponse, error) {
	return ListCuentas(db, empresaID, q)
}

// GetCuenta devuelve una cuenta por ID validando que pertenezca a la empresa.
func GetCuenta(db *gorm.DB, idCuenta, empresaID int64) (*CuentaFiadoResponse, error) {
	cuenta, err := GetCuentaByID(db, idCuenta, empresaID)
	if err != nil {
		return nil, err
	}
	if cuenta == nil || cuenta.IDCuenta == 0 {
		return nil, errors.New("cuenta no encontrada")
	}
	return cuenta, nil
}

// EditCuenta actualiza límite de crédito, observaciones y QR de cobro.
func EditCuenta(db *gorm.DB, idCuenta int64, req *ActualizarCuentaRequest) error {
	// Verificar que la cuenta existe y pertenece a la empresa
	cuenta, err := GetCuentaByID(db, idCuenta, req.EmpresaID)
	if err != nil {
		return err
	}
	if cuenta == nil || cuenta.IDCuenta == 0 {
		return errors.New("cuenta no encontrada")
	}
	return UpdateCuenta(db, idCuenta, req)
}

// ─── Movimientos ──────────────────────────────────────────────────────────────

// RegistrarFiado orquesta el registro de un nuevo fiado.
//
// Lógica de resolución del cliente:
//  - Si llega id_cuenta > 0: usa esa cuenta directamente.
//  - Si no: busca el tercero por cédula en cfg_terceros.
//    - Si no existe: lo crea (nombre, cédula, teléfono, dirección son obligatorios).
//    - fn_registrar_fiado crea la cuenta si el tercero no tiene una aún.

func RegistrarFiado(db *gorm.DB, req *RegistrarFiadoRequest) (*FiadoMovResult, error) {

	var idTercero int64

	if req.IDCuenta > 0 {
		// Cliente identificado por cuenta — leer el id_tercero
		cuenta, err := GetCuentaByID(db, req.IDCuenta, req.EmpresaID)
		if err != nil {
			return nil, err
		}
		if cuenta == nil || cuenta.IDCuenta == 0 {
			return nil, errors.New("cuenta no encontrada")
		}
		idTercero = cuenta.IDTercero
	} else {
		// Cliente identificado por cédula
		if strings.TrimSpace(req.Cedula) == "" {
			return nil, errors.New("debe proporcionar id_cuenta o la cédula del cliente")
		}

		tercerosList, err := terceros.FindTerceroByDocument(db, req.Cedula)
		if err != nil {
			return nil, err
		}

		if len(tercerosList) > 0 {
			idTercero = int64(tercerosList[0].ID)
		} else {
			// Validar campos obligatorios para el nuevo cliente
			if strings.TrimSpace(req.NombreCliente) == "" {
				return nil, errors.New("nombre_cliente es obligatorio para registrar un cliente nuevo")
			}

			estado := true

			nuevoTercero := &terceros.TerceroRequest{
				IdTipoPersona:      1,
				IdTipoDocumento:    2,
				NumeroDocumento:    req.Cedula,
				Dv:                 "",
				PrimerNombre:       req.NombreCliente,
				SegundoNombre:      "",
				PrimerApellido:     req.ApellidoCliente,
				SegundoApellido:    "",
				RazonSocial:        "",
				RepresentanteLegal: "",
				IdGenero:           3,
				Telefono:           req.Telefono,
				Telefono_2:         "",
				Email:              "",
				PaginaWeb:          "",
				Direccion:          req.Direccion,
				IdPais:             1,
				IdDepartamento:     10,
				IdCiudad:           429,
				IdZona:             1,
				Estado:             &estado,
				UserID:             req.UserID,
			}

			newID, err := terceros.CreateTercero(db, nuevoTercero)
			if err != nil {
				return nil, err
			}
			idTercero = newID
		}
	}

	return CallRegistrarFiado(db, req, idTercero)
}

// RegistrarAbono registra un pago sobre una cuenta existente.
func RegistrarAbono(db *gorm.DB, req *RegistrarAbonoRequest) (*AbonoMovResult, error) {
	// Verificar que la cuenta existe y pertenece a la empresa
	cuenta, err := GetCuentaByID(db, req.IDCuenta, req.EmpresaID)
	if err != nil {
		return nil, err
	}
	if cuenta == nil || cuenta.IDCuenta == 0 {
		return nil, errors.New("cuenta no encontrada")
	}
	return CallRegistrarAbono(db, req)
}

// AnularMovimiento anula un movimiento de fiado o abono.
func AnularMovimiento(db *gorm.DB, idMovimiento int64, req *AnularMovimientoRequest) error {
	if strings.TrimSpace(req.Motivo) == "" {
		return errors.New("el motivo de anulación es obligatorio")
	}
	return CallAnularMovimiento(db, idMovimiento, req)
}

// GetHistorial retorna el historial de movimientos de una cuenta.
func GetHistorial(
	db *gorm.DB,
	idCuenta, empresaID int64,
	fechaInicio, fechaFin string,
) ([]MovimientoFiadoResponse, error) {

	// Validar que la cuenta pertenece a la empresa antes de exponer el historial
	cuenta, err := GetCuentaByID(db, idCuenta, empresaID)
	if err != nil {
		return nil, err
	}
	if cuenta == nil || cuenta.IDCuenta == 0 {
		return nil, errors.New("cuenta no encontrada")
	}

	return ListHistorialCuenta(db, idCuenta, empresaID, fechaInicio, fechaFin)
}

// GetResumen retorna el resumen de cartera para el dashboard.
func GetResumen(db *gorm.DB, empresaID int64, sedeID *int64) (*ResumenCarteraResponse, error) {
	return GetResumenCartera(db, empresaID, sedeID)
}
