package seguridad

import (
	"gorm.io/gorm"
)

func GetCompanysAccess(db *gorm.DB, id int) ([]CompanyAccess, error) {

	var results []CompanyAccess

	err := db.
		Table("seguridad.cfg_empresas_roles re").
		Joins("INNER JOIN configuracion.cfg_empresas e ON e.id = re.id_empresa").
		Select(`
			re.id,
			e.id as id_empresa,
			e.razon_social,
			e.descripcion,
			e.usa_sedes
			`).
		Where("re.id_rol = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, err

}

func GetSedesAccess(db *gorm.DB, id int) ([]SedesAccess, error) {

	var results []SedesAccess

	err := db.
		Table("seguridad.cfg_sedes_roles sr").
		Joins("INNER JOIN configuracion.cfg_sedes AS s on s.id = sr.id_sede").
		Select(`
				sr.id,
                sr.id_empresa,
                sr.id_sede,
                s.nombre`).
		Where("sr.id_rol = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}
