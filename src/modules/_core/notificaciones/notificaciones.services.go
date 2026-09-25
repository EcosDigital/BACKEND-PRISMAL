package notificaciones

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sort"

	"github.com/ecosistema/core/src/modules/_core/roles"
	"github.com/ecosistema/core/src/shared/integraciones/firebase"
	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)

// horas de vigencia de la sesion (igual al token de refresco del login).
// pasado ese tiempo sin abrir Prismar, el dispositivo no recibe push.
const vigenciaSesionHoras = 24

// consultar los tipos de notificacion disponibles
func FilterTiposNotificacion(db *gorm.DB) ([]TipoNotificacionResponse, error) {
	return ListTiposNotificacion(db)
}

// registrar el token del navegador del usuario
func RegisterDispositivoPush(db *gorm.DB, req *DispositivoPushRequest) (int64, error) {
	return UpsertDispositivoPush(db, req)
}

// dar de baja el token del navegador del usuario
func BajaDispositivoPush(db *gorm.DB, req *BajaDispositivoRequest) error {
	return UpdateBajaDispositivoPush(db, req)
}

// consultar los destinos disponibles (usuario, rol, modulo)
func FilterOrigenesDestinatario(db *gorm.DB) ([]OrigenDestinatarioResponse, error) {
	return ListOrigenesDestinatario(db)
}

// registrar una notificacion: resuelve a quienes les llega y guarda el
// mensaje junto con su reparto. cuando va dirigida a un modulo, el modulo
// del mensaje es ese mismo destino, para que el frontend no lo mande dos veces.
func RegisterNotificacion(db *gorm.DB, req *NotificacionRequest) (int64, error) {

	if req.IDOrigen != OrigenUsuario &&
		req.IDOrigen != OrigenRol &&
		req.IDOrigen != OrigenModulo {
		return 0, errors.New("el destino de la notificación no es válido")
	}

	if req.IDOrigen == OrigenModulo {
		req.IDModuloRef = &req.IDDestino
	}

	//resolver a quienes les llega
	idsUsuarios, err := resolverDestinatarios(db, req.SedeID, req.IDOrigen, req.IDDestino)
	if err != nil {
		return 0, err
	}

	destinatarios := make([]destinatario, 0, len(idsUsuarios))
	for _, idUsuario := range idsUsuarios {
		destinatarios = append(destinatarios, destinatario{IDUsuario: idUsuario, IDOrigen: req.IDOrigen})
	}

	return guardarYEnviar(db, req, destinatarios)

}

// registrar una notificacion para varios destinos a la vez (ej. dos usuarios
// y un rol). cada persona recibe una sola notificacion y un solo push aunque
// llegue por varios destinos. el id_origen e id_destino del req no se usan.
func RegisterNotificacionMultiple(db *gorm.DB, req *NotificacionRequest, destinos []Destino) (int64, error) {

	//primero usuarios, luego roles y modulos: quien llegue por varios
	//caminos queda registrado como "Usuario"
	ordenados := make([]Destino, len(destinos))
	copy(ordenados, destinos)
	sort.SliceStable(ordenados, func(i, j int) bool {
		return ordenados[i].IDOrigen < ordenados[j].IDOrigen
	})

	destinatarios := make([]destinatario, 0)
	vistos := make(map[int64]bool)

	for _, d := range ordenados {

		//un destino que no se pueda resolver (ej. un usuario inactivo) no
		//debe dejar sin aviso a los demas
		idsUsuarios, err := resolverDestinatarios(db, req.SedeID, d.IDOrigen, d.IDDestino)
		if err != nil {
			logging.Error.Printf("destino de notificación omitido (origen %d, id %d): %v", d.IDOrigen, d.IDDestino, err)
			continue
		}

		for _, idUsuario := range idsUsuarios {
			if vistos[idUsuario] {
				continue
			}
			vistos[idUsuario] = true
			destinatarios = append(destinatarios, destinatario{IDUsuario: idUsuario, IDOrigen: d.IDOrigen})
		}
	}

	return guardarYEnviar(db, req, destinatarios)

}

// guardar el mensaje y su reparto en una sola transaccion y enviar el push
func guardarYEnviar(db *gorm.DB, req *NotificacionRequest, destinatarios []destinatario) (int64, error) {

	if len(destinatarios) <= 0 {
		return 0, errors.New("no se encontraron destinatarios para esta notificación")
	}

	//el mensaje y su reparto se guardan juntos, o no se guarda ninguno
	var idNotificacion int64

	err := db.Transaction(func(tx *gorm.DB) error {

		newID, err := CreateNotificacion(tx, req)
		if err != nil {
			return err
		}

		if err := CreateDestinatarios(tx, newID, destinatarios); err != nil {
			return err
		}

		idNotificacion = newID

		return nil
	})

	if err != nil {
		return 0, err
	}

	idsUsuarios := make([]int64, 0, len(destinatarios))
	for _, d := range destinatarios {
		idsUsuarios = append(idsUsuarios, d.IDUsuario)
	}

	//el push es el aviso en pantalla; la notificacion ya quedo guardada
	enviarPush(db, idsUsuarios, req.Titulo, req.Mensaje)

	return idNotificacion, nil

}

