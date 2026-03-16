package search

import "gorm.io/gorm"

func FilterDynamic(db *gorm.DB, data SearchDynamics) ([]interface{}, error) {
	return ListDynamics(db, data)
}
