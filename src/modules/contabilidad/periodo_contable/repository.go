package periodo_contable

import (
	"time"

	"gorm.io/gorm"
)

// ─── Helpers de formateo ──────────────────────────────────────────────────────

func fmtFecha(t time.Time) string {
	return t.Format("02/01/2006")
}

func fmtFechaHora(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}

// ─── Referencias / Catálogos ──────────────────────────────────────────────────

// GetCatalogosPeriodos retorna todos los catálogos de referencia en una sola consulta.
func GetCatalogosPeriodos(db *gorm.DB) (*CatalogosResponse, error) {
	var estadosPeriodo []RefEstadoPeriodo
	var modulos []RefModuloContable
	var estadosModulo []RefEstadoModulo

	if err := db.Model(&RefEstadoPeriodo{}).Order("id ASC").Find(&estadosPeriodo).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&RefModuloContable{}).Where("is_active = ?", true).Order("id ASC").Find(&modulos).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&RefEstadoModulo{}).Order("id ASC").Find(&estadosModulo).Error; err != nil {
		return nil, err
	}

	return &CatalogosResponse{
		EstadosPeriodo: estadosPeriodo,
		Modulos:        modulos,
		EstadosModulo:  estadosModulo,
	}, nil
}

// ─── Crear períodos ───────────────────────────────────────────────────────────

// CrearPeriodosDelAño inserta los 12 períodos mensuales para un año fiscal.
// Debe ejecutarse dentro de una transacción.
// Todos los períodos nacen con id_estado = EstadoPeriodoCerrado (patrón seguro).

