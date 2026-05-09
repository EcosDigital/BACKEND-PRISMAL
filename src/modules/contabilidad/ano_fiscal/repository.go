package ano_fiscal

import (
	"time"

	"gorm.io/gorm"
)

// ─── Helper: formateo de fechas ───────────────────────────────────────────────

func fmtFecha(t time.Time) string {
	return t.Format("02/01/2006")
}

func fmtFechaHora(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}

func toResponse(af *AñoFiscal, estado string) AñoFiscalResponse {
	return AñoFiscalResponse{
		ID:        af.ID,
		Year:      af.Year,
		IDEstado:  af.IDEstado,
		Estado:    estado,
		CreatedAt: fmtFecha(af.CreatedAt),
		UpdatedAt: fmtFechaHora(af.UpdatedAt),
	}
}

// ─── Referencias ──────────────────────────────────────────────────────────────

// GetEstadosAñoFiscal retorna el catálogo de estados activos.
// GORM resuelve la tabla via TableName() definido en el struct.
func GetEstadosAñoFiscal(db *gorm.DB) ([]RefEstadoAñoFiscal, error) {
	var results []RefEstadoAñoFiscal

	err := db.
		Model(&RefEstadoAñoFiscal{}).
		Select("id, nombre, is_active").
		Where("is_active = ?", true).
		Order("id ASC").
		Find(&results).Error

	if results == nil {
		results = []RefEstadoAñoFiscal{}
	}
	return results, err
}

// resolverNombreEstado obtiene el nombre de un estado por su ID.
func resolverNombreEstado(db *gorm.DB, idEstado int) (string, error) {
	var estado RefEstadoAñoFiscal

	err := db.
		Model(&RefEstadoAñoFiscal{}).
		Select("id, nombre").
		Where("id = ?", idEstado).
		First(&estado).Error

	if err != nil {
		return "", err
	}
	return estado.Nombre, nil
}

// ─── Crear ────────────────────────────────────────────────────────────────────

