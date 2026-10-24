package orquestador

import (
	"fmt"
	"strings"
	"time"

	"github.com/ecosistema/core/src/modules/agentes/configuracion"
	"github.com/ecosistema/core/src/modules/agentes/herramientas"
	"github.com/ecosistema/core/src/shared/integraciones/gemini"
	"gorm.io/gorm"
)

// Contexto identifica el negocio y la sede desde donde se pregunta.
type Contexto struct {
	DB         *gorm.DB
	SedeID     int64
	TenantSlug string
}

// HerramientaRegistrarPedido le permite a un agente que atiende clientes
// registrar un pedido a domicilio, que llega a Pedidos a Domicilio.
func HerramientaRegistrarPedido(ctx Contexto, agenteNombre string) gemini.Herramienta {

	return gemini.Herramienta{
		Nombre: "registrar_pedido",
		Descripcion: "Registra un pedido a domicilio cuando el cliente ya confirmó los productos y dio su nombre, " +
			"teléfono y dirección. El precio y la disponibilidad los valida el sistema.",
		Parametros: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"nombre":     map[string]interface{}{"type": "string", "description": "Nombre del cliente"},
				"telefono":   map[string]interface{}{"type": "string", "description": "Teléfono del cliente"},
				"direccion":  map[string]interface{}{"type": "string", "description": "Dirección de entrega"},
				"referencia": map[string]interface{}{"type": "string", "description": "Referencia para encontrar la dirección (opcional)"},
				"productos": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"producto": map[string]interface{}{"type": "string", "description": "Nombre o código del producto"},
							"cantidad": map[string]interface{}{"type": "integer"},
						},
						"required": []string{"producto", "cantidad"},
					},
				},
			},
			"required": []string{"nombre", "telefono", "direccion", "productos"},
		},
		Ejecutar: func(args map[string]interface{}) (interface{}, error) {
			texto := func(k string) string { s, _ := args[k].(string); return s }

			pedido := herramientas.PedidoCliente{
				Nombre: texto("nombre"), Telefono: texto("telefono"),
				Direccion: texto("direccion"), Referencia: texto("referencia"),
				Agente: agenteNombre,
			}

			lista, _ := args["productos"].([]interface{})
			for _, item := range lista {
				m, _ := item.(map[string]interface{})
				nombre, _ := m["producto"].(string)
				cantidad, _ := m["cantidad"].(float64) // los números del JSON llegan como float64
				pedido.Productos = append(pedido.Productos, herramientas.ProductoPedido{Producto: nombre, Cantidad: int(cantidad)})
			}

			registrado, err := herramientas.RegistrarPedido(ctx.DB, ctx.TenantSlug, ctx.SedeID, pedido)
			if err != nil {
				// el error es un motivo para el cliente, no una falla del agente
				return map[string]interface{}{"registrado": false, "motivo": err.Error()}, nil
			}
			return map[string]interface{}{"registrado": true, "pedido": registrado}, nil
		},
	}
}

type moduloInstalado struct {
	ID     int
	Codigo string
	Nombre string
}

// HerramientaConsultarEncargado le permite a un agente (Eco o un agente de
// chat) pasarle una pregunta al agente encargado de un módulo instalado.
// sinEncargado se marca si se preguntó por un módulo sin encargado activo.
// paraClientes indica que quien pregunta atiende a clientes finales: el
// encargado no revela cantidades de existencias, solo si hay disponible.
func HerramientaConsultarEncargado(ctx Contexto, excluirAgenteID int, sinEncargado *bool, paraClientes bool) (gemini.Herramienta, error) {

	var modulos []moduloInstalado
	if err := ctx.DB.Table("configuracion.ref_modulos_tenant").
		Select("id, codigo, nombre").Scan(&modulos).Error; err != nil {
		return gemini.Herramienta{}, err
	}

	nombres := make([]string, 0, len(modulos))
	for _, m := range modulos {
		nombres = append(nombres, m.Nombre)
	}

	return gemini.Herramienta{
		Nombre: "consultar_encargado",
		Descripcion: "Pregunta al agente encargado de un módulo del negocio por datos reales de ese módulo " +
			"(por ejemplo existencias, precios, movimientos de inventario). Úsala siempre que la pregunta " +
			"necesite datos del negocio; nunca los inventes.",
		Parametros: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"modulo":   map[string]interface{}{"type": "string", "enum": nombres, "description": "Módulo al que pertenece la pregunta"},
				"pregunta": map[string]interface{}{"type": "string", "description": "La pregunta completa para el encargado"},
			},
			"required": []string{"modulo", "pregunta"},
		},
		Ejecutar: func(args map[string]interface{}) (interface{}, error) {
			nombreModulo, _ := args["modulo"].(string)
			pregunta, _ := args["pregunta"].(string)

			var modulo *moduloInstalado
			for i := range modulos {
				if strings.EqualFold(modulos[i].Nombre, nombreModulo) {
					modulo = &modulos[i]
				}
			}
			if modulo == nil {
				return "ese módulo no está instalado en este negocio", nil
			}

			agentes, err := configuracion.FilterAgentes(ctx.DB)
			if err != nil {
				return nil, err
			}

			for _, a := range agentes {
				if a.IsActive && a.IDModuloRef != nil && *a.IDModuloRef == modulo.ID && a.ID != excluirAgenteID {
					return consultarEncargado(ctx, a, *modulo, pregunta, paraClientes)
				}
			}

			if sinEncargado != nil {
				*sinEncargado = true
			}
			return map[string]interface{}{
				"sin_encargado": true,
				"modulo":        modulo.Nombre,
			}, nil
		},
	}, nil
}

