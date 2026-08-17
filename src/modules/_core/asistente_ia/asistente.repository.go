package asistente

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// defaultModel es el modelo servido hoy por el endpoint de IA de cada
// empresa. Si más adelante hace falta variar el modelo por cliente, se
// puede promover a un tercer campo configurable en la ficha de empresa,
// igual que ia_endpoint/ia_token.
const defaultModel = "qwen3:8b"

// defaultKeepAlive le pide a Ollama mantener el modelo cargado en memoria
// 2 horas después de responder, para que preguntas seguidas durante el
// día no paguen el costo de recargarlo. Se reinicia con cada solicitud
// nueva, así que en la práctica no se descarga mientras haya uso activo,
// y se libera solo en horas de inactividad real (de noche, fines de
// semana).
const defaultKeepAlive = "2h"

// aiHTTPClient usa el transporte por defecto de Go (misma resolución de
// DNS/IPv4/IPv6 y negociación de protocolo que usaría cualquier cliente
// HTTP normal, incluido curl) — solo se le da más tiempo de espera,
// porque una recarga del modelo en el servidor de IA puede tardar más
// que una solicitud HTTP típica.
var aiHTTPClient = &http.Client{
	Timeout: 180 * time.Second,
}

// AskAI envía el contexto + la pregunta del usuario al endpoint de IA
// configurado en la ficha de la empresa (formato Ollama de un solo
// "prompt") y devuelve el texto de la respuesta.
//
// El endpoint se usa exactamente como está guardado en la base de datos,
// sin agregarle ni quitarle ninguna ruta — si está mal configurado, el
// error devuelto lo indica explícitamente en vez de intentar adivinar.
func AskAI(endpoint, token, context, pregunta string) (string, error) {

	prompt := fmt.Sprintf("%s\n\n---\n\nPregunta del usuario: %s\n\nRespuesta:", context, pregunta)

	reqBody := ollamaGenerateRequest{
		Model:     defaultModel,
		Prompt:    prompt,
		Think:     false,
		Stream:    false,
		KeepAlive: defaultKeepAlive,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error preparando solicitud al asistente de IA: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("el endpoint de IA configurado para esta empresa no es una URL válida (%q): %v", endpoint, err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := aiHTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("no se pudo contactar al endpoint de IA configurado (%q): %v", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("el endpoint de IA configurado (%q) respondió con error (%d): %s", endpoint, resp.StatusCode, string(raw))
	}

	var ollamaResp ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("la respuesta de %q no tiene el formato esperado: %v", endpoint, err)
	}

	if strings.TrimSpace(ollamaResp.Response) == "" {
		return "", fmt.Errorf("el endpoint de IA configurado (%q) no devolvió ninguna respuesta", endpoint)
	}

	return ollamaResp.Response, nil
}
