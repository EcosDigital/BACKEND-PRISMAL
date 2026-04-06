package horarios

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Inserta registro (plantilla horario)
func CreateHorario(db *gorm.DB, req *HorarioRequest) (int64, error) {

	data := map[string]interface{}{
		"codigo":             req.Codigo,
		"nombre":             req.Nombre,
		"descripcion":        req.Descripcion,
		"minutos_tolerancia": req.MinutosTolerancia,
		"activo":             true,
		"fecha_creacion":     time.Now(),
	}

	if err := db.Table("control_asistencia.cfg_horarios_asistencia").Create(&data).Error; err != nil {
		return 0, err
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil

}

// bloques de horario (detalle)
func CreateBloques(db *gorm.DB, idHorario int64, bloques []BloqueHorarioRequest) error {

	rows := make([]map[string]interface{}, 0, len(bloques))

	for _, b := range bloques {
		rows = append(rows, map[string]interface{}{
			"id_horario":     idHorario,
			"hora_entrada":   b.HoraEntrada,
			"hora_salida":    b.HoraSalida,
			"orden":          b.Orden,
			"activo":         true,
			"fecha_creacion": time.Now(),
		})
	}

	return db.Table("control_asistencia.cfg_horarios_detalle").Create(&rows).Error

}

// registros recientes
func FindLastHorarios(db *gorm.DB) ([]HorarioListResponse, error) {

	var results []HorarioListResponse

	err := db.
		Table("control_asistencia.cfg_horarios_asistencia").
		Select(`
			id,
			codigo,
			nombre,
			minutos_tolerancia,
			activo,
			fecha_creacion::text AS fecha_creacion`).
		Order("fecha_creacion DESC").
		Limit(10).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

// registro por ID
func FindHorarioByID(db *gorm.DB, id int64) (*horarioRow, error) {

	var results *horarioRow

	err := db.
		Table("control_asistencia.cfg_horarios_asistencia").
		Select(`
			id,
			codigo,
			nombre,
			descripcion,
			minutos_tolerancia,
			activo,
			fecha_creacion::text AS fecha_creacion`).
		Where("id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results.ID == 0 {
		return nil, nil
	}

	return results, nil

}

// bloques por ID horario
func FindBloquesByHorario(db *gorm.DB, idHorario int64) ([]bloqueRow, error) {

	var rows []bloqueRow

	err := db.
		Table("control_asistencia.cfg_horarios_detalle").
		Select("id, id_horario, hora_entrada::text AS hora_entrada, hora_salida::text AS hora_salida, orden, activo").
		Where("id_horario = ? AND activo = true", idHorario).
		Order("orden ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// cambiar estado del horario
func ChangeHorarioEstado(db *gorm.DB, id int64, activo bool) error {

	result := db.
		Table("control_asistencia.cfg_horarios_asistencia").
		Where("id = ?", id).
		Update("activo", activo)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("horario con id %d no encontrado", id)
	}

	return nil
}

// EditHorario actualiza solo los campos básicos de la cabecera del horario.
// Los bloques (cfg_horarios_detalle) son inmutables y nunca se tocan aquí.
func UpdateHorarioData(db *gorm.DB, id int64, req *HorarioUpdateRequest) error {

	result := db.
		Table("control_asistencia.cfg_horarios_asistencia").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"nombre":             req.Nombre,
			"descripcion":        req.Descripcion,
			"minutos_tolerancia": req.MinutosTolerancia,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("horario con id %d no encontrado", id)
	}

	return nil
}

// Retorna true si ya existe otro horario con ese código.
func FindHorarioByCodigo(db *gorm.DB, codigo string, excludeID int64) (bool, error) {

	var count int64

	err := db.
		Table("control_asistencia.cfg_horarios_asistencia").
		Where("codigo = ? AND id <> ?", codigo, excludeID).
		Count(&count).Error

	return count > 0, err
}

// ─── cfg_horarios_asignacion ──────────────────────────────────────────────────
func InsertAsignacion(db *gorm.DB, req *AsignacionHorarioRequest) error {

	now := time.Now()

	// Cerrar asignación activa previa del tercero si existe
	if err := db.
		Table("control_asistencia.cfg_empleado_horario").
		Where("id_empleado = ? AND activo = true", req.IDTercero).
		Updates(map[string]interface{}{
			"activo":    false,
			"fecha_fin": now.Format("2006-01-02"),
		}).Error; err != nil {
		return err
	}

	// Insertar nueva asignación
	data := map[string]interface{}{
		"id_empleado":    req.IDTercero,
		"id_horario":     req.IDHorario,
		"fecha_inicio":   req.FechaInicio,
		"fecha_fin":      nil,
		"activo":         true,
		"observaciones":  req.Observaciones,
		"registrado_por": req.RegistradoPor,
		"fecha_creacion": now,
	}

	return db.Table("control_asistencia.cfg_empleado_horario").Create(&data).Error
}

// FindAsignacionesByHorario devuelve las asignaciones activas de un horario,
// haciendo JOIN con cfg_terceros para obtener nombre y documento del empleado.
func FindAsignacionesByHorario(db *gorm.DB, idHorario int64) ([]AsignacionListResponse, error) {

	var rows []AsignacionListResponse

	err := db.
		Table("control_asistencia.cfg_empleado_horario eh").
		Select(`
			eh.id,
			eh.id_empleado,
			TRIM(COALESCE(t.primer_nombre, '') || ' ' || COALESCE(t.segundo_nombre, '') || ' ' || COALESCE(t.primer_apellido, '') || ' ' || COALESCE(t.segundo_apellido, '')) AS nombre_tercero,
			COALESCE(t.numero_documento, '')                               AS numero_documento,
			eh.fecha_inicio::text                                           AS fecha_inicio,
			COALESCE(eh.observaciones, '')                                 AS observaciones`).
		Joins("LEFT JOIN configuracion.cfg_terceros t ON t.id = eh.id_empleado").
		Where("eh.id_horario = ? AND eh.activo = true", idHorario).
		Order("eh.fecha_creacion DESC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	if rows == nil {
		rows = []AsignacionListResponse{}
	}

	return rows, nil
}

// DeleteAsignacion desactiva una asignación por su ID (soft delete).
// Retorna error si no existe o ya estaba inactiva.
// DeleteAsignacion elimina físicamente una asignación por su ID.
func DeleteAsignacion(db *gorm.DB, idAsignacion int64) error {

	result := db.
		Table("control_asistencia.cfg_empleado_horario").
		Where("id = ?", idAsignacion).
		Delete(nil)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("asignación con id %d no encontrada", idAsignacion)
	}

	return nil
}

func FindAsignacionActivaByTercero(db *gorm.DB, idTercero int64) (bool, error) {

	var count int64

	err := db.
		Table("control_asistencia.cfg_empleado_horario").
		Where("id_tercero = ? AND activo = true", idTercero).
		Count(&count).Error

	return count > 0, err
}
