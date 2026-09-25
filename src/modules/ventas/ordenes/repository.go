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

		if err := tx.Raw(`
			INSERT INTO ventas.mov_ordenes
				(id_estado, id_mesa, identificador_mesa, valor_subtotal, valor_total, observaciones, created_by, id_empresa, id_sede)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			RETURNING id
		`, idEstadoSolicitado, req.IDMesa, req.IdentificadorMesa, subtotal, subtotal, req.Observaciones, req.UserID, req.EmpresaID, req.SedeID).
			Scan(&idOrden).Error; err != nil {
			return err
		}

		if idOrden == 0 {
			return errors.New("no se pudo obtener el id de la orden creada")
		}

		for _, linea := range req.Detalle {
			if err := tx.Exec(`
				INSERT INTO ventas.mov_ordenes_detalle
					(id_orden, id_articulo_origen, codigo, nombre, cantidad, precio_unitario, observacion)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, idOrden, linea.IDArticuloOrigen, linea.Codigo, linea.Nombre, linea.Cantidad, linea.PrecioUnitario, linea.Observacion).Error; err != nil {
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

// consultar las ordenes que estan en un estado, de la mas reciente a la mas
// antigua. consulta nueva: va con el ORM, no con SQL directo como las de arriba
func ListOrdenesByEstado(db *gorm.DB, codigoEstado string) ([]OrdenListItem, error) {

	results := make([]OrdenListItem, 0)

	err := db.
		Table("ventas.mov_ordenes o").
		Select(`
			o.id,
			o.identificador_mesa,
			e.codigo AS estado_codigo,
			e.nombre AS estado_nombre,
			o.valor_total,
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI') AS created_at`).
		Joins("INNER JOIN ventas.ref_estado_orden e ON e.id = o.id_estado").
		Where("e.codigo = ?", codigoEstado).
		Order("o.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil

}

// consultar el estado actual de la orden y los datos para su aviso
func GetOrdenAviso(db *gorm.DB, idOrden int64) ([]ordenAviso, error) {

	var results []ordenAviso

	err := db.
		Table("ventas.mov_ordenes o").
		Select("e.codigo as estado_codigo, o.identificador_mesa, o.id_empresa, o.id_sede").
		Joins("INNER JOIN ventas.ref_estado_orden e ON e.id = o.id_estado").
		Where("o.id = ?", idOrden).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if results == nil {
		results = []ordenAviso{}
	}

	return results, nil

}
