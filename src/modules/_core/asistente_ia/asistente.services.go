package asistente

import (
	"fmt"

	"github.com/ecosistema/core/src/modules/agentes/orquestador"
	"github.com/ecosistema/core/src/shared/integraciones/gemini"
	"gorm.io/gorm"
)

// modeloAsistente es el modelo de Gemini que usa el asistente de Prismar.
const modeloAsistente = "gemini-3.7-flash"

// reglaEncargados le explica a Eco cuándo consultar a un encargado y qué
// decir si no hay uno.
const reglaEncargados = `

## Agentes encargados
Si la pregunta necesita datos reales del negocio (existencias, precios, movimientos, etc.), usa la herramienta consultar_encargado con el módulo que corresponde y entrega su respuesta con tus palabras.
Si la herramienta indica sin_encargado, responde: "Actualmente no tengo forma de responder eso, porque no hay un agente encargado del módulo X." y explica cómo crear uno:
1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Agentes IA**.
3. Haz clic en el ítem **Agentes**.
4. Presiona **+**, ponle un nombre, escribe su tarea en **Instrucciones** y en **Encargado de** elige el módulo X.
5. Presiona **Guardar**.`

// AskCopilot arma el contexto (identidad + manuales del tema) y consulta a
// Gemini, que puede pasarle la pregunta al agente encargado del módulo
// cuando necesita datos reales del negocio. Además resuelve si la pregunta
// corresponde a un tema con pantalla asociada, para ofrecer un acceso
// directo junto con la respuesta.
func AskCopilot(db *gorm.DB, sedeID int64, pregunta string) (ChatResponse, error) {

	context, err := BuildContext(pregunta)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error preparando el contexto del asistente: %v", err)
	}

	sinEncargado := false
	consultar, err := orquestador.HerramientaConsultarEncargado(orquestador.Contexto{DB: db, SedeID: sedeID}, 0, &sinEncargado, false)
	if err != nil {
		return ChatResponse{}, err
	}

	respuesta, err := gemini.GenerarConHerramientas(modeloAsistente, context+reglaEncargados,
		[]gemini.Mensaje{{EsIA: false, Texto: pregunta}}, []gemini.Herramienta{consultar})
	if err != nil {
		return ChatResponse{}, err
	}

	enlace := ResolveEnlace(pregunta)
	if sinEncargado {
		enlace = &EnlaceSugerido{Label: "Ir a Agentes", Path: "/agentes"}
	}

	return ChatResponse{
		Respuesta: respuesta,
		Enlace:    enlace,
	}, nil
}
