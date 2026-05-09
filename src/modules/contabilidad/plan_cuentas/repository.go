package plan_cuentas

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ─── SELECT base reutilizable ─────────────────────────────────────────────────

const selectCuentaFull = `
	c.id,
	c.codigo_cuenta,
	c.nombre_cuenta,
	c.id_cuenta_padre,
	COALESCE(p.codigo_cuenta || ' - ' || p.nombre_cuenta, '') AS cuenta_padre,
	c.id_naturaleza,
	n.nombre  AS naturaleza,
	c.id_tipo_cuenta,
	COALESCE(t.nombre, '') AS tipo_cuenta,
	c.id_nivel_cuenta,
	nv.nombre AS nivel_cuenta,
	c.permite_movimientos,
	c.requiere_tercero,
	c.requiere_centro_costo,
	c.is_active`

const joinsCuenta = `
	INNER JOIN contabilidad.ref_naturaleza_contable n  ON n.id  = c.id_naturaleza
	INNER JOIN contabilidad.ref_nivel_cuenta        nv ON nv.id = c.id_nivel_cuenta
	LEFT  JOIN contabilidad.cfg_tipo_cuentas        t  ON t.id  = c.id_tipo_cuenta
	LEFT  JOIN contabilidad.cfg_cuentas_contables   p  ON p.id  = c.id_cuenta_padre`

// ─── Listado ──────────────────────────────────────────────────────────────────

func ListCuentas(db *gorm.DB, empresaID int64) ([]CuentaContableResponse, error) {
	var results []CuentaContableResponse

	err := db.
		Table("contabilidad.cfg_cuentas_contables c").
		Select(selectCuentaFull).
		Joins(joinsCuenta).
		Where("c.id_empresa = ?", empresaID).
		Order("c.id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []CuentaContableResponse{}
	}
	return results, nil
}

// ─── Detalle por ID ───────────────────────────────────────────────────────────

