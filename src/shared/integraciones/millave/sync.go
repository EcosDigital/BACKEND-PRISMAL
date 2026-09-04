// Package millave sincroniza el catálogo público de cada tenant con las
// tablas centrales de la base admin (integraciones.mi_llave_negocios,
// mi_llave_articulos, mi_llave_horarios) que va a leer la app Mi Llave.
//
// Es event-driven (se llama justo después de guardar algo en el tenant,
// no hay job periódico) e idempotente: cada llamada recalcula el estado
// completo del catálogo dado, no solo lo que cambió puntualmente — así
// nunca queda desactualizado sin importar qué se editó.
package millave

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ecosistema/core/src/database"
	"gorm.io/gorm"
)

const maxIntentos = 3
const esperaEntreIntentos = 200 * time.Millisecond

// ─── Lectura del contexto del catálogo dentro del tenant ────────────────────

type catalogoContexto struct {
	IDCatalogo     int64
	IDSede         int64
	CatalogoActivo bool
	SedeNombre     string
	SedeDireccion  string
	SedeTelefono   string
	SedeLat        *float64
	SedeLon        *float64
	SedeImagenURL  string
	EmpresaLogoURL string
}

func leerContextoCatalogo(tenantDB *gorm.DB, idCatalogo int64) (*catalogoContexto, error) {

	var ctx catalogoContexto

	err := tenantDB.Raw(`
		SELECT
			c.id                          AS id_catalogo,
			c.id_sede                     AS id_sede,
			c.is_active                   AS catalogo_activo,
			s.nombre                      AS sede_nombre,
			s.direccion                   AS sede_direccion,
			s.telefono                    AS sede_telefono,
			s.geolocalizacion_lat         AS sede_lat,
			s.geolocalizacion_lon         AS sede_lon,
			COALESCE(s.imagen_url, '')    AS sede_imagen_url,
			COALESCE(e.logo_url, '')      AS empresa_logo_url
		FROM inventario.cfg_catalogos c
		JOIN configuracion.cfg_sedes s ON s.id = c.id_sede
		JOIN configuracion.cfg_empresas e ON e.id = c.id_empresa
		WHERE c.id = ?
	`, idCatalogo).Scan(&ctx).Error

	if err != nil {
		return nil, err
	}
	if ctx.IDCatalogo == 0 {
		return nil, fmt.Errorf("catálogo %d no encontrado", idCatalogo)
	}

	return &ctx, nil
}

type articuloElegible struct {
	IDArticuloOrigen int64
	Codigo           string
	Nombre           string
	Descripcion      string
	ImagenURL        string
	PrecioPublico    float64
}

