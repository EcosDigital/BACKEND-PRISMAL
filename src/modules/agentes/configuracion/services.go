package configuracion

import (
	"errors"
	"net/url"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/integraciones/whatsapp"
	"gorm.io/gorm"
)

// ConectarWhatsApp conecta al agente el número que el dueño eligió en la
// ventana de Meta y devuelve ese número para mostrarlo en pantalla.
func ConectarWhatsApp(db *gorm.DB, id int64, req *ConectarWhatsAppRequest) (string, error) {

	agentes, err := ListAgenteByID(db, id)
	if err != nil {
		return "", err
	}
	if len(agentes) == 0 {
		return "", errors.New("agente no encontrado")
	}
	if agentes[0].IDModuloRef != nil {
		return "", errors.New("solo un agente sin módulo a cargo puede atender WhatsApp")
	}

	sede, err := GetSedeAgente(db, id)
	if err != nil {
		return "", err
	}
	if sede == 0 {
		return "", errors.New("el agente debe pertenecer a una sede")
	}

	ocupado, err := ExistsTelefonoEnOtroAgente(db, req.IDTelefonoMeta, id)
	if err != nil {
		return "", err
	}
	if ocupado {
		return "", errors.New("ese número ya está conectado a otro agente")
	}

	llave, err := whatsapp.ObtenerLlave(req.Codigo)
	if err != nil {
		return "", err
	}

	telefono, err := whatsapp.ObtenerTelefono(req.IDTelefonoMeta, llave)
	if err != nil {
		return "", err
	}

	if err := whatsapp.SuscribirCuenta(req.IDCuentaMeta, llave); err != nil {
		return "", err
	}

	if !req.Coexistencia {
		if err := whatsapp.RegistrarNumero(req.IDTelefonoMeta, llave); err != nil {
			return "", err
		}
	}

	direccion := "https://" + req.TenantSlug + "." + core.Cfg.Tenant_domain + "/api/v1/agentes/chat/whatsapp?tenant=" + url.QueryEscape(req.TenantSlug)
	if err := whatsapp.ConfigurarAvisos(req.IDTelefonoMeta, llave, direccion); err != nil {
		return "", err
	}

	llaveCifrada, err := whatsapp.CifrarLlave(llave)
	if err != nil {
		return "", err
	}

	err = GuardarWhatsAppAgente(db, id, ConexionWhatsApp{
		Telefono:       telefono,
		IDTelefonoMeta: req.IDTelefonoMeta,
		IDCuentaMeta:   req.IDCuentaMeta,
		LlaveCifrada:   llaveCifrada,
		Conectado:      true,
	}, req.UserID)

	return telefono, err

}

// defaultModelo es el modelo de Gemini que usa el agente cuando no se
// elige otro al configurarlo.
const defaultModelo = "gemini-3.7-flash"

func RegisterAgente(db *gorm.DB, req *AgenteRequest) (int64, error) {

	if req.Modelo == "" {
		req.Modelo = defaultModelo
	}

	if err := validarEncargado(db, req, 0); err != nil {
		return 0, err
	}

	return CreateAgente(db, req)

}

func FilterAgentes(db *gorm.DB) ([]AgenteResponse, error) {
	return ListAgentes(db)
}

func FilterAgenteByID(db *gorm.DB, id int64) ([]AgenteResponse, error) {
	return ListAgenteByID(db, id)
}

func FilterModulosInstalados(db *gorm.DB) ([]ModuloResponse, error) {
	return ListModulosInstalados(db)
}

func EditAgente(db *gorm.DB, id int64, req *AgenteRequest) (int64, error) {

	//verificar si existe el registro por ID
	exists, err := ListAgenteByID(db, id)
	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		return 0, errors.New("No se encontraron registros para actualizar")
	}

	if req.Modelo == "" {
		req.Modelo = defaultModelo
	}

	if err := validarEncargado(db, req, id); err != nil {
		return 0, err
	}

	return UpdateAgente(db, id, req)

}

// solo puede haber un agente encendido encargado de cada módulo, para que
// el asistente sepa siempre a quién consultar
func validarEncargado(db *gorm.DB, req *AgenteRequest, id int64) error {

	if req.IDModuloRef == nil || req.IsActive == nil || !*req.IsActive {
		return nil
	}

	ocupado, err := ExistsEncargadoActivo(db, *req.IDModuloRef, id)
	if err != nil {
		return err
	}

	if ocupado {
		return errors.New("Ya hay otro agente encendido encargado de este módulo")
	}

	return nil

}
