package espacio_comercial

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/utils"
	"gorm.io/gorm"
)

// ─── Config ─────────────────────────────────────────────────────────────────

// GetConfig retorna la configuración de la empresa, creándola si no existe.
func GetConfig(db *gorm.DB, empresaID, userID int64) (*ConfigResponse, error) {
	config, err := FindConfigByEmpresa(db, empresaID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return CreateConfig(db, empresaID, userID)
	}
	return config, nil
}

// EditConfig actualiza frecuencia y estado activo, creando la fila si aún no existe.
func EditConfig(db *gorm.DB, req *ActualizarConfigRequest) (*ConfigResponse, error) {
	if _, err := GetConfig(db, req.EmpresaID, req.UserID); err != nil {
		return nil, err
	}
	if err := UpdateConfig(db, req.EmpresaID, req); err != nil {
		return nil, err
	}
	return FindConfigByEmpresa(db, req.EmpresaID)
}

// ─── Playlists ──────────────────────────────────────────────────────────────

// GetPlaylists devuelve las playlists de la empresa con sus videos.
func GetPlaylists(db *gorm.DB, empresaID int64) ([]PlaylistResponse, error) {
	playlists, err := ListPlaylists(db, empresaID)
	if err != nil {
		return nil, err
	}
	for i := range playlists {
		videos, err := ListVideosByPlaylist(db, playlists[i].ID)
		if err != nil {
			return nil, err
		}
		playlists[i].Videos = videos
	}
	return playlists, nil
}

// NuevaPlaylist crea una playlist nueva.
func NuevaPlaylist(db *gorm.DB, req *CrearPlaylistRequest) (int64, error) {
	return CreatePlaylist(db, req)
}

// EditPlaylist actualiza una playlist validando que pertenezca a la empresa.
func EditPlaylist(db *gorm.DB, id int64, req *ActualizarPlaylistRequest) error {
	playlist, err := GetPlaylistByID(db, id, req.EmpresaID)
	if err != nil {
		return err
	}
	if playlist == nil {
		return errors.New("playlist no encontrada")
	}
	return UpdatePlaylist(db, id, req.EmpresaID, req)
}

// RemovePlaylist elimina una playlist validando que pertenezca a la empresa.
func RemovePlaylist(db *gorm.DB, id, empresaID int64) error {
	playlist, err := GetPlaylistByID(db, id, empresaID)
	if err != nil {
		return err
	}
	if playlist == nil {
		return errors.New("playlist no encontrada")
	}
	return DeletePlaylist(db, id, empresaID)
}

// ─── Videos de playlist ───────────────────────────────────────────────────────

// NuevoVideo agrega un video a una playlist validando que pertenezca a la empresa.
func NuevoVideo(db *gorm.DB, playlistID int64, req *AgregarVideoRequest) (int64, error) {
	playlist, err := GetPlaylistByID(db, playlistID, req.EmpresaID)
	if err != nil {
		return 0, err
	}
	if playlist == nil {
		return 0, errors.New("playlist no encontrada")
	}
	return AddVideo(db, playlistID, req)
}

// EditVideo actualiza un video validando que pertenezca a la empresa.
func EditVideo(db *gorm.DB, videoID int64, req *ActualizarVideoRequest) error {
	empresaDueno, err := GetVideoPlaylistEmpresa(db, videoID)
	if err != nil {
		return err
	}
	if empresaDueno == 0 || empresaDueno != req.EmpresaID {
		return errors.New("video no encontrado")
	}
	return UpdateVideo(db, videoID, req)
}

// RemoveVideo elimina un video validando que pertenezca a la empresa.
func RemoveVideo(db *gorm.DB, videoID, empresaID int64) error {
	empresaDueno, err := GetVideoPlaylistEmpresa(db, videoID)
	if err != nil {
		return err
	}
	if empresaDueno == 0 || empresaDueno != empresaID {
		return errors.New("video no encontrado")
	}
	return DeleteVideo(db, videoID)
}

