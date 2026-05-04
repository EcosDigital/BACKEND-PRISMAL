package scrapping

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ─── Constantes de código de estado ──────────────────────────────────────────
// Deben coincidir exactamente con la columna `codigo` de las tablas de referencia.
// ref_estado_job:  001=Pendiente 002=Ejecutando 003=Completado 004=Fallido
// ref_estado_lead: 001=Nuevo 002=Contactado 003=Calificado 004=Descartado

const (
	EstadoJobPendiente  = "001"
	EstadoJobEjecutando = "002"
	EstadoJobCompletado = "003"
	EstadoJobFallido    = "004"

	EstadoLeadNuevo      = "001"
	EstadoLeadContactado = "002"
	EstadoLeadCalificado = "003"
	EstadoLeadDescartado = "004"
)

// ─── Jobs ─────────────────────────────────────────────────────────────────────

// CreateScrapingJob inserta un job en estado Pendiente y retorna su ID.
func CreateScrapingJob(db *gorm.DB, req *ScrapingJobRequest, userID, empresaID, sedeID int64) (int64, error) {

	limite := req.Limite
	if limite == 0 {
		limite = 50
	}

	idEstadoPendiente, err := resolveEstadoJob(db, EstadoJobPendiente)
	if err != nil {
		return 0, err
	}

	data := map[string]interface{}{
		"keyword":    req.Keyword,
		"ciudad":     req.Ciudad,
		"limite":     limite,
		"id_estado":  idEstadoPendiente,
		"id_empresa": empresaID,
		"id_sede":    sedeID,
		"created_at": time.Now(),
		"created_by": userID,
	}

	tx := db.Table("prospeccion.cfg_scraping_jobs").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateJobEstado cambia el estado del job.
// Usar las constantes EstadoJobXxx definidas en este archivo.
func UpdateJobEstado(db *gorm.DB, jobID int64, codigoEstado string) error {

	idEstado, err := resolveEstadoJob(db, codigoEstado)
	if err != nil {
		return err
	}

	return db.
		Table("prospeccion.cfg_scraping_jobs").
		Where("id = ?", jobID).
		Update("id_estado", idEstado).Error
}

// UpdateJobCounters persiste los contadores finales y marca finished_at.
func UpdateJobCounters(db *gorm.DB, jobID int64, c JobCounters, codigoEstado string) error {

	idEstado, err := resolveEstadoJob(db, codigoEstado)
	if err != nil {
		return err
	}

	now := time.Now()
	return db.
		Table("prospeccion.cfg_scraping_jobs").
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"id_estado":         idEstado,
			"total_encontrados": c.Encontrados,
			"total_guardados":   c.Guardados,
			"total_duplicados":  c.Duplicados,
			"total_errores":     c.Errores,
			"finished_at":       now,
		}).Error
}

// UpdateJobError persiste el error global cuando el job falla completamente.
func UpdateJobError(db *gorm.DB, jobID int64, mensaje string) error {

	idEstado, err := resolveEstadoJob(db, EstadoJobFallido)
	if err != nil {
		return err
	}

	now := time.Now()
	return db.
		Table("prospeccion.cfg_scraping_jobs").
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"id_estado":     idEstado,
			"error_mensaje": mensaje,
			"finished_at":   now,
		}).Error
}

// ListScrapingJobs retorna el historial de jobs de una empresa.
func ListScrapingJobs(db *gorm.DB, empresaID int64) ([]JobListResponse, error) {

	var results []JobListResponse

	err := db.
		Table("prospeccion.cfg_scraping_jobs j").
		Select(`
			j.id,
			j.keyword,
			j.ciudad,
			j.limite,
			e.nombre            AS estado,
			j.total_encontrados,
			j.total_guardados,
			j.total_duplicados,
			j.total_errores,
			j.created_at::TEXT  AS created_at,
			j.finished_at::TEXT AS finished_at`).
		Joins("INNER JOIN prospeccion.ref_estado_job e ON e.id = j.id_estado").
		Where("j.id_empresa = ?", empresaID).
		Order("j.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []JobListResponse{}
	}

	return results, nil
}

// ─── Leads ────────────────────────────────────────────────────────────────────

// CheckLeadDuplicado llama a fn_check_lead_duplicado en PostgreSQL.
func CheckLeadDuplicado(db *gorm.DB, telefono, nombre, ciudad string, empresaID int64) (*DuplicadoCheck, error) {

	var result struct {
		EsDuplicado     bool   `gorm:"column:es_duplicado"`
		Motivo          string `gorm:"column:motivo"`
		IDLeadExistente *int64 `gorm:"column:id_lead_existente"`
	}

	err := db.Raw(`
		SELECT * FROM prospeccion.fn_check_lead_duplicado(?, ?, ?, ?)`,
		nullIfEmpty(telefono), nombre, ciudad, empresaID,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &DuplicadoCheck{
		EsDuplicado:     result.EsDuplicado,
		Motivo:          result.Motivo,
		IDLeadExistente: result.IDLeadExistente,
	}, nil
}

// CreateLead inserta un lead ya validado (no duplicado).
func CreateLead(db *gorm.DB, lead *LeadRaw, jobID, userID, empresaID, sedeID int64) error {

	idEstadoNuevo, err := resolveEstadoLead(db, EstadoLeadNuevo)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"id_job":         jobID,
		"nombre":         lead.Nombre,
		"tipo_negocio":   lead.TipoNegocio,
		"ciudad":         lead.Ciudad,
		"direccion":      lead.Direccion,
		"telefono":       nullIfEmpty(lead.Telefono),
		"sitio_web":      nullIfEmpty(lead.SitioWeb),
		"rating":         nullIfZeroFloat(lead.Rating),
		"total_reviews":  nullIfZeroInt(lead.TotalReviews),
		"maps_url":       nullIfEmpty(lead.MapsURL),
		"id_estado_lead": idEstadoNuevo,
		"id_empresa":     empresaID,
		"id_sede":        sedeID,
		"created_at":     time.Now(),
		"created_by":     userID,
	}

	return db.Table("prospeccion.prospeccion_leads").Create(&data).Error
}

