package notificaciones

import (
	"errors"
	"time"

	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)

// consultar los eventos que el tenant puede configurar, agrupados por modulo.
// antes se copia del catalogo admin lo que corresponda a sus modulos; si eso
// falla se muestra lo que ya estaba copiado, para no dejar la pantalla vacia
func FilterEventosConfigurables(db *gorm.DB) ([]ModuloEventosResponse, error) {

	if err := sincronizarEventos(db); err != nil {
		logging.Error.Printf("error sincronizando eventos de notificación: %v", err)
	}

	rows, err := ListEventosTenant(db)
	if err != nil {
		return nil, err
	}

	results := make([]ModuloEventosResponse, 0)
	posicion := make(map[int64]int)

	for _, r := range rows {

		i, ok := posicion[r.IDModulo]
		if !ok {
			results = append(results, ModuloEventosResponse{
				IDModulo: r.IDModulo,
				Codigo:   r.CodigoModulo,
				Nombre:   r.NombreModulo,
				Eventos:  []EventoResponse{},
			})
			i = len(results) - 1
			posicion[r.IDModulo] = i
		}

		results[i].Eventos = append(results[i].Eventos, EventoResponse{
			Codigo:      r.Codigo,
			Nombre:      r.Nombre,
			Descripcion: r.Descripcion,
			Tipo:        r.Tipo,
		})
	}

	return results, nil

}

// copiar a la base del tenant los eventos del catalogo admin de los modulos
// que tiene contratados. solo se escribe lo que falta o cambio: si todo esta
// igual no se toca ninguna fila. el modulo se cruza por id_ref y el tipo por
// nombre: los ids locales cambian de un tenant a otro
func sincronizarEventos(db *gorm.DB) error {

	modulos, err := ListModulosTenant(db)
	if err != nil {
		return err
	}

	idModuloLocal := make(map[int64]int64, len(modulos))
	idsModulosRef := make([]int64, 0, len(modulos))

	for _, m := range modulos {
		idModuloLocal[m.IDRef] = m.ID
		idsModulosRef = append(idsModulosRef, m.IDRef)
	}

	tipos, err := ListTiposNotificacion(db)
	if err != nil {
		return err
	}

	idTipoLocal := make(map[string]int, len(tipos))
	for _, t := range tipos {
		idTipoLocal[t.Nombre] = t.ID
	}

	catalogo := []eventoCatalogo{}
	if len(idsModulosRef) > 0 {
		catalogo, err = ListEventosCatalogoByModulos(idsModulosRef)
		if err != nil {
			return err
		}
	}

	locales, err := ListEventosLocales(db)
	if err != nil {
		return err
	}

	existentes := make(map[int64]eventoLocal, len(locales))
	for _, l := range locales {
		existentes[l.IDRef] = l
	}

	nuevos := make([]map[string]interface{}, 0)
	cambiados := make(map[int64]map[string]interface{})
	vigentes := make([]int64, 0, len(catalogo))

	for _, e := range catalogo {

		idTipo, ok := idTipoLocal[e.Tipo]
		if !ok {
			logging.Error.Printf("evento %s omitido: el tipo '%s' no existe en este tenant", e.Codigo, e.Tipo)
			continue
		}

		vigentes = append(vigentes, e.ID)

		esperado := eventoLocal{
			IDRef:              e.ID,
			Codigo:             e.Codigo,
			Nombre:             e.Nombre,
			Descripcion:        e.Descripcion,
			IDModuloRef:        idModuloLocal[e.IDModulo],
			IDTipoNotificacion: idTipo,
			OrdenLista:         e.OrdenLista,
			IsActive:           true,
		}

		datos := map[string]interface{}{
			"codigo":               esperado.Codigo,
			"nombre":               esperado.Nombre,
			"descripcion":          esperado.Descripcion,
			"id_modulo_ref":        esperado.IDModuloRef,
			"id_tipo_notificacion": esperado.IDTipoNotificacion,
			"orden_lista":          esperado.OrdenLista,
			"is_active":            true,
		}

		actual, existe := existentes[e.ID]

		//no existe en el tenant: se registra
		if !existe {
			datos["id_ref"] = e.ID
			nuevos = append(nuevos, datos)
			continue
		}

		//existe y cambio en el catalogo (o estaba inactivo): se actualiza.
		//si esta igual no se toca
		if !eventoIgual(actual, esperado) {
			datos["updated_at"] = time.Now()
			cambiados[e.ID] = datos
		}
	}

	//lo nuevo, lo cambiado y la limpieza van juntos, o no se aplica nada
	return db.Transaction(func(tx *gorm.DB) error {

		if err := CreateEventosTenant(tx, nuevos); err != nil {
			return err
		}

		for idRef, datos := range cambiados {
			if err := UpdateEventoTenant(tx, idRef, datos); err != nil {
				return err
			}
		}

		return UpdateDesactivarEventosTenant(tx, vigentes)
	})

}

// comparar el evento copiado en el tenant con el que dice el catalogo
func eventoIgual(a eventoLocal, b eventoLocal) bool {

	descA, descB := "", ""
	if a.Descripcion != nil {
		descA = *a.Descripcion
	}
	if b.Descripcion != nil {
		descB = *b.Descripcion
	}

	return a.Codigo == b.Codigo &&
		a.Nombre == b.Nombre &&
		descA == descB &&
		a.IDModuloRef == b.IDModuloRef &&
		a.IDTipoNotificacion == b.IDTipoNotificacion &&
		a.OrdenLista == b.OrdenLista &&
		a.IsActive == b.IsActive
}

