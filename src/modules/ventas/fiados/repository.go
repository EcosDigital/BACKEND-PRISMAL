package fiados

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// ─── Cuentas ──────────────────────────────────────────────────────────────────

// FindCuentaByTercero devuelve la cuenta de fiado asociada a un tercero
// en una empresa. Retorna (id_cuenta, encontrado, error).

func FindCuentaByTercero(db *gorm.DB, idTercero, empresaID int64) (int64, bool, error) {
	var id int64
	err := db.Raw(
		`SELECT id FROM fiados.cfg_cuentas_fiado
		 WHERE id_tercero = ? AND id_empresa = ?
		 LIMIT 1`,
		idTercero, empresaID,
	).Scan(&id).Error
	if err != nil {
		return 0, false, err
	}
	return id, id > 0, nil
}

// ListCuentas ejecuta la vista v_cuentas filtrando por empresa.
// Soporta búsqueda opcional por nombre, cédula o teléfono (q).

func ListCuentas(db *gorm.DB, empresaID int64, q string) ([]CuentaFiadoResponse, error) {
	var results []CuentaFiadoResponse

	query := db.Raw(
		`SELECT * FROM fiados.v_cuentas
		 WHERE id_empresa = ?
		   AND (? = '' OR
		        nombre_cliente ILIKE '%' || ? || '%' OR
		        cedula         ILIKE '%' || ? || '%' OR
		        telefono       ILIKE '%' || ? || '%')
		 ORDER BY
		   CASE estado_codigo WHEN '01' THEN 1 WHEN '02' THEN 2 ELSE 3 END,
		   nombre_cliente ASC`,
		empresaID, q, q, q, q,
	)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	if results == nil {
		results = []CuentaFiadoResponse{}
	}
	return results, nil
}

// GetCuentaByID retorna una cuenta desde la vista v_cuentas por id_cuenta.
func GetCuentaByID(db *gorm.DB, idCuenta, empresaID int64) (*CuentaFiadoResponse, error) {
	var result CuentaFiadoResponse
	err := db.Raw(
		`SELECT * FROM fiados.v_cuentas
		 WHERE id_cuenta = ? AND id_empresa = ?`,
		idCuenta, empresaID,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCuenta actualiza límite de crédito, observaciones y QR de cobro.
func UpdateCuenta(db *gorm.DB, idCuenta int64, req *ActualizarCuentaRequest) error {
	data := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": req.UserID,
	}

	if req.LimiteCredito != nil {
		data["limite_credito"] = req.LimiteCredito
	}
	if req.Observaciones != "" {
		data["observaciones"] = req.Observaciones
	}
	if req.QrCobroURL != "" {
		data["qr_cobro_url"] = req.QrCobroURL
	}

	tx := db.Table("fiados.cfg_cuentas_fiado").
		Where("id = ? AND id_empresa = ?", idCuenta, req.EmpresaID).
		Updates(data)

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Movimientos ──────────────────────────────────────────────────────────────

// CallRegistrarFiado ejecuta la función fn_registrar_fiado de PostgreSQL.
func CallRegistrarFiado(db *gorm.DB, req *RegistrarFiadoRequest, idTercero int64) (*FiadoMovResult, error) {
	var result FiadoMovResult

	fecha := req.Fecha
	if fecha == "" {
		fecha = time.Now().Format("2006-01-02")
	}

	err := db.Raw(
		`SELECT * FROM fiados.fn_registrar_fiado(?, ?, ?, ?, ?, ?, ?, ?)`,
		idTercero,
		req.IDComprobante,
		req.Valor,
		req.Observaciones,
		fecha,
		req.UserID,
		req.EmpresaID,
		req.SedeID,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CallRegistrarAbono ejecuta fn_registrar_abono de PostgreSQL.
func CallRegistrarAbono(db *gorm.DB, req *RegistrarAbonoRequest) (*AbonoMovResult, error) {
	var result AbonoMovResult

	fecha := req.Fecha
	if fecha == "" {
		fecha = time.Now().Format("2006-01-02")
	}

	err := db.Raw(
		`SELECT * FROM fiados.fn_registrar_abono(?, ?, ?, ?, ?, ?, ?, ?)`,
		req.IDCuenta,
		req.IDComprobante,
		req.Valor,
		req.Observaciones,
		fecha,
		req.UserID,
		req.EmpresaID,
		req.SedeID,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CallAnularMovimiento ejecuta fn_anular_mov_fiado de PostgreSQL.
func CallAnularMovimiento(db *gorm.DB, idMovimiento int64, req *AnularMovimientoRequest) error {
	return db.Exec(
		`SELECT fiados.fn_anular_mov_fiado(?, ?, ?, ?)`,
		idMovimiento,
		req.Motivo,
		req.UserID,
		req.EmpresaID,
	).Error
}

// ListHistorialCuenta ejecuta fn_historial_cuenta para una cuenta dada.
func ListHistorialCuenta(
	db *gorm.DB,
	idCuenta, empresaID int64,
	fechaInicio, fechaFin string,
) ([]MovimientoFiadoResponse, error) {

	var results []MovimientoFiadoResponse

	// La función acepta NULL para los rangos de fecha
	var fi, ff interface{}
	if fechaInicio != "" {
		fi = fechaInicio
	}
	if fechaFin != "" {
		ff = fechaFin
	}

	err := db.Raw(
		`SELECT * FROM fiados.fn_historial_cuenta(?, ?, ?, ?)`,
		idCuenta, empresaID, fi, ff,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []MovimientoFiadoResponse{}
	}
	return results, nil
}

// GetResumenCartera ejecuta fn_resumen_cartera para el dashboard.
func GetResumenCartera(db *gorm.DB, empresaID int64, sedeID *int64) (*ResumenCarteraResponse, error) {
	var result ResumenCarteraResponse
	err := db.Raw(
		`SELECT * FROM fiados.fn_resumen_cartera(?, ?)`,
		empresaID, sedeID,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
