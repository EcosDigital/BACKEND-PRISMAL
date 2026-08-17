package asistente

import (
	"errors"
	"fmt"

	empresa "github.com/ecosistema/core/src/modules/_core/gestion_empresa"
	"gorm.io/gorm"
)

// AskCopilot resuelve la configuración de IA de la empresa que pregunta,
// arma el contexto (identidad + manuales) y consulta al modelo. Además
// resuelve si la pregunta corresponde a un tema con pantalla asociada,
// para ofrecer un acceso directo junto con la respuesta.
func AskCopilot(db *gorm.DB, empresaID int64, pregunta string) (ChatResponse, error) {

	empresas, err := empresa.FilterEmpresaByID(db, empresaID)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error consultando la empresa: %v", err)
	}

	if len(empresas) == 0 {
		return ChatResponse{}, errors.New("empresa no encontrada")
	}

	emp := empresas[0]

	if emp.IaEndpoint == "" || emp.IaToken == "" {
		return ChatResponse{}, errors.New("el asistente de IA no está configurado para esta empresa; configúralo en Configuración → Empresas")
	}

	context, err := BuildContext(pregunta)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("error preparando el contexto del asistente: %v", err)
	}

	respuesta, err := AskAI(emp.IaEndpoint, emp.IaToken, context, pregunta)
	if err != nil {
		return ChatResponse{}, err
	}

	return ChatResponse{
		Respuesta: respuesta,
		Enlace:    ResolveEnlace(pregunta),
	}, nil
}
