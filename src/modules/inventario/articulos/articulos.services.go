package articulos

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func RegisterArticulo(db *gorm.DB, req *ArticuloRequest) (int64, error) {
	return CreateArticulo(db, req)
}

func FilterArticulos(db *gorm.DB, empresaID int64) ([]ArticuloResponse, error) {
	return ListArticulos(db, empresaID)
}

func FilterArticuloByID(db *gorm.DB, id int64) (*ArticuloResponseFull, error) {
	return ListArticuloByID(db, id)
}

func EditArticulo(db *gorm.DB, id int64, req *ArticuloUpdateRequest) (int64, error) {

	exists, err := ListArticuloByID(db, id)
	if err != nil {
		return 0, err
	}
	if exists == nil || exists.ID == 0 {
		return 0, errors.New("no se encontró artículo con este ID")
	}

	return UpdateArticulo(db, req, id)
}

func FilterSearchArticulos(db *gorm.DB, nombreCodigo string, empresaID int64) ([]ArticuloResponse, error) {
	return SearchArticulos(db, nombreCodigo, empresaID)
}

// ─── Carga masiva ─────────────────────────────────────────────────────────────

// ImportArticulos procesa el plano de carga masiva:
//  1. Carga el cache de referencias (1 query por tabla de referencia).
//  2. Valida y hace upsert fila a fila dentro de una transacción.
//  3. Retorna el resumen agregado y el detalle por fila.
//
// Si una fila falla NO se hace rollback de las demás — cada fila es independiente.
// Esto permite que el usuario corrija solo las filas con error.
func ImportArticulos(
	db *gorm.DB,
	req *ImportArticulosRequest,
	userID int64,
	empresaID int64,
	sedeID int64,
) (*ImportArticulosResponse, error) {

	// ── 1. Validaciones previas ───────────────────────────────────────────────

	if len(req.Filas) == 0 {
		return nil, errors.New("el plano no contiene filas")
	}

	if len(req.Filas) > 5000 {
		return nil, fmt.Errorf("el plano excede el límite de 5000 filas (recibidas: %d)", len(req.Filas))
	}

	// ── 2. Cargar cache de referencias una sola vez ───────────────────────────

	cache, err := LoadRefCache(db, empresaID)
	if err != nil {
		return nil, fmt.Errorf("error cargando referencias: %w", err)
	}

	// ── 3. Procesar fila a fila ───────────────────────────────────────────────

	resp := &ImportArticulosResponse{
		Detalle: make([]ImportRowResult, 0, len(req.Filas)),
	}

	for i, fila := range req.Filas {
		filaNum := i + 2 // fila 1 = encabezado en el Excel → datos desde fila 2

		// Validaciones básicas antes de tocar la BD
		if validErr := validatePlanoRow(&fila, filaNum); validErr != nil {
			resp.Errores++
			resp.Detalle = append(resp.Detalle, ImportRowResult{
				Fila:   filaNum,
				Codigo: fila.Codigo,
				Accion: "error",
				Error:  validErr.Error(),
			})
			continue
		}

		// Upsert
		accion, upsertErr := BulkUpsertArticulo(db, &fila, cache, userID, empresaID, sedeID)
		if upsertErr != nil {
			resp.Errores++
			resp.Detalle = append(resp.Detalle, ImportRowResult{
				Fila:   filaNum,
				Codigo: fila.Codigo,
				Accion: "error",
				Error:  upsertErr.Error(),
			})
			continue
		}

		if accion == "creado" {
			resp.Creados++
		} else {
			resp.Actualizados++
		}

		resp.Detalle = append(resp.Detalle, ImportRowResult{
			Fila:   filaNum,
			Codigo: fila.Codigo,
			Accion: accion,
		})
	}

	return resp, nil
}

// validatePlanoRow realiza validaciones de negocio sobre una fila del plano
// antes de intentar el upsert en la BD.
func validatePlanoRow(fila *PlanoArticuloRow, filaNum int) error {
	var msgs []string

	if strings.TrimSpace(fila.Codigo) == "" {
		msgs = append(msgs, "CODIGO es obligatorio")
	}

	if len(strings.TrimSpace(fila.Nombre)) < 3 {
		msgs = append(msgs, "NOMBRE debe tener al menos 3 caracteres")
	}

	if strings.TrimSpace(fila.IDTipoArticuloNombre) == "" {
		msgs = append(msgs, "ID_TIPO_ARTICULO es obligatorio")
	}

	if strings.TrimSpace(fila.IDUnidadMediaNombre) == "" {
		msgs = append(msgs, "ID_UNIDAD_MEDIA es obligatoria")
	}

	estadoNorm := strings.ToLower(strings.TrimSpace(fila.Estado))
	if estadoNorm != "activo" && estadoNorm != "inactivo" &&
		estadoNorm != "true" && estadoNorm != "false" && estadoNorm != "1" && estadoNorm != "0" {
		msgs = append(msgs, fmt.Sprintf("ESTADO '%s' no válido (use Activo o Inactivo)", fila.Estado))
	}

	if len(msgs) > 0 {
		return fmt.Errorf("fila %d — %s", filaNum, strings.Join(msgs, "; "))
	}

	return nil
}
