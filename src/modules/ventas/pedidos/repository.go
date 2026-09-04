package pedidos

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ecosistema/core/src/database"
)

// CodigoRechazado es el código de ref_estado_pedido para "Rechazado" —
// exige motivo_rechazo al pasar un pedido a este estado.
const CodigoRechazado = "008"

func ListPedidos(idTenant int64) ([]PedidoListItem, error) {

	results := make([]PedidoListItem, 0)

	err := database.GormDB.Raw(`
		SELECT
			p.id,
			p.id_referencia_millave,
			e.codigo                         AS estado_codigo,
			e.nombre                         AS estado_nombre,
			COALESCE(ent.cliente_nombre, '') AS cliente_nombre,
			p.valor_total,
			TO_CHAR(p.created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM integraciones.mov_pedidos p
		JOIN integraciones.ref_estado_pedido e ON e.id = p.id_estado
		LEFT JOIN integraciones.mov_pedidos_entrega ent ON ent.id_pedido = p.id
		WHERE p.id_tenant = ?
		ORDER BY p.created_at DESC
	`, idTenant).Scan(&results).Error

	return results, err
}

func GetPedidoDetalle(idTenant, idPedido int64) (*PedidoDetalle, error) {

	var cabecera PedidoDetalle

	err := database.GormDB.Raw(`
		SELECT
			p.id,
			p.id_referencia_millave,
			e.codigo    AS estado_codigo,
			e.nombre    AS estado_nombre,
			e.es_final  AS estado_es_final,
			p.valor_subtotal,
			p.valor_domicilio,
			p.valor_total,
			p.observaciones,
			p.motivo_rechazo,
			TO_CHAR(p.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
			TO_CHAR(p.updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at
		FROM integraciones.mov_pedidos p
		JOIN integraciones.ref_estado_pedido e ON e.id = p.id_estado
		WHERE p.id = ? AND p.id_tenant = ?
	`, idPedido, idTenant).Scan(&cabecera).Error

	if err != nil {
		return nil, err
	}
	if cabecera.ID == 0 {
		return nil, errors.New("pedido no encontrado")
	}

	detalle := make([]PedidoDetalleItem, 0)
	if err := database.GormDB.Raw(`
		SELECT id, codigo, nombre, cantidad, precio_unitario, subtotal, observacion
		FROM integraciones.mov_pedidos_detalle
		WHERE id_pedido = ?
		ORDER BY id
	`, idPedido).Scan(&detalle).Error; err != nil {
		return nil, err
	}
	cabecera.Detalle = detalle

	var entrega PedidoEntrega
	if err := database.GormDB.Raw(`
		SELECT cliente_nombre, cliente_telefono, tipo_documento, cliente_documento,
			direccion_entrega, geolocalizacion_lat, geolocalizacion_lon, referencia
		FROM integraciones.mov_pedidos_entrega
		WHERE id_pedido = ?
	`, idPedido).Scan(&entrega).Error; err != nil {
		return nil, err
	}
	cabecera.Entrega = entrega

	return &cabecera, nil
}

func UpdateEstadoPedido(idTenant, idPedido int64, codigo, motivo string) error {

	var actual struct {
		IDEstado int64
		EsFinal  bool
	}
	err := database.GormDB.Raw(`
		SELECT p.id_estado, e.es_final
		FROM integraciones.mov_pedidos p
		JOIN integraciones.ref_estado_pedido e ON e.id = p.id_estado
		WHERE p.id = ? AND p.id_tenant = ?
	`, idPedido, idTenant).Scan(&actual).Error
	if err != nil {
		return err
	}
	if actual.IDEstado == 0 {
		return errors.New("pedido no encontrado")
	}
	if actual.EsFinal {
		return errors.New("el pedido ya está en un estado final y no se puede cambiar")
	}

	if codigo == CodigoRechazado && strings.TrimSpace(motivo) == "" {
		return errors.New("motivo_rechazo es obligatorio al rechazar un pedido")
	}

	var nuevoEstado struct {
		ID int64
	}
	err = database.GormDB.Raw(`
		SELECT id FROM integraciones.ref_estado_pedido WHERE codigo = ?
	`, codigo).Scan(&nuevoEstado).Error
	if err != nil {
		return err
	}
	if nuevoEstado.ID == 0 {
		return fmt.Errorf("código de estado '%s' no existe", codigo)
	}

	return database.GormDB.Exec(`
		UPDATE integraciones.mov_pedidos
		SET id_estado = ?, motivo_rechazo = COALESCE(NULLIF(?, ''), motivo_rechazo), updated_at = NOW()
		WHERE id = ? AND id_tenant = ?
	`, nuevoEstado.ID, motivo, idPedido, idTenant).Error
}
