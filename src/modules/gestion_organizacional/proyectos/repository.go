package proyectos

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

// ─── Conversores a Response ───────────────────────────────────────────────────
func toProyectoResponse(p *Proyecto, estado, colorEstado, prioridad, colorPrioridad string) ProyectoResponse {
	resp := ProyectoResponse{
		ID:             p.ID,
		Codigo:         p.Codigo,
		Nombre:         p.Nombre,
		Descripcion:    p.Descripcion,
		FechaInicio:    fmtFechaISO(p.FechaInicio),
		IDEstado:       p.IDEstado,
		Estado:         estado,
		ColorEstado:    colorEstado,
		IDPrioridad:    p.IDPrioridad,
		Prioridad:      prioridad,
		ColorPrioridad: colorPrioridad,
		IDResponsable:  p.IDResponsable,
		CreadoEn:       fmtFecha(p.CreadoEn),
		ActualizadoEn:  fmtFechaHora(p.ActualizadoEn),
	}
	if p.FechaFin != nil {
		s := fmtFechaISO(*p.FechaFin)
		resp.FechaFin = &s
	}
	return resp
}

func toParteInteresadaResponse(p *ParteInteresada) ParteInteresadaResponse {
	return ParteInteresadaResponse{
		ID:            p.ID,
		IDProyecto:    p.IDProyecto,
		Nombre:        p.Nombre,
		Rol:           p.Rol,
		Organizacion:  p.Organizacion,
		Correo:        p.Correo,
		Telefono:      p.Telefono,
		Observaciones: p.Observaciones,
		CreadoEn:      fmtFecha(p.CreadoEn),
		ActualizadoEn: fmtFechaHora(p.ActualizadoEn),
	}
}

func toComentarioResponse(c *Comentario, evidencias []EvidenciaResponse) ComentarioResponse {
	if evidencias == nil {
		evidencias = []EvidenciaResponse{}
	}
	return ComentarioResponse{
		ID:             c.ID,
		IDProyecto:     c.IDProyecto,
		FechaActividad: fmtFechaISO(c.FechaActividad),
		Comentario:     c.Comentario,
		CreadoPor:      c.CreadoPor,
		CreadoEn:       fmtFechaHora(c.CreadoEn),
		Evidencias:     evidencias,
	}
}

