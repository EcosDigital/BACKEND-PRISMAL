package search

import (
	"encoding/json"
	"fmt"

	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)

func ListDynamics(db *gorm.DB, data SearchDynamics) ([]interface{}, error) {

	query := `SELECT * FROM configuracion.qry_busquedas_dinamicas(
		operacion => $1,
		nombre_codigo => $2,
		id_tipo_registro => $3,
		fecha_inicial => $4,
		fecha_final => $5,
		id_estado_registro => $6,
		active => $7,
		aplicar_limit => $8
	)`

	// ── Preparar parámetros — nil cuando no vienen ────────────────
	var nombreCodigo interface{}
	if data.NombreCodigo != "" {
		nombreCodigo = data.NombreCodigo
	}

	var idTipoRegistro interface{}
	if data.Filters.IdTipoRegistro != nil {
		idTipoRegistro = *data.Filters.IdTipoRegistro
	}

	var fechaInicial interface{}
	if data.Filters.FechaInicial != nil {
		fechaInicial = *data.Filters.FechaInicial
	}

	var fechaFinal interface{}
	if data.Filters.FechaFinal != nil {
		fechaFinal = *data.Filters.FechaFinal
	}

	var idEstadoRegistro interface{}
	if data.Filters.IdEstadoRegistro != nil {
		idEstadoRegistro = *data.Filters.IdEstadoRegistro
	}

	var active interface{}
	if data.Filters.Active != nil {
		active = *data.Filters.Active
	}

	var aplicarLimit interface{}
	if data.Filters.AplicarLimit != nil {
		aplicarLimit = *data.Filters.AplicarLimit
	}

	// ── Ejecutar con GORM Raw ─────────────────────────────────────
	rawDB := db.Raw(
		query,
		data.Operacion,
		nombreCodigo,
		idTipoRegistro,
		fechaInicial,
		fechaFinal,
		idEstadoRegistro,
		active,
		aplicarLimit,
	)

	rows, err := rawDB.Rows()

	if err != nil {
		logging.Error.Printf("error ejecutando query: %v", err)
		return nil, fmt.Errorf("error ejecutando query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		logging.Error.Printf("error obteniendo columnas: %v", err)
		return nil, fmt.Errorf("error obteniendo columnas: %w", err)
	}

	var results []interface{}

	// ── Iterar filas ──────────────────────────────────────────────
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			logging.Error.Printf("error escaneando fila: %v", err)
			return nil, fmt.Errorf("error escaneando la fila: %w", err)
		}

		// ── Caso: una sola columna JSON (función PostgreSQL retorna JSON) ──
		if len(columns) == 1 {
			val := values[0]

			if b, ok := val.([]byte); ok {
				// Intentar como objeto con clave qry_busquedas_dinamicas
				var jsonWrapper map[string]interface{}
				if err := json.Unmarshal(b, &jsonWrapper); err == nil {
					if resultado, exists := jsonWrapper["qry_busquedas_dinamicas"]; exists && len(jsonWrapper) == 1 {
						if resultado == nil {
							return []interface{}{}, nil
						}
						if arr, ok := resultado.([]interface{}); ok {
							return arr, nil
						}
						return []interface{}{resultado}, nil
					}
				}

				// Intentar como array JSON
				var jsonData []map[string]interface{}
				if err := json.Unmarshal(b, &jsonData); err == nil {
					for _, item := range jsonData {
						results = append(results, item)
					}
					return results, nil
				}

				// Intentar como objeto JSON simple
				var jsonObj map[string]interface{}
				if err := json.Unmarshal(b, &jsonObj); err == nil {
					if resultado, exists := jsonObj["qry_busquedas_dinamicas"]; exists && len(jsonObj) == 1 {
						if resultado == nil {
							return []interface{}{}, nil
						}
						if arr, ok := resultado.([]interface{}); ok {
							return arr, nil
						}
						return []interface{}{resultado}, nil
					}
					results = append(results, jsonObj)
					return results, nil
				}
			}
		}

		// ── Caso general: múltiples columnas ──────────────────────
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		logging.Error.Printf("error iterando filas: %v", err)
		return nil, fmt.Errorf("error iterando filas: %w", err)
	}

	if len(results) == 0 {
		return []interface{}{}, nil
	}

	return results, nil
}
