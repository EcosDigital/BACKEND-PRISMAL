// Package firebase envía las notificaciones push a los navegadores de los
// usuarios a través de Firebase Cloud Messaging.
//
// El cliente se inicializa una sola vez al arrancar el servidor (Init) y se
// reutiliza en cada envío. Si las credenciales no cargan, el servidor sigue
// funcionando: Enviar devuelve error y la notificación igual queda guardada
// en la base de datos para verse en la campanita.
package firebase

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"
	"google.golang.org/api/option"
)

// tope de tokens por envio que acepta Firebase
const maxTokensPorEnvio = 500

var client *messaging.Client

// Init carga las credenciales y deja listo el cliente de envio
func Init() error {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil,
		option.WithCredentialsFile(core.Cfg.Firebase_credentials))
	if err != nil {
		return err
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return err
	}

	client = messagingClient

	return nil

}

// Enviar manda la notificacion a los tokens indicados y devuelve los tokens
// que Firebase reporto como invalidos o vencidos, para darlos de baja.
func Enviar(tokens []string, titulo string, mensaje string) ([]string, error) {

	if client == nil {
		return nil, errors.New("firebase no está inicializado")
	}

	if len(tokens) == 0 {
		return []string{}, nil
	}

	ctx := context.Background()
	invalidos := make([]string, 0)

	//firebase acepta maximo 500 tokens por envio
	for inicio := 0; inicio < len(tokens); inicio += maxTokensPorEnvio {

		fin := inicio + maxTokensPorEnvio
		if fin > len(tokens) {
			fin = len(tokens)
		}

		lote := tokens[inicio:fin]

		respuesta, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
			Tokens: lote,
			Notification: &messaging.Notification{
				Title: titulo,
				Body:  mensaje,
			},
		})

		if err != nil {
			return invalidos, err
		}

		//si al menos un token funciono, el mensaje esta bien y los rechazos
		//son culpa de cada token. si fallaron todos, el sospechoso es el
		//mensaje y no se da de baja ningun dispositivo.
		huboExito := respuesta.SuccessCount > 0

		for i, resultado := range respuesta.Responses {
			if resultado.Success {
				continue
			}

			//navegador que ya no existe: el token murio, se da de baja
			if messaging.IsUnregistered(resultado.Error) {
				invalidos = append(invalidos, lote[i])
				continue
			}

			//token de otro proyecto de firebase: nunca va a funcionar
			if messaging.IsSenderIDMismatch(resultado.Error) {
				invalidos = append(invalidos, lote[i])
				continue
			}

			//token corrupto, pero solo si el mensaje demostro estar bien
			if huboExito && messaging.IsInvalidArgument(resultado.Error) {
				invalidos = append(invalidos, lote[i])
				continue
			}

			logging.Error.Printf("error enviando push: %v", resultado.Error)
		}
	}

	return invalidos, nil

}