func GetCuentaByID(db *gorm.DB, id int64) (*CuentaContableResponseFull, error) {
	var result CuentaContableResponseFull

	err := db.
		Table("contabilidad.cfg_cuentas_contables c").
		Select(selectCuentaFull).
		Joins(joinsCuenta).
		Where("c.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ─── Crear ────────────────────────────────────────────────────────────────────

func CreateCuenta(db *gorm.DB, req *CuentaContableRequest) (int64, error) {
	data := map[string]interface{}{
		"codigo_cuenta":         strings.TrimSpace(req.CodigoCuenta),
		"nombre_cuenta":         strings.TrimSpace(req.NombreCuenta),
		"id_cuenta_padre":       req.IDCuentaPadre,
		"id_naturaleza":         req.IDNaturaleza,
		"id_tipo_cuenta":        req.IDTipoCuenta,
		"id_nivel_cuenta":       req.IDNivelCuenta,
		"permite_movimientos":   req.PermiteMovimientos,
		"requiere_tercero":      req.RequiereTercero,
		"requiere_centro_costo": req.RequiereCentroCosto,
		"is_active":             req.IsActive,
		"created_at":            time.Now(),
		"created_by":            req.UserID,
		"id_empresa":            req.EmpresaID,
		"id_sede":               req.SedeID,
	}

	tx := db.Table("contabilidad.cfg_cuentas_contables").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

// ─── Actualizar ───────────────────────────────────────────────────────────────

func UpdateCuenta(db *gorm.DB, id int64, req *CuentaContableUpdateRequest) (int64, error) {
	data := map[string]interface{}{
		"nombre_cuenta":         strings.TrimSpace(req.NombreCuenta),
		"id_cuenta_padre":       req.IDCuentaPadre,
		"id_naturaleza":         req.IDNaturaleza,
		"id_tipo_cuenta":        req.IDTipoCuenta,
		"id_nivel_cuenta":       req.IDNivelCuenta,
		"permite_movimientos":   req.PermiteMovimientos,
		"requiere_tercero":      req.RequiereTercero,
		"requiere_centro_costo": req.RequiereCentroCosto,
		"is_active":             req.IsActive,
		"updated_at":            time.Now(),
		"updated_by":            req.UserID,
	}

	tx := db.Table("contabilidad.cfg_cuentas_contables").
		Where("id = ? AND id_empresa = ?", id, req.EmpresaID).
		Updates(data)

	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}
	return id, nil
}

// ─── Validaciones ─────────────────────────────────────────────────────────────

// ExistsCodigoCuenta verifica si ya existe el código en la empresa (excluyendo un ID dado)
func ExistsCodigoCuenta(db *gorm.DB, codigo string, empresaID int64, excludeID int64) (bool, error) {
	var count int64
	query := db.Table("contabilidad.cfg_cuentas_contables").
		Where("codigo_cuenta = ? AND id_empresa = ?", strings.TrimSpace(codigo), empresaID)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// HasHijos verifica si una cuenta tiene cuentas hijas (para validar permite_movimientos)
func HasHijos(db *gorm.DB, id int64) (bool, error) {
	var count int64
	err := db.Table("contabilidad.cfg_cuentas_contables").
		Where("id_cuenta_padre = ?", id).
		Count(&count).Error
	return count > 0, err
}

// ─── Referencias ─────────────────────────────────────────────────────────────

func GetNaturalezas(db *gorm.DB) ([]RefNaturaleza, error) {
	var results []RefNaturaleza
	err := db.Raw(`SELECT id, nombre FROM contabilidad.ref_naturaleza_contable ORDER BY id`).
		Scan(&results).Error
	return results, err
}

func GetTiposCuenta(db *gorm.DB) ([]RefTipoCuenta, error) {
	var results []RefTipoCuenta
	err := db.Raw(`SELECT id, nombre FROM contabilidad.cfg_tipo_cuentas WHERE is_active = true ORDER BY id`).
		Scan(&results).Error
	return results, err
}

func GetNivelesCuenta(db *gorm.DB) ([]RefNivelCuenta, error) {
	var results []RefNivelCuenta
	err := db.Raw(`SELECT id, nombre FROM contabilidad.ref_nivel_cuenta ORDER BY id`).
		Scan(&results).Error
	return results, err
}

// ─── LoadCuentasRefCache ──────────────────────────────────────────────────────
//
// Carga de una sola vez todos los maps nombre→id necesarios para resolver
// el plano PUC sin N+1 queries a la BD.
//
// Maps que construye:
//   - Naturalezas:       "activo" → id, "pasivo" → id ...
//   - TiposCuenta:       "inventario" → id, "caja" → id ...
//   - NivelesCuenta:     "nivel 1" → id  Y  "1" → id  (ambos formatos)
//   - CodigosExistentes: "1105" → id  (códigos ya registrados en la empresa)
//
// CodigosExistentes se actualiza en memoria al insertar cada cuenta nueva,
// permitiendo que las filas siguientes del mismo plano la usen como padre.
func LoadCuentasRefCache(db *gorm.DB, empresaID int64) (*cuentasRefCache, error) {
	cache := &cuentasRefCache{
		Naturalezas:       make(map[string]int),
		TiposCuenta:       make(map[string]int),
		NivelesCuenta:     make(map[string]int),
		CodigosExistentes: make(map[string]int),
	}

	type row struct {
		ID     int
		Nombre string
	}

	// ── Naturalezas ───────────────────────────────────────────────────────────
	var nats []row
	if err := db.Raw(
		`SELECT id, nombre FROM contabilidad.ref_naturaleza_contable`,
	).Scan(&nats).Error; err != nil {
		return nil, err
	}
	for _, n := range nats {
		cache.Naturalezas[strings.ToLower(strings.TrimSpace(n.Nombre))] = n.ID
	}

	// ── Tipos de cuenta ───────────────────────────────────────────────────────
	var tipos []row
	if err := db.Raw(
		`SELECT id, nombre FROM contabilidad.cfg_tipo_cuentas WHERE is_active = true`,
	).Scan(&tipos).Error; err != nil {
		return nil, err
	}
	for _, t := range tipos {
		cache.TiposCuenta[strings.ToLower(strings.TrimSpace(t.Nombre))] = t.ID
	}

	// ── Niveles — indexar por "nivel 1" Y por "1" ─────────────────────────────
	// El plano puede traer el número puro (cuando Excel escribe un entero)
	// o el texto "Nivel N". Indexamos ambas formas con el mismo ID.
	var niveles []row
	if err := db.Raw(
		`SELECT id, nombre FROM contabilidad.ref_nivel_cuenta ORDER BY id`,
	).Scan(&niveles).Error; err != nil {
		return nil, err
	}
	for i, nv := range niveles {
		keyTexto := strings.ToLower(strings.TrimSpace(nv.Nombre)) // "nivel 1"
		keyNumero := fmt.Sprintf("%d", i+1)                       // "1"
		cache.NivelesCuenta[keyTexto] = nv.ID
		cache.NivelesCuenta[keyNumero] = nv.ID
	}

	// ── Códigos existentes de la empresa ──────────────────────────────────────
	type codigoRow struct {
		ID     int
		Codigo string
	}
	var codigos []codigoRow
	if err := db.Raw(
		`SELECT id, codigo_cuenta AS codigo
		 FROM contabilidad.cfg_cuentas_contables
		 WHERE id_empresa = ?`,
		empresaID,
	).Scan(&codigos).Error; err != nil {
		return nil, err
	}
	for _, c := range codigos {
		cache.CodigosExistentes[strings.ToLower(strings.TrimSpace(c.Codigo))] = c.ID
	}

	return cache, nil
}

// GetCuentasPadre retorna cuentas que pueden ser padre:
// excluye las de Nivel 4 (son hojas) y opcionalmente la cuenta que se está editando
func GetCuentasPadre(db *gorm.DB, empresaID int64, excludeID int64) ([]RefCuentaPadre, error) {
	var results []RefCuentaPadre
	query := db.Table("contabilidad.cfg_cuentas_contables c").
		Select(`c.id, c.codigo_cuenta, c.nombre_cuenta, nv.nombre AS nivel_cuenta`).
		Joins(`INNER JOIN contabilidad.ref_nivel_cuenta nv ON nv.id = c.id_nivel_cuenta`).
		Where("c.id_empresa = ? AND c.is_active = true AND nv.nombre != 'Nivel 4'", empresaID)

	if excludeID > 0 {
		query = query.Where("c.id != ?", excludeID)
	}

	err := query.Order("c.codigo_cuenta ASC").Scan(&results).Error
	if results == nil {
		results = []RefCuentaPadre{}
	}
	return results, err
}

// ─── Carga masiva: upsert de una fila ────────────────────────────────────────

// BulkUpsertCuenta inserta o actualiza una cuenta contable a partir de una fila del plano.
// Retorna la acción ejecutada ("creada"|"actualizada") y el error si ocurre.
//
// Lógica de resolución:
//   - CUENTA_PADRE  → busca en cache.CodigosExistentes por código
//   - ID_NATURALEZA → busca en cache.Naturalezas por nombre (case-insensitive)
//   - ID_TIPO_CUENTA → busca en cache.TiposCuenta por nombre (opcional)
//   - NIVEL_CUENTA  → busca en cache.NivelesCuenta por "1"|"2"|"3"|"4" o "Nivel N"
//   - Si el CODIGO_CUENTA ya existe → UPDATE (solo nombre, config, no el código)
//   - Si no existe → INSERT
func BulkUpsertCuenta(
	db *gorm.DB,
	fila *PlanoCuentaRow,
	cache *cuentasRefCache,
	userID int64,
	empresaID int64,
	sedeID int64,
) (accion string, err error) {

	// ── Resolver naturaleza ────────────────────────────────────────────────────
	naturalezaID, ok := cache.Naturalezas[strings.ToLower(strings.TrimSpace(fila.IDNaturaleza))]
	if !ok {
		return "error", fmt.Errorf("naturaleza '%s' no reconocida", fila.IDNaturaleza)
	}

	// ── Resolver nivel ─────────────────────────────────────────────────────────
	nivelKey := strings.ToLower(strings.TrimSpace(fila.NivelCuenta))
	nivelID, ok := cache.NivelesCuenta[nivelKey]
	if !ok {
		return "error", fmt.Errorf("nivel de cuenta '%s' no reconocido (use 1, 2, 3, 4 o 'Nivel N')", fila.NivelCuenta)
	}

	// ── Resolver tipo de cuenta (opcional) ────────────────────────────────────
	var tipoCuentaID *int
	if strings.TrimSpace(fila.IDTipoCuenta) != "" {
		id, ok := cache.TiposCuenta[strings.ToLower(strings.TrimSpace(fila.IDTipoCuenta))]
		if !ok {
			return "error", fmt.Errorf("tipo de cuenta '%s' no reconocido", fila.IDTipoCuenta)
		}
		tipoCuentaID = &id
	}

	// ── Resolver cuenta padre (opcional — busca por CÓDIGO) ───────────────────
	var cuentaPadreID *int
	if strings.TrimSpace(fila.CuentaPadre) != "" {
		padreID, ok := cache.CodigosExistentes[strings.ToLower(strings.TrimSpace(fila.CuentaPadre))]
		if !ok {
			return "error", fmt.Errorf("cuenta padre con código '%s' no encontrada (¿está en el plano antes que esta fila?)", fila.CuentaPadre)
		}
		cuentaPadreID = &padreID
	}

	// ── Resolver booleanos (Si/No, si/no, SI/NO) ──────────────────────────────
	parseBool := func(s string) bool {
		return strings.ToLower(strings.TrimSpace(s)) == "si"
	}
	permiteMovimientos := parseBool(fila.PermiteMovimientos)
	requiereTercero := parseBool(fila.RequiereTerceros)
	requiereCentroCosto := parseBool(fila.RequiereCentrosCosto)
	isActive := strings.ToLower(strings.TrimSpace(fila.Estado)) != "inactivo"

	now := time.Now()

	// ── ¿Existe el código? → UPDATE o INSERT ──────────────────────────────────
	existingID, found := cache.CodigosExistentes[strings.ToLower(strings.TrimSpace(fila.CodigoCuenta))]

	if found {
		// UPDATE — el código no se modifica
		data := map[string]interface{}{
			"nombre_cuenta":         strings.TrimSpace(fila.Nombre),
			"id_cuenta_padre":       cuentaPadreID,
			"id_naturaleza":         naturalezaID,
			"id_tipo_cuenta":        tipoCuentaID,
			"id_nivel_cuenta":       nivelID,
			"permite_movimientos":   permiteMovimientos,
			"requiere_tercero":      requiereTercero,
			"requiere_centro_costo": requiereCentroCosto,
			"is_active":             isActive,
			"updated_at":            now,
			"updated_by":            userID,
		}
		tx := db.Table("contabilidad.cfg_cuentas_contables").
			Where("id = ?", existingID).Updates(data)
		if tx.Error != nil {
			return "error", tx.Error
		}
		return "actualizada", nil
	}

	// INSERT
	data := map[string]interface{}{
		"codigo_cuenta":         strings.TrimSpace(fila.CodigoCuenta),
		"nombre_cuenta":         strings.TrimSpace(fila.Nombre),
		"id_cuenta_padre":       cuentaPadreID,
		"id_naturaleza":         naturalezaID,
		"id_tipo_cuenta":        tipoCuentaID,
		"id_nivel_cuenta":       nivelID,
		"permite_movimientos":   permiteMovimientos,
		"requiere_tercero":      requiereTercero,
		"requiere_centro_costo": requiereCentroCosto,
		"is_active":             isActive,
		"created_at":            now,
		"created_by":            userID,
		"id_empresa":            empresaID,
		"id_sede":               sedeID,
	}
	tx := db.Table("contabilidad.cfg_cuentas_contables").Create(&data)
	if tx.Error != nil {
		return "error", tx.Error
	}

	// Agregar el nuevo código al cache para que las filas siguientes
	// del mismo plano puedan usarlo como cuenta padre
	var newID int
	if err := db.Raw("SELECT lastval()").Scan(&newID).Error; err == nil {
		cache.CodigosExistentes[strings.ToLower(strings.TrimSpace(fila.CodigoCuenta))] = newID
	}

	return "creada", nil
}

// ─── Exportación ──────────────────────────────────────────────────────────────

// GetCuentasExport retorna todas las cuentas de la empresa con los campos
// requeridos para el informe Excel. El campo created_at se formatea en la
// query como 'DD/MM/YYYY HH24:MI' para facilitar la lectura directa en Excel.
func GetCuentasExport(db *gorm.DB, empresaID int64) ([]CuentaExportRow, error) {
	var results []CuentaExportRow

	err := db.
		Table("contabilidad.cfg_cuentas_contables c").
		Select(`
			c.codigo_cuenta,
			c.nombre_cuenta,
			COALESCE(p.codigo_cuenta || ' - ' || p.nombre_cuenta, '') AS cuenta_padre,
			n.nombre  AS naturaleza,
			COALESCE(t.nombre, '')  AS tipo_cuenta,
			nv.nombre AS nivel_cuenta,
			c.permite_movimientos,
			c.requiere_tercero,
			c.requiere_centro_costo,
			c.is_active,
			TO_CHAR(c.created_at, 'DD/MM/YYYY HH24:MI') AS created_at
		`).
		Joins(joinsCuenta).
		Where("c.id_empresa = ?", empresaID).
		Order("c.codigo_cuenta ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []CuentaExportRow{}
	}
	return results, nil
}
