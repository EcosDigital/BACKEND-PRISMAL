// Tiempo real de órdenes: a diferencia de "pedidos" (una sola BD admin, un
// solo listener que demuxa por id_tenant), cada tenant tiene su propia base
// física (prismar_<tenant>) y Postgres LISTEN/NOTIFY no cruza bases de
// datos — hace falta un listener independiente por tenant, abierto
// perezosamente al primer suscriptor SSE de ese tenant.
package ordenes

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/lib/pq"
)

const canalOrdenesNuevas = "ordenes_nuevas"

// Notification es el payload que manda pg_notify por el canal ordenes_nuevas.
type Notification struct {
	IDOrden int64 `json:"id_orden"`
}

// Hub reparte cada Notification a las conexiones SSE abiertas de UN tenant.
// A diferencia de pedidos.Hub no hace falta demuxar por tenant: cada Hub
// aquí ya está atado a la conexión LISTEN de un solo tenant.
type Hub struct {
	mu   sync.Mutex
	subs map[chan Notification]struct{}
}

func newHub() *Hub {
	return &Hub{subs: make(map[chan Notification]struct{})}
}

func (h *Hub) subscribe() (chan Notification, func()) {
	ch := make(chan Notification, 8)

	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()

	unsubscribe := func() {
		h.mu.Lock()
		delete(h.subs, ch)
		h.mu.Unlock()
		close(ch)
	}

	return ch, unsubscribe
}

func (h *Hub) publish(n Notification) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for ch := range h.subs {
		select {
		case ch <- n:
		default:
			// Suscriptor lento: no bloquea al resto, la orden también
			// queda disponible al refrescar la pantalla por el endpoint normal.
		}
	}
}

// tenantHubs guarda, por slug de tenant, el Hub ya creado — y garantiza que
// el listener de ese tenant se abra una sola vez.
var (
	tenantHubsMu sync.Mutex
	tenantHubs   = make(map[string]*Hub)
)

// SubscribeOrdenes suscribe una conexión SSE a los avisos de órdenes nuevas
// del tenant dado, abriendo su listener si todavía no existe.
func SubscribeOrdenes(tenantSlug string) (chan Notification, func()) {
	tenantHubsMu.Lock()
	hub, exists := tenantHubs[tenantSlug]
	if !exists {
		hub = newHub()
		tenantHubs[tenantSlug] = hub
		go startTenantListener(tenantSlug, hub)
	}
	tenantHubsMu.Unlock()

	return hub.subscribe()
}

// startTenantListener abre una conexión LISTEN permanente contra la base
// propia del tenant (prismar_<tenant>). Se llama una sola vez por tenant,
// al primer suscriptor SSE. Simplificación consciente de esta primera
// versión: el listener no se cierra por inactividad.
func startTenantListener(tenantSlug string, hub *Hub) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=prismar_%s sslmode=disable",
		core.Cfg.Db_host,
		core.Cfg.Db_port,
		core.Cfg.Db_user,
		core.Cfg.Db_pass,
		tenantSlug,
	)

	listener := pq.NewListener(dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logging.Error.Printf("ordenes (%s): evento de LISTEN: %v", tenantSlug, err)
		}
	})

	if err := listener.Listen(canalOrdenesNuevas); err != nil {
		logging.Error.Printf("ordenes (%s): no se pudo escuchar el canal %s: %v", tenantSlug, canalOrdenesNuevas, err)
		return
	}

	logging.Info.Printf("ordenes (%s): escuchando avisos de órdenes nuevas (canal %s) 🔔", tenantSlug, canalOrdenesNuevas)

	for n := range listener.Notify {
		if n == nil {
			continue
		}

		var payload Notification
		if err := json.Unmarshal([]byte(n.Extra), &payload); err != nil {
			logging.Error.Printf("ordenes (%s): aviso con payload inválido: %v", tenantSlug, err)
			continue
		}

		hub.publish(payload)
	}
}
