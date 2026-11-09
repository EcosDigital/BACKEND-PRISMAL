package configuracion

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// registrar agente
func CreateAgente(db *gorm.DB, req *AgenteRequest) (int64, error) {

	data := map[string]interface{}{
		"nombre":        req.Nombre,
		"instrucciones": req.Instrucciones,
		"modelo":        req.Modelo,
		"is_active":     req.IsActive,
		"id_modulo_ref": req.IDModuloRef,
		"created_at":    time.Now(),
		"created_by":    req.UserID,
		"id_empresa":    req.EmpresaID,
		"id_sede":       req.SedeID,
	}

	err := db.
		Table("agentes.cfg_agentes").
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

// consultar todos los agentes
func ListAgentes(db *gorm.DB) ([]AgenteResponse, error) {

	var results []AgenteResponse

	err := db.
		Table("agentes.cfg_agentes a").
		Select(`
			a.id,
			a.nombre,
			a.instrucciones,
			a.modelo,
			a.is_active,
			a.id_modulo_ref,
			COALESCE(m.nombre, '') as modulo`).
		Joins("LEFT JOIN configuracion.ref_modulos_tenant m ON m.id = a.id_modulo_ref").
		Order("a.id DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []AgenteResponse{}
	}

	return results, nil

}

// consultar agente por id
func ListAgenteByID(db *gorm.DB, id int64) ([]AgenteResponse, error) {

	var results []AgenteResponse

	err := db.
		Table("agentes.cfg_agentes a").
		Select(`
			a.id,
			a.nombre,
			a.instrucciones,
			a.modelo,
			a.is_active,
			a.id_modulo_ref,
			COALESCE(m.nombre, '') as modulo`).
		Joins("LEFT JOIN configuracion.ref_modulos_tenant m ON m.id = a.id_modulo_ref").
		Where("a.id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []AgenteResponse{}
	}

	return results, nil

}

// actualizar agente
func UpdateAgente(db *gorm.DB, id int64, req *AgenteRequest) (int64, error) {

	data := map[string]interface{}{
		"nombre":        req.Nombre,
		"instrucciones": req.Instrucciones,
		"modelo":        req.Modelo,
		"is_active":     req.IsActive,
		"id_modulo_ref": req.IDModuloRef,
		"update_at":     time.Now(),
		"update_by":     req.UserID,
	}

	err := db.Table("agentes.cfg_agentes").
		Where("id = ?", id).
		Updates(data)

	if err.Error != nil {
		return 0, err.Error
	}

	if err.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil

}

// consultar los módulos instalados en el negocio
func ListModulosInstalados(db *gorm.DB) ([]ModuloResponse, error) {

	var results []ModuloResponse

	err := db.
		Table("configuracion.ref_modulos_tenant").
		Select("id, nombre").
		Order("nombre ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ModuloResponse{}
	}

	return results, nil

}

// verificar si ya hay otro agente encendido encargado de ese módulo
func ExistsEncargadoActivo(db *gorm.DB, idModuloRef int64, excluirID int64) (bool, error) {

	var count int64

	err := db.
		Table("agentes.cfg_agentes").
		Where("id_modulo_ref = ? AND is_active = true AND id <> ?", idModuloRef, excluirID).
		Count(&count).Error

	return count > 0, err

}

// consultar la sede a la que pertenece el agente
func GetSedeAgente(db *gorm.DB, id int64) (int64, error) {

	var sede int64

	err := db.
		Table("agentes.cfg_agentes").
		Select("COALESCE(id_sede, 0)").
		Where("id = ?", id).
		Scan(&sede).Error

	return sede, err

}

// guardar los datos de WhatsApp del agente
func GuardarWhatsAppAgente(db *gorm.DB, id int64, conexion ConexionWhatsApp, userID int64) error {

	conexionJSON, err := json.Marshal(conexion)
	if err != nil {
		return err
	}

	resultado := db.Exec(`
		UPDATE agentes.cfg_agentes
		SET details = jsonb_set(details, '{whatsapp}', ?::jsonb), update_at = NOW(), update_by = ?
		WHERE id = ?`,
		string(conexionJSON), userID, id)

	if resultado.Error != nil {
		return resultado.Error
	}

	if resultado.RowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

// consultar el agente que atiende un número de WhatsApp conectado
func FilterAgentePorTelefonoMeta(db *gorm.DB, idTelefonoMeta string) ([]AgenteWhatsApp, error) {

	filtro, err := json.Marshal(map[string]interface{}{
		"whatsapp": map[string]interface{}{"id_telefono_meta": idTelefonoMeta, "conectado": true},
	})
	if err != nil {
		return nil, err
	}

	var results []AgenteWhatsApp

	err = db.
		Table("agentes.cfg_agentes").
		Select(`
			id, nombre, instrucciones, modelo, is_active,
			id_empresa AS empresa_id,
			COALESCE(id_sede, 0) AS sede_id,
			details->'whatsapp'->>'llave_cifrada' AS llave_cifrada`).
		Where("details @> ?::jsonb", string(filtro)).
		Scan(&results).Error

	return results, err

}

// verificar si el número ya está conectado a otro agente
func ExistsTelefonoEnOtroAgente(db *gorm.DB, idTelefonoMeta string, excluirID int64) (bool, error) {

	filtro, err := json.Marshal(map[string]interface{}{
		"whatsapp": map[string]interface{}{"id_telefono_meta": idTelefonoMeta, "conectado": true},
	})
	if err != nil {
		return false, err
	}

	var count int64

	err = db.
		Table("agentes.cfg_agentes").
		Where("details @> ?::jsonb AND id <> ?", string(filtro), excluirID).
		Count(&count).Error

	return count > 0, err

}
