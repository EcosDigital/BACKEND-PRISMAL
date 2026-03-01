package plan_cuentas

import (
	"time"

	"github.com/ecosistema/core/src/database"
)

func CreateCuentaContable(req *CuentaContableRequest) (int64, error) {

	data := map[string]interface{}{
		"codigo_cuenta":         req.CodigoCuenta,
		"nombre_cuenta":         req.NombreCuenta,
		"id_cuenta_padre":       req.IDCuentaPadre,
		"id_naturaleza":         req.IDNaturaleza,
		"id_tipo_cuenta":        req.IDTipoCuenta,
		"id_nivel_cuenta":       req.IDNivelCuenta,
		"permite_movimientos":   req.PermiteMovimientos,
		"requiere_tercero":      req.ReqTercero,
		"requiere_centro_costo": req.ReqCentroCosto,
		"is_active":             req.IsActive,
		"created_at":            time.Now(),
		"created_by":            req.UserID,
		"id_empresa":            req.EmpresaID,
		"id_sede":               req.SedeID,
	}

	err := database.GormDB.
		Table("contabilidad.cfg_cuentas_contables").
		Create(&data)

	if err.Error != nil {
		return 0, err.Error
	}

	id, ok := data["id"].(int64)
	if !ok {
		return 0, nil
	}

	return id, nil

}
