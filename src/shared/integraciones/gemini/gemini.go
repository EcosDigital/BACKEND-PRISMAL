package gemini

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ecosistema/core/src/core"
)

const baseURL = "https://generativelanguage.googleapis.com/v1beta/models/"

// Mensaje es un turno de la conversación que se le envía al modelo.
type Mensaje struct {
	EsIA  bool
	Texto string
}

// Herramienta es una consulta que el modelo puede pedir mientras arma su
// respuesta. Parametros sigue el formato JSON Schema de Gemini.
type Herramienta struct {
	Nombre      string
	Descripcion string
	Parametros  map[string]interface{}
	Ejecutar    func(args map[string]interface{}) (interface{}, error)
}

type part struct {
	Text             string            `json:"text,omitempty"`
	FunctionCall     *functionCall     `json:"functionCall,omitempty"`
	FunctionResponse *functionResponse `json:"functionResponse,omitempty"`
	// Gemini 3 exige devolver esta firma junto a cada consulta que pidió
	ThoughtSignature string `json:"thoughtSignature,omitempty"`
}

type functionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type functionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

// thinkingLevel "low": los modelos Gemini 3.x razonan al máximo por
// defecto, lo que en un chat hace que cada respuesta tarde más de un
// minuto. Para conversar basta con razonamiento bajo.
type thinkingConfig struct {
	ThinkingLevel string `json:"thinkingLevel"`
}

type generationConfig struct {
	ThinkingConfig thinkingConfig `json:"thinkingConfig"`
}

type tool struct {
	FunctionDeclarations []map[string]interface{} `json:"functionDeclarations"`
}

// toolConfig con mode "NONE" obliga al modelo a responder con texto aunque
// tenga herramientas disponibles.
type toolConfig struct {
	FunctionCallingConfig struct {
		Mode string `json:"mode"`
	} `json:"functionCallingConfig"`
}

type generateRequest struct {
	SystemInstruction *content          `json:"system_instruction,omitempty"`
	Contents          []content         `json:"contents"`
	Tools             []tool            `json:"tools,omitempty"`
	ToolConfig        *toolConfig       `json:"toolConfig,omitempty"`
	GenerationConfig  *generationConfig `json:"generationConfig,omitempty"`
}

type generateResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

// modeloRespaldo se usa cuando el modelo del agente sigue saturado en
// Google (503/429) después de reintentar.
const modeloRespaldo = "gemini-3.5-flash"

// maxRondasHerramientas evita que el modelo encadene consultas sin fin.
const maxRondasHerramientas = 4

// errSaturado indica que Google no pudo atender la solicitud en ese
// momento (saturación, límite de uso o demora); vale la pena reintentar.
var errSaturado = errors.New("Gemini está saturado en este momento, intenta de nuevo en unos segundos")

// Generar envía las instrucciones del agente y la conversación (en orden,
// del mensaje más viejo al más nuevo) al modelo indicado y devuelve el
// texto de la respuesta.
func Generar(modelo, instrucciones string, conversacion []Mensaje) (string, error) {

	resp, err := generarContenido(modelo, instrucciones, armarContents(conversacion), nil, false)
	if err != nil {
		return "", err
	}

	return textoDe(resp)
}

// GenerarConHerramientas funciona como Generar, pero le ofrece al modelo
// herramientas. Si el modelo pide una, se ejecuta, se le devuelve el
// resultado y se le vuelve a preguntar, hasta que responda con texto.
func GenerarConHerramientas(modelo, instrucciones string, conversacion []Mensaje, herramientas []Herramienta) (string, error) {

	contents := armarContents(conversacion)
	porNombre := map[string]Herramienta{}
	var decls []map[string]interface{}
	for _, h := range herramientas {
		porNombre[h.Nombre] = h
		decls = append(decls, map[string]interface{}{
			"name": h.Nombre, "description": h.Descripcion, "parameters": h.Parametros,
		})
	}

	for ronda := 0; ronda <= maxRondasHerramientas; ronda++ {

		// en la última ronda ya no se permiten más consultas: el modelo debe
		// responder con lo que ya averiguó
		resp, err := generarContenido(modelo, instrucciones, contents, decls, ronda == maxRondasHerramientas)
		if err != nil {
			return "", err
		}

		var llamadas []part
		for _, p := range resp.Parts {
			if p.FunctionCall != nil {
				llamadas = append(llamadas, p)
			}
		}

		if len(llamadas) == 0 {
			return textoDe(resp)
		}

		// se devuelve el turno del modelo tal cual (con su firma) y luego
		// el resultado de cada consulta
		contents = append(contents, resp)
		resultados := content{Role: "user"}
		for _, l := range llamadas {
			h, ok := porNombre[l.FunctionCall.Name]
			var salida interface{}
			if !ok {
				salida = "herramienta no disponible"
			} else if r, err := h.Ejecutar(l.FunctionCall.Args); err != nil {
				salida = "no se pudo consultar: " + err.Error()
			} else {
				salida = r
			}
			resultados.Parts = append(resultados.Parts, part{FunctionResponse: &functionResponse{
				Name: l.FunctionCall.Name, Response: map[string]interface{}{"resultado": salida},
			}})
		}
		contents = append(contents, resultados)
	}

	return "", errors.New("el asistente no logró armar la respuesta, intenta preguntarlo de otra forma")
}