func CrearPeriodosDelAño(db *gorm.DB, req *GenerarPeriodosRequest, year int) ([]PeriodoContableResponse, error) {
	now := time.Now()
	periodos := make([]PeriodoContable, 12)

	for mes := 1; mes <= 12; mes++ {
		fechaInicio := time.Date(year, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
		// Último día del mes
		fechaFin := fechaInicio.AddDate(0, 1, -1)

		periodos[mes-1] = PeriodoContable{
			IDAñoFiscal: req.IDAñoFiscal,
			EmpresaID:   req.EmpresaID,
			NumeroMes:   mes,
			NombreMes:   NombreMes[mes],
			FechaInicio: fechaInicio,
			FechaFin:    fechaFin,
			IDEstado:    EstadoPeriodoCerrado,
			CreatedBy:   req.UserID,
			UpdatedBy:   req.UserID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	if err := db.Model(&PeriodoContable{}).Create(&periodos).Error; err != nil {
		return nil, err
	}

	// Crear controles de módulos por defecto para cada período
	for _, p := range periodos {
		if err := crearControlesModulosDefault(db, p.ID, req.EmpresaID, req.UserID); err != nil {
			return nil, err
		}
	}

	return ListarPeriodosPorAño(db, req.IDAñoFiscal, req.EmpresaID)
}

// crearControlesModulosDefault crea el registro de control para cada módulo activo
// al inicializar un período, todos en estado Cerrado por defecto.
func crearControlesModulosDefault(db *gorm.DB, idPeriodo int64, empresaID int64, userID int64) error {
	var modulos []RefModuloContable
	if err := db.Session(&gorm.Session{NewDB: true}).
		Model(&RefModuloContable{}).
		Where("is_active = ?", true).
		Find(&modulos).Error; err != nil {
		return err
	}

	now := time.Now()
	controles := make([]ControlModuloPeriodo, len(modulos))
	for i, m := range modulos {
		controles[i] = ControlModuloPeriodo{
			IDPeriodo: idPeriodo,
			IDModulo:  m.ID,
			IDEstado:  EstadoModuloCerrado,
			EmpresaID: empresaID,
			UpdatedBy: userID,
			UpdatedAt: now,
			CreatedAt: now,
		}
	}

	return db.Session(&gorm.Session{NewDB: true}).
		Model(&ControlModuloPeriodo{}).
		Create(&controles).Error
}

// ─── Listar ───────────────────────────────────────────────────────────────────

// ListarPeriodosPorAño retorna todos los períodos de un año fiscal con estado y módulos.
func ListarPeriodosPorAño(db *gorm.DB, idAñoFiscal int64, empresaID int64) ([]PeriodoContableResponse, error) {
	// Query principal: períodos con estado y año resueltos
	type rowPeriodo struct {
		PeriodoContable
		Estado      string `gorm:"column:estado"`
		ColorEstado string `gorm:"column:color_estado"`
		AñoFiscal   int    `gorm:"column:año_fiscal"`
	}

	var rows []rowPeriodo
	err := db.
		Model(&PeriodoContable{}).
		Select(
			"cfg_periodos_contables.id, cfg_periodos_contables.id_año_fiscal, "+
				"cfg_periodos_contables.empresa_id, cfg_periodos_contables.numero_mes, "+
				"cfg_periodos_contables.nombre_mes, cfg_periodos_contables.fecha_inicio, "+
				"cfg_periodos_contables.fecha_fin, cfg_periodos_contables.id_estado, "+
				"cfg_periodos_contables.notas, cfg_periodos_contables.created_by, "+
				"cfg_periodos_contables.updated_by, cfg_periodos_contables.created_at, "+
				"cfg_periodos_contables.updated_at, "+
				"ep.nombre AS estado, ep.color AS color_estado, af.year AS año_fiscal",
		).
		Joins("INNER JOIN contabilidad.ref_estado_periodo ep ON ep.id = cfg_periodos_contables.id_estado").
		Joins("INNER JOIN contabilidad.cfg_años_fiscales af ON af.id = cfg_periodos_contables.id_año_fiscal").
		Where("cfg_periodos_contables.id_año_fiscal = ? AND cfg_periodos_contables.empresa_id = ?", idAñoFiscal, empresaID).
		Order("cfg_periodos_contables.numero_mes ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// Resolver controles de módulos para cada período
	results := make([]PeriodoContableResponse, len(rows))
	for i, r := range rows {
		modulos, err := getControlesModulos(db, r.ID)
		if err != nil {
			return nil, err
		}
		results[i] = toPeriodoResponse(&r.PeriodoContable, r.Estado, r.ColorEstado, r.AñoFiscal, modulos)
	}

	return results, nil
}

// getControlesModulos retorna los controles de módulos de un período con JOINs resueltos.
func getControlesModulos(db *gorm.DB, idPeriodo int64) ([]ControlModuloPeriodoResponse, error) {
	type rowControl struct {
		IDModulo     int    `gorm:"column:id_modulo"`
		CodigoModulo string `gorm:"column:codigo_modulo"`
		NombreModulo string `gorm:"column:nombre_modulo"`
		IDEstado     int    `gorm:"column:id_estado"`
		Estado       string `gorm:"column:estado"`
		PermiteOp    bool   `gorm:"column:permite_op"`
		Color        string `gorm:"column:color"`
		Notas        string `gorm:"column:notas"`
	}

	var rows []rowControl
	err := db.Session(&gorm.Session{NewDB: true}).
		Select(
			"cmp.id_modulo, rm.codigo AS codigo_modulo, rm.nombre AS nombre_modulo, "+
				"cmp.id_estado, em.nombre AS estado, em.permite_op, em.color, cmp.notas",
		).
		Table("contabilidad.cfg_control_modulos_periodo cmp").
		Joins("INNER JOIN contabilidad.ref_modulo_contable rm ON rm.id = cmp.id_modulo").
		Joins("INNER JOIN contabilidad.ref_estado_modulo em ON em.id = cmp.id_estado").
		Where("cmp.id_periodo = ?", idPeriodo).
		Order("rm.id ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make([]ControlModuloPeriodoResponse, len(rows))
	for i, r := range rows {
		result[i] = ControlModuloPeriodoResponse{
			IDModulo:     r.IDModulo,
			CodigoModulo: r.CodigoModulo,
			NombreModulo: r.NombreModulo,
			IDEstado:     r.IDEstado,
			Estado:       r.Estado,
			PermiteOp:    r.PermiteOp,
			Color:        r.Color,
			Notas:        r.Notas,
		}
	}

	return result, nil
}

// toPeriodoResponse construye el DTO público desde el modelo crudo.
func toPeriodoResponse(
	p *PeriodoContable,
	estado string,
	colorEstado string,
	añoFiscal int,
	modulos []ControlModuloPeriodoResponse,
) PeriodoContableResponse {
	return PeriodoContableResponse{
		ID:          p.ID,
		IDAñoFiscal: p.IDAñoFiscal,
		AñoFiscal:   añoFiscal,
		EmpresaID:   p.EmpresaID,
		NumeroMes:   p.NumeroMes,
		NombreMes:   p.NombreMes,
		FechaInicio: fmtFecha(p.FechaInicio),
		FechaFin:    fmtFecha(p.FechaFin),
		IDEstado:    p.IDEstado,
		Estado:      estado,
		ColorEstado: colorEstado,
		Notas:       p.Notas,
		Modulos:     modulos,
		CreatedAt:   fmtFecha(p.CreatedAt),
		UpdatedAt:   fmtFechaHora(p.UpdatedAt),
	}
}

// ─── Obtener por ID ───────────────────────────────────────────────────────────

// GetPeriodoByID obtiene un período con estado y módulos resueltos. Retorna nil, nil si no existe.
func GetPeriodoByID(db *gorm.DB, id int64, empresaID int64) (*PeriodoContableResponse, error) {
	type rowPeriodo struct {
		PeriodoContable
		Estado      string `gorm:"column:estado"`
		ColorEstado string `gorm:"column:color_estado"`
		AñoFiscal   int    `gorm:"column:año_fiscal"`
	}

	var row rowPeriodo
	err := db.
		Model(&PeriodoContable{}).
		Select(
			"cfg_periodos_contables.id, cfg_periodos_contables.id_año_fiscal, "+
				"cfg_periodos_contables.empresa_id, cfg_periodos_contables.numero_mes, "+
				"cfg_periodos_contables.nombre_mes, cfg_periodos_contables.fecha_inicio, "+
				"cfg_periodos_contables.fecha_fin, cfg_periodos_contables.id_estado, "+
				"cfg_periodos_contables.notas, cfg_periodos_contables.created_by, "+
				"cfg_periodos_contables.updated_by, cfg_periodos_contables.created_at, "+
				"cfg_periodos_contables.updated_at, "+
				"ep.nombre AS estado, ep.color AS color_estado, af.year AS año_fiscal",
		).
		Joins("INNER JOIN contabilidad.ref_estado_periodo ep ON ep.id = cfg_periodos_contables.id_estado").
		Joins("INNER JOIN contabilidad.cfg_años_fiscales af ON af.id = cfg_periodos_contables.id_año_fiscal").
		Where("cfg_periodos_contables.id = ? AND cfg_periodos_contables.empresa_id = ?", id, empresaID).
		First(&row).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	modulos, err := getControlesModulos(db, row.ID)
	if err != nil {
		return nil, err
	}

	resp := toPeriodoResponse(&row.PeriodoContable, row.Estado, row.ColorEstado, row.AñoFiscal, modulos)
	return &resp, nil
}

// getRawPeriodoByID obtiene el modelo crudo para validaciones internas.
func getRawPeriodoByID(db *gorm.DB, id int64, empresaID int64) (*PeriodoContable, error) {
	var result PeriodoContable
	err := db.
		Model(&PeriodoContable{}).
		Where("id = ? AND empresa_id = ?", id, empresaID).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ─── Validaciones ─────────────────────────────────────────────────────────────

// ExistenPeriodosParaAño verifica si ya se generaron los períodos de un año fiscal.
func ExistenPeriodosParaAño(db *gorm.DB, idAñoFiscal int64, empresaID int64) (bool, error) {
	var count int64
	err := db.Model(&PeriodoContable{}).
		Where("id_año_fiscal = ? AND empresa_id = ?", idAñoFiscal, empresaID).
		Count(&count).Error
	return count > 0, err
}

// GetPeriodoPorFecha retorna el período activo para una fecha específica de la empresa.
// Usado por el servicio de validación centralizada de operaciones.
func GetPeriodoPorFecha(db *gorm.DB, empresaID int64, fecha time.Time) (*PeriodoContable, error) {
	var result PeriodoContable
	err := db.Model(&PeriodoContable{}).
		Where("empresa_id = ? AND fecha_inicio <= ? AND fecha_fin >= ?", empresaID, fecha, fecha).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// GetControlModulo retorna el control de un módulo específico en un período.
func GetControlModulo(db *gorm.DB, idPeriodo int64, idModulo int) (*ControlModuloPeriodo, error) {
	var result ControlModuloPeriodo
	err := db.Model(&ControlModuloPeriodo{}).
		Where("id_periodo = ? AND id_modulo = ?", idPeriodo, idModulo).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ─── Cambio de estado de período ─────────────────────────────────────────────

// SetEstadoPeriodo actualiza el estado de un período contable.
func SetEstadoPeriodo(db *gorm.DB, id int64, idEstado int, notas string, updatedBy int64) error {
	return db.
		Model(&PeriodoContable{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"id_estado":  idEstado,
			"notas":      notas,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

// CerrarPeriodosAbiertosPorAño cierra todos los períodos abiertos de un año fiscal
// excepto el especificado. Se ejecuta en transacción al abrir un período.
func CerrarPeriodosAbiertosPorAño(db *gorm.DB, idAñoFiscal int64, empresaID int64, exceptID int64, updatedBy int64) error {
	q := db.Model(&PeriodoContable{}).
		Where("id_año_fiscal = ? AND empresa_id = ? AND id_estado = ?",
			idAñoFiscal, empresaID, EstadoPeriodoAbierto)

	if exceptID > 0 {
		q = q.Where("id != ?", exceptID)
	}

	return q.Updates(map[string]interface{}{
		"id_estado":  EstadoPeriodoCerrado,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}).Error
}

// ─── Cambio de estado de módulo en período ────────────────────────────────────

// SetEstadoModuloPeriodo actualiza el estado operacional de un módulo en un período.
func SetEstadoModuloPeriodo(db *gorm.DB, idPeriodo int64, idModulo int, idEstado int, notas string, updatedBy int64) error {
	return db.
		Model(&ControlModuloPeriodo{}).
		Where("id_periodo = ? AND id_modulo = ?", idPeriodo, idModulo).
		Updates(map[string]interface{}{
			"id_estado":  idEstado,
			"notas":      notas,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

// ─── Validación centralizada de operaciones ───────────────────────────────────

// PuedeOperarModulo verifica si un módulo puede ejecutar operaciones en la fecha dada.
// Esta función es el núcleo del motor de control transaccional de PRISMAR.
func PuedeOperarModulo(db *gorm.DB, empresaID int64, fecha time.Time, idModulo int) (bool, string, int64, error) {
	// 1. Buscar el período que contiene la fecha
	periodo, err := GetPeriodoPorFecha(db, empresaID, fecha)
	if err != nil {
		return false, "Error consultando el período contable", 0, err
	}
	if periodo == nil {
		return false, "No existe un período contable configurado para la fecha indicada", 0, nil
	}

	// 2. Validar estado del período
	switch periodo.IDEstado {
	case EstadoPeriodoCerrado:
		return false, "El período contable está Cerrado", periodo.ID, nil
	case EstadoPeriodoBloqueado:
		return false, "El período contable está Bloqueado. No se permiten operaciones", periodo.ID, nil
	}

	// 3. Validar estado del módulo dentro del período
	control, err := GetControlModulo(db, periodo.ID, idModulo)
	if err != nil {
		return false, "Error consultando el control de módulo", periodo.ID, err
	}
	if control == nil {
		return false, "El módulo no está configurado para este período", periodo.ID, nil
	}

	// 4. Consultar si el estado del módulo permite operaciones
	var estadoModulo RefEstadoModulo
	if err := db.Model(&RefEstadoModulo{}).Where("id = ?", control.IDEstado).First(&estadoModulo).Error; err != nil {
		return false, "Error consultando estado de módulo", periodo.ID, err
	}

	if !estadoModulo.PermiteOp {
		return false, "El módulo '" + estadoModulo.Nombre + "' no permite operaciones en este período", periodo.ID, nil
	}

	return true, "", periodo.ID, nil
}
