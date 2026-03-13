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

// ─── Gestiones ────────────────────────────────────────────────────────────────
func RegisterGestion(db *gorm.DB, ticketID int64, req *GestionRequest) (int64, error) {

	// Verificar que el ticket exista
	var count int64
	if err := db.
		Table("gestiones.cfg_tickets_soporte").
		Where("id = ?", ticketID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("el ticket no existe")
	}

	return CreateGestion(db, ticketID, req)
}

func FilterGestiones(db *gorm.DB, ticketID int64) ([]GestionResponse, error) {

	// Verificar que el ticket exista
	var count int64
	if err := db.
		Table("gestiones.cfg_tickets_soporte").
		Where("id = ?", ticketID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("el ticket no existe")
	}

	return ListGestiones(db, ticketID)
}
