package reports

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func ListReportByID(db *gorm.DB, id int) ([]InformeDB, error) {

	var results []InformeDB

	err := db.
		Table("reportes.cfg_informes_generales").
		Select(`
			id,
			funcion_sql,
			is_active`).
		Where("id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []InformeDB{}
	}

	return results, nil

}

func ExecuteReport(
	db *gorm.DB,
	functioName string,
	params map[string]interface{},
) ([]map[string]interface{}, error) {

	var results []map[string]interface{}

	//construir argumentos nombrados
	args := []string{}
	values := []interface{}{}

	i := 1
	for key, value := range params {
		args = append(args, fmt.Sprintf("%s := ?", key))
		values = append(values, value)
		i++
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s(%s)",
		functioName,
		strings.Join(args, ", "),
	)

	err := db.Raw(query, values...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil

}
