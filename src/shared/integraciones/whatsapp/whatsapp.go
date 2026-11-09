package whatsapp

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/ecosistema/core/src/core"
)

const baseURL = "https://graph.facebook.com/v27.0/"

var clienteHTTP = &http.Client{Timeout: 20 * time.Second}

// EnviarTexto le manda un mensaje de texto a un cliente por WhatsApp, desde
// el número del negocio. Los errores nunca incluyen la llave.
func EnviarTexto(idTelefonoMeta, llave, telefonoCliente, texto string) error {

	cuerpo, err := json.Marshal(map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                telefonoCliente,
		"type":              "text",
		"text":              map[string]string{"body": texto},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+idTelefonoMeta+"/messages", bytes.NewReader(cuerpo))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+llave)
	req.Header.Set("Content-Type", "application/json")

	resp, err := clienteHTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Meta rechazó el mensaje (código %d)", resp.StatusCode)
	}

	return nil
}

// ObtenerLlave cambia el código temporal que entrega la ventana de Meta por
// la llave con la que Prismar enviará mensajes por ese número.
func ObtenerLlave(codigo string) (string, error) {

	parametros := url.Values{}
	parametros.Set("client_id", core.Cfg.Meta_app_id)
	parametros.Set("client_secret", core.Cfg.Meta_app_secret)
	parametros.Set("code", codigo)

	resp, err := clienteHTTP.Get(baseURL + "oauth/access_token?" + parametros.Encode())
	if err != nil {
		// el error original trae la dirección completa, con el secreto
		return "", errors.New("no se pudo conectar con Meta")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Meta rechazó la conexión (código %d)", resp.StatusCode)
	}

	var datos struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&datos); err != nil || datos.AccessToken == "" {
		return "", errors.New("Meta no entregó la llave")
	}

	return datos.AccessToken, nil
}

// llamarMeta hace un pedido a Meta. Los errores nunca incluyen la llave.
func llamarMeta(metodo, ruta, llave string, cuerpo map[string]interface{}, destino interface{}) error {

	var contenido bytes.Buffer
	if cuerpo != nil {
		if err := json.NewEncoder(&contenido).Encode(cuerpo); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(metodo, baseURL+ruta, &contenido)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+llave)
	req.Header.Set("Content-Type", "application/json")

	resp, err := clienteHTTP.Do(req)
	if err != nil {
		return errors.New("no se pudo conectar con Meta")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Meta rechazó el pedido (código %d)", resp.StatusCode)
	}

	if destino != nil {
		return json.NewDecoder(resp.Body).Decode(destino)
	}

	return nil
}

// ObtenerTelefono consulta qué número eligió el dueño en la ventana de Meta.
func ObtenerTelefono(idTelefonoMeta, llave string) (string, error) {

	var datos struct {
		Telefono string `json:"display_phone_number"`
	}

	err := llamarMeta(http.MethodGet, idTelefonoMeta+"?fields=display_phone_number", llave, nil, &datos)

	return datos.Telefono, err
}

// SuscribirCuenta le pide a Meta que le avise a Prismar lo que pase en esa cuenta.
func SuscribirCuenta(idCuentaMeta, llave string) error {
	return llamarMeta(http.MethodPost, idCuentaMeta+"/subscribed_apps", llave, nil, nil)
}

// RegistrarNumero activa un número nuevo. No se usa si el número sigue
// funcionando en la app de WhatsApp Business del celular.
func RegistrarNumero(idTelefonoMeta, llave string) error {

	pin, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return errors.New("no se pudo generar el pin")
	}

	return llamarMeta(http.MethodPost, idTelefonoMeta+"/register", llave, map[string]interface{}{
		"messaging_product": "whatsapp",
		"pin":               fmt.Sprintf("%06d", pin.Int64()),
	}, nil)
}

// ConfigurarAvisos le indica a Meta a qué dirección avisar los mensajes de
// ese número. La dirección incluye el negocio al que pertenece.
func ConfigurarAvisos(idTelefonoMeta, llave, direccion string) error {
	return llamarMeta(http.MethodPost, idTelefonoMeta, llave, map[string]interface{}{
		"webhook_configuration": map[string]string{
			"override_callback_uri": direccion,
			"verify_token":          core.Cfg.Meta_verify_token,
		},
	}, nil)
}