func leerArticulosElegibles(tenantDB *gorm.DB, idCatalogo int64) ([]articuloElegible, error) {

	results := make([]articuloElegible, 0)

	err := tenantDB.Raw(`
		SELECT
			ca.id                        AS id_articulo_origen,
			a.codigo,
			a.nombre,
			COALESCE(a.descripcion, '')  AS descripcion,
			COALESCE(a.imagen_url, '')   AS imagen_url,
			ca.precio_publico
		FROM inventario.cfg_catalogo_articulos ca
		JOIN inventario.cfg_articulos a ON a.id = ca.id_articulo
		WHERE ca.id_catalogo = ? AND ca.aplica_domicilio = true
	`, idCatalogo).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

type horarioDia struct {
	DiaSemana    int
	Abierto      bool
	HoraApertura *string
	HoraCierre   *string
}

func leerHorarios(tenantDB *gorm.DB, idSede int64) ([]horarioDia, error) {

	results := make([]horarioDia, 0)

	err := tenantDB.Raw(`
		SELECT
			dia_semana,
			abierto,
			TO_CHAR(hora_apertura, 'HH24:MI') AS hora_apertura,
			TO_CHAR(hora_cierre, 'HH24:MI')   AS hora_cierre
		FROM configuracion.cfg_sede_horarios
		WHERE id_sede = ?
	`, idSede).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// ─── Resolución del tenant en la base admin ──────────────────────────────────

// ResolverTenant resuelve el id y el flag delivery_app de un tenant en la
// base admin a partir de su slug (dominio). Se usa desde cualquier módulo
// que necesite pasar de tenantSlug -> id_tenant, no solo desde este paquete.
func ResolverTenant(tenantSlug string) (idTenant int64, deliveryApp bool, err error) {

	var row struct {
		ID          int64
		DeliveryApp bool
	}

	err = database.GormDB.
		Table("configuracion.cfg_tenants").
		Select("id, delivery_app").
		Where("dominio = ?", tenantSlug).
		Take(&row).Error

	if err != nil {
		return 0, false, fmt.Errorf("no se pudo resolver el tenant '%s' en la base admin: %w", tenantSlug, err)
	}

	return row.ID, row.DeliveryApp, nil
}

// ─── Escritura central ───────────────────────────────────────────────────────

func eliminarNegocio(tx *gorm.DB, idTenant, idSede int64) error {
	return tx.Exec(`
		DELETE FROM integraciones.mi_llave_negocios
		WHERE id_tenant = ? AND id_sede_origen = ?
	`, idTenant, idSede).Error
}

func upsertNegocio(tx *gorm.DB, idTenant int64, ctx *catalogoContexto) (int64, error) {

	logoURL := ctx.SedeImagenURL
	if logoURL == "" {
		logoURL = ctx.EmpresaLogoURL
	}

	var idNegocio int64

	err := tx.Raw(`
		INSERT INTO integraciones.mi_llave_negocios
			(id_tenant, id_sede_origen, id_catalogo_origen, nombre, direccion, telefono,
			 geolocalizacion_lat, geolocalizacion_lon, logo_url, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		ON CONFLICT (id_tenant, id_sede_origen) DO UPDATE SET
			id_catalogo_origen  = EXCLUDED.id_catalogo_origen,
			nombre              = EXCLUDED.nombre,
			direccion           = EXCLUDED.direccion,
			telefono            = EXCLUDED.telefono,
			geolocalizacion_lat = EXCLUDED.geolocalizacion_lat,
			geolocalizacion_lon = EXCLUDED.geolocalizacion_lon,
			logo_url            = EXCLUDED.logo_url,
			updated_at          = NOW()
		RETURNING id
	`,
		idTenant, ctx.IDSede, ctx.IDCatalogo, ctx.SedeNombre, ctx.SedeDireccion, ctx.SedeTelefono,
		ctx.SedeLat, ctx.SedeLon, logoURL,
	).Scan(&idNegocio).Error

	if err != nil {
		return 0, err
	}

	return idNegocio, nil
}

func reconciliarHorarios(tx *gorm.DB, idNegocio int64, dias []horarioDia) error {

	if err := tx.Exec(`DELETE FROM integraciones.mi_llave_horarios WHERE id_negocio = ?`, idNegocio).Error; err != nil {
		return err
	}

	for _, d := range dias {
		err := tx.Exec(`
			INSERT INTO integraciones.mi_llave_horarios
				(id_negocio, dia_semana, abierto, hora_apertura, hora_cierre)
			VALUES (?, ?, ?, ?, ?)
		`, idNegocio, d.DiaSemana, d.Abierto, d.HoraApertura, d.HoraCierre).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func reconciliarArticulos(tx *gorm.DB, idNegocio int64, articulos []articuloElegible) error {

	idsVigentes := make([]int64, 0, len(articulos))
	for _, a := range articulos {
		idsVigentes = append(idsVigentes, a.IDArticuloOrigen)
	}

	// Borra lo que ya no debe estar (despublicado, ya no aplica a domicilio, etc.)
	del := tx.Table("integraciones.mi_llave_articulos").Where("id_negocio = ?", idNegocio)
	if len(idsVigentes) > 0 {
		del = del.Where("id_articulo_origen NOT IN ?", idsVigentes)
	}
	if err := del.Delete(nil).Error; err != nil {
		return err
	}

	// Upsert de cada artículo vigente
	for _, a := range articulos {
		err := tx.Exec(`
			INSERT INTO integraciones.mi_llave_articulos
				(id_negocio, id_articulo_origen, codigo, nombre, descripcion, imagen_url, precio_publico, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
			ON CONFLICT (id_negocio, id_articulo_origen) DO UPDATE SET
				codigo         = EXCLUDED.codigo,
				nombre         = EXCLUDED.nombre,
				descripcion    = EXCLUDED.descripcion,
				imagen_url     = EXCLUDED.imagen_url,
				precio_publico = EXCLUDED.precio_publico,
				updated_at     = NOW()
		`, idNegocio, a.IDArticuloOrigen, a.Codigo, a.Nombre, a.Descripcion, a.ImagenURL, a.PrecioPublico).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// ─── Punto de entrada ────────────────────────────────────────────────────────

// SyncCatalogo reconcilia en la base central el negocio (cabecera + horarios)
// y todos los artículos actualmente elegibles del catálogo indicado. Se debe
// llamar después de cualquier alta/edición/baja de un artículo publicado, o
// al activar/desactivar el catálogo completo.
//
// Reintenta unas pocas veces ante fallas transitorias; si sigue fallando,
// retorna el error para que el caller marque el estado como "no sincronizado"
// y lo informe en pantalla — el guardado en el tenant YA fue exitoso antes de
// llegar aquí, esto solo afecta el reflejo en Mi Llave.
func SyncCatalogo(tenantDB *gorm.DB, tenantSlug string, idCatalogo int64) error {

	var ultimoError error

	for intento := 1; intento <= maxIntentos; intento++ {
		if err := sincronizarUnaVez(tenantDB, tenantSlug, idCatalogo); err != nil {
			ultimoError = err
			if intento < maxIntentos {
				time.Sleep(esperaEntreIntentos)
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("no se pudo sincronizar con Mi Llave tras %d intentos: %w", maxIntentos, ultimoError)
}

// SyncCatalogoPorSede resuelve qué catálogo tiene esa sede (a lo sumo uno,
// UNIQUE por sede) y sincroniza. No hace nada, sin error, si esa sede
// todavía no tiene ningún catálogo creado — no es una falla, simplemente
// no hay nada que reflejar en Mi Llave todavía.
func SyncCatalogoPorSede(tenantDB *gorm.DB, tenantSlug string, idSede int64) error {

	var idCatalogo int64

	err := tenantDB.
		Table("inventario.cfg_catalogos").
		Select("id").
		Where("id_sede = ?", idSede).
		Row().
		Scan(&idCatalogo)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return SyncCatalogo(tenantDB, tenantSlug, idCatalogo)
}

func sincronizarUnaVez(tenantDB *gorm.DB, tenantSlug string, idCatalogo int64) error {

	idTenant, deliveryApp, err := ResolverTenant(tenantSlug)
	if err != nil {
		return err
	}

	ctx, err := leerContextoCatalogo(tenantDB, idCatalogo)
	if err != nil {
		return err
	}

	elegible := deliveryApp && ctx.CatalogoActivo

	if !elegible {
		return eliminarNegocio(database.GormDB, idTenant, ctx.IDSede)
	}

	articulos, err := leerArticulosElegibles(tenantDB, idCatalogo)
	if err != nil {
		return err
	}

	dias, err := leerHorarios(tenantDB, ctx.IDSede)
	if err != nil {
		return err
	}

	return database.GormDB.Transaction(func(tx *gorm.DB) error {

		idNegocio, err := upsertNegocio(tx, idTenant, ctx)
		if err != nil {
			return err
		}
		if idNegocio == 0 {
			return errors.New("no se pudo resolver el id del negocio central")
		}

		if err := reconciliarHorarios(tx, idNegocio, dias); err != nil {
			return err
		}

		return reconciliarArticulos(tx, idNegocio, articulos)
	})
}
