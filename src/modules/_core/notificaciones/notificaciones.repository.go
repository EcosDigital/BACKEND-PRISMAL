package notificaciones

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// consultar los tipos de notificacion disponibles
func ListTiposNotificacion(db *gorm.DB) ([]TipoNotificacionResponse, error) {

	var results []TipoNotificacionResponse

	err := db.
		Table("notificaciones.ref_tipo_notificacion").
		Select("id, nombre").
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []TipoNotificacionResponse{}
	}

	return results, nil

}

// registrar el token del navegador, o reactivarlo si ya existia.
// el token pertenece al navegador, no a la persona: si otro usuario
// entra en el mismo equipo, el token se reasigna a quien esta entrando.
func UpsertDispositivoPush(db *gorm.DB, req *DispositivoPushRequest) (int64, error) {

	data := map[string]interface{}{
		"id_usuario": req.UserID,
		"token":      req.Token,
		"navegador":  req.Navegador,
		"activo":     true,
		"created_at": time.Now(),
	}

	err := db.
		Table("notificaciones.cfg_dispositivos_push").
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "token"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"id_usuario": req.UserID,
				"navegador":  req.Navegador,
				"activo":     true,
				"updated_at": time.Now(),
			}),
		}).
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

// consultar los destinos disponibles (usuario, rol, modulo)
func ListOrigenesDestinatario(db *gorm.DB) ([]OrigenDestinatarioResponse, error) {

	var results []OrigenDestinatarioResponse

	err := db.
		Table("notificaciones.ref_origen_destinatario").
		Select("id, nombre").
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []OrigenDestinatarioResponse{}
	}

	return results, nil

}

// registrar el mensaje de la notificacion
func CreateNotificacion(db *gorm.DB, req *NotificacionRequest) (int64, error) {

	row := NotificacionRow{
		Titulo:             req.Titulo,
		Mensaje:            req.Mensaje,
		IDModuloRef:        req.IDModuloRef,
		IDTipoNotificacion: req.IDTipoNotificacion,
		EnviadoPor:         req.UserID,
		IDEmpresa:          req.EmpresaID,
		IDSede:             req.SedeID,
		CreatedAt:          time.Now(),
	}

	if err := db.
		Table("notificaciones.mov_notificaciones").
		Create(&row).Error; err != nil {
		return 0, err
	}

	if row.ID == 0 {
		return 0, errors.New("no se pudo obtener el id de la notificación")
	}

	return row.ID, nil

}

// dar de baja el token (cierre de sesion o el usuario desactiva las
// notificaciones). solo afecta los dispositivos del propio usuario.
func UpdateBajaDispositivoPush(db *gorm.DB, req *BajaDispositivoRequest) error {

	result := db.
		Table("notificaciones.cfg_dispositivos_push").
		Where("token = ? AND id_usuario = ?", req.Token, req.UserID).
		Updates(map[string]interface{}{
			"activo":     false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

// consultar un usuario activo por su id
func ListUsuarioActivoByID(db *gorm.DB, idUsuario int64) ([]int64, error) {

	var results []int64

	err := db.
		Table("seguridad.cfg_usuarios").
		Select("id").
		Where("id = ? AND activo = ?", idUsuario, true).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []int64{}
	}

	return results, nil

}

// consultar los usuarios activos que pertenecen a los roles indicados
func ListUsuariosByRoles(db *gorm.DB, idsRoles []int64) ([]int64, error) {

	var results []int64

	err := db.
		Table("seguridad.cfg_usuarios").
		Distinct("id").
		Where("id_rol IN ? AND activo = ?", idsRoles, true).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []int64{}
	}

	return results, nil

}

// consultar el id_ref (id del catalogo admin) de un modulo del tenant
func ListIDRefModulo(db *gorm.DB, idModulo int64) ([]int64, error) {

	var results []int64

	err := db.
		Table("configuracion.ref_modulos_tenant").
		Select("id_ref").
		Where("id = ?", idModulo).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []int64{}
	}

	return results, nil

}

// consultar el arbol de permisos de los roles configurados en una sede
func ListConfigRolesBySede(db *gorm.DB, idSede int64) ([]ConfigRolModulos, error) {

	var results []ConfigRolModulos

	err := db.
		Table("seguridad.cfg_sedes_roles").
		Select("id_rol, json_modules").
		Where("id_sede = ?", idSede).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ConfigRolModulos{}
	}

	return results, nil

}

// registrar el reparto de la notificacion. el ON CONFLICT evita duplicar
// a una misma persona que llegue por dos caminos (ej. dos roles distintos).
func CreateDestinatarios(db *gorm.DB, idNotificacion int64, idsUsuarios []int64, idOrigen int) error {

	rows := make([]map[string]interface{}, 0, len(idsUsuarios))

	for _, idUsuario := range idsUsuarios {
		rows = append(rows, map[string]interface{}{
			"id_notificacion": idNotificacion,
			"id_usuario":      idUsuario,
			"id_origen":       idOrigen,
			"leido":           false,
			"created_at":      time.Now(),
		})
	}

	return db.
		Table("notificaciones.mov_destinatarios").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&rows).Error

}

