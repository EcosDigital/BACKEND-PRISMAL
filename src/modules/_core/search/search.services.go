package search

import (
	"database/sql")

func FilterDynamic(db *sql.DB, data SearchDynamics) ([]interface{}, error){
	return ListDynamics(db, data)
}