package ordenes

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// CreateOrdenController POST /ventas/ordenes
func CreateOrdenController(c *fiber.Ctx) error {
	var req OrdenRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	newID, err := RegisterOrden(db, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Orden registrada...",
		"id":      newID,
	})
}

// FindOrdenesController GET /ventas/ordenes
func FindOrdenesController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterOrdenes(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": results})
}

// FindOrdenByIDController GET /ventas/ordenes/:id
func FindOrdenByIDController(c *fiber.Ctx) error {
	idOrden, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := FilterOrdenDetalle(db, idOrden)
	if err != nil {
		if err.Error() == "orden no encontrada" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// ChangeEstadoOrdenController PUT /ventas/ordenes/:id/estado
func ChangeEstadoOrdenController(c *fiber.Ctx) error {
	idOrden, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CambiarEstadoOrdenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := CambiarEstadoOrden(db, idOrden, &req); err != nil {
		if err.Error() == "orden no encontrada" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "estado actualizado"})
}

// StreamOrdenesController GET /ventas/ordenes/stream
// Canal SSE: mientras la pantalla esté abierta, avisa en tiempo real cuando
// llega una orden nueva de este tenant. Si no hay nadie escuchando, la
// orden no se pierde — queda disponible por FindOrdenesController.
func StreamOrdenesController(c *fiber.Ctx) error {
	tenantSlug := middlewares.GetTenantSlug(c)
	if tenantSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "el canal de órdenes requiere modo multi-tenant",
		})
	}

	ch, unsubscribe := SubscribeOrdenes(tenantSlug)

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
				fmt.Fprintf(w, "event: orden_nueva\ndata: {\"id_orden\":%d}\n\n", n.IDOrden)
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
