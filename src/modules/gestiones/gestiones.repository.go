package gestiones

import (
	"time"

	"gorm.io/gorm"
)

// CreateTicket inserta un nuevo ticket en la base de datos
func CreateTicket(db *gorm.DB, req *TicketRequest) (int64, error) {

	data := map[string]interface{}{
		"titulo":            req.Titulo,
		"descripcion":       req.Descripcion,
		"id_nivel":          req.IDNivel,
		"id_estado":         1,
		"id_origen":         req.IDOrigen,
		"id_tenant":         nil,
		"fecha_entrega":     nil,
		"contacto_nombre":   req.ContactoNombre,
		"contacto_email":    req.ContactoEmail,
		"contacto_telefono": req.ContactoTelefono,
		"created_at":        time.Now(),
		"created_by":        req.UserID,
	}

	tx := db.
		Table("gestiones.cfg_tickets_soporte").
		Create(&data)

	if tx.Error != nil {
		return 0, tx.Error
	}

	var id int64
	if err := db.Raw("SELECT lastval()").Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

// ListTicketByID busca un ticket por su ID con JOINs a nivel, estado y origen
func ListTicketByID(db *gorm.DB, id int64) ([]TicketResponse, error) {

	var results []TicketResponse

	err := db.
		Table("gestiones.cfg_tickets_soporte t").
		Select(`
			t.id,
			t.titulo,
			t.descripcion,
			t.id_nivel,
			n.nombre        AS nivel,
			n.color_hex     AS color_nivel,
			t.id_estado,
			e.nombre        AS estado,
			e.color_hex     AS color_estado,
			t.id_origen,
			o.nombre        AS origen,
			t.id_tenant,
			t.fecha_entrega,
			t.contacto_nombre,
			t.contacto_email,
			t.contacto_telefono,
			t.created_at
		`).
		Joins("INNER JOIN gestiones.cfg_niveles_caso n ON n.id = t.id_nivel").
		Joins("INNER JOIN gestiones.cfg_estados_ticket e ON e.id = t.id_estado").
		Joins("INNER JOIN gestiones.ref_origenes_ticket o ON o.id = t.id_origen").
		Where("t.id = ?", id).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func ListTickets(db *gorm.DB, f *TicketFiltros, userID int64) ([]TicketResponse, error) {

	var results []TicketResponse

	query := db.
		Table("gestiones.cfg_tickets_soporte t").
		Select(`
            t.id,
            t.titulo,
            t.descripcion,
            t.id_nivel,
            n.nombre    AS nivel,
            n.color_hex AS color_nivel,
            t.id_estado,
            e.nombre    AS estado,
            e.color_hex AS color_estado,
            t.id_origen,
            o.nombre    AS origen,
            t.id_tenant,
            t.fecha_entrega,
            t.contacto_nombre,
            t.contacto_email,
            t.contacto_telefono,
            t.created_at
        `).
		Joins("INNER JOIN gestiones.cfg_niveles_caso n ON n.id = t.id_nivel").
		Joins("INNER JOIN gestiones.cfg_estados_ticket e ON e.id = t.id_estado").
		Joins("INNER JOIN gestiones.ref_origenes_ticket o ON o.id = t.id_origen").
		Order("t.id DESC")

	// Filtros opcionales
	if f.IDEstado != 0 {
		query = query.Where("t.id_estado = ?", f.IDEstado)
	}
	if f.IDNivel != 0 {
		query = query.Where("t.id_nivel = ?", f.IDNivel)
	}
	if f.IDOrigen != 0 {
		query = query.Where("t.id_origen = ?", f.IDOrigen)
	}
	if f.IDColaborador != 0 {
		query = query.
			Joins("INNER JOIN gestiones.cfg_ticket_colaboradores tc ON tc.id_ticket = t.id").
			Where("tc.id_colaborador = ?", f.IDColaborador)
	}
	if f.Numero != "" {
		query = query.Where("CAST(t.id AS TEXT) = ?", f.Numero)
	}
	if f.SoloAsignados {
		// Filtra por el userID extraído del JWT en el controller
		query = query.
			Joins("INNER JOIN gestiones.cfg_ticket_colaboradores tca ON tca.id_ticket = t.id").
			Where("tca.id_colaborador = ?", userID)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	if results == nil {
		results = []TicketResponse{}
	}

	return results, nil
}
