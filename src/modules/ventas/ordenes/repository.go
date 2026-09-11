package ordenes

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// idEstadoMesaOcupada es el id fijo de "Ocupada" en
// configuracion.ref_estado_mesa, sembrado por 027_init.up.sql
// (1=Disponible, 2=Ocupada).
const idEstadoMesaOcupada = 2

// CreateOrden registra la cabecera y las líneas de la orden en una sola
// transacción y, si la orden va contra una mesa real (no ficticia), marca
// esa mesa como "Ocupada".
func CreateOrden(db *gorm.DB, req *OrdenRequest) (int64, error) {

	var idOrden int64

	err := db.Transaction(func(tx *gorm.DB) error {

		var idEstadoSolicitado int64
		if err := tx.Raw(`SELECT id FROM ventas.ref_estado_orden WHERE codigo = '001'`).Scan(&idEstadoSolicitado).Error; err != nil {
			return err
		}

		var subtotal float64
		for _, linea := range req.Detalle {
			subtotal += float64(linea.Cantidad) * linea.PrecioUnitario
		}

		cabecera := map[string]interface{}{
			"id_estado":          idEstadoSolicitado,
			"id_mesa":            req.IDMesa,
			"identificador_mesa": req.IdentificadorMesa,
			"valor_subtotal":     subtotal,
			"valor_total":        subtotal,
			"observaciones":      req.Observaciones,
			"created_by":         req.UserID,
			"id_empresa":         req.EmpresaID,
			"id_sede":            req.SedeID,
		}

		if err := tx.Table("ventas.mov_ordenes").Create(&cabecera).Error; err != nil {
			return err
		}

		id, ok := cabecera["id"].(int64)
		if !ok {
			return errors.New("no se pudo obtener el id de la orden creada")
		}
		idOrden = id

		for _, linea := range req.Detalle {
			detalle := map[string]interface{}{
				"id_orden":           id,
				"id_articulo_origen": linea.IDArticuloOrigen,
				"codigo":             linea.Codigo,
				"nombre":             linea.Nombre,
				"cantidad":           linea.Cantidad,
				"precio_unitario":    linea.PrecioUnitario,
				"observacion":        linea.Observacion,
			}
			if err := tx.Table("ventas.mov_ordenes_detalle").Create(&detalle).Error; err != nil {
				return err
			}
		}

		// La mesa pasa a "Ocupada" solo cuando es una mesa real; una mesa
		// ficticia (id_mesa nulo) no tiene fila en cfg_mesas que actualizar.
		if req.IDMesa != nil {
			if err := tx.Table("configuracion.cfg_mesas").
				Where("id = ?", *req.IDMesa).
				Update("id_estado", idEstadoMesaOcupada).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return idOrden, nil
}

func ListOrdenes(db *gorm.DB) ([]OrdenListItem, error) {

	results := make([]OrdenListItem, 0)

	err := db.Raw(`
		SELECT
			o.id,
			o.identificador_mesa,
			e.codigo AS estado_codigo,
			e.nombre AS estado_nombre,
			o.valor_total,
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM ventas.mov_ordenes o
		JOIN ventas.ref_estado_orden e ON e.id = o.id_estado
		ORDER BY o.created_at DESC
	`).Scan(&results).Error

	return results, err
}

func GetOrdenDetalle(db *gorm.DB, idOrden int64) (*OrdenDetalle, error) {

	var cabecera OrdenDetalle

	err := db.Raw(`
		SELECT
			o.id,
			o.identificador_mesa,
			e.codigo   AS estado_codigo,
			e.nombre   AS estado_nombre,
			e.es_final AS estado_es_final,
			o.valor_subtotal,
			o.valor_total,
			o.observaciones,
			o.motivo_rechazo,
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
			TO_CHAR(o.update_at, 'YYYY-MM-DD HH24:MI') AS updated_at
		FROM ventas.mov_ordenes o
		JOIN ventas.ref_estado_orden e ON e.id = o.id_estado
		WHERE o.id = ?
	`, idOrden).Scan(&cabecera).Error

	if err != nil {
		return nil, err
	}
	if cabecera.ID == 0 {
		return nil, errors.New("orden no encontrada")
	}

	detalle := make([]OrdenDetalleItem, 0)
	if err := db.Raw(`
		SELECT id, codigo, nombre, cantidad, precio_unitario, subtotal, observacion
		FROM ventas.mov_ordenes_detalle
		WHERE id_orden = ?
		ORDER BY id
	`, idOrden).Scan(&detalle).Error; err != nil {
		return nil, err
	}
	cabecera.Detalle = detalle

	return &cabecera, nil
}

func UpdateEstadoOrden(db *gorm.DB, idOrden int64, codigo, motivo string) error {

	var actual struct {
		IDEstado int64
		EsFinal  bool
	}
	err := db.Raw(`
		SELECT o.id_estado, e.es_final
		FROM ventas.mov_ordenes o
		JOIN ventas.ref_estado_orden e ON e.id = o.id_estado
		WHERE o.id = ?
	`, idOrden).Scan(&actual).Error
	if err != nil {
		return err
	}
	if actual.IDEstado == 0 {
		return errors.New("orden no encontrada")
	}
	if actual.EsFinal {
		return errors.New("la orden ya está en un estado final y no se puede cambiar")
	}

	var nuevoEstado struct {
		ID int64
	}
	err = db.Raw(`SELECT id FROM ventas.ref_estado_orden WHERE codigo = ?`, codigo).Scan(&nuevoEstado).Error
	if err != nil {
		return err
	}
	if nuevoEstado.ID == 0 {
		return fmt.Errorf("código de estado '%s' no existe", codigo)
	}

	return db.Exec(`
		UPDATE ventas.mov_ordenes
		SET id_estado = ?, motivo_rechazo = COALESCE(NULLIF(?, ''), motivo_rechazo), update_at = NOW()
		WHERE id = ?
	`, nuevoEstado.ID, motivo, idOrden).Error
}