// enviar el push a los dispositivos vigentes de los destinatarios. si algo
// falla no se devuelve error: la notificacion ya quedo guardada y el usuario
// la vera en la campanita cuando entre.
func enviarPush(db *gorm.DB, idsUsuarios []int64, titulo string, mensaje string) {

	tokens, err := ListTokensByUsuarios(db, idsUsuarios, vigenciaSesionHoras)
	if err != nil {
		logging.Error.Printf("error consultando tokens push: %v", err)
		return
	}

	if len(tokens) == 0 {
		return
	}

	invalidos, err := firebase.Enviar(tokens, titulo, mensaje)
	if err != nil {
		logging.Error.Printf("error enviando notificación push: %v", err)
	}

	if err := UpdateBajaTokens(db, invalidos); err != nil {
		logging.Error.Printf("error dando de baja tokens vencidos: %v", err)
	}

}

// resolver los usuarios que reciben la notificacion segun el destino elegido
func resolverDestinatarios(db *gorm.DB, idSede int64, idOrigen int, idDestino int64) ([]int64, error) {

	switch idOrigen {

	case OrigenUsuario:
		exists, err := ListUsuarioActivoByID(db, idDestino)
		if err != nil {
			return nil, err
		}

		if len(exists) <= 0 {
			return nil, errors.New("no se encontró un usuario activo con ese ID")
		}

		return exists, nil

	case OrigenRol:
		return ListUsuariosByRoles(db, []int64{idDestino})

	case OrigenModulo:
		idsRoles, err := rolesConAccesoAModulo(db, idSede, idDestino)
		if err != nil {
			return nil, err
		}

		if len(idsRoles) <= 0 {
			return []int64{}, nil
		}

		return ListUsuariosByRoles(db, idsRoles)

	}

	return nil, errors.New("el destino de la notificación no es válido")

}

// roles de la sede cuyo arbol de permisos incluye el modulo indicado.
// el arbol se guarda como JSON en seguridad.cfg_sedes_roles, por eso se
// recorre en Go: en el primer nivel estan los modulos, y el id de cada uno
// corresponde al id_ref del catalogo, no al id local del tenant.
func rolesConAccesoAModulo(db *gorm.DB, idSede int64, idModulo int64) ([]int64, error) {

	refs, err := ListIDRefModulo(db, idModulo)
	if err != nil {
		return nil, err
	}

	if len(refs) <= 0 {
		return nil, errors.New("no se encontró el módulo indicado")
	}

	idRefModulo := refs[0]

	configs, err := ListConfigRolesBySede(db, idSede)
	if err != nil {
		return nil, err
	}

	idsRoles := make([]int64, 0, len(configs))
	vistos := make(map[int64]bool, len(configs))

	for _, cfg := range configs {

		if vistos[cfg.IDRol] {
			continue
		}

		var items []roles.Item
		if err := json.Unmarshal([]byte(cfg.JsonModules), &items); err != nil {
			logging.Error.Printf("error al deserializar json_modules del rol %d: %v", cfg.IDRol, err)
			continue
		}

		for _, item := range items {
			if item.ID != nil && int64(*item.ID) == idRefModulo {
				idsRoles = append(idsRoles, cfg.IDRol)
				vistos[cfg.IDRol] = true
				break
			}
		}
	}

	return idsRoles, nil

}

// consultar el historial de notificaciones del usuario y cuantas tiene sin leer
func FilterNotificacionesByUsuario(db *gorm.DB, idUsuario int64, idSede int64) ([]NotificacionUsuarioResponse, int64, error) {

	results, err := ListNotificacionesByUsuario(db, idUsuario, idSede)
	if err != nil {
		return nil, 0, err
	}

	noLeidas, err := CountNoLeidasByUsuario(db, idUsuario, idSede)
	if err != nil {
		return nil, 0, err
	}

	return results, noLeidas, nil

}

// marcar una notificacion como leida
func MarcarLeida(db *gorm.DB, idNotificacion int64, idUsuario int64) error {

	err := UpdateLeidoByUsuario(db, idNotificacion, idUsuario)

	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("no se encontró esa notificación para este usuario")
	}

	return err

}

// marcar como leidas todas las notificaciones pendientes del usuario
func MarcarLeidasTodas(db *gorm.DB, idUsuario int64, idSede int64) (int64, error) {
	return UpdateLeidoTodasByUsuario(db, idUsuario, idSede)
}

// usuarios activos para el selector de destino
func FilterUsuariosDestino(db *gorm.DB) ([]DestinoResponse, error) {
	return ListUsuariosDestino(db)
}

// roles para el selector de destino
func FilterRolesDestino(db *gorm.DB) ([]DestinoResponse, error) {
	return ListRolesDestino(db)
}
