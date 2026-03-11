package gestiones

import (
	"errors"

	"gorm.io/gorm"
)

// RegisterTicket valida y crea un nuevo ticket
func RegisterTicket(db *gorm.DB, req *TicketRequest) (int64, error) {

	newID, err := CreateTicket(db, req)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func FilterTickets(db *gorm.DB, f *TicketFiltros, userID int64) ([]TicketResponse, error) {
	return ListTickets(db, f, userID)
}

// FilterTicketByID busca un ticket por ID y valida que exista
func FilterTicketByID(db *gorm.DB, id int64) ([]TicketResponse, error) {

	results, err := ListTicketByID(db, id)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no se encontró el ticket solicitado")
	}

	return results, nil
}