// consultar los tokens vigentes de los usuarios indicados. se excluyen los
// dispositivos sin actividad reciente: si el usuario no abre Prismar hace
// mas de ese tiempo, su sesion ya vencio y no debe recibir notificaciones.
func ListTokensByUsuarios(db *gorm.DB, idsUsuarios []int64, horasVigencia int) ([]string, error) {

	var results []string

	err := db.
		Table("notificaciones.cfg_dispositivos_push").
		Select("token").
		Where("id_usuario IN ? AND activo = ?", idsUsuarios, true).
		Where("COALESCE(updated_at, created_at) >= NOW() - make_interval(hours => ?)", horasVigencia).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []string{}
	}

	return results, nil

}

// dar de baja los tokens que firebase reporto como invalidos
func UpdateBajaTokens(db *gorm.DB, tokens []string) error {

	if len(tokens) == 0 {
		return nil
	}

	return db.
		Table("notificaciones.cfg_dispositivos_push").
		Where("token IN ?", tokens).
		Updates(map[string]interface{}{
			"activo":     false,
			"updated_at": time.Now(),
		}).Error

}

// consultar las notificaciones que le llegaron a un usuario en su sede.
// se incluyen las que no tienen sede asignada, para no ocultar avisos
// generales de la empresa.
func ListNotificacionesByUsuario(db *gorm.DB, idUsuario int64, idSede int64) ([]NotificacionUsuarioResponse, error) {

	var results []NotificacionUsuarioResponse

	err := db.
		Table("notificaciones.mov_destinatarios d").
		Select(`
			n.id,
			n.titulo,
			n.mensaje,
			t.nombre as tipo,
			COALESCE(m.nombre, '') as modulo,
			o.nombre as origen,
			d.leido,
			n.created_at::text as created_at`).
		Joins("INNER JOIN notificaciones.mov_notificaciones n ON n.id = d.id_notificacion").
		Joins("INNER JOIN notificaciones.ref_tipo_notificacion t ON t.id = n.id_tipo_notificacion").
		Joins("INNER JOIN notificaciones.ref_origen_destinatario o ON o.id = d.id_origen").
		Joins("LEFT JOIN configuracion.ref_modulos_tenant m ON m.id = n.id_modulo_ref").
		Where("d.id_usuario = ?", idUsuario).
		Where("n.id_sede = ? OR n.id_sede IS NULL", idSede).
		Order("n.created_at DESC").
		Limit(maxHistorial).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []NotificacionUsuarioResponse{}
	}

	return results, nil

}

// contar las notificaciones sin leer de un usuario en su sede
func CountNoLeidasByUsuario(db *gorm.DB, idUsuario int64, idSede int64) (int64, error) {

	var total int64

	err := db.
		Table("notificaciones.mov_destinatarios d").
		Joins("INNER JOIN notificaciones.mov_notificaciones n ON n.id = d.id_notificacion").
		Where("d.id_usuario = ? AND d.leido = ?", idUsuario, false).
		Where("n.id_sede = ? OR n.id_sede IS NULL", idSede).
		Count(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil

}

// marcar una notificacion como leida para un usuario. es idempotente: si ya
// estaba leida no falla, y COALESCE conserva la hora de la primera lectura.
// si no afecta ninguna fila, esa notificacion no es de este usuario.
func UpdateLeidoByUsuario(db *gorm.DB, idNotificacion int64, idUsuario int64) error {

	result := db.
		Table("notificaciones.mov_destinatarios").
		Where("id_notificacion = ? AND id_usuario = ?", idNotificacion, idUsuario).
		Updates(map[string]interface{}{
			"leido":    true,
			"leido_at": gorm.Expr("COALESCE(leido_at, NOW())"),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

// marcar como leidas todas las notificaciones pendientes del usuario en su
// sede, con el mismo criterio con el que las ve en el historial.
func UpdateLeidoTodasByUsuario(db *gorm.DB, idUsuario int64, idSede int64) (int64, error) {

	notificacionesVisibles := db.
		Table("notificaciones.mov_notificaciones").
		Select("id").
		Where("id_sede = ? OR id_sede IS NULL", idSede)

	result := db.
		Table("notificaciones.mov_destinatarios").
		Where("id_usuario = ? AND leido = ?", idUsuario, false).
		Where("id_notificacion IN (?)", notificacionesVisibles).
		Updates(map[string]interface{}{
			"leido":    true,
			"leido_at": time.Now(),
		})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil

}

// usuarios activos para el selector de destino, como "Nombre (correo)" para
// distinguir personas con el mismo nombre
func ListUsuariosDestino(db *gorm.DB) ([]DestinoResponse, error) {

	var results []DestinoResponse

	err := db.
		Table("seguridad.cfg_usuarios u").
		Select(`
			u.id,
			CONCAT(
				CASE
					WHEN t.id_tipo_persona = 1 THEN
						TRIM(CONCAT_WS(' ', t.primer_nombre, t.segundo_nombre, t.primer_apellido, t.segundo_apellido))
					ELSE t.razon_social
				END,
				' (', u.email, ')'
			) as nombre`).
		Joins("INNER JOIN configuracion.cfg_terceros t ON t.id = u.id_tercero").
		Where("u.activo = ?", true).
		Order("nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []DestinoResponse{}
	}

	return results, nil

}

// roles para el selector de destino
func ListRolesDestino(db *gorm.DB) ([]DestinoResponse, error) {

	var results []DestinoResponse

	err := db.
		Table("seguridad.cfg_roles_usuario").
		Select("id, nombre").
		Order("nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []DestinoResponse{}
	}

	return results, nil

}
