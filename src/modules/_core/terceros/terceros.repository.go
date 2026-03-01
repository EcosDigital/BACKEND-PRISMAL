package terceros

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/ecosistema/core/src/database"
	"gorm.io/gorm"
)

// valida existencia tercero by documento
func FindTerceroByDocument(db *gorm.DB, numero_documento string) ([]TerceroResponse, error) {

	var tercero []TerceroResponse

	err := db.
		Table("configuracion.cfg_terceros").
		Select(`
			id,
			numero_documento,
			primer_nombre as nombre_completo,
			email,
			telefono,
			'' as tipo_documento,
			estado
		`).
		Where("numero_documento = ?", numero_documento).
		Scan(&tercero).Error

	if err != nil {
		return nil, err
	}

	return tercero, err

}

func CreateTercero(db *gorm.DB, req *TerceroRequest) (int64, error) {

	data := map[string]interface{}{
		"id_tipo_persona":     req.IdTipoPersona,
		"id_tipo_documento":   req.IdTipoDocumento,
		"numero_documento":    req.NumeroDocumento,
		"dv":                  req.Dv,
		"primer_nombre":       req.PrimerNombre,
		"segundo_nombre":      req.SegundoNombre,
		"primer_apellido":     req.PrimerApellido,
		"segundo_apellido":    req.SegundoApellido,
		"razon_social":        req.RazonSocial,
		"representante_legal": req.RepresentanteLegal,
		"id_genero":           req.IdGenero,
		"telefono":            req.Telefono,
		"telefono_2":          req.Telefono_2,
		"email":               req.Email,
		"pagina_web":          req.PaginaWeb,
		"direccion":           req.Direccion,
		"id_pais":             req.IdPais,
		"id_departamento":     req.IdDepartamento,
		"id_ciudad":           req.IdCiudad,
		"id_zona":             req.IdZona,
		"estado":              req.Estado,
		"created_at":          time.Now(),
		"created_by":          req.UserID,
	}

	err := db.
		Table("configuracion.cfg_terceros").
		Create(&data).Error

	if err != nil {
		return 0, err
	}

	var id_reg int64
	database.GormDB.Raw("SELECT lastval()").Scan(&id_reg)

	return id_reg, nil

}

func AddClaseTercero(db *gorm.DB, idTercero int64, IdClase int, IdUser int64) error {

	data := map[string]interface{}{
		"id_tercero": idTercero,
		"id_clase":   IdClase,
		"created_by": IdUser,
		"created_at": time.Now(),
	}

	err := db.
		Table("configuracion.cfg_clase_terceros").
		Create(&data).Error

	return err

}

func AddClasesTercero(db *gorm.DB, idTercero int64, clases []int, idUser int64) error {
	for _, IdClase := range clases {
		err := AddClaseTercero(db, idTercero, IdClase, idUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveClaseTerceros(db *gorm.DB, idTercero int64) error {

	err := db.
		Table("configuracion.cfg_clase_terceros").
		Where("id_tercero = ?", idTercero).
		Delete(nil).Error

	if err != nil {
		return err
	}

	return nil

}

func FindClasesByTerceroId(db *gorm.DB, idTercero int64) ([]int, error) {

	var clases []int

	err := db.
		Table("configuracion.cfg_clase_terceros").
		Select("id_clase").
		Where("id_tercero = ?", idTercero).
		Scan(&clases).Error

	if err != nil {
		return nil, err
	}

	return clases, nil

}

func FindLastTercero(db *gorm.DB) ([]TerceroResponse, error) {

	var results []TerceroResponse

	err := db.
		Table("configuracion.cfg_terceros t").
		Select(`
			t.id,
			t.numero_documento,
			CASE
				WHEN t.id_tipo_persona = 1 THEN
					TRIM(CONCAT_WS(' ', t.primer_nombre, t.segundo_nombre, t.primer_apellido, t.segundo_apellido))
				WHEN t.id_tipo_persona = 2 THEN
					t.razon_social
				ELSE ''	
			END as nombre_completo,
			t.email,
			t.telefono,
			td.nombre as tipo_documento,
			t.estado`).
		Joins("INNER JOIN configuracion.ref_tipo_documento td ON td.id = t.id_tipo_documento").
		Order("id ASC").Limit(10).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func ListTerceroById(db *gorm.DB, id int64) (*TerceroResponseFull, error) {

	var result TerceroResponseFull

	err := db.
		Table("configuracion.cfg_terceros").
		Where("id = ?", id).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	clases, err := FindClasesByTerceroId(db, id)
	if err != nil {
		return nil, err
	}

	result.ClaseTercero = clases

	return &result, err

}

func UpdateTerceroData(db *gorm.DB, id int64, req *TerceroUpdateRequest) (int64, error) {

	data := map[string]interface{}{
		"primer_nombre":       req.PrimerNombre,
		"segundo_nombre":      req.SegundoNombre,
		"primer_apellido":     req.PrimerApellido,
		"segundo_apellido":    req.SegundoApellido,
		"razon_social":        req.RazonSocial,
		"representante_legal": req.RepresentanteLegal,
		"id_genero":           req.IdGenero,
		"telefono":            req.Telefono,
		"telefono_2":          req.Telefono_2,
		"email":               req.Email,
		"pagina_web":          req.PaginaWeb,
		"direccion":           req.Direccion,
		"id_pais":             req.IdPais,
		"id_departamento":     req.IdDepartamento,
		"id_ciudad":           req.Idciudad,
		"id_zona":             req.IdZona,
		"estado":              req.Estado,
		"created_at":          time.Now(),
		"created_by":          req.UserID,
	}

	tx := db.
		Table("configuracion.cfg_terceros").
		Where("id = ?", id).
		Updates(data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil

}

func updateTerceroIdentity(db *sql.DB, id int, req *TerceroUpdateIndenty) (int, error) {

	var jsonResult []byte
	var id_usuario = 1
	var dv sql.NullString

	if req.Dv == "" || req.Dv == "0" || strings.TrimSpace(req.Dv) == "" {
		dv = sql.NullString{Valid: false} // se enviará como NULL
	} else {
		dv = sql.NullString{String: req.Dv, Valid: true}
	}

	query := `
        SELECT configuracion.qry_terceros(
            operacion => $1,
            id_tercero_p => $2,
			id_usuario_p => $3,
            id_tipo_persona_p => $4,
            id_tipo_documento_p => $5,
            numero_documento_p => $6,
            dv_p => $7
        )`

	err := db.QueryRow(
		query,
		5,
		id,
		id_usuario,
		req.IdTipoPersona,
		req.IdTipoDocumento,
		req.NumeroDocumento,
		dv,
	).Scan(&jsonResult)

	if err != nil {
		return 0, err
	}

	var result struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(jsonResult, &result); err != nil {
		return 0, err
	}

	return result.ID, nil

}
