package documentacion_sop

import (
	"time"

	"gorm.io/gorm"
)

// ─── Helpers: formateo de fechas ──────────────────────────────────────────────

func fmtFecha(t time.Time) string {
	return t.Format("02/01/2006")
}

func fmtFechaHora(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}

func fmtFechaISO(t time.Time) string {
	return t.Format("2006-01-02")
}

// toResponse convierte un SOP crudo más los campos de JOIN al DTO público.
func toResponse(s *SOP, categoria, estado, estadoColor string) SOPResponse {
	resp := SOPResponse{
		ID:            s.ID,
		Code:          s.Code,
		Title:         s.Title,
		IDCategoria:   s.IDCategoria,
		Categoria:     categoria,
		Area:          s.Area,
		IDEstado:      s.IDEstado,
		Estado:        estado,
		EstadoColor:   estadoColor,
		Content:       s.Content,
		Version:       s.Version,
		EffectiveDate: fmtFechaISO(s.EffectiveDate),
		CreatedAt:     fmtFecha(s.CreatedAt),
		UpdatedAt:     fmtFechaHora(s.UpdatedAt),
	}

	if s.ReviewDate != nil {
		f := fmtFechaISO(*s.ReviewDate)
		resp.ReviewDate = &f
	}

	return resp
}

// ─── Referencias ──────────────────────────────────────────────────────────────

// GetCategoriasSOP retorna el catálogo de categorías activas.
func GetCategoriasSOP(db *gorm.DB) ([]RefCategoriaSOP, error) {
	var results []RefCategoriaSOP

	err := db.
		Model(&RefCategoriaSOP{}).
		Select("id, nombre, is_active").
		Where("is_active = ?", true).
		Order("id ASC").
		Find(&results).Error

	if results == nil {
		results = []RefCategoriaSOP{}
	}
	return results, err
}

// GetEstadosSOP retorna el catálogo de estados activos.
func GetEstadosSOP(db *gorm.DB) ([]RefEstadoSOP, error) {
	var results []RefEstadoSOP

	err := db.
		Model(&RefEstadoSOP{}).
		Select("id, nombre, color, is_active").
		Where("is_active = ?", true).
		Order("id ASC").
		Find(&results).Error

	if results == nil {
		results = []RefEstadoSOP{}
	}
	return results, err
}

// ─── Crear ────────────────────────────────────────────────────────────────────

// CreateSOP inserta un nuevo SOP y retorna el registro completo con JOINs.
// GORM popula s.ID con el BIGSERIAL asignado por PostgreSQL al hacer Create.
func CreateSOP(db *gorm.DB, req *CreateSOPRequest) (*SOPResponse, error) {
	now := time.Now()

	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return nil, ErrFechaEfectivaInvalida
	}

	var reviewDate *time.Time
	if req.ReviewDate != nil {
		rd, err := time.Parse("2006-01-02", *req.ReviewDate)
		if err != nil {
			return nil, ErrFechaRevisionInvalida
		}
		reviewDate = &rd
	}

	s := SOP{
		Code:          req.Code,
		Title:         req.Title,
		IDCategoria:   req.IDCategoria,
		Area:          req.Area,
		IDEstado:      EstadoDraft, // todo SOP nace como Draft
		Content:       req.Content,
		Version:       req.Version,
		EffectiveDate: effectiveDate,
		ReviewDate:    reviewDate,
		IsActive:      true,
		EmpresaID:     req.EmpresaID,
		CreatedBy:     req.UserID,
		UpdatedBy:     req.UserID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := db.Model(&s).Create(&s).Error; err != nil {
		return nil, err
	}

	return GetSOPByID(db, s.ID, req.EmpresaID)
}

// ─── Listar ───────────────────────────────────────────────────────────────────

// ListarSOPs retorna los SOPs activos de la empresa aplicando los filtros opcionales.
// Los filtros de valor cero (0 para int, "" para string) se omiten.

func ListarSOPs(db *gorm.DB, empresaID int64, filters SOPFilters) ([]SOPResponse, error) {
	type row struct {
		SOP
		Categoria   string `gorm:"column:categoria"`
		Estado      string `gorm:"column:estado"`
		EstadoColor string `gorm:"column:estado_color"`
	}

	var rows []row

	q := db.
		Model(&SOP{}).
		Select("documentacion.cfg_sop.*, "+
			"c.nombre AS categoria, "+
			"e.nombre AS estado, "+
			"e.color  AS estado_color").
		Joins("INNER JOIN documentacion.ref_categoria_sop c ON c.id = documentacion.cfg_sop.id_categoria").
		Joins("INNER JOIN documentacion.ref_estado_sop    e ON e.id = documentacion.cfg_sop.id_estado").
		Where("documentacion.cfg_sop.empresa_id = ?", empresaID).
		Where("documentacion.cfg_sop.is_active  = ?", true)

	if filters.IDCategoria > 0 {
		q = q.Where("documentacion.cfg_sop.id_categoria = ?", filters.IDCategoria)
	}

	if filters.IDEstado > 0 {
		q = q.Where("documentacion.cfg_sop.id_estado = ?", filters.IDEstado)
	}

	if filters.Area != "" {
		q = q.Where("documentacion.cfg_sop.area ILIKE ?", "%"+filters.Area+"%")
	}

	if filters.Search != "" {
		pattern := "%" + filters.Search + "%"
		q = q.Where(
			"documentacion.cfg_sop.title ILIKE ? OR documentacion.cfg_sop.code ILIKE ?",
			pattern, pattern,
		)
	}

	err := q.Order("documentacion.cfg_sop.created_at DESC").Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]SOPResponse, len(rows))
	for i, r := range rows {
		results[i] = toResponse(&r.SOP, r.Categoria, r.Estado, r.EstadoColor)
	}
	return results, nil
}

