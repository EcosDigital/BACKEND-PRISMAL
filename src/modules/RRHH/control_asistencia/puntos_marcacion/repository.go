package control_asistencia

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterPuntoMarcacion inserta un nuevo punto usando
func CreatePuntoMarcacion(db *gorm.DB, req *PuntoMarcacionRequest) (int64, error) {

	radioMetros := 80
	if req.RadioMetros > 0 {
		radioMetros = req.RadioMetros
	}

	data := map[string]interface{}{
		"id_sede":        req.SedeID,
		"id_empresa":     req.EmpresaID,
		"id_tipo_punto":  req.IDTipoPunto,
		"nombre":         req.Nombre,
		"descripcion":    req.Descripcion,
		"token_qr":       uuid.NewString(),
		"latitud":        req.Latitud,
		"longitud":       req.Longitud,
		"radio_metros":   radioMetros,
		"activo":         true,
		"fecha_creacion": time.Now(),
	}

	tx := db.Table("control_asistencia.cfg_puntos_marcacion").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil

}

// // ListPuntosMarcacion devuelve todos los puntos
func FindPuntosMarcacion(db *gorm.DB, empresaID int64) ([]PuntoMarcacionListResponse, error) {

	var results []PuntoMarcacionListResponse

	err := db.
		Table("control_asistencia.cfg_puntos_marcacion p").
		Select(`
			p.id,
			p.id_sede,
			t.nombre                        AS tipo_punto,
			p.nombre,
			p.token_qr,
			p.latitud,
			p.longitud,
			p.radio_metros,
			p.activo,
			p.fecha_creacion::text          AS fecha_creacion`).
		Joins("INNER JOIN control_asistencia.ref_tipo_punto_marcacion t ON t.id = p.id_tipo_punto").
		Where("p.id_empresa = ?", empresaID).
		Order("p.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []PuntoMarcacionListResponse{}
	}

	return results, nil
}

// GetPuntoMarcacionByID devuelve el detalle completo de un punto por su ID.
func FindPuntoMarcacionByID(db *gorm.DB, id int64) (*PuntoMarcacionResponse, error) {

	var result PuntoMarcacionResponse

	err := db.
		Table("control_asistencia.cfg_puntos_marcacion p").
		Select(`
			p.id,
			p.id_sede,
			p.id_tipo_punto,
			t.nombre                        AS tipo_punto,
			p.nombre,
			COALESCE(p.descripcion, '')     AS descripcion,
			p.token_qr,
			p.latitud,
			p.longitud,
			p.radio_metros,
			p.activo,
			p.fecha_creacion::text          AS fecha_creacion`).
		Joins("INNER JOIN control_asistencia.ref_tipo_punto_marcacion t ON t.id = p.id_tipo_punto").
		Where("p.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	if result.ID == 0 {
		return nil, nil
	}

	return &result, nil

}

// EditPuntoMarcacion actualiza los campos modificables de un punto.
// El token_qr es inmutable y nunca se incluye en el SET.
func EditPuntoMarcacion(db *gorm.DB, id int64, req *UpdatePuntoMarcacionRequest) (*PuntoMarcacionResponse, error) {

	// Construir SET dinámico solo con los campos enviados
	updates := map[string]interface{}{}

	if req.SedeID > 0 {
		updates["id_sede"] = req.SedeID
	}
	if req.IDTipoPunto > 0 {
		updates["id_tipo_punto"] = req.IDTipoPunto
	}
	if req.Nombre != "" {
		updates["nombre"] = req.Nombre
	}
	if req.Descripcion != "" {
		updates["descripcion"] = req.Descripcion
	}
	if req.Latitud != 0 {
		updates["latitud"] = req.Latitud
	}
	if req.Longitud != 0 {
		updates["longitud"] = req.Longitud
	}
	if req.RadioMetros > 0 {
		updates["radio_metros"] = req.RadioMetros
	}

	if len(updates) == 0 {
		return FindPuntoMarcacionByID(db, id)
	}

	if err := db.
		Table("control_asistencia.cfg_puntos_marcacion").
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	return FindPuntoMarcacionByID(db, id)
}

