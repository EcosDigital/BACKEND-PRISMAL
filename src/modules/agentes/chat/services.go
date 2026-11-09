package chat

import (
	"errors"
	"fmt"

	"github.com/ecosistema/core/src/modules/agentes/configuracion"
	"github.com/ecosistema/core/src/modules/agentes/orquestador"
	"github.com/ecosistema/core/src/shared/integraciones/gemini"
	"github.com/ecosistema/core/src/shared/integraciones/whatsapp"
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

// datosAgente es lo que necesita un agente para responder, venga el mensaje
// de la pantalla de Prismar o de WhatsApp.
type datosAgente struct {
	ID            int
	Nombre        string
	Instrucciones string
	Modelo        string
	ParaClientes  bool
}

// generarRespuesta arma la conversación y consulta a Gemini.
func generarRespuesta(db *gorm.DB, agente datosAgente, sedeID int64, tenantSlug string, chatID map[string]interface{}, mensaje string) (string, error) {

	limite := limiteHistorial
	if agente.ParaClientes {
		limite = limiteHistorialClientes
	}

	historial, err := ListUltimosMensajes(db, chatID, limite)
	if err != nil {
		return "", err
	}

	conversacion := append(historial, gemini.Mensaje{EsIA: false, Texto: mensaje})

	instrucciones := fmt.Sprintf("Te llamas %s.\n\n%s\n\n"+
		"Si te preguntan por datos del negocio (productos, precios, disponibilidad, etc.), usa la herramienta "+
		"consultar_encargado; nunca inventes esos datos.", agente.Nombre, agente.Instrucciones)

	// un agente sin "Encargado de" atiende a clientes finales: no revela
	// cantidades de inventario
	if agente.ParaClientes {
		instrucciones += reglasClientes
	}

	ctx := orquestador.Contexto{DB: db, SedeID: sedeID, TenantSlug: tenantSlug}

	// el agente puede consultar a los encargados de módulo, pero no a sí mismo
	consultar, err := orquestador.HerramientaConsultarEncargado(ctx, agente.ID, nil, agente.ParaClientes)
	if err != nil {
		return "", err
	}

	herramientasAgente := []gemini.Herramienta{consultar}
	if agente.ParaClientes {
		herramientasAgente = append(herramientasAgente, orquestador.HerramientaRegistrarPedido(ctx, agente.Nombre))
	}

	return gemini.GenerarConHerramientas(agente.Modelo, instrucciones, conversacion, herramientasAgente)

}

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

	respuesta, err := generarRespuesta(db, datosAgente{
		ID:            agente.ID,
		Nombre:        agente.Nombre,
		Instrucciones: agente.Instrucciones,
		Modelo:        agente.Modelo,
		ParaClientes:  agente.IDModuloRef == nil,
	}, req.SedeID, req.TenantSlug, chatID, req.Mensaje)
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

// detalleChatWhatsApp identifica el chat de un agente con un cliente.
func detalleChatWhatsApp(idAgente int, telefono string) map[string]interface{} {
	return map[string]interface{}{"id_agente": idAgente, "origen": "whatsapp", "telefono": telefono}
}

// detalleMensajeWhatsApp agrega quién escribió y el identificador de Meta.
func detalleMensajeWhatsApp(chatID map[string]interface{}, autor, idMensajeMeta string) map[string]interface{} {
	detalle := map[string]interface{}{"autor": autor}
	for k, v := range chatID {
		detalle[k] = v
	}
	if idMensajeMeta != "" {
		detalle["id_mensaje_meta"] = idMensajeMeta
	}
	return detalle
}

// ProcesarMensajeWhatsApp recibe lo que escribió un cliente y le responde
// por WhatsApp, salvo que el agente esté apagado o en pausa en ese chat.
func ProcesarMensajeWhatsApp(db *gorm.DB, tenantSlug string, msg MensajeWhatsApp) error {

	agentes, err := configuracion.FilterAgentePorTelefonoMeta(db, msg.IDTelefonoMeta)
	if err != nil {
		return err
	}
	if len(agentes) == 0 {
		return nil
	}
	agente := agentes[0]

	repetido, err := ExistsMensajeMeta(db, msg.IDMensajeMeta)
	if err != nil {
		return err
	}
	if repetido {
		return nil
	}

	chatID := detalleChatWhatsApp(agente.ID, msg.Telefono)
	detalleCliente := detalleMensajeWhatsApp(chatID, "cliente", msg.IDMensajeMeta)

	pausado, err := ExistsPausaDueno(db, agente.ID, msg.Telefono)
	if err != nil {
		return err
	}

	// si el agente no debe responder, igual se guarda lo que dijo el cliente
	if !agente.IsActive || pausado {
		return CreateMensaje(db, msg.Texto, false, detalleCliente, agente.EmpresaID)
	}

	respuesta, err := generarRespuesta(db, datosAgente{
		ID:            agente.ID,
		Nombre:        agente.Nombre,
		Instrucciones: agente.Instrucciones,
		Modelo:        agente.Modelo,
		ParaClientes:  true,
	}, agente.SedeID, tenantSlug, chatID, msg.Texto)
	if err != nil {
		return err
	}

	llave, err := whatsapp.DescifrarLlave(agente.LlaveCifrada)
	if err != nil {
		return err
	}

	if err := CreateMensaje(db, msg.Texto, false, detalleCliente, agente.EmpresaID); err != nil {
		return err
	}

	if err := whatsapp.EnviarTexto(msg.IDTelefonoMeta, llave, msg.Telefono, respuesta); err != nil {
		return err
	}

	return CreateMensaje(db, respuesta, true, detalleMensajeWhatsApp(chatID, "agente", ""), agente.EmpresaID)

}

// RegistrarMensajeDueno guarda lo que el dueño escribió desde su celular.
// Con eso el agente sabe qué se habló y se calla un rato en ese chat.
func RegistrarMensajeDueno(db *gorm.DB, msg MensajeWhatsApp) error {

	agentes, err := configuracion.FilterAgentePorTelefonoMeta(db, msg.IDTelefonoMeta)
	if err != nil {
		return err
	}
	if len(agentes) == 0 {
		return nil
	}
	agente := agentes[0]

	repetido, err := ExistsMensajeMeta(db, msg.IDMensajeMeta)
	if err != nil || repetido {
		return err
	}

	chatID := detalleChatWhatsApp(agente.ID, msg.Telefono)

	return CreateMensaje(db, msg.Texto, true, detalleMensajeWhatsApp(chatID, "dueno", msg.IDMensajeMeta), agente.EmpresaID)

}