// ─── Obtener por ID ───────────────────────────────────────────────────────────

// GetSOPByID obtiene un SOP activo por su ID con los catálogos resueltos.
// Retorna nil, nil si no existe o fue eliminado lógicamente.
func GetSOPByID(db *gorm.DB, id int64, empresaID int64) (*SOPResponse, error) {
	type row struct {
		SOP
		Categoria   string `gorm:"column:categoria"`
		Estado      string `gorm:"column:estado"`
		EstadoColor string `gorm:"column:estado_color"`
	}

	var r row

	err := db.
		Model(&SOP{}).
		Select("documentacion.cfg_sop.*, "+
			"c.nombre AS categoria, "+
			"e.nombre AS estado, "+
			"e.color  AS estado_color").
		Joins("INNER JOIN documentacion.ref_categoria_sop c ON c.id = documentacion.cfg_sop.id_categoria").
		Joins("INNER JOIN documentacion.ref_estado_sop    e ON e.id = documentacion.cfg_sop.id_estado").
		Where("documentacion.cfg_sop.id         = ?", id).
		Where("documentacion.cfg_sop.empresa_id = ?", empresaID).
		Where("documentacion.cfg_sop.is_active  = ?", true).
		First(&r).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := toResponse(&r.SOP, r.Categoria, r.Estado, r.EstadoColor)
	return &resp, nil
}

// getRawSOPByID obtiene el modelo crudo (sin JOIN) para validaciones internas.
func getRawSOPByID(db *gorm.DB, id int64, empresaID int64) (*SOP, error) {
	var result SOP

	err := db.
		Model(&SOP{}).
		Where("id = ? AND empresa_id = ? AND is_active = ?", id, empresaID, true).
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

// ExisteCode verifica si ya existe un SOP activo con ese código en la empresa.
// excludeID permite ignorar el propio registro en actualizaciones.
func ExisteCode(db *gorm.DB, code string, empresaID int64, excludeID int64) (bool, error) {
	var count int64

	q := db.
		Model(&SOP{}).
		Where("code = ? AND empresa_id = ? AND is_active = ?", code, empresaID, true)

	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}

	err := q.Count(&count).Error
	return count > 0, err
}

// ExisteCategoria verifica que el id_categoria exista y esté activo en el catálogo.
func ExisteCategoria(db *gorm.DB, idCategoria int) (bool, error) {
	var count int64
	err := db.
		Model(&RefCategoriaSOP{}).
		Where("id = ? AND is_active = ?", idCategoria, true).
		Count(&count).Error
	return count > 0, err
}

// ExisteEstado verifica que el id_estado exista y esté activo en el catálogo.
func ExisteEstado(db *gorm.DB, idEstado int) (bool, error) {
	var count int64
	err := db.
		Model(&RefEstadoSOP{}).
		Where("id = ? AND is_active = ?", idEstado, true).
		Count(&count).Error
	return count > 0, err
}

// ─── Actualizar ───────────────────────────────────────────────────────────────

// UpdateSOP modifica los campos editables de un SOP existente.
func UpdateSOP(db *gorm.DB, id int64, req *UpdateSOPRequest) (*SOPResponse, error) {
	effectiveDate, err := time.Parse("2006-01-02", req.EffectiveDate)
	if err != nil {
		return nil, ErrFechaEfectivaInvalida
	}

	fields := map[string]interface{}{
		"title":          req.Title,
		"id_categoria":   req.IDCategoria,
		"area":           req.Area,
		"id_estado":      req.IDEstado,
		"content":        req.Content,
		"version":        req.Version,
		"effective_date": effectiveDate,
		"review_date":    nil,
		"updated_by":     req.UserID,
		"updated_at":     time.Now(),
	}

	if req.ReviewDate != nil {
		rd, err := time.Parse("2006-01-02", *req.ReviewDate)
		if err != nil {
			return nil, ErrFechaRevisionInvalida
		}
		fields["review_date"] = rd
	}

	err = db.
		Model(&SOP{}).
		Where("id = ? AND empresa_id = ? AND is_active = ?", id, req.EmpresaID, true).
		Updates(fields).Error

	if err != nil {
		return nil, err
	}

	return GetSOPByID(db, id, req.EmpresaID)
}

// ─── Eliminar (soft delete) ───────────────────────────────────────────────────

// SoftDeleteSOP realiza la eliminación lógica del SOP.
// Marca is_active = false sin borrar el registro de la base de datos.
func SoftDeleteSOP(db *gorm.DB, id int64, empresaID int64, updatedBy int64) error {
	return db.
		Model(&SOP{}).
		Where("id = ? AND empresa_id = ? AND is_active = ?", id, empresaID, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}