// consultar a quien le llega un evento en la sede, separado en usuarios y roles
func FilterDestinatariosEvento(db *gorm.DB, codigo string, idSede int64) (*DestinatariosEventoResponse, error) {

	idEvento, err := buscarEvento(db, codigo)
	if err != nil {
		return nil, err
	}

	rows, err := ListDestinosEvento(db, idEvento, idSede)
	if err != nil {
		return nil, err
	}

	result := &DestinatariosEventoResponse{
		Usuarios: []int64{},
		Roles:    []int64{},
	}

	for _, r := range rows {
		switch r.IDOrigen {
		case OrigenUsuario:
			result.Usuarios = append(result.Usuarios, r.IDDestino)
		case OrigenRol:
			result.Roles = append(result.Roles, r.IDDestino)
		}
	}

	return result, nil

}

// guardar a quien le llega un evento en la sede. la lista enviada reemplaza
// la anterior; si llega vacia el evento queda sin destinatarios en esa sede
func SaveDestinatariosEvento(db *gorm.DB, codigo string, req *DestinatariosEventoRequest) error {

	if req.SedeID <= 0 {
		return errors.New("no se pudo identificar la sede")
	}

	idEvento, err := buscarEvento(db, codigo)
	if err != nil {
		return err
	}

	usuarios := sinRepetidos(req.Usuarios)
	roles := sinRepetidos(req.Roles)

	if len(usuarios) > 0 {
		existentes, err := ListUsuariosActivosByIDs(db, usuarios)
		if err != nil {
			return err
		}

		if len(existentes) != len(usuarios) {
			return errors.New("uno o más usuarios no existen o están inactivos")
		}
	}

	if len(roles) > 0 {
		existentes, err := ListRolesByIDs(db, roles)
		if err != nil {
			return err
		}

		if len(existentes) != len(roles) {
			return errors.New("uno o más roles no existen")
		}
	}

	rows := make([]map[string]interface{}, 0, len(usuarios)+len(roles))

	agregar := func(idOrigen int, ids []int64) {
		for _, id := range ids {
			rows = append(rows, map[string]interface{}{
				"id_evento":  idEvento,
				"id_origen":  idOrigen,
				"id_destino": id,
				"id_empresa": req.EmpresaID,
				"id_sede":    req.SedeID,
				"created_by": req.UserID,
				"created_at": time.Now(),
			})
		}
	}

	agregar(OrigenUsuario, usuarios)
	agregar(OrigenRol, roles)

	//la lista vieja se borra y la nueva se guarda juntas, o no cambia nada
	return db.Transaction(func(tx *gorm.DB) error {

		if err := DeleteDestinosEvento(tx, idEvento, req.SedeID); err != nil {
			return err
		}

		return CreateDestinosEvento(tx, rows)
	})

}

// buscar el id local de un evento activo por su codigo
func buscarEvento(db *gorm.DB, codigo string) (int64, error) {

	ids, err := ListIDEventoByCodigo(db, codigo)
	if err != nil {
		return 0, err
	}

	if len(ids) <= 0 {
		return 0, errors.New("el evento no existe o no está disponible")
	}

	return ids[0], nil

}

// quitar ids repetidos conservando el orden
func sinRepetidos(ids []int64) []int64 {

	vistos := make(map[int64]bool, len(ids))
	results := make([]int64, 0, len(ids))

	for _, id := range ids {
		if vistos[id] {
			continue
		}
		vistos[id] = true
		results = append(results, id)
	}

	return results

}

// usuarios activos con su rol para la pantalla de Parametros generales
func FilterUsuariosDestinoDetalle(db *gorm.DB) ([]UsuarioDestinoDetalle, error) {
	return ListUsuariosDestinoDetalle(db)
}

// enviar la notificacion automatica de un evento a quienes esten configurados
// en la sede. la llama cualquier modulo (ej. ventas al tener una orden lista)
// y nunca le devuelve error: si el evento no esta configurado o no tiene
// destinatarios no se envia nada, y lo que falle queda en el log. puede
// correr en segundo plano, por eso se protege de un panic
func NotificarEvento(db *gorm.DB, codigoEvento string, idSede int64, idEmpresa int64, titulo string, mensaje string) {

	defer func() {
		if r := recover(); r != nil {
			logging.Error.Printf("error inesperado notificando el evento %s: %v", codigoEvento, r)
		}
	}()

	eventos, err := ListEventoActivoByCodigo(db, codigoEvento)
	if err != nil {
		logging.Error.Printf("error consultando el evento %s: %v", codigoEvento, err)
		return
	}

	//el tenant no tiene el evento (no ha abierto Parametros generales o no
	//tiene el modulo): no hay a quien avisar
	if len(eventos) <= 0 {
		return
	}

	evento := eventos[0]

	rows, err := ListDestinosEvento(db, evento.ID, idSede)
	if err != nil {
		logging.Error.Printf("error consultando destinatarios del evento %s: %v", codigoEvento, err)
		return
	}

	//nadie configurado en esta sede: no se envia nada
	if len(rows) <= 0 {
		return
	}

	destinos := make([]Destino, 0, len(rows))
	for _, r := range rows {
		destinos = append(destinos, Destino{IDOrigen: r.IDOrigen, IDDestino: r.IDDestino})
	}

	idModulo := evento.IDModuloRef

	//sin UserID: la notificacion queda como enviada por el sistema
	req := &NotificacionRequest{
		Titulo:             titulo,
		Mensaje:            mensaje,
		IDTipoNotificacion: evento.IDTipoNotificacion,
		IDModuloRef:        &idModulo,
		EmpresaID:          idEmpresa,
		SedeID:             idSede,
	}

	if _, err := RegisterNotificacionMultiple(db, req, destinos); err != nil {
		logging.Error.Printf("error enviando la notificación del evento %s: %v", codigoEvento, err)
	}

}
