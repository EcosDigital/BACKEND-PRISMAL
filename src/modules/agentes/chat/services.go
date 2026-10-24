package chat

import (
	"errors"
	"fmt"

	"github.com/ecosistema/core/src/modules/agentes/configuracion"
	"github.com/ecosistema/core/src/modules/agentes/orquestador"
	"github.com/ecosistema/core/src/shared/integraciones/gemini"
	"gorm.io/gorm"
)

// limiteHistorial es cuántos mensajes anteriores del chat se le envían
// al modelo como memoria de la conversación.
const limiteHistorial = 3

// limiteHistorialClientes es la memoria de los agentes que atienden clientes
// (sin "Encargado de"): tomar un pedido completo toma varios mensajes.
const limiteHistorialClientes = 10

// reglasClientes se agregan a los agentes que atienden clientes finales
// (sin "Encargado de"): no revelan inventario y saben tomar pedidos.
const reglasClientes = `

Nunca le digas al cliente cuántas unidades hay en inventario; di solo si está disponible o agotado.

Tú atiendes a los clientes del negocio. Si te preguntan cómo usar el sistema Prismar (menús, pantallas, configuración), no lo expliques: di con amabilidad que tú te encargas de atender pedidos y que para eso está Eco, el asistente de Prismar.

Cuando el cliente confirme que quiere hacer un pedido:
1. Repítele los productos, cantidades y el total, y pide confirmación.
2. Pídele nombre, teléfono y dirección de entrega (y una referencia, opcional).
3. Usa registrar_pedido. Si responde con un motivo de error, explícaselo al cliente con amabilidad.
4. Si queda registrado, dale el número de pedido y dile que pronto lo contactarán.
Nunca digas que un pedido quedó registrado si registrar_pedido no lo confirmó.`

// EnviarMensaje recibe el mensaje de un usuario de Prismar que prueba el
// agente, arma la conversación (últimos mensajes + el nuevo), consulta a
// Gemini y guarda ambos mensajes.
func EnviarMensaje(db *gorm.DB, req *ChatRequest) (ChatResponse, error) {

	agentes, err := configuracion.FilterAgenteByID(db, req.IDAgente)
	if err != nil {
		return ChatResponse{}, err
	}

	if len(agentes) == 0 {
		return ChatResponse{}, errors.New("agente no encontrado")
	}

	agente := agentes[0]

	if !agente.IsActive {
		return ChatResponse{}, errors.New("el agente está apagado; enciéndelo en su configuración")
	}

	// identifica el chat: este agente con este usuario de Prismar, en esta
	// sesión de la pantalla (recargar o cambiar de agente empieza otra)
	chatID := map[string]interface{}{
		"id_agente":  agente.ID,
		"origen":     "prismar",
		"id_usuario": req.UserID,
		"id_sesion":  req.IDSesion,
	}

	limite := limiteHistorial
	if agente.IDModuloRef == nil {
		limite = limiteHistorialClientes
	}

	historial, err := ListUltimosMensajes(db, chatID, limite)
	if err != nil {
		return ChatResponse{}, err
	}

	conversacion := append(historial, gemini.Mensaje{EsIA: false, Texto: req.Mensaje})

	instrucciones := fmt.Sprintf("Te llamas %s.\n\n%s\n\n"+
		"Si te preguntan por datos del negocio (productos, precios, disponibilidad, etc.), usa la herramienta "+
		"consultar_encargado; nunca inventes esos datos.", agente.Nombre, agente.Instrucciones)

	// un agente sin "Encargado de" atiende a clientes finales: no revela
	// cantidades de inventario
	paraClientes := agente.IDModuloRef == nil
	if paraClientes {
		instrucciones += reglasClientes
	}

	ctx := orquestador.Contexto{DB: db, SedeID: req.SedeID, TenantSlug: req.TenantSlug}

	// el agente puede consultar a los encargados de módulo, pero no a sí mismo
	consultar, err := orquestador.HerramientaConsultarEncargado(ctx, agente.ID, nil, paraClientes)
	if err != nil {
		return ChatResponse{}, err
	}

	herramientasAgente := []gemini.Herramienta{consultar}
	if paraClientes {
		herramientasAgente = append(herramientasAgente, orquestador.HerramientaRegistrarPedido(ctx, agente.Nombre))
	}

	respuesta, err := gemini.GenerarConHerramientas(agente.Modelo, instrucciones, conversacion, herramientasAgente)
	if err != nil {
		return ChatResponse{}, err
	}

	if err := CreateMensaje(db, req.Mensaje, false, chatID, req.EmpresaID); err != nil {
		return ChatResponse{}, err
	}

	if err := CreateMensaje(db, respuesta, true, chatID, req.EmpresaID); err != nil {
		return ChatResponse{}, err
	}

	return ChatResponse{Respuesta: respuesta}, nil

}
