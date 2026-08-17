package asistente

type ChatRequest struct {
	Pregunta string `json:"pregunta" validate:"required,min=2,max=2000"`
}

type ChatResponse struct {
	Respuesta string          `json:"respuesta"`
	Enlace    *EnlaceSugerido `json:"enlace,omitempty"`
}

// ollamaGenerateRequest/ollamaGenerateResponse representan el contrato HTTP
// real del servidor de IA: POST {endpoint}/api/generate (Ollama, endpoint
// de generación con un único "prompt", no de chat con "messages").
type ollamaGenerateRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	Think     bool   `json:"think"`
	Stream    bool   `json:"stream"`
	KeepAlive string `json:"keep_alive,omitempty"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}
