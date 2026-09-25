package notificaciones

import (
	"time"

	"github.com/ecosistema/core/src/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// consultar los modulos contratados por el tenant
func ListModulosTenant(db *gorm.DB) ([]moduloTenant, error) {

	var results []moduloTenant

	err := db.
		Table("configuracion.ref_modulos_tenant").
		Select("id, id_ref").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []moduloTenant{}
	}

	return results, nil

}

// consultar en la base admin los eventos activos de los modulos indicados
// (ids del catalogo admin). usa la conexion global a la base admin, no la
// del tenant: el catalogo maestro solo vive alli
func ListEventosCatalogoByModulos(idsModulos []int64) ([]eventoCatalogo, error) {

	var results []eventoCatalogo

	err := database.GormDB.
		Table("configuracion.cfg_eventos_notificacion e").
		Select("e.id, e.codigo, e.nombre, e.descripcion, e.id_modulo, t.nombre as tipo, e.orden_lista").
		Joins("INNER JOIN notificaciones.ref_tipo_notificacion t ON t.id = e.id_tipo_notificacion").
		Where("e.id_modulo IN ? AND e.is_active = ?", idsModulos, true).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []eventoCatalogo{}
	}

	return results, nil

}

// consultar los eventos ya copiados en el tenant, activos o no
func ListEventosLocales(db *gorm.DB) ([]eventoLocal, error) {

	var results []eventoLocal

	err := db.
		Table("notificaciones.ref_evento_notificacion").
		Select("id_ref, codigo, nombre, descripcion, id_modulo_ref, id_tipo_notificacion, orden_lista, is_active").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []eventoLocal{}
	}

	return results, nil

}

// registrar los eventos que el tenant todavia no tiene. el ON CONFLICT cubre
// el caso de dos personas abriendo la pantalla al mismo tiempo: el segundo
// insert del mismo evento se ignora en vez de fallar
func CreateEventosTenant(db *gorm.DB, rows []map[string]interface{}) error {

	if len(rows) == 0 {
		return nil
	}

	return db.
		Table("notificaciones.ref_evento_notificacion").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id_ref"}},
			DoNothing: true,
		}).
		Create(&rows).Error

}

// actualizar un evento del tenant que cambio en el catalogo admin
func UpdateEventoTenant(db *gorm.DB, idRef int64, data map[string]interface{}) error {

	return db.
		Table("notificaciones.ref_evento_notificacion").
		Where("id_ref = ?", idRef).
		Updates(data).Error

}

// desactivar los eventos locales que ya no vienen del catalogo (ej. el
// tenant dejo de tener el modulo). con la lista vacia se desactivan todos
func UpdateDesactivarEventosTenant(db *gorm.DB, idsRefVigentes []int64) error {

	query := db.
		Table("notificaciones.ref_evento_notificacion").
		Where("is_active = ?", true)

	if len(idsRefVigentes) > 0 {
		query = query.Where("id_ref NOT IN ?", idsRefVigentes)
	}

	return query.Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error

}

// consultar los eventos activos del tenant con el modulo al que pertenecen
func ListEventosTenant(db *gorm.DB) ([]eventoModuloRow, error) {

	var results []eventoModuloRow

	err := db.
		Table("notificaciones.ref_evento_notificacion e").
		Select(`
			m.id as id_modulo,
			m.codigo as codigo_modulo,
			m.nombre as nombre_modulo,
			e.codigo,
			e.nombre,
			COALESCE(e.descripcion, '') as descripcion,
			t.nombre as tipo`).
		Joins("INNER JOIN configuracion.ref_modulos_tenant m ON m.id = e.id_modulo_ref").
		Joins("INNER JOIN notificaciones.ref_tipo_notificacion t ON t.id = e.id_tipo_notificacion").
		Where("e.is_active = ?", true).
		Order("m.nombre ASC, e.orden_lista ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []eventoModuloRow{}
	}

	return results, nil

}

// consultar el id local de un evento activo por su codigo
func ListIDEventoByCodigo(db *gorm.DB, codigo string) ([]int64, error) {

	var results []int64

	err := db.
		Table("notificaciones.ref_evento_notificacion").
		Select("id").
		Where("codigo = ? AND is_active = ?", codigo, true).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []int64{}
	}

	return results, nil

}

// consultar los destinos configurados para un evento en una sede
func ListDestinosEvento(db *gorm.DB, idEvento int64, idSede int64) ([]eventoDestinoRow, error) {

	var results []eventoDestinoRow

	err := db.
		Table("notificaciones.cfg_evento_destinatarios").
		Select("id_origen, id_destino").
		Where("id_evento = ? AND id_sede = ?", idEvento, idSede).
		Order("id_origen ASC, id_destino ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []eventoDestinoRow{}
	}

	return results, nil

}

// consultar cuales de los usuarios indicados existen y estan activos
func ListUsuariosActivosByIDs(db *gorm.DB, ids []int64) ([]int64, error) {

	var results []int64

	err := db.
		Table("seguridad.cfg_usuarios").
		Select("id").
		Where("id IN ? AND activo = ?", ids, true).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []int64{}
	}

	return results, nil

}

// consultar cuales de los roles indicados existen
func ListRolesByIDs(db *gorm.DB, ids []int64) ([]int64, error) {

	var results []int64

	err := db.
		Table("seguridad.cfg_roles_usuario").
		Select("id").
		Where("id IN ?", ids).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []int64{}
	}

	return results, nil

}

// borrar los destinos de un evento en una sede, para reemplazarlos
func DeleteDestinosEvento(db *gorm.DB, idEvento int64, idSede int64) error {

	return db.
		Table("notificaciones.cfg_evento_destinatarios").
		Where("id_evento = ? AND id_sede = ?", idEvento, idSede).
		Delete(&eventoDestinoRow{}).Error

}

// registrar los destinos de un evento en una sede
func CreateDestinosEvento(db *gorm.DB, rows []map[string]interface{}) error {

	if len(rows) == 0 {
		return nil
	}

	return db.
		Table("notificaciones.cfg_evento_destinatarios").
		Create(&rows).Error

}

// consultar los usuarios activos con su rol, para el selector y la tabla
// de Parametros generales
func ListUsuariosDestinoDetalle(db *gorm.DB) ([]UsuarioDestinoDetalle, error) {

	var results []UsuarioDestinoDetalle

	err := db.
		Table("seguridad.cfg_usuarios u").
		Select(`
			u.id,
			CASE
				WHEN t.id_tipo_persona = 1 THEN
					TRIM(CONCAT_WS(' ', t.primer_nombre, t.segundo_nombre, t.primer_apellido, t.segundo_apellido))
				ELSE t.razon_social
			END as nombre,
			u.email,
			u.id_rol,
			r.nombre as rol`).
		Joins("INNER JOIN configuracion.cfg_terceros t ON t.id = u.id_tercero").
		Joins("INNER JOIN seguridad.cfg_roles_usuario r ON r.id = u.id_rol").
		Where("u.activo = ?", true).
		Order("nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []UsuarioDestinoDetalle{}
	}

	return results, nil

}

// consultar un evento activo del tenant por su codigo, con su modulo y tipo
func ListEventoActivoByCodigo(db *gorm.DB, codigo string) ([]eventoActivo, error) {

	var results []eventoActivo

	err := db.
		Table("notificaciones.ref_evento_notificacion").
		Select("id, id_modulo_ref, id_tipo_notificacion").
		Where("codigo = ? AND is_active = ?", codigo, true).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []eventoActivo{}
	}

	return results, nil

}