// ListLeads retorna leads paginados con filtros opcionales.
func ListLeads(db *gorm.DB, f LeadFilterRequest, empresaID int64) ([]LeadResponse, error) {

	page := f.Page
	if page <= 0 {
		page = 1
	}
	limit := f.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	q := db.
		Table("prospeccion.prospeccion_leads l").
		Select(`
			l.id,
			l.id_job,
			l.nombre,
			COALESCE(l.tipo_negocio,  '') AS tipo_negocio,
			l.ciudad,
			COALESCE(l.direccion,     '') AS direccion,
			COALESCE(l.telefono,      '') AS telefono,
			COALESCE(l.sitio_web,     '') AS sitio_web,
			COALESCE(l.rating,         0) AS rating,
			COALESCE(l.total_reviews,  0) AS total_reviews,
			COALESCE(l.maps_url,      '') AS maps_url,
			e.nombre                      AS estado_lead,
			l.created_at::TEXT            AS created_at`).
		Joins("INNER JOIN prospeccion.ref_estado_lead e ON e.id = l.id_estado_lead").
		Where("l.id_empresa = ?", empresaID)

	if f.IDJob != nil {
		q = q.Where("l.id_job = ?", *f.IDJob)
	}
	if f.Ciudad != "" {
		q = q.Where("l.ciudad ILIKE ?", "%"+f.Ciudad+"%")
	}
	if f.IDEstado != nil {
		q = q.Where("l.id_estado_lead = ?", *f.IDEstado)
	}

	var results []LeadResponse
	err := q.Order("l.id DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []LeadResponse{}
	}

	return results, nil
}

// ListLeadsForExport retorna todos los leads sin paginación para el Excel.
func ListLeadsForExport(db *gorm.DB, f LeadExportRequest, empresaID int64) ([]LeadExportRow, error) {

	q := db.
		Table("prospeccion.prospeccion_leads l").
		Select(`
			l.nombre,
			COALESCE(l.tipo_negocio,  '') AS tipo_negocio,
			l.ciudad,
			COALESCE(l.direccion,     '') AS direccion,
			COALESCE(l.telefono,      '') AS telefono,
			COALESCE(l.sitio_web,     '') AS sitio_web,
			COALESCE(l.rating,         0) AS rating,
			COALESCE(l.total_reviews,  0) AS total_reviews,
			COALESCE(l.maps_url,      '') AS maps_url,
			e.nombre                      AS estado`).
		Joins("INNER JOIN prospeccion.ref_estado_lead e ON e.id = l.id_estado_lead").
		Where("l.id_empresa = ?", empresaID)

	if f.IDJob != nil {
		q = q.Where("l.id_job = ?", *f.IDJob)
	}
	if f.Ciudad != "" {
		q = q.Where("l.ciudad ILIKE ?", "%"+f.Ciudad+"%")
	}
	if f.IDEstado != nil {
		q = q.Where("l.id_estado_lead = ?", *f.IDEstado)
	}

	var results []LeadExportRow
	err := q.Order("l.nombre ASC").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []LeadExportRow{}
	}

	return results, nil
}

// ─── Helpers privados ─────────────────────────────────────────────────────────

// resolveEstadoJob obtiene el id de ref_estado_job por su codigo consecutivo.
func resolveEstadoJob(db *gorm.DB, codigo string) (int, error) {
	var id int
	err := db.Raw(
		"SELECT id FROM prospeccion.ref_estado_job WHERE codigo = ? LIMIT 1", codigo,
	).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("ref_estado_job: codigo '%s' no encontrado", codigo)
	}
	return id, nil
}

// resolveEstadoLead obtiene el id de ref_estado_lead por su codigo consecutivo.
func resolveEstadoLead(db *gorm.DB, codigo string) (int, error) {
	var id int
	err := db.Raw(
		"SELECT id FROM prospeccion.ref_estado_lead WHERE codigo = ? LIMIT 1", codigo,
	).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("ref_estado_lead: codigo '%s' no encontrado", codigo)
	}
	return id, nil
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullIfZeroFloat(f float64) interface{} {
	if f == 0 {
		return nil
	}
	return f
}

func nullIfZeroInt(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

// Asegura que sql se importe si algún caller lo necesita en el futuro.
var _ = sql.ErrNoRows
