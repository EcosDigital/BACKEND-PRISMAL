package chat

import (
	"encoding/json"

	"github.com/ecosistema/core/src/shared/integraciones/gemini"
	"gorm.io/gorm"
)

// registrar mensaje (de la persona o de la IA)
func CreateMensaje(db *gorm.DB, mensaje string, esIA bool, detalle map[string]interface{}, empresaID int64) error {

	detalleJSON, err := json.Marshal(detalle)
	if err != nil {
		return err
	}

	return db.Exec(`
		INSERT INTO agentes.mov_mensajes (mensaje, es_ia, detalle, created_at, id_empresa)
		VALUES (?, ?, ?::jsonb, NOW(), ?)`,
		mensaje, esIA, string(detalleJSON), empresaID).Error

}

// consultar los últimos mensajes de un chat (filtrado por detalle),
// devueltos del más viejo al más nuevo
func ListUltimosMensajes(db *gorm.DB, filtro map[string]interface{}, limite int) ([]gemini.Mensaje, error) {

	filtroJSON, err := json.Marshal(filtro)
	if err != nil {
		return nil, err
	}

	var rows []struct {
		Mensaje string
		EsIA    bool `gorm:"column:es_ia"`
	}

	err = db.Raw(`
		SELECT mensaje, es_ia
		FROM agentes.mov_mensajes
		WHERE detalle @> ?::jsonb
		ORDER BY created_at DESC
		LIMIT ?`,
		string(filtroJSON), limite).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	mensajes := make([]gemini.Mensaje, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		mensajes = append(mensajes, gemini.Mensaje{EsIA: rows[i].EsIA, Texto: rows[i].Mensaje})
	}

	return mensajes, nil

}