// ─── Ads ────────────────────────────────────────────────────────────────────

// GetAds devuelve los ads de la empresa con su URL pública resuelta.
func GetAds(db *gorm.DB, empresaID int64) ([]AdResponse, error) {
	ads, err := ListAds(db, empresaID)
	if err != nil {
		return nil, err
	}
	for i := range ads {
		ads[i].URL = buildAdURL(ads[i].ArchivoPath)
	}
	return ads, nil
}

// SubirAd valida y guarda el archivo de un video publicitario, y crea su registro.
func SubirAd(
	db *gorm.DB,
	empresaID, userID int64,
	nombre string,
	orden int,
	file *multipart.FileHeader,
) (*AdResponse, error) {

	if err := utils.ValidateVideoWithMaxSize(file, MaxAdFileSize, AllowedAdExtensions); err != nil {
		return nil, err
	}

	filename, err := utils.SaveVideo(file, "espacio_comercial_ads")
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	mimeType := AllowedAdExtensions[ext]

	id, err := CreateAd(db, empresaID, userID, nombre, filename, orden, file.Size, mimeType)
	if err != nil {
		_ = utils.DeleteVideo("espacio_comercial_ads", filename)
		return nil, err
	}

	ad, err := GetAdByID(db, id, empresaID)
	if err != nil {
		return nil, err
	}
	ad.URL = buildAdURL(ad.ArchivoPath)
	return ad, nil
}

// EditAd actualiza nombre, orden y estado activo de un ad.
func EditAd(db *gorm.DB, id int64, req *ActualizarAdRequest) error {
	ad, err := GetAdByID(db, id, req.EmpresaID)
	if err != nil {
		return err
	}
	if ad == nil {
		return errors.New("ad no encontrado")
	}
	return UpdateAd(db, id, req.EmpresaID, req)
}

// RemoveAd elimina un ad (registro + archivo físico).
func RemoveAd(db *gorm.DB, id, empresaID int64) error {
	archivoPath, err := DeleteAd(db, id, empresaID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("ad no encontrado")
		}
		return err
	}
	if err := utils.DeleteVideo("espacio_comercial_ads", archivoPath); err != nil {
		logging.Error.Printf("No se pudo eliminar el archivo de ad %s: %v", archivoPath, err)
	}
	return nil
}

// ─── Público (pantalla de reproducción) ──────────────────────────────────────

// GetPlayback arma el payload consumido por la pantalla pública sin login.
// El tenant ya llega aislado por el header X-Tenant-ID (resuelto por
// subdominio), así que no requiere ningún token adicional. Solo incluye
// playlists activas con al menos un video, y ads activos, para que la
// pantalla nunca reciba datos que no pueda reproducir.
func GetPlayback(db *gorm.DB) (*PublicPlaybackResponse, error) {
	config, err := FindFirstActiveConfig(db)
	if err != nil {
		return nil, err
	}
	if config == nil || !config.Activo {
		return nil, errors.New("espacio comercial no encontrado o inactivo")
	}

	playlists, err := GetPlaylists(db, config.EmpresaID)
	if err != nil {
		return nil, err
	}
	activePlaylists := make([]PlaylistResponse, 0, len(playlists))
	for _, p := range playlists {
		if p.Activo && len(p.Videos) > 0 {
			activePlaylists = append(activePlaylists, p)
		}
	}

	ads, err := GetAds(db, config.EmpresaID)
	if err != nil {
		return nil, err
	}
	activeAds := make([]AdResponse, 0, len(ads))
	for _, a := range ads {
		if a.Activo {
			activeAds = append(activeAds, a)
		}
	}

	return &PublicPlaybackResponse{
		FrecuenciaAds: config.FrecuenciaAds,
		Playlists:     activePlaylists,
		Ads:           activeAds,
	}, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func buildAdURL(filename string) string {
	return fmt.Sprintf("%s/uploads/espacio_comercial_ads/%s", core.Cfg.Backend_public_url, filename)
}
