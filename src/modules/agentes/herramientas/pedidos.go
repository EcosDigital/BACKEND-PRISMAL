package herramientas

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ecosistema/core/src/database"
	"github.com/ecosistema/core/src/shared/integraciones/millave"
	"gorm.io/gorm"
)

type ProductoPedido struct {
	Producto string
	Cantidad int
}

type PedidoCliente struct {
	Nombre     string
	Telefono   string
	Direccion  string
	Referencia string
	Productos  []ProductoPedido
	Agente     string
}

type LineaPedido struct {
	Nombre         string  `json:"nombre"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

type PedidoRegistrado struct {
	IDPedido  int64         `json:"numero_pedido"`
	Total     float64       `json:"total"`
	Productos []LineaPedido `json:"productos"`
}

type productoCatalogo struct {
	IDArticulo      int64
	Codigo          string
	Nombre          string
	PrecioPublico   float64
	AplicaDomicilio bool
	Existencias     float64
}

// RegistrarPedido crea un pedido a domicilio en la base central (la misma
// de los pedidos de Mi Llave). Precio, domicilio y disponibilidad se toman
// siempre del catálogo y el inventario de la sede, nunca de lo que diga el
// modelo. El valor del domicilio va en 0 y el documento del cliente vacío.
func RegistrarPedido(db *gorm.DB, tenantSlug string, sedeID int64, p PedidoCliente) (*PedidoRegistrado, error) {

	if strings.TrimSpace(p.Nombre) == "" || strings.TrimSpace(p.Telefono) == "" || strings.TrimSpace(p.Direccion) == "" {
		return nil, errors.New("faltan datos del cliente: nombre, teléfono y dirección son obligatorios")
	}
	if len(p.Productos) == 0 {
		return nil, errors.New("el pedido no tiene productos")
	}

	idTenant, deliveryApp, err := millave.ResolverTenant(tenantSlug)
	if err != nil {
		return nil, err
	}
	if !deliveryApp {
		return nil, errors.New("este negocio no tiene activos los pedidos a domicilio")
	}

	var idNegocio int64
	database.GormDB.Raw(`SELECT id FROM integraciones.mi_llave_negocios WHERE id_tenant = ? AND id_sede_origen = ?`,
		idTenant, sedeID).Scan(&idNegocio)
	if idNegocio == 0 {
		return nil, errors.New("esta sede no está habilitada para recibir domicilios")
	}

	registro := &PedidoRegistrado{}
	var lineas []productoCatalogo
	var cantidades []int

	for _, item := range p.Productos {
		if item.Cantidad <= 0 {
			return nil, fmt.Errorf("la cantidad de «%s» debe ser mayor a cero", item.Producto)
		}
		prod, err := buscarEnCatalogo(db, sedeID, item.Producto)
		if err != nil {
			return nil, err
		}
		if !prod.AplicaDomicilio {
			return nil, fmt.Errorf("«%s» no está disponible para domicilio", prod.Nombre)
		}
		if prod.Existencias < float64(item.Cantidad) {
			return nil, fmt.Errorf("no hay suficientes unidades disponibles de «%s»", prod.Nombre)
		}
		lineas = append(lineas, prod)
		cantidades = append(cantidades, item.Cantidad)
		registro.Total += prod.PrecioPublico * float64(item.Cantidad)
		registro.Productos = append(registro.Productos, LineaPedido{prod.Nombre, item.Cantidad, prod.PrecioPublico})
	}

	err = database.GormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`
			INSERT INTO integraciones.mov_pedidos
				(id_negocio, id_tenant, id_estado, valor_subtotal, valor_domicilio, valor_total, observaciones, created_at)
			SELECT ?, ?, e.id, ?, 0, ?, ?, NOW()
			FROM integraciones.ref_estado_pedido e WHERE e.codigo = '001'
			RETURNING id`,
			idNegocio, idTenant, registro.Total, registro.Total,
			strings.TrimSpace(fmt.Sprintf("Pedido tomado por el agente %s (chat). %s", p.Agente, p.Referencia)),
		).Scan(&registro.IDPedido).Error; err != nil {
			return err
		}

		for i, prod := range lineas {
			if err := tx.Exec(`
				INSERT INTO integraciones.mov_pedidos_detalle
					(id_pedido, id_articulo_origen, codigo, nombre, cantidad, precio_unitario)
				VALUES (?, ?, ?, ?, ?, ?)`,
				registro.IDPedido, prod.IDArticulo, prod.Codigo, prod.Nombre, cantidades[i], prod.PrecioPublico,
			).Error; err != nil {
				return err
			}
		}

		return tx.Exec(`
			INSERT INTO integraciones.mov_pedidos_entrega
				(id_pedido, cliente_nombre, cliente_telefono, tipo_documento, cliente_documento, direccion_entrega, referencia)
			VALUES (?, ?, ?, '', '', ?, ?)`,
			registro.IDPedido, p.Nombre, p.Telefono, p.Direccion, p.Referencia,
		).Error
	})
	if err != nil {
		return nil, fmt.Errorf("no se pudo registrar el pedido: %v", err)
	}

	return registro, nil
}

// buscarEnCatalogo ubica un producto publicado en el catálogo de la sede por
// código o nombre. Si coincide con varios, pide precisar.
func buscarEnCatalogo(db *gorm.DB, sedeID int64, texto string) (productoCatalogo, error) {

	var rows []productoCatalogo
	err := db.Raw(`
		SELECT a.id AS id_articulo, a.codigo, a.nombre, ca.precio_publico, ca.aplica_domicilio,
		       COALESCE(SUM(e.stock_actual), 0) AS existencias
		FROM inventario.cfg_catalogos c
		JOIN inventario.cfg_catalogo_articulos ca ON ca.id_catalogo = c.id
		JOIN inventario.cfg_articulos a ON a.id = ca.id_articulo AND a.is_active
		LEFT JOIN inventario.inv_existencias_articulos e ON e.id_articulo = a.id
		     AND e.id_bodega IN (SELECT id FROM inventario.cfg_bodegas WHERE id_sede = ? AND is_active)
		WHERE c.id_sede = ? AND c.is_active AND (a.codigo ILIKE ? OR a.nombre ILIKE ?)
		GROUP BY a.id, a.codigo, a.nombre, ca.precio_publico, ca.aplica_domicilio`,
		sedeID, sedeID, texto, "%"+texto+"%").Scan(&rows).Error
	if err != nil {
		return productoCatalogo{}, err
	}

	switch len(rows) {
	case 0:
		return productoCatalogo{}, fmt.Errorf("no encontré «%s» en el catálogo", texto)
	case 1:
		return rows[0], nil
	}

	for _, r := range rows {
		if strings.EqualFold(r.Nombre, texto) || strings.EqualFold(r.Codigo, texto) {
			return r, nil
		}
	}

	nombres := make([]string, 0, len(rows))
	for _, r := range rows {
		nombres = append(nombres, r.Nombre)
	}
	return productoCatalogo{}, fmt.Errorf("«%s» coincide con varios productos (%s); pide al cliente que precise cuál", texto, strings.Join(nombres, ", "))
}
