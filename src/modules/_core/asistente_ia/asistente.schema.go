package asistente

type ChatRequest struct {
	Pregunta string `json:"pregunta" validate:"required,min=2,max=2000"`
}

type ChatResponse struct {
	Respuesta string          `json:"respuesta"`
	Enlace    *EnlaceSugerido `json:"enlace,omitempty"`
}
