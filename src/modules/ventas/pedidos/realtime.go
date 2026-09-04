// Tiempo real de pedidos: una conexión LISTEN permanente contra la base
// admin (canal "pedidos_nuevos", disparado por el trigger de 026_init.up.sql)
// y un Hub en memoria que reparte cada aviso a las conexiones SSE abiertas
// del tenant dueño del pedido.
package pedidos

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/lib/pq"
)

const canalPedidosNuevos = "pedidos_nuevos"

// Notification es el payload que manda pg_notify por el canal pedidos_nuevos.
type Notification struct {
	IDTenant int64 `json:"id_tenant"`
	IDPedido int64 `json:"id_pedido"`
}

// Hub reparte cada Notification a las conexiones SSE abiertas del tenant
// correspondiente. Si el tenant no tiene ninguna pantalla abierta en ese
// momento, el aviso simplemente no tiene a quién llegarle — no se pierde
// nada porque el pedido ya quedó guardado en la base.
type Hub struct {
	mu   sync.Mutex
	subs map[int64]map[chan Notification]struct{}
}

func NewHub() *Hub {
	return &Hub{subs: make(map[int64]map[chan Notification]struct{})}
}

// Subscribe registra un canal nuevo para el tenant dado. El caller debe
// invocar la función de limpieza retornada cuando cierre la conexión SSE.
func (h *Hub) Subscribe(idTenant int64) (chan Notification, func()) {
	ch := make(chan Notification, 8)

	h.mu.Lock()
	if h.subs[idTenant] == nil {
		h.subs[idTenant] = make(map[chan Notification]struct{})
	}
	h.subs[idTenant][ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		delete(h.subs[idTenant], ch)
		if len(h.subs[idTenant]) == 0 {
			delete(h.subs, idTenant)
		}
		h.mu.Unlock()
		close(ch)
	}

	return ch, unsubscribe
}

func (h *Hub) publish(n Notification) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for ch := range h.subs[n.IDTenant] {
		select {
		case ch <- n:
		default:
			// Suscriptor lento: no bloquea al resto, el aviso también
			// queda disponible al refrescar la pantalla por el endpoint normal.
		}
	}
}

// PedidosHub es el único Hub del proceso.
var PedidosHub = NewHub()

// StartListener abre la conexión LISTEN permanente contra la base admin.
// Debe llamarse una sola vez, solo en el proceso multi-tenant (main.go).
func StartListener() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
		core.Cfg.Db_name,
	)

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logging.Error.Printf("pedidos: evento de LISTEN: %v", err)
		}
	})

	if err := listener.Listen(canalPedidosNuevos); err != nil {
		logging.Error.Printf("pedidos: no se pudo escuchar el canal %s: %v", canalPedidosNuevos, err)
		return
	}

	logging.Info.Printf("pedidos: escuchando avisos de pedidos nuevos (canal %s) 🔔", canalPedidosNuevos)

	go func() {
		for n := range listener.Notify {
			if n == nil {
				continue
			}

			var payload Notification
			if err := json.Unmarshal([]byte(n.Extra), &payload); err != nil {
				logging.Error.Printf("pedidos: aviso con payload inválido: %v", err)
				continue
			}

			PedidosHub.publish(payload)
		}
	}()
}
