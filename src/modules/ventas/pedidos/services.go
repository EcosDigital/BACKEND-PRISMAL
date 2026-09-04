package pedidos

// GetPedidos devuelve el listado de pedidos del tenant.
func GetPedidos(idTenant int64) ([]PedidoListItem, error) {
	return ListPedidos(idTenant)
}

// GetPedido devuelve el detalle completo de un pedido, validando que
// pertenezca al tenant.
func GetPedido(idTenant, idPedido int64) (*PedidoDetalle, error) {
	return GetPedidoDetalle(idTenant, idPedido)
}

// CambiarEstadoPedido valida y aplica el cambio de estado de un pedido.
func CambiarEstadoPedido(idTenant, idPedido int64, req *CambiarEstadoRequest) error {
	return UpdateEstadoPedido(idTenant, idPedido, req.CodigoEstado, req.MotivoRechazo)
}
