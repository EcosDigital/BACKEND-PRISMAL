package pedidos

import (
	"bufio"
	"fmt"
	"strconv"
	"time"

	"github.com/ecosistema/core/src/shared/integraciones/millave"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// resolverTenantActual obtiene el id_tenant de la petición actual, o responde
// el error correspondiente si el proceso no está en modo multi-tenant.
func resolverTenantActual(c *fiber.Ctx) (int64, bool) {
	tenantSlug := middlewares.GetTenantSlug(c)
	if tenantSlug == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pedidos de domicilio requiere modo multi-tenant",
		})
		return 0, false
	}

	idTenant, _, err := millave.ResolverTenant(tenantSlug)
	if err != nil {
		logging.Error.Printf("pedidos: %v", err)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		return 0, false
	}

	return idTenant, true
}

// FindPedidosController GET /ventas/pedidos
func FindPedidosController(c *fiber.Ctx) error {
	idTenant, ok := resolverTenantActual(c)
	if !ok {
		return nil
	}

	results, err := GetPedidos(idTenant)
	if err != nil {
		logging.Error.Printf("FindPedidos error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// FindPedidoByIDController GET /ventas/pedidos/:id
func FindPedidoByIDController(c *fiber.Ctx) error {
	idTenant, ok := resolverTenantActual(c)
	if !ok {
		return nil
	}

	idPedido, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	result, err := GetPedido(idTenant, idPedido)
	if err != nil {
		if err.Error() == "pedido no encontrado" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		logging.Error.Printf("FindPedidoByID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// ChangeEstadoPedidoController PUT /ventas/pedidos/:id/estado
func ChangeEstadoPedidoController(c *fiber.Ctx) error {
	idTenant, ok := resolverTenantActual(c)
	if !ok {
		return nil
	}

	idPedido, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CambiarEstadoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	if err := CambiarEstadoPedido(idTenant, idPedido, &req); err != nil {
		if err.Error() == "pedido no encontrado" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "estado actualizado"})
}

// StreamPedidosController GET /ventas/pedidos/stream
// Canal SSE: mientras la pantalla esté abierta, avisa en tiempo real cuando
// llega un pedido nuevo de este tenant. Si no hay nadie escuchando, el
// pedido no se pierde — queda disponible por FindPedidosController.
func StreamPedidosController(c *fiber.Ctx) error {
	idTenant, ok := resolverTenantActual(c)
	if !ok {
		return nil
	}

	ch, unsubscribe := PedidosHub.Subscribe(idTenant)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer unsubscribe()

		heartbeat := time.NewTicker(25 * time.Second)
		defer heartbeat.Stop()

		for {
			select {
			case n, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "event: pedido_nuevo\ndata: {\"id_pedido\":%d}\n\n", n.IDPedido)
				if err := w.Flush(); err != nil {
					return
				}
			case <-heartbeat.C:
				fmt.Fprint(w, ": ping\n\n")
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	}))

	return nil
}
