package ordenes

import (
	"fmt"

	"github.com/ecosistema/core/src/modules/_core/notificaciones"
	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)

// codigo de "Listo para servir" en ventas.ref_estado_orden y el evento de
// notificacion que dispara (configuracion.cfg_eventos_notificacion)
const (
	codigoEstadoListo = "003"
	eventoOrdenLista  = "EV-001"
)

func RegisterOrden(db *gorm.DB, req *OrdenRequest) (int64, error) {
	return CreateOrden(db, req)
}

func FilterOrdenes(db *gorm.DB) ([]OrdenListItem, error) {
	return ListOrdenes(db)
}

// ordenes filtradas por estado (ej. las que estan listas para servir)
func FilterOrdenesByEstado(db *gorm.DB, codigoEstado string) ([]OrdenListItem, error) {
	return ListOrdenesByEstado(db, codigoEstado)
}

func FilterOrdenDetalle(db *gorm.DB, id int64) (*OrdenDetalle, error) {
	return GetOrdenDetalle(db, id)
}

// cambiar el estado de la orden. si entra a "Listo para servir" se avisa a
// los configurados en su sede; el aviso nunca frena el cambio de estado
func CambiarEstadoOrden(db *gorm.DB, id int64, req *CambiarEstadoOrdenRequest) error {

	//el estado de antes se lee primero: solo se avisa cuando la orden entra
	//a "Listo para servir", no si ya estaba ahi. si esta lectura falla, el
	//cambio de estado se hace igual y solo se pierde el aviso
	ordenes, err := GetOrdenAviso(db, id)
	if err != nil {
		logging.Error.Printf("error consultando la orden %d para su aviso: %v", id, err)
	}

	if err := UpdateEstadoOrden(db, id, req.CodigoEstado, req.MotivoRechazo); err != nil {
		return err
	}

	if len(ordenes) > 0 &&
		req.CodigoEstado == codigoEstadoListo &&
		ordenes[0].EstadoCodigo != codigoEstadoListo {
		avisarOrdenLista(db, id, ordenes[0], req.SedeID)
	}

	return nil
}

// avisar en segundo plano que la orden esta lista: el mesero no espera el
// envio del push y, si algo falla, la orden ya quedo actualizada
func avisarOrdenLista(db *gorm.DB, idOrden int64, orden ordenAviso, idSedeUsuario int64) {

	//la sede de la orden; si no tiene, la de quien la marco como lista
	idSede := idSedeUsuario
	if orden.IDSede != nil && *orden.IDSede > 0 {
		idSede = *orden.IDSede
	}

	titulo := "Orden lista para servir"
	mensaje := fmt.Sprintf("La orden #%d (%s) está lista para servir", idOrden, orden.IdentificadorMesa)

	go notificaciones.NotificarEvento(db, eventoOrdenLista, idSede, orden.IDEmpresa, titulo, mensaje)
}