// ChangePuntoMarcacionEstado activa o inactiva un punto de marcación.
func ChangePuntoMarcacionEstado(db *gorm.DB, id int64, activo bool) error {

	result := db.
		Table("control_asistencia.cfg_puntos_marcacion").
		Where("id = ?", id).
		Update("activo", activo)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("punto de marcación con id %d no encontrado", id)
	}

	return nil
}

// ExisteNombrePuntoMarcacion verifica unicidad del nombre, excluyendo un ID.
func FindPuntoMarcacionByNombre(db *gorm.DB, nombre string, excludeID int64) (bool, error) {

	var count int64
	err := db.
		Table("control_asistencia.cfg_puntos_marcacion").
		Where("nombre = ? AND id <> ?", nombre, excludeID).
		Count(&count).Error

	return count > 0, err
}

func FindTiposPuntoMarcacion(db *gorm.DB) ([]RefTipoPuntoMarcacion, error) {

	var results []RefTipoPuntoMarcacion

	err := db.
		Table("control_asistencia.ref_tipo_punto_marcacion").
		Select("id, codigo, nombre, COALESCE(descripcion, '') AS descripcion").
		Where("activo = true").
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []RefTipoPuntoMarcacion{}
	}

	return results, nil
}

// ======= MARCACIONES (EMPLEADOS SCANN QR) ========= //

// FindPuntoByToken busca un punto de marcación activo por su token QR.
// Retorna nil si no existe o está inactivo.
func FindPuntoByToken(db *gorm.DB, token string) (*PuntoMarcacionValidacion, error) {

	var result PuntoMarcacionValidacion

	err := db.
		Table("control_asistencia.cfg_puntos_marcacion").
		Select("id, latitud, longitud, radio_metros, activo").
		Where("token_qr = ? AND activo = true", token).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	if result.ID == 0 {
		return nil, nil
	}

	return &result, nil

}

// FindOrigenQR obtiene el id del origen "QR" desde ref_origen_marcacion.
// Si no existe, retorna 0 y nil — el service decide si forzar fallo o no.
func FindOrigenQR(db *gorm.DB) (int64, error) {

	var id int64

	err := db.
		Table("control_asistencia.ref_origen_marcacion").
		Select("id").
		Where("UPPER(nombre) = 'QR' AND activo = true").
		Scan(&id).Error

	if err != nil {
		return 0, err
	}

	return id, nil
}

// ExisteMarcacionReciente verifica si el empleado ya registró una marcación
// en el mismo punto dentro de los últimos segundosTolerados segundos.
// Previene el fraude de doble marcación en corto tiempo.
func ExisteMarcacionReciente(db *gorm.DB, idEmpleado int64, idPunto int64, segundosTolerados int) (bool, error) {

	var count int64

	err := db.
		Table("control_asistencia.mov_marcaciones").
		Where(`
			id_empleado        = ? AND
			id_punto_marcacion = ? AND
			fecha_hora         >= NOW() - (? * INTERVAL '1 second')
		`, idEmpleado, idPunto, segundosTolerados).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// InsertMarcacion inserta el evento de marcación en mov_marcaciones.
// Retorna el ID generado y la fecha_hora del registro.
func InsertMarcacion(db *gorm.DB, req *MarcacionRequest, idPunto int64, idOrigen int64) (int64, time.Time, error) {

	data := map[string]interface{}{
		"id_empleado":        req.IDEmpleado,
		"id_punto_marcacion": idPunto,
		"latitud":            req.Latitud,
		"longitud":           req.Longitud,
		"precision_metros":   req.Precision,
		"direccion_ip":       req.DireccionIP,
		"user_agent":         req.UserAgent,
		"tipo_evento":        "ASISTENCIA",
		"origen":             idOrigen,
		"fecha_hora":         time.Now(),
		"fecha_creacion":     time.Now(),
	}

	tx := db.Table("control_asistencia.mov_marcaciones").Create(&data)
	if tx.Error != nil {
		return 0, time.Time{}, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, time.Time{}, err
	}

	// Recuperar la fecha_hora exacta que quedó en BD
	var fechaHora time.Time
	if err := db.
		Table("control_asistencia.mov_marcaciones").
		Select("fecha_hora").
		Where("id = ?", id).
		Scan(&fechaHora).Error; err != nil {
		return id, time.Now(), nil // no fatal si falla el fetch de fecha
	}

	return id, fechaHora, nil
}
