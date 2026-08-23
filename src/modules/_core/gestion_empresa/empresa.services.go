package empresa

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/integraciones/millave"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/utils"
	"gorm.io/gorm"
)

// sincronizarSedeConCentral resincroniza el catálogo de esa sede (si tiene uno)
// hacia la base central. El guardado en el tenant ya fue exitoso antes de
// llegar aquí — el resultado solo decide qué advertencia mostrar.
func sincronizarSedeConCentral(db *gorm.DB, tenantSlug string, idSede int64) string {

	if tenantSlug == "" {
		return ""
	}

	if err := millave.SyncCatalogoPorSede(db, tenantSlug, idSede); err != nil {
		logging.Error.Printf("sincronización con Mi Llave falló (sede %d): %v", idSede, err)
		return "Los datos se guardaron, pero no se pudo sincronizar con Mi Llave. Intenta de nuevo."
	}

	return ""
}

// MaxImagenSedeSize límite de peso para la imagen/logo de una sede (5MB).
const MaxImagenSedeSize = 5 * 1024 * 1024

func RegisterEmpresa(db *gorm.DB, req *EmpresaRequest) (int64, error) {

	//verificar existencia
	exists, err := ListEmpresaByNit(db, req.Nit)

	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateEmpresa(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func FilterLastEmpresa(db *gorm.DB) ([]EmpresaResponse, error) {
	return ListEmpresaLast(db)
}

func FilterEmpresaByID(db *gorm.DB, id int64) ([]EmpresaResponseFull, error) {
	return ListEmpresaByID(db, id)
}

func EditEmpresa(db *gorm.DB, id int64, req *EmpresaUpdateRequest) (int64, error) {

	//verificar registro por ID
	exist, err := ListEmpresaByID(db, id)
	if err != nil {
		return 0, err
	}

	if len(exist) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No Se encontraron resultados de este registro")
	}

	//actualizar registro
	regID, err := UpdateEmpresa(db, id, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

// SEDES
func FilterLsatSede(db *gorm.DB) ([]SedeResponse, error) {
	return ListSedeLast(db)
}

func FilterSedeByID(db *gorm.DB, id int64) ([]SedeResponseFull, error) {
	return ListSedeByID(db, id)
}

func RegisterSede(db *gorm.DB, req *SedeRequest) (int64, error) {

	//verificar existencia
	exists, err := ListSedeByCode(db, req.Codigo)

	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateSede(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func EditSede(db *gorm.DB, tenantSlug string, id int64, req *SedeUpdateRequest) (int64, string, error) {

	//verificar registro por ID
	exist, err := ListSedeByID(db, id)
	if err != nil {
		return 0, "", err
	}

	if len(exist) <= 0 {
		// Existe registro con número de documento
		return 0, "", errors.New("No Se encontraron resultados de este registro")
	}

	//actualizar registro
	regID, err := UpdateSede(db, id, req)
	if err != nil {
		return 0, "", err
	}

	syncWarning := sincronizarSedeConCentral(db, tenantSlug, id)

	return regID, syncWarning, nil

}

// ─── Imagen de sede ──────────────────────────────────────────────────────────

// SetSedeImagen valida, guarda y asocia una nueva imagen a la sede. Si ya tenía
// una imagen, el archivo físico anterior se borra DESPUÉS de que la BD quede
// actualizada con la nueva URL.
func SetSedeImagen(db *gorm.DB, tenantSlug string, id int64, file *multipart.FileHeader) (string, string, error) {

	oldImagenURL, exists, err := GetSedeImagenURLByID(db, id)
	if err != nil {
		return "", "", err
	}
	if !exists {
		return "", "", errors.New("no se encontró sede con este ID")
	}

	if err := utils.ValidateImageWithMaxSize(file, MaxImagenSedeSize); err != nil {
		return "", "", err
	}

	filename, err := utils.SaveImage(file)
	if err != nil {
		return "", "", fmt.Errorf("error guardando la imagen: %w", err)
	}

	newImagenURL := fmt.Sprintf("%s/uploads/%s", core.Cfg.Backend_public_url, filename)

	if err := UpdateSedeImagen(db, id, newImagenURL); err != nil {
		_ = utils.DeleteImage(filename)
		return "", "", fmt.Errorf("error actualizando la sede: %w", err)
	}

	if oldImagenURL != "" {
		oldFilename := utils.ExtractFilenameFromUrl(oldImagenURL)
		if err := utils.DeleteImage(oldFilename); err != nil {
			logging.Error.Printf("no se pudo borrar la imagen anterior de la sede %d (%s): %v", id, oldFilename, err)
		}
	}

	syncWarning := sincronizarSedeConCentral(db, tenantSlug, id)

	return newImagenURL, syncWarning, nil
}

// ─── Horario de atención por sede ───────────────────────────────────────────

func FilterSedeHorarios(db *gorm.DB, idSede int64) ([]SedeHorarioResponse, error) {
	return ListSedeHorarios(db, idSede)
}

// SaveSedeHorario crea o actualiza el horario de un día puntual de la sede.
func SaveSedeHorario(db *gorm.DB, tenantSlug string, idSede int64, diaSemana int, req *SedeHorarioRequest) (*SedeHorarioResponse, string, error) {

	if diaSemana < 0 || diaSemana > 6 {
		return nil, "", errors.New("día de la semana inválido (debe estar entre 0 y 6)")
	}

	// hora_cierre < hora_apertura es válido (cruza medianoche, ej. abre 5pm
	// cierra 2am) — solo se exige que ambas horas estén presentes.
	if req.Abierto != nil && *req.Abierto {
		if req.HoraApertura == nil || *req.HoraApertura == "" || req.HoraCierre == nil || *req.HoraCierre == "" {
			return nil, "", errors.New("debes indicar hora de apertura y de cierre para un día abierto")
		}
	}

	result, err := UpsertSedeHorario(db, idSede, diaSemana, req)
	if err != nil {
		return nil, "", err
	}

	syncWarning := sincronizarSedeConCentral(db, tenantSlug, idSede)

	return result, syncWarning, nil
}
