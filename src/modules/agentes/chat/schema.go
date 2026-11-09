package chat

type ChatRequest struct {
	IDAgente   int64  `json:"id_agente" validate:"required,gt=0"`
	Mensaje    string `json:"mensaje"   validate:"required,min=1,max=2000"`
	IDSesion   string `json:"id_sesion" validate:"required,max=64"`
	UserID     int64  `json:"user_id"   validate:"omitempty"`
	EmpresaID  int64  `json:"empresa_id" validate:"omitempty"`
	SedeID     int64  `json:"sede_id"   validate:"omitempty"`
	TenantSlug string `json:"-"`
}

type ChatResponse struct {
	Respuesta string `json:"respuesta"`
}

// mensajeMeta es un mensaje dentro del aviso que envía Meta.
type mensajeMeta struct {
	From string `json:"from"`
	To   string `json:"to"`
	ID   string `json:"id"`
	Type string `json:"type"`
	Text struct {
		Body string `json:"body"`
	} `json:"text"`
}

// avisoMeta es lo que Meta envía cuando pasa algo en un número.
type avisoMeta struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Metadata struct {
					IDTelefono string `json:"phone_number_id"`
				} `json:"metadata"`
				Messages      []mensajeMeta `json:"messages"`
				MessageEchoes []mensajeMeta `json:"message_echoes"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// MensajeWhatsApp es un mensaje de WhatsApp ya leído del aviso de Meta.
type MensajeWhatsApp struct {
	IDTelefonoMeta string
	IDMensajeMeta  string
	Telefono       string
	Texto          string
}
