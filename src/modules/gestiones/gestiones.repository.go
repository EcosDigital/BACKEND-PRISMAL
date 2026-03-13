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

// ─── Gestiones ────────────────────────────────────────────────────────────────
func CreateGestion(db *gorm.DB, ticketID int64, req *GestionRequest) (int64, error) {

	var newID int64

	err := db.Transaction(func(tx *gorm.DB) error {

		// 1. Leer el estado actual del ticket (para guardarlo como estado_anterior)
		var idEstadoActual int
		if err := tx.
			Table("gestiones.cfg_tickets_soporte").
			Select("id_estado").
			Where("id = ?", ticketID).
			Scan(&idEstadoActual).Error; err != nil {
			return err
		}

		// 2. Insertar la gestión en el historial
		gestion := map[string]interface{}{
			"id_ticket":          ticketID,
			"comentario":         req.Comentario,
			"id_estado_anterior": idEstadoActual,
			"id_estado_nuevo":    req.IDEstadoNuevo,
			"created_at":         time.Now(),
			"created_by":         req.UserID,
		}

		if err := tx.
			Table("gestiones.mov_ticket_historial").
			Create(&gestion).Error; err != nil {
			return err
		}

		if err := tx.Raw("SELECT lastval()").Scan(&newID).Error; err != nil {
			return err
		}

		// 3. Actualizar el estado del ticket solo si cambió
		if req.IDEstadoNuevo != idEstadoActual {
			if err := tx.
				Table("gestiones.cfg_tickets_soporte").
				Where("id = ?", ticketID).
				Updates(map[string]interface{}{
					"id_estado": req.IDEstadoNuevo,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func ListGestiones(db *gorm.DB, ticketID int64) ([]GestionResponse, error) {

	var results []GestionResponse

	err := db.
		Table("gestiones.mov_ticket_historial h").
		Select(`
			h.id,
			h.id_ticket,
			h.comentario,
			h.id_estado_anterior,
			ea.nombre    AS estado_anterior,
			ea.color_hex AS color_anterior,
			h.id_estado_nuevo,
			en.nombre    AS estado_nuevo,
			en.color_hex AS color_nuevo,
			h.created_by AS id_autor,
			CASE
				WHEN t.id_tipo_persona = 1 THEN
					TRIM(CONCAT_WS(' ', t.primer_nombre, t.segundo_nombre, t.primer_apellido, t.segundo_apellido))
				WHEN t.id_tipo_persona = 2 THEN
					t.razon_social
				ELSE ''	
			END AS autor,
			h.created_at
		`).
		Joins("LEFT JOIN gestiones.cfg_estados_ticket ea ON ea.id = h.id_estado_anterior").
		Joins("LEFT JOIN gestiones.cfg_estados_ticket en ON en.id = h.id_estado_nuevo").
		Joins("LEFT JOIN configuracion.cfg_terceros t ON t.id = h.created_by").
		Where("h.id_ticket = ?", ticketID).
		Order("h.created_at ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []GestionResponse{}
	}

	return results, nil
}
