package chat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// responde la validación inicial con la que Meta comprueba la dirección
func VerificarWhatsAppController(c *fiber.Ctx) error {

	token := core.Cfg.Meta_verify_token

	if token != "" && c.Query("hub.mode") == "subscribe" && c.Query("hub.verify_token") == token {
		return c.SendString(c.Query("hub.challenge"))
	}

	return c.SendStatus(http.StatusForbidden)
}

// confirma que el aviso de verdad lo envió Meta
func firmaValida(c *fiber.Ctx) bool {

	if core.Cfg.Meta_app_secret == "" {
		return false
	}

	calculo := hmac.New(sha256.New, []byte(core.Cfg.Meta_app_secret))
	calculo.Write(c.Body())
	esperada := "sha256=" + hex.EncodeToString(calculo.Sum(nil))

	return hmac.Equal([]byte(c.Get("X-Hub-Signature-256")), []byte(esperada))
}

// recibe el aviso y le contesta a Meta de inmediato; el agente trabaja aparte
func RecibirWhatsAppController(c *fiber.Ctx) error {

	if !firmaValida(c) {
		return c.SendStatus(http.StatusUnauthorized)
	}

	var aviso avisoMeta
	if err := json.Unmarshal(c.Body(), &aviso); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	tenant := strings.Clone(middlewares.GetTenantSlug(c))

	go procesarAviso(db, tenant, aviso)

	return c.SendStatus(http.StatusOK)
}

// procesarAviso reparte lo que escribieron los clientes y el dueño
func procesarAviso(db *gorm.DB, tenant string, aviso avisoMeta) {

	for _, entrada := range aviso.Entry {
		for _, cambio := range entrada.Changes {

			idTelefono := cambio.Value.Metadata.IDTelefono

			for _, m := range cambio.Value.Messages {
				if m.Type != "text" {
					continue
				}
				msg := MensajeWhatsApp{IDTelefonoMeta: idTelefono, IDMensajeMeta: m.ID, Telefono: m.From, Texto: m.Text.Body}
				if err := ProcesarMensajeWhatsApp(db, tenant, msg); err != nil {
					logging.Error.Printf("error atendiendo mensaje de WhatsApp: %v", err)
				}
			}

			for _, m := range cambio.Value.MessageEchoes {
				if m.Type != "text" {
					continue
				}
				msg := MensajeWhatsApp{IDTelefonoMeta: idTelefono, IDMensajeMeta: m.ID, Telefono: m.To, Texto: m.Text.Body}
				if err := RegistrarMensajeDueno(db, msg); err != nil {
					logging.Error.Printf("error guardando mensaje del dueño: %v", err)
				}
			}
		}
	}
}

// enviar mensaje al agente y devolver su respuesta
func EnviarMensajeController(c *fiber.Ctx) error {
	var req ChatRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))
	req.TenantSlug = middlewares.GetTenantSlug(c)

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//invocar services
	resultado, err := EnviarMensaje(db, &req)
	if err != nil {

		logging.Error.Printf("error en chat del agente: %v", err)

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(resultado)
}
