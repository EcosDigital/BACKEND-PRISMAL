package plan_cuentas

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func FilterCuentas(db *gorm.DB, empresaID int64) ([]CuentaContableResponse, error) {
	return ListCuentas(db, empresaID)
}

func FilterCuentaByID(db *gorm.DB, id int64) (*CuentaContableResponseFull, error) {
	result, err := GetCuentaByID(db, id)
	if err != nil {
		return nil, err
	}
	if result == nil || result.ID == 0 {
		return nil, errors.New("cuenta contable no encontrada")
	}
	return result, nil
}

// RegisterCuenta crea una nueva cuenta contable aplicando todas las validaciones de negocio.
func RegisterCuenta(db *gorm.DB, req *CuentaContableRequest) (int64, error) {

	// ── 1. Código único por empresa ───────────────────────────────────────────
	exists, err := ExistsCodigoCuenta(db, req.CodigoCuenta, req.EmpresaID, 0)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("ya existe una cuenta con el código '%s' en esta empresa",
			strings.TrimSpace(req.CodigoCuenta))
	}

	// ── 2. Validar cuenta padre si se especificó ──────────────────────────────
	if req.IDCuentaPadre != nil {
		padre, err := GetCuentaByID(db, int64(*req.IDCuentaPadre))
		if err != nil || padre == nil || padre.ID == 0 {
			return 0, errors.New("la cuenta padre especificada no existe")
		}
		// El nivel de la cuenta hija debe ser mayor que el del padre
		if req.IDNivelCuenta <= padre.IDNivelCuenta {
			return 0, fmt.Errorf(
				"el nivel de la cuenta (%d) debe ser mayor al nivel de la cuenta padre (%d)",
				req.IDNivelCuenta, padre.IDNivelCuenta,
			)
		}
	}

	return CreateCuenta(db, req)
}

// EditCuenta actualiza una cuenta contable aplicando validaciones de negocio.
func EditCuenta(db *gorm.DB, id int64, req *CuentaContableUpdateRequest) (int64, error) {

	// ── 1. Verificar que existe ───────────────────────────────────────────────
	current, err := GetCuentaByID(db, id)
	if err != nil || current == nil || current.ID == 0 {
		return 0, errors.New("cuenta contable no encontrada")
	}

	// ── 2. Validar cuenta padre ───────────────────────────────────────────────
	if req.IDCuentaPadre != nil {
		// No puede ser su propio padre
		if int64(*req.IDCuentaPadre) == id {
			return 0, errors.New("una cuenta no puede ser su propio padre")
		}
		padre, err := GetCuentaByID(db, int64(*req.IDCuentaPadre))
		if err != nil || padre == nil || padre.ID == 0 {
			return 0, errors.New("la cuenta padre especificada no existe")
		}
		if req.IDNivelCuenta <= padre.IDNivelCuenta {
			return 0, fmt.Errorf(
				"el nivel de la cuenta (%d) debe ser mayor al nivel de la cuenta padre (%d)",
				req.IDNivelCuenta, padre.IDNivelCuenta,
			)
		}
	}

	// ── 3. Si permite_movimientos = true, verificar que no tenga hijos ────────
	if req.PermiteMovimientos != nil && *req.PermiteMovimientos {
		tieneHijos, err := HasHijos(db, id)
		if err != nil {
			return 0, err
		}
		if tieneHijos {
			return 0, errors.New("no se puede habilitar movimientos en una cuenta que tiene subcuentas")
		}
	}

	return UpdateCuenta(db, id, req)
}

// GetReferencias devuelve todas las referencias necesarias para el formulario
func GetReferencias(db *gorm.DB, empresaID int64, excludeID int64) (map[string]interface{}, error) {
	naturalezas, err := GetNaturalezas(db)
	if err != nil {
		return nil, err
	}

	tipos, err := GetTiposCuenta(db)
	if err != nil {
		return nil, err
	}

	niveles, err := GetNivelesCuenta(db)
	if err != nil {
		return nil, err
	}

	padres, err := GetCuentasPadre(db, empresaID, excludeID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"naturalezas":   naturalezas,
		"tipos_cuenta":  tipos,
		"niveles":       niveles,
		"cuentas_padre": padres,
	}, nil
}

// ─── Carga masiva ─────────────────────────────────────────────────────────────