// CreateAñoFiscal inserta un nuevo año fiscal y retorna el registro completo.
// GORM popula af.ID con el SERIAL asignado por PostgreSQL al hacer Create.
func CreateAñoFiscal(db *gorm.DB, req *CreateAñoFiscalRequest) (*AñoFiscalResponse, error) {
	now := time.Now()

	af := AñoFiscal{
		Year:      req.Year,
		IDEstado:  EstadoCerrado,
		EmpresaID: req.EmpresaID,
		CreatedBy: req.UserID,
		UpdatedBy: req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Model(&af) usa TableName() para resolver la tabla correctamente.
	if err := db.Model(&af).Create(&af).Error; err != nil {
		return nil, err
	}

	nombreEstado, err := resolverNombreEstado(db, af.IDEstado)
	if err != nil {
		return nil, err
	}

	resp := toResponse(&af, nombreEstado)
	return &resp, nil
}

// ─── Listar ───────────────────────────────────────────────────────────────────

// ListarAñosFiscales retorna todos los años fiscales de la empresa con JOIN al catálogo.
func ListarAñosFiscales(db *gorm.DB, empresaID int64) ([]AñoFiscalResponse, error) {
	var rows []struct {
		AñoFiscal
		Estado string `gorm:"column:estado"`
	}

	err := db.
		Model(&AñoFiscal{}).
		Select("cfg_años_fiscales.id, cfg_años_fiscales.year, cfg_años_fiscales.id_estado, "+
			"cfg_años_fiscales.empresa_id, cfg_años_fiscales.created_by, cfg_años_fiscales.updated_by, "+
			"cfg_años_fiscales.created_at, cfg_años_fiscales.updated_at, "+
			"e.nombre AS estado").
		Joins("INNER JOIN contabilidad.ref_estado_año_fiscal e ON e.id = cfg_años_fiscales.id_estado").
		Where("cfg_años_fiscales.empresa_id = ?", empresaID).
		Order("cfg_años_fiscales.year DESC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	results := make([]AñoFiscalResponse, len(rows))
	for i, r := range rows {
		results[i] = toResponse(&r.AñoFiscal, r.Estado)
	}
	return results, nil
}

// ─── Obtener por ID ───────────────────────────────────────────────────────────

// GetAñoFiscalByID obtiene un año fiscal por su ID con el nombre del estado resuelto.
// Retorna nil, nil si no existe.
func GetAñoFiscalByID(db *gorm.DB, id int64, empresaID int64) (*AñoFiscalResponse, error) {
	var row struct {
		AñoFiscal
		Estado string `gorm:"column:estado"`
	}

	err := db.
		Model(&AñoFiscal{}).
		Select("cfg_años_fiscales.id, cfg_años_fiscales.year, cfg_años_fiscales.id_estado, "+
			"cfg_años_fiscales.empresa_id, cfg_años_fiscales.created_by, cfg_años_fiscales.updated_by, "+
			"cfg_años_fiscales.created_at, cfg_años_fiscales.updated_at, "+
			"e.nombre AS estado").
		Joins("INNER JOIN contabilidad.ref_estado_año_fiscal e ON e.id = cfg_años_fiscales.id_estado").
		Where("cfg_años_fiscales.id = ? AND cfg_años_fiscales.empresa_id = ?", id, empresaID).
		First(&row).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := toResponse(&row.AñoFiscal, row.Estado)
	return &resp, nil
}

// getRawAñoFiscalByID obtiene el modelo crudo (sin JOIN) para validaciones internas.
func getRawAñoFiscalByID(db *gorm.DB, id int64, empresaID int64) (*AñoFiscal, error) {
	var result AñoFiscal

	err := db.
		Model(&AñoFiscal{}).
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

// ExisteYear verifica si ya existe un año fiscal con ese año en la empresa.
func ExisteYear(db *gorm.DB, year int, empresaID int64, excludeID int64) (bool, error) {
	var count int64

	q := db.
		Model(&AñoFiscal{}).
		Where("year = ? AND empresa_id = ?", year, empresaID)

	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}

	err := q.Count(&count).Error
	return count > 0, err
}

// GetAñoFiscalAbierto retorna el año fiscal con id_estado = EstadoAbierto de la empresa.
func GetAñoFiscalAbierto(db *gorm.DB, empresaID int64) (*AñoFiscal, error) {
	var result AñoFiscal

	err := db.
		Model(&AñoFiscal{}).
		Where("empresa_id = ? AND id_estado = ?", empresaID, EstadoAbierto).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

// ─── Cambio de estado ─────────────────────────────────────────────────────────

// SetEstadoAñoFiscal actualiza el id_estado de un año fiscal.
// Usa Model(&AñoFiscal{}) para que GORM resuelva la tabla via TableName()
// y genere un UPDATE correcto sin ambigüedad en la cláusula WHERE.
func SetEstadoAñoFiscal(db *gorm.DB, id int64, idEstado int, updatedBy int64) error {
	return db.
		Model(&AñoFiscal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"id_estado":  idEstado,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

// CerrarTodosLosAbiertosExcepto cierra todos los años fiscales abiertos de la empresa
// excepto el que se está abriendo. Siempre se ejecuta dentro de una transacción.
func CerrarTodosLosAbiertosExcepto(db *gorm.DB, empresaID int64, exceptID int64, updatedBy int64) error {
	q := db.
		Model(&AñoFiscal{}).
		Where("empresa_id = ? AND id_estado = ?", empresaID, EstadoAbierto)

	if exceptID > 0 {
		q = q.Where("id != ?", exceptID)
	}

	return q.Updates(map[string]interface{}{
		"id_estado":  EstadoCerrado,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}).Error
}

// ─── Placeholder: validaciones contables ──────────────────────────────────────

// TieneAsientosPendientes verifica si el año fiscal tiene asientos pendientes.
// Placeholder activo — retorna false hasta que el módulo de asientos esté disponible.
//
// TODO: implementar con ORM cuando exista el modelo Asiento:
//
//	db.Model(&Asiento{}).
//	    Where("id_año_fiscal = ? AND id_estado_asiento = ?", añoFiscalID, EstadoAsientoPendiente).
//	    Count(&count)
func TieneAsientosPendientes(db *gorm.DB, añoFiscalID int64) (bool, error) {
	_ = db
	_ = añoFiscalID
	return false, nil
}