func toEvidenciaResponse(e *Evidencia, tipoNombre string) EvidenciaResponse {
	return EvidenciaResponse{
		ID:            e.ID,
		IDComentario:  e.IDComentario,
		NombreArchivo: e.NombreArchivo,
		RutaArchivo:   e.RutaArchivo,
		IDTipo:        e.IDTipo,
		TipoEvidencia: tipoNombre,
		SubidoEn:      fmtFechaHora(e.SubidoEn),
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// CATÁLOGOS
// ═══════════════════════════════════════════════════════════════════════════════

func ObtenerEstadosProyecto(db *gorm.DB) ([]RefEstadoProyecto, error) {
	var resultados []RefEstadoProyecto
	err := db.Model(&RefEstadoProyecto{}).
		Where("es_activo = ?", true).
		Order("id ASC").
		Find(&resultados).Error
	if resultados == nil {
		resultados = []RefEstadoProyecto{}
	}
	return resultados, err
}

func ObtenerPrioridadesProyecto(db *gorm.DB) ([]RefPrioridadProyecto, error) {
	var resultados []RefPrioridadProyecto
	err := db.Model(&RefPrioridadProyecto{}).
		Where("es_activo = ?", true).
		Order("id ASC").
		Find(&resultados).Error
	if resultados == nil {
		resultados = []RefPrioridadProyecto{}
	}
	return resultados, err
}

func ObtenerTiposEvidencia(db *gorm.DB) ([]RefTipoEvidencia, error) {
	var resultados []RefTipoEvidencia
	err := db.Model(&RefTipoEvidencia{}).
		Where("es_activo = ?", true).
		Order("id ASC").
		Find(&resultados).Error
	if resultados == nil {
		resultados = []RefTipoEvidencia{}
	}
	return resultados, err
}

// ═══════════════════════════════════════════════════════════════════════════════
// PROYECTOS
// ═══════════════════════════════════════════════════════════════════════════════

// InsertarProyecto persiste un nuevo proyecto y retorna el DTO con JOINs resueltos.
func InsertarProyecto(db *gorm.DB, req *CrearProyectoRequest) (*ProyectoResponse, error) {
	ahora := time.Now()

	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicio)
	if err != nil {
		return nil, ErrFechaInicioInvalida
	}

	var fechaFin *time.Time
	if req.FechaFin != nil {
		ff, err := time.Parse("2006-01-02", *req.FechaFin)
		if err != nil {
			return nil, ErrFechaFinInvalida
		}
		fechaFin = &ff
	}

	p := Proyecto{
		Codigo:         req.Codigo,
		Nombre:         req.Nombre,
		Descripcion:    req.Descripcion,
		FechaInicio:    fechaInicio,
		FechaFin:       fechaFin,
		IDEstado:       EstadoPlaneacion, // todo proyecto nace en Planeación
		IDPrioridad:    req.IDPrioridad,
		IDResponsable:  req.IDResponsable,
		EsActivo:       true,
		IDEmpresa:      req.IDEmpresa,
		CreadoPor:      req.IDUsuario,
		ActualizadoPor: req.IDUsuario,
		CreadoEn:       ahora,
		ActualizadoEn:  ahora,
	}

	if err := db.Model(&p).Create(&p).Error; err != nil {
		return nil, err
	}

	return BuscarProyectoPorID(db, p.ID, req.IDEmpresa)
}

// ListarProyectos retorna proyectos activos de la empresa con filtros opcionales.
func ListarProyectos(db *gorm.DB, idEmpresa int64, filtros FiltrosProyecto) ([]ProyectoResponse, error) {
	type fila struct {
		Proyecto
		Estado         string `gorm:"column:estado"`
		ColorEstado    string `gorm:"column:color_estado"`
		Prioridad      string `gorm:"column:prioridad"`
		ColorPrioridad string `gorm:"column:color_prioridad"`
	}

	var filas []fila

	q := db.Model(&Proyecto{}).
		Select("organizacion.cfg_proyectos.*, "+
			"ep.nombre AS estado, "+
			"ep.color  AS color_estado, "+
			"pr.nombre AS prioridad, "+
			"pr.color  AS color_prioridad").
		Joins("INNER JOIN organizacion.ref_estado_proyecto    ep ON ep.id = organizacion.cfg_proyectos.id_estado").
		Joins("INNER JOIN organizacion.ref_prioridad_proyecto pr ON pr.id = organizacion.cfg_proyectos.id_prioridad").
		Where("organizacion.cfg_proyectos.id_empresa = ?", idEmpresa).
		Where("organizacion.cfg_proyectos.es_activo  = ?", true)

	if filtros.IDEstado > 0 {
		q = q.Where("organizacion.cfg_proyectos.id_estado = ?", filtros.IDEstado)
	}
	if filtros.IDPrioridad > 0 {
		q = q.Where("organizacion.cfg_proyectos.id_prioridad = ?", filtros.IDPrioridad)
	}
	if filtros.IDResponsable > 0 {
		q = q.Where("organizacion.cfg_proyectos.id_responsable = ?", filtros.IDResponsable)
	}
	if filtros.Busqueda != "" {
		patron := "%" + filtros.Busqueda + "%"
		q = q.Where(
			"organizacion.cfg_proyectos.nombre ILIKE ? OR organizacion.cfg_proyectos.codigo ILIKE ?",
			patron, patron,
		)
	}

	err := q.Order("organizacion.cfg_proyectos.created_at DESC").Scan(&filas).Error
	if err != nil {
		return nil, err
	}

	resultados := make([]ProyectoResponse, len(filas))
	for i, f := range filas {
		resultados[i] = toProyectoResponse(&f.Proyecto, f.Estado, f.ColorEstado, f.Prioridad, f.ColorPrioridad)
	}
	return resultados, nil
}

// BuscarProyectoPorID retorna un proyecto activo con JOINs. nil,nil si no existe.
func BuscarProyectoPorID(db *gorm.DB, id int64, idEmpresa int64) (*ProyectoResponse, error) {
	type fila struct {
		Proyecto
		Estado         string `gorm:"column:estado"`
		ColorEstado    string `gorm:"column:color_estado"`
		Prioridad      string `gorm:"column:prioridad"`
		ColorPrioridad string `gorm:"column:color_prioridad"`
	}

	var f fila

	err := db.Model(&Proyecto{}).
		Select("organizacion.cfg_proyectos.*, "+
			"ep.nombre AS estado, "+
			"ep.color  AS color_estado, "+
			"pr.nombre AS prioridad, "+
			"pr.color  AS color_prioridad").
		Joins("INNER JOIN organizacion.ref_estado_proyecto    ep ON ep.id = organizacion.cfg_proyectos.id_estado").
		Joins("INNER JOIN organizacion.ref_prioridad_proyecto pr ON pr.id = organizacion.cfg_proyectos.id_prioridad").
		Where("organizacion.cfg_proyectos.id         = ?", id).
		Where("organizacion.cfg_proyectos.id_empresa = ?", idEmpresa).
		Where("organizacion.cfg_proyectos.es_activo  = ?", true).
		First(&f).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	resp := toProyectoResponse(&f.Proyecto, f.Estado, f.ColorEstado, f.Prioridad, f.ColorPrioridad)
	return &resp, nil
}

// buscarProyectoCrudo obtiene el modelo sin JOIN para validaciones internas.
func buscarProyectoCrudo(db *gorm.DB, id int64, idEmpresa int64) (*Proyecto, error) {
	var resultado Proyecto
	err := db.Model(&Proyecto{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, idEmpresa, true).
		First(&resultado).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &resultado, nil
}

func ExisteCodigoProyecto(db *gorm.DB, codigo string, idEmpresa int64, excluirID int64) (bool, error) {
	var total int64
	q := db.Model(&Proyecto{}).
		Where("codigo = ? AND id_empresa = ? AND es_activo = ?", codigo, idEmpresa, true)
	if excluirID > 0 {
		q = q.Where("id != ?", excluirID)
	}
	err := q.Count(&total).Error
	return total > 0, err
}

func ExisteEstadoProyecto(db *gorm.DB, id int) (bool, error) {
	var total int64
	err := db.Model(&RefEstadoProyecto{}).
		Where("id = ? AND es_activo = ?", id, true).Count(&total).Error
	return total > 0, err
}

func ExistePrioridad(db *gorm.DB, id int) (bool, error) {
	var total int64
	err := db.Model(&RefPrioridadProyecto{}).
		Where("id = ? AND es_activo = ?", id, true).Count(&total).Error
	return total > 0, err
}

// ActualizarProyecto modifica los campos editables de un proyecto.
func ActualizarProyecto(db *gorm.DB, id int64, req *ActualizarProyectoRequest) (*ProyectoResponse, error) {
	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicio)
	if err != nil {
		return nil, ErrFechaInicioInvalida
	}

	campos := map[string]interface{}{
		"nombre":          req.Nombre,
		"descripcion":     req.Descripcion,
		"fecha_inicio":    fechaInicio,
		"fecha_fin":       nil,
		"id_estado":       req.IDEstado,
		"id_prioridad":    req.IDPrioridad,
		"id_responsable":  req.IDResponsable,
		"actualizado_por": req.IDUsuario,
		"actualizado_en":  time.Now(),
	}

	if req.FechaFin != nil {
		ff, err := time.Parse("2006-01-02", *req.FechaFin)
		if err != nil {
			return nil, ErrFechaFinInvalida
		}
		campos["fecha_fin"] = ff
	}

	err = db.Model(&Proyecto{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, req.IDEmpresa, true).
		Updates(campos).Error
	if err != nil {
		return nil, err
	}

	return BuscarProyectoPorID(db, id, req.IDEmpresa)
}

// EliminarProyectoLogico marca es_activo = false.
func EliminarProyectoLogico(db *gorm.DB, id int64, idEmpresa int64, idUsuario int64) error {
	return db.Model(&Proyecto{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, idEmpresa, true).
		Updates(map[string]interface{}{
			"es_activo":       false,
			"actualizado_por": idUsuario,
			"actualizado_en":  time.Now(),
		}).Error
}

// ═══════════════════════════════════════════════════════════════════════════════
// PARTES INTERESADAS
// ═══════════════════════════════════════════════════════════════════════════════

func InsertarParteInteresada(db *gorm.DB, req *CrearParteInteresadaRequest) (*ParteInteresadaResponse, error) {
	ahora := time.Now()
	p := ParteInteresada{
		IDProyecto:     req.IDProyecto,
		Nombre:         req.Nombre,
		Rol:            req.Rol,
		Organizacion:   req.Organizacion,
		Correo:         req.Correo,
		Telefono:       req.Telefono,
		Observaciones:  req.Observaciones,
		EsActivo:       true,
		IDEmpresa:      req.IDEmpresa,
		CreadoPor:      req.IDUsuario,
		ActualizadoPor: req.IDUsuario,
		CreadoEn:       ahora,
		ActualizadoEn:  ahora,
	}
	if err := db.Model(&p).Create(&p).Error; err != nil {
		return nil, err
	}
	resp := toParteInteresadaResponse(&p)
	return &resp, nil
}

func ListarPartesInteresadas(db *gorm.DB, idProyecto int64, idEmpresa int64) ([]ParteInteresadaResponse, error) {
	var filas []ParteInteresada
	err := db.Model(&ParteInteresada{}).
		Where("id_proyecto = ? AND id_empresa = ? AND es_activo = ?", idProyecto, idEmpresa, true).
		Order("creado_en ASC").
		Find(&filas).Error
	if err != nil {
		return nil, err
	}
	resultados := make([]ParteInteresadaResponse, len(filas))
	for i, f := range filas {
		resultados[i] = toParteInteresadaResponse(&f)
	}
	return resultados, nil
}

func BuscarParteInteresadaPorID(db *gorm.DB, id int64, idEmpresa int64) (*ParteInteresadaResponse, error) {
	var p ParteInteresada
	err := db.Model(&ParteInteresada{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, idEmpresa, true).
		First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	resp := toParteInteresadaResponse(&p)
	return &resp, nil
}

func buscarParteInteresadaCruda(db *gorm.DB, id int64, idEmpresa int64) (*ParteInteresada, error) {
	var p ParteInteresada
	err := db.Model(&ParteInteresada{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, idEmpresa, true).
		First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func ActualizarParteInteresada(db *gorm.DB, id int64, req *ActualizarParteInteresadaRequest) (*ParteInteresadaResponse, error) {
	err := db.Model(&ParteInteresada{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, req.IDEmpresa, true).
		Updates(map[string]interface{}{
			"nombre":          req.Nombre,
			"rol":             req.Rol,
			"organizacion":    req.Organizacion,
			"correo":          req.Correo,
			"telefono":        req.Telefono,
			"observaciones":   req.Observaciones,
			"actualizado_por": req.IDUsuario,
			"actualizado_en":  time.Now(),
		}).Error
	if err != nil {
		return nil, err
	}
	return BuscarParteInteresadaPorID(db, id, req.IDEmpresa)
}

func EliminarParteInteresadaLogica(db *gorm.DB, id int64, idEmpresa int64, idUsuario int64) error {
	return db.Model(&ParteInteresada{}).
		Where("id = ? AND id_empresa = ? AND es_activo = ?", id, idEmpresa, true).
		Updates(map[string]interface{}{
			"es_activo":       false,
			"actualizado_por": idUsuario,
			"actualizado_en":  time.Now(),
		}).Error
}

// ═══════════════════════════════════════════════════════════════════════════════
// COMENTARIOS (BITÁCORA)
// ═══════════════════════════════════════════════════════════════════════════════

func InsertarComentario(db *gorm.DB, req *CrearComentarioRequest) (*ComentarioResponse, error) {
	fechaActividad, err := time.Parse("2006-01-02", req.FechaActividad)
	if err != nil {
		return nil, ErrFechaActividadInvalida
	}

	c := Comentario{
		IDProyecto:     req.IDProyecto,
		FechaActividad: fechaActividad,
		Comentario:     req.Comentario,
		IDEmpresa:      req.IDEmpresa,
		CreadoPor:      req.IDUsuario,
		CreadoEn:       time.Now(),
	}

	if err := db.Model(&c).Create(&c).Error; err != nil {
		return nil, err
	}

	resp := toComentarioResponse(&c, nil)
	return &resp, nil
}

// ListarComentarios retorna la bitácora de un proyecto ordenada por fecha de actividad.
// Cada comentario incluye sus evidencias resueltas.
func ListarComentarios(db *gorm.DB, idProyecto int64, idEmpresa int64) ([]ComentarioResponse, error) {
	var filas []Comentario
	err := db.Model(&Comentario{}).
		Where("id_proyecto = ? AND id_empresa = ?", idProyecto, idEmpresa).
		Order("fecha_actividad DESC, created_at DESC").
		Find(&filas).Error
	if err != nil {
		return nil, err
	}

	resultados := make([]ComentarioResponse, len(filas))
	for i, c := range filas {
		evidencias, err := buscarEvidenciasPorComentario(db, c.ID)
		if err != nil {
			return nil, err
		}
		resultados[i] = toComentarioResponse(&c, evidencias)
	}
	return resultados, nil
}

func buscarComentarioCrudo(db *gorm.DB, id int64, idEmpresa int64) (*Comentario, error) {
	var c Comentario
	err := db.Model(&Comentario{}).
		Where("id = ? AND id_empresa = ?", id, idEmpresa).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// ═══════════════════════════════════════════════════════════════════════════════
// EVIDENCIAS
// ═══════════════════════════════════════════════════════════════════════════════

func InsertarEvidencia(db *gorm.DB, req *CrearEvidenciaRequest) (*EvidenciaResponse, error) {
	e := Evidencia{
		IDComentario:  req.IDComentario,
		NombreArchivo: req.NombreArchivo,
		RutaArchivo:   req.RutaArchivo,
		IDTipo:        req.IDTipo,
		IDEmpresa:     req.IDEmpresa,
		SubidoPor:     req.IDUsuario,
		SubidoEn:      time.Now(),
	}

	if err := db.Model(&e).Create(&e).Error; err != nil {
		return nil, err
	}

	tipoNombre, err := resolverNombreTipo(db, e.IDTipo)
	if err != nil {
		return nil, err
	}

	resp := toEvidenciaResponse(&e, tipoNombre)
	return &resp, nil
}

// buscarEvidenciasPorComentario es un helper interno usado al cargar la bitácora.
func buscarEvidenciasPorComentario(db *gorm.DB, idComentario int64) ([]EvidenciaResponse, error) {
	type fila struct {
		Evidencia
		TipoNombre string `gorm:"column:tipo_nombre"`
	}

	var filas []fila
	err := db.Model(&Evidencia{}).
		Select("organizacion.bit_evidencias.*, t.nombre AS tipo_nombre").
		Joins("INNER JOIN organizacion.ref_tipo_evidencia t ON t.id = organizacion.bit_evidencias.id_tipo").
		Where("organizacion.bit_evidencias.id_comentario = ?", idComentario).
		Order("organizacion.bit_evidencias.created_at ASC").
		Scan(&filas).Error
	if err != nil {
		return nil, err
	}

	resultados := make([]EvidenciaResponse, len(filas))
	for i, f := range filas {
		resultados[i] = toEvidenciaResponse(&f.Evidencia, f.TipoNombre)
	}
	return resultados, nil
}

// EliminarEvidenciaFisica borra físicamente la evidencia (no hay soft delete en bitácora).
func EliminarEvidenciaFisica(db *gorm.DB, id int64, idEmpresa int64) error {
	return db.Where("id = ? AND id_empresa = ?", id, idEmpresa).
		Delete(&Evidencia{}).Error
}

func buscarEvidenciaCruda(db *gorm.DB, id int64, idEmpresa int64) (*Evidencia, error) {
	var e Evidencia
	err := db.Model(&Evidencia{}).
		Where("id = ? AND id_empresa = ?", id, idEmpresa).
		First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func ExisteTipoEvidencia(db *gorm.DB, id int) (bool, error) {
	var total int64
	err := db.Model(&RefTipoEvidencia{}).
		Where("id = ? AND es_activo = ?", id, true).Count(&total).Error
	return total > 0, err
}

func resolverNombreTipo(db *gorm.DB, id int) (string, error) {
	var t RefTipoEvidencia
	err := db.Model(&RefTipoEvidencia{}).Where("id = ?", id).First(&t).Error
	if err != nil {
		return "", err
	}
	return t.Nombre, nil
}