// ImportCuentas procesa el plano PUC:
//
//  1. Carga el cache de referencias (4 queries totales).
//  2. Procesa las filas en ORDEN — esto es crítico porque una cuenta
//     puede referenciar como padre a otra cuenta del mismo plano
//     que todavía no existía en la BD. El cache se actualiza en memoria
//     al insertar cada cuenta nueva para que las filas siguientes la encuentren.
//  3. Cada fila es independiente: si una falla las demás continúan.
func ImportCuentas(
	db *gorm.DB,
	req *ImportCuentasRequest,
	userID int64,
	empresaID int64,
	sedeID int64,
) (*ImportCuentasResponse, error) {

	if len(req.Filas) == 0 {
		return nil, errors.New("el plano no contiene filas")
	}
	if len(req.Filas) > 10000 {
		return nil, fmt.Errorf("el plano excede el límite de 10.000 filas (recibidas: %d)", len(req.Filas))
	}

	cache, err := LoadCuentasRefCache(db, empresaID)
	if err != nil {
		return nil, fmt.Errorf("error cargando referencias: %w", err)
	}

	resp := &ImportCuentasResponse{
		Detalle: make([]ImportCuentaRowResult, 0, len(req.Filas)),
	}

	for i, fila := range req.Filas {
		filaNum := i + 2 // fila 1 = encabezado → datos desde fila 2

		if validErr := validatePlanoCuentaRow(&fila, filaNum); validErr != nil {
			resp.Errores++
			resp.Detalle = append(resp.Detalle, ImportCuentaRowResult{
				Fila: filaNum, Codigo: fila.CodigoCuenta, Accion: "error", Error: validErr.Error(),
			})
			continue
		}

		accion, upsertErr := BulkUpsertCuenta(db, &fila, cache, userID, empresaID, sedeID)
		if upsertErr != nil {
			resp.Errores++
			resp.Detalle = append(resp.Detalle, ImportCuentaRowResult{
				Fila: filaNum, Codigo: fila.CodigoCuenta, Accion: "error", Error: upsertErr.Error(),
			})
			continue
		}

		if accion == "creada" {
			resp.Creadas++
		} else {
			resp.Actualizadas++
		}
		resp.Detalle = append(resp.Detalle, ImportCuentaRowResult{
			Fila: filaNum, Codigo: fila.CodigoCuenta, Accion: accion,
		})
	}

	return resp, nil
}

// validatePlanoCuentaRow valida los campos obligatorios y formatos antes de tocar la BD.
func validatePlanoCuentaRow(fila *PlanoCuentaRow, filaNum int) error {
	var msgs []string

	if strings.TrimSpace(fila.CodigoCuenta) == "" {
		msgs = append(msgs, "CODIGO_CUENTA es obligatorio")
	}
	if len(strings.TrimSpace(fila.Nombre)) < 3 {
		msgs = append(msgs, "NOMBRE debe tener al menos 3 caracteres")
	}
	if strings.TrimSpace(fila.IDNaturaleza) == "" {
		msgs = append(msgs, "ID_NATURALEZA es obligatorio")
	}
	if strings.TrimSpace(fila.NivelCuenta) == "" {
		msgs = append(msgs, "NIVEL_CUENTA es obligatorio")
	}

	pm := strings.ToLower(strings.TrimSpace(fila.PermiteMovimientos))
	if pm != "si" && pm != "no" {
		msgs = append(msgs, fmt.Sprintf("PERMITE MOVIMIENTOS '%s' no válido (use Si o No)", fila.PermiteMovimientos))
	}
	rt := strings.ToLower(strings.TrimSpace(fila.RequiereTerceros))
	if rt != "si" && rt != "no" {
		msgs = append(msgs, fmt.Sprintf("REQUIERE_TERCEROS '%s' no válido (use Si o No)", fila.RequiereTerceros))
	}
	rc := strings.ToLower(strings.TrimSpace(fila.RequiereCentrosCosto))
	if rc != "si" && rc != "no" {
		msgs = append(msgs, fmt.Sprintf("REQUIERE_CENTROS_COSTO '%s' no válido (use Si o No)", fila.RequiereCentrosCosto))
	}

	estado := strings.ToLower(strings.TrimSpace(fila.Estado))
	if estado != "activo" && estado != "inactivo" {
		msgs = append(msgs, fmt.Sprintf("ESTADO '%s' no válido (use Activo o Inactivo)", fila.Estado))
	}

	if len(msgs) > 0 {
		return fmt.Errorf("fila %d — %s", filaNum, strings.Join(msgs, "; "))
	}
	return nil
}

// ─── Exportación ──────────────────────────────────────────────────────────────

// ExportCuentas obtiene todas las cuentas de la empresa listas para exportar.
// La capa de servicio es un pass-through limpio hacia el repositorio; si en el
// futuro se necesitara filtrar o enriquecer la lista, se haría aquí.
func ExportCuentas(db *gorm.DB, empresaID int64) ([]CuentaExportRow, error) {
	rows, err := GetCuentasExport(db, empresaID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