// consultarEncargado hace que el agente encargado responda la pregunta
// usando únicamente las consultas de su módulo.
func consultarEncargado(ctx Contexto, agente configuracion.AgenteResponse, modulo moduloInstalado, pregunta string, paraClientes bool) (interface{}, error) {

	consultas := herramientasDeModulo(ctx, modulo.Codigo)
	if len(consultas) == 0 {
		return fmt.Sprintf("El encargado de %s todavía no tiene consultas disponibles para responder esto.", modulo.Nombre), nil
	}

	instrucciones := fmt.Sprintf(`Te llamas %s y eres el agente encargado del módulo %s de este negocio.

%s

Reglas:
- Responde solo con datos que obtengas de tus consultas. Nunca inventes cifras, precios ni existencias.
- Si una consulta no devuelve resultados, dilo claramente.
- Si un artículo no tiene precio público, di que no tiene precio publicado en el catálogo.
- Hoy es %s.`, agente.Nombre, modulo.Nombre, agente.Instrucciones, time.Now().Format("2006-01-02"))

	if paraClientes {
		instrucciones += `
- Quien pregunta atiende a clientes del negocio: NO digas cantidades de existencias ni de movimientos. Di solo si el producto está disponible o agotado, su precio y si tiene domicilio.`
	}

	respuesta, err := gemini.GenerarConHerramientas(agente.Modelo, instrucciones,
		[]gemini.Mensaje{{EsIA: false, Texto: pregunta}}, consultas)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"encargado": agente.Nombre, "respuesta": respuesta}, nil
}

// herramientasDeModulo devuelve las consultas disponibles para cada módulo.
func herramientasDeModulo(ctx Contexto, codigo string) []gemini.Herramienta {
	switch codigo {
	case "MD-004": // Inventario
		return herramientasInventario(ctx)
	}
	return nil
}

func herramientasInventario(ctx Contexto) []gemini.Herramienta {
	texto := func(args map[string]interface{}, k string) string { s, _ := args[k].(string); return s }

	return []gemini.Herramienta{
		{
			Nombre: "ver_catalogo",
			Descripcion: "Lista los productos publicados en el catálogo de la sede con precio, domicilio y existencias. " +
				"Úsala para recomendar u ofrecer productos cuando no preguntan por uno en particular.",
			Parametros: map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Ejecutar: func(a map[string]interface{}) (interface{}, error) {
				return herramientas.ListarCatalogo(ctx.DB, ctx.SedeID)
			},
		},
		{
			Nombre:      "buscar_articulos",
			Descripcion: "Busca artículos por nombre o código y devuelve existencias en la sede, precio público y si tiene domicilio.",
			Parametros: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{"texto": map[string]interface{}{"type": "string"}},
				"required":   []string{"texto"},
			},
			Ejecutar: func(a map[string]interface{}) (interface{}, error) {
				return herramientas.BuscarArticulos(ctx.DB, ctx.SedeID, texto(a, "texto"))
			},
		},
		{
			Nombre:      "movimientos_articulo",
			Descripcion: "Resume por tipo (entradas, bajas, traslados, despachos) cuántas veces y qué cantidad se movió un artículo entre dos fechas.",
			Parametros: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"texto": map[string]interface{}{"type": "string"},
					"desde": map[string]interface{}{"type": "string", "description": "Fecha AAAA-MM-DD"},
					"hasta": map[string]interface{}{"type": "string", "description": "Fecha AAAA-MM-DD"},
				},
				"required": []string{"texto", "desde", "hasta"},
			},
			Ejecutar: func(a map[string]interface{}) (interface{}, error) {
				return herramientas.MovimientosArticulo(ctx.DB, ctx.SedeID, texto(a, "texto"), texto(a, "desde"), texto(a, "hasta"))
			},
		},
		{
			Nombre:      "articulos_bajo_minimo",
			Descripcion: "Lista los artículos con existencias en o por debajo del mínimo de su bodega.",
			Parametros:  map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
			Ejecutar: func(a map[string]interface{}) (interface{}, error) {
				return herramientas.ArticulosBajoMinimo(ctx.DB, ctx.SedeID)
			},
		},
	}
}
