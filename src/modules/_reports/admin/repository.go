package reports

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func CreateConnection(db *gorm.DB, req *ReportServerConnectionRequest, userID string) (int64, error) {

	data := map[string]interface{}{
		"name":       req.Name,
		"base_url":   req.BaseURL,
		"api_key":    req.APIKey,
		"is_active":  false,
		"created_by": userID,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	tx := db.Table("administracion.report_server_connections").Create(&data)
	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

func ListConnections(db *gorm.DB, userID string) ([]ReportServerConnectionResponse, error) {

	var results []ReportServerConnectionResponse

	err := db.
		Table("administracion.report_server_connections").
		Select(`
			id,
			name,
			base_url,
			is_active,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at`).
		Where("created_by = ?", userID).
		Order("id ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ReportServerConnectionResponse{}
	}

	return results, nil
}

func ListConnectionByID(db *gorm.DB, id int64, userID string) (*ReportServerConnectionResponseFull, error) {

	var result ReportServerConnectionResponseFull

	err := db.
		Table("administracion.report_server_connections").
		Select(`
			id,
			name,
			base_url,
			api_key,
			is_active,
			created_by,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
			TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at`).
		Where("id = ? AND created_by = ?", id, userID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateConnection(db *gorm.DB, id int64, req *ReportServerConnectionUpdateRequest, userID string) (int64, error) {

	data := map[string]interface{}{
		"name":       req.Name,
		"base_url":   req.BaseURL,
		"api_key":    req.APIKey,
		"updated_at": time.Now(),
	}

	tx := db.
		Table("administracion.report_server_connections").
		Where("id = ? AND created_by = ?", id, userID).
		Updates(data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil
}

func DeleteConnection(db *gorm.DB, id int64, userID string) (int64, error) {

	tx := db.
		Table("administracion.report_server_connections").
		Where("id = ? AND created_by = ?", id, userID).
		Delete(nil)

	if tx.Error != nil {
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil
}

// SetActiveConnection marks the given connection as active and deactivates all others
// for this user. Uses a transaction to guarantee consistency.
func SetActiveConnection(db *gorm.DB, id int64, userID string) (int64, error) {

	err := db.Transaction(func(tx *gorm.DB) error {

		// Deactivate all connections for this user
		if err := tx.
			Table("administracion.report_server_connections").
			Where("created_by = ?", userID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Activate the selected one (must belong to this user)
		res := tx.
			Table("administracion.report_server_connections").
			Where("id = ? AND created_by = ?", id, userID).
			Update("is_active", true)

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
