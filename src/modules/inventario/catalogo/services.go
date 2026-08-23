package catalogo

import (
	"errors"

	"github.com/ecosistema/core/src/shared/integraciones/millave"
	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)

// sincronizarConCentral corre la sincronización con la base central y, según
// el resultado, marca sync_ok en el tenant. El guardado en el tenant ya fue
// exitoso antes de llegar aquí — esto solo decide qué advertencia mostrar.
// Retorna un mensaje de advertencia no vacío si la sincronización falló.
func sincronizarConCentral(db *gorm.DB, tenantSlug string, idCatalogo int64) string {

	if tenantSlug == "" {
		// Modo single-tenant: no hay concepto de tenant para sincronizar.
		return ""
	}

	err := millave.SyncCatalogo(db, tenantSlug, idCatalogo)

	if err := SetSyncOkForCatalogo(db, idCatalogo, err == nil); err != nil {
		logging.Error.Printf("no se pudo actualizar sync_ok del catálogo %d: %v", idCatalogo, err)
	}

	if err != nil {
		logging.Error.Printf("sincronización con Mi Llave falló (catálogo %d): %v", idCatalogo, err)
		return "El artículo se guardó, pero no se pudo sincronizar con Mi Llave. Intenta de nuevo."
	}

	return ""
}

func RegisterCatalogo(db *gorm.DB, req *CatalogoRequest) (int64, error) {

	exists, err := ExistsCatalogoBySede(db, req.IDSede)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("esta sede ya tiene un catálogo creado")
	}

	codigo, err := NextCatalogoCodigo(db)
	if err != nil {
		return 0, err
	}

	return CreateCatalogo(db, req, codigo)
}

func FilterCatalogos(db *gorm.DB, empresaID int64) ([]CatalogoResponse, error) {
	return ListCatalogos(db, empresaID)
}

func FilterCatalogoByID(db *gorm.DB, id int64) (*CatalogoResponseFull, error) {
	return ListCatalogoByID(db, id)
}

// EditCatalogo actualiza el catálogo y, si cambia su estado (activo/inactivo),
// dispara la sincronización con la base central — activar/desactivar el
// catálogo entero afecta a todos sus artículos de golpe.
func EditCatalogo(db *gorm.DB, tenantSlug string, id int64, req *CatalogoUpdateRequest) (int64, string, error) {

	exists, err := ListCatalogoByID(db, id)
	if err != nil {
		return 0, "", err
	}
	if exists == nil || exists.ID == 0 {
		return 0, "", errors.New("no se encontró catálogo con este ID")
	}

	uptID, err := UpdateCatalogo(db, req, id)
	if err != nil {
		return 0, "", err
	}

	syncWarning := sincronizarConCentral(db, tenantSlug, id)

	return uptID, syncWarning, nil
}

// ─── Detalle: artículos publicados dentro del catálogo ──────────────────────

func RegisterCatalogoArticulo(db *gorm.DB, tenantSlug string, catalogoID int64, req *CatalogoArticuloRequest) (int64, string, error) {

	catalogo, err := ListCatalogoByID(db, catalogoID)
	if err != nil {
		return 0, "", err
	}
	if catalogo == nil || catalogo.ID == 0 {
		return 0, "", errors.New("no se encontró catálogo con este ID")
	}

	yaPublicado, err := ExistsCatalogoArticulo(db, catalogoID, req.IDArticulo)
	if err != nil {
		return 0, "", err
	}
	if yaPublicado {
		return 0, "", errors.New("este artículo ya está publicado en este catálogo")
	}

	newID, err := CreateCatalogoArticulo(db, catalogoID, req)
	if err != nil {
		return 0, "", err
	}

	syncWarning := sincronizarConCentral(db, tenantSlug, catalogoID)

	return newID, syncWarning, nil
}

func FilterCatalogoArticulos(db *gorm.DB, catalogoID int64) ([]CatalogoArticuloResponse, error) {
	return ListCatalogoArticulos(db, catalogoID)
}

func EditCatalogoArticulo(db *gorm.DB, tenantSlug string, id int64, req *CatalogoArticuloUpdateRequest) (int64, string, error) {

	exists, err := ListCatalogoArticuloByID(db, id)
	if err != nil {
		return 0, "", err
	}
	if exists == nil || exists.ID == 0 {
		return 0, "", errors.New("no se encontró publicación con este ID")
	}

	idCatalogo, err := GetIDCatalogoByArticuloID(db, id)
	if err != nil {
		return 0, "", err
	}

	uptID, err := UpdateCatalogoArticulo(db, req, id)
	if err != nil {
		return 0, "", err
	}

	syncWarning := sincronizarConCentral(db, tenantSlug, idCatalogo)

	return uptID, syncWarning, nil
}

func RemoveCatalogoArticulo(db *gorm.DB, tenantSlug string, id int64) (string, error) {

	exists, err := ListCatalogoArticuloByID(db, id)
	if err != nil {
		return "", err
	}
	if exists == nil || exists.ID == 0 {
		return "", errors.New("no se encontró publicación con este ID")
	}

	idCatalogo, err := GetIDCatalogoByArticuloID(db, id)
	if err != nil {
		return "", err
	}

	if err := DeleteCatalogoArticulo(db, id); err != nil {
		return "", err
	}

	syncWarning := sincronizarConCentral(db, tenantSlug, idCatalogo)

	return syncWarning, nil
}