func armarContents(conversacion []Mensaje) []content {

	contents := make([]content, 0, len(conversacion))
	for _, m := range conversacion {
		role := "user"
		if m.EsIA {
			role = "model"
		}
		contents = append(contents, content{Role: role, Parts: []part{{Text: m.Texto}}})
	}

	return contents
}

// generarContenido hace la solicitud a Gemini y devuelve el turno del
// modelo. Si Google está saturado, reintenta una vez con el mismo modelo
// y luego con el modelo de respaldo.
func generarContenido(modelo, instrucciones string, contents []content, decls []map[string]interface{}, soloTexto bool) (content, error) {

	if core.Cfg.Gemini_api_key == "" {
		return content{}, errors.New("la llave de Gemini no está configurada")
	}

	reqBody := generateRequest{
		SystemInstruction: &content{Parts: []part{{Text: instrucciones}}},
		Contents:          contents,
		GenerationConfig:  &generationConfig{ThinkingConfig: thinkingConfig{ThinkingLevel: "low"}},
	}

	if len(decls) > 0 {
		reqBody.Tools = []tool{{FunctionDeclarations: decls}}
		if soloTexto {
			reqBody.ToolConfig = &toolConfig{}
			reqBody.ToolConfig.FunctionCallingConfig.Mode = "NONE"
		}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return content{}, fmt.Errorf("error preparando solicitud a Gemini: %v", err)
	}

	intentos := []string{modelo, modelo, modeloRespaldo}

	for i, m := range intentos {
		resp, err := llamar(m, bodyBytes)
		if err == nil || !errors.Is(err, errSaturado) || i == len(intentos)-1 {
			return resp, err
		}
		time.Sleep(time.Second)
	}

	return content{}, errSaturado
}

// llamar hace una sola solicitud a Gemini con el modelo indicado.
func llamar(modelo string, bodyBytes []byte) (content, error) {

	httpReq, err := http.NewRequest(http.MethodPost, baseURL+modelo+":generateContent", bytes.NewReader(bodyBytes))
	if err != nil {
		return content{}, fmt.Errorf("error preparando solicitud a Gemini: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", core.Cfg.Gemini_api_key)

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return content{}, errSaturado
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusServiceUnavailable,
		resp.StatusCode == http.StatusTooManyRequests,
		resp.StatusCode == http.StatusInternalServerError:
		return content{}, errSaturado
	case resp.StatusCode != http.StatusOK:
		return content{}, fmt.Errorf("Gemini respondió con error (%d)", resp.StatusCode)
	}

	var genResp generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return content{}, fmt.Errorf("la respuesta de Gemini no tiene el formato esperado: %v", err)
	}

	if len(genResp.Candidates) == 0 {
		return content{}, errors.New("Gemini no devolvió ninguna respuesta")
	}

	turno := genResp.Candidates[0].Content
	if turno.Role == "" {
		turno.Role = "model"
	}

	return turno, nil
}

// textoDe une el texto de un turno del modelo.
func textoDe(turno content) (string, error) {

	var sb strings.Builder
	for _, p := range turno.Parts {
		sb.WriteString(p.Text)
	}

	respuesta := strings.TrimSpace(sb.String())
	if respuesta == "" {
		return "", errors.New("Gemini no devolvió ninguna respuesta")
	}

	return respuesta, nil
}
