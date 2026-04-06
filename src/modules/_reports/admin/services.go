package reports

import (
	"errors"

	"gorm.io/gorm"
)

func RegisterConnection(db *gorm.DB, req *ReportServerConnectionRequest, userID string) (int64, error) {
	return CreateConnection(db, req, userID)
}

func FilterConnections(db *gorm.DB, userID string) ([]ReportServerConnectionResponse, error) {
	return ListConnections(db, userID)
}

func FilterConnectionByID(db *gorm.DB, id int64, userID string) (*ReportServerConnectionResponseFull, error) {
	return ListConnectionByID(db, id, userID)
}

func EditConnection(db *gorm.DB, id int64, req *ReportServerConnectionUpdateRequest, userID string) (int64, error) {

	exists, err := ListConnectionByID(db, id, userID)
	if err != nil {
		return 0, err
	}
	if exists == nil || exists.ID == 0 {
		return 0, errors.New("no se encontró conexión con este ID")
	}

	return UpdateConnection(db, id, req, userID)
}

func RemoveConnection(db *gorm.DB, id int64, userID string) (int64, error) {

	exists, err := ListConnectionByID(db, id, userID)
	if err != nil {
		return 0, err
	}
	if exists == nil || exists.ID == 0 {
		return 0, errors.New("no se encontró conexión con este ID")
	}

	return DeleteConnection(db, id, userID)
}

func ActivateConnection(db *gorm.DB, id int64, userID string) (int64, error) {

	exists, err := ListConnectionByID(db, id, userID)
	if err != nil {
		return 0, err
	}
	if exists == nil || exists.ID == 0 {
		return 0, errors.New("no se encontró conexión con este ID")
	}

	return SetActiveConnection(db, id, userID)
}
