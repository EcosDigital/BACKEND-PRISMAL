package espacio_comercial

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// ─── Config ─────────────────────────────────────────────────────────────────

// FindConfigByEmpresa devuelve la fila de configuración de una empresa, si existe.
func FindConfigByEmpresa(db *gorm.DB, empresaID int64) (*ConfigResponse, error) {
	var result ConfigResponse
	err := db.Raw(
		`SELECT id, id_empresa, frecuencia_ads, activo
		 FROM espacio_comercial.cfg_espacio_comercial
		 WHERE id_empresa = ?`,
		empresaID,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, nil
	}
	return &result, nil
}

// CreateConfig crea la fila de configuración inicial de una empresa.
func CreateConfig(db *gorm.DB, empresaID, userID int64) (*ConfigResponse, error) {
	err := db.Exec(
		`INSERT INTO espacio_comercial.cfg_espacio_comercial
			(id_empresa, frecuencia_ads, activo, created_at, created_by)
		 VALUES (?, 3, TRUE, ?, ?)`,
		empresaID, time.Now(), userID,
	).Error
	if err != nil {
		return nil, err
	}

	return FindConfigByEmpresa(db, empresaID)
}

// UpdateConfig actualiza frecuencia y estado activo de la configuración.
func UpdateConfig(db *gorm.DB, empresaID int64, req *ActualizarConfigRequest) error {
	tx := db.Exec(
		`UPDATE espacio_comercial.cfg_espacio_comercial
		 SET frecuencia_ads = ?, activo = ?, updated_at = ?, updated_by = ?
		 WHERE id_empresa = ?`,
		req.FrecuenciaAds, req.Activo, time.Now(), req.UserID, empresaID,
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// FindFirstActiveConfig resuelve la configuración a mostrar en la pantalla
// pública. El tenant ya está aislado por su propia base de datos (resuelta
// por subdominio vía X-Tenant-ID), así que basta con tomar la config activa
// de esa base — no se necesita ningún token adicional en la URL.
func FindFirstActiveConfig(db *gorm.DB) (*ConfigResponse, error) {
	var result ConfigResponse
	err := db.Raw(
		`SELECT id, id_empresa, frecuencia_ads, activo
		 FROM espacio_comercial.cfg_espacio_comercial
		 WHERE activo = TRUE
		 ORDER BY id ASC
		 LIMIT 1`,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, nil
	}
	return &result, nil
}

// ─── Playlists ──────────────────────────────────────────────────────────────

// ListPlaylists devuelve las playlists de una empresa, ordenadas por "orden".
func ListPlaylists(db *gorm.DB, empresaID int64) ([]PlaylistResponse, error) {
	var results []PlaylistResponse
	err := db.Raw(
		`SELECT id, nombre, orden, activo
		 FROM espacio_comercial.cfg_playlists
		 WHERE id_empresa = ?
		 ORDER BY orden ASC, id ASC`,
		empresaID,
	).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []PlaylistResponse{}
	}
	return results, nil
}

// GetPlaylistByID retorna una playlist validando que pertenezca a la empresa.
func GetPlaylistByID(db *gorm.DB, id, empresaID int64) (*PlaylistResponse, error) {
	var result PlaylistResponse
	err := db.Raw(
		`SELECT id, nombre, orden, activo
		 FROM espacio_comercial.cfg_playlists
		 WHERE id = ? AND id_empresa = ?`,
		id, empresaID,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, nil
	}
	return &result, nil
}

// CreatePlaylist inserta una nueva playlist y retorna su ID.
func CreatePlaylist(db *gorm.DB, req *CrearPlaylistRequest) (int64, error) {
	var id int64
	err := db.Raw(
		`INSERT INTO espacio_comercial.cfg_playlists
			(id_empresa, nombre, orden, activo, created_at, created_by)
		 VALUES (?, ?, ?, TRUE, ?, ?)
		 RETURNING id`,
		req.EmpresaID, req.Nombre, req.Orden, time.Now(), req.UserID,
	).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdatePlaylist actualiza nombre, orden y estado activo de una playlist.
func UpdatePlaylist(db *gorm.DB, id, empresaID int64, req *ActualizarPlaylistRequest) error {
	tx := db.Exec(
		`UPDATE espacio_comercial.cfg_playlists
		 SET nombre = ?, orden = ?, activo = ?, updated_at = ?, updated_by = ?
		 WHERE id = ? AND id_empresa = ?`,
		req.Nombre, req.Orden, req.Activo, time.Now(), req.UserID, id, empresaID,
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeletePlaylist elimina una playlist (sus videos se borran en cascada).
func DeletePlaylist(db *gorm.DB, id, empresaID int64) error {
	tx := db.Exec(
		`DELETE FROM espacio_comercial.cfg_playlists WHERE id = ? AND id_empresa = ?`,
		id, empresaID,
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Videos de playlist ───────────────────────────────────────────────────────

// ListVideosByPlaylist devuelve los videos de una playlist, ordenados.
func ListVideosByPlaylist(db *gorm.DB, playlistID int64) ([]PlaylistVideoResponse, error) {
	var results []PlaylistVideoResponse
	err := db.Raw(
		`SELECT id, playlist_id, youtube_url, orden
		 FROM espacio_comercial.cfg_playlist_videos
		 WHERE playlist_id = ?
		 ORDER BY orden ASC, id ASC`,
		playlistID,
	).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []PlaylistVideoResponse{}
	}
	return results, nil
}

// AddVideo agrega un video a una playlist y retorna su ID.
func AddVideo(db *gorm.DB, playlistID int64, req *AgregarVideoRequest) (int64, error) {
	var id int64
	err := db.Raw(
		`INSERT INTO espacio_comercial.cfg_playlist_videos
			(playlist_id, youtube_url, orden, created_at)
		 VALUES (?, ?, ?, ?)
		 RETURNING id`,
		playlistID, req.YoutubeURL, req.Orden, time.Now(),
	).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetVideoPlaylistEmpresa retorna el id_empresa dueño del video (join con la
// playlist), para validar propiedad antes de editar/eliminar.
func GetVideoPlaylistEmpresa(db *gorm.DB, videoID int64) (int64, error) {
	var empresaID int64
	err := db.Raw(
		`SELECT p.id_empresa
		 FROM espacio_comercial.cfg_playlist_videos v
		 JOIN espacio_comercial.cfg_playlists p ON p.id = v.playlist_id
		 WHERE v.id = ?`,
		videoID,
	).Scan(&empresaID).Error
	if err != nil {
		return 0, err
	}
	return empresaID, nil
}

// UpdateVideo actualiza la URL y el orden de un video.
func UpdateVideo(db *gorm.DB, videoID int64, req *ActualizarVideoRequest) error {
	tx := db.Exec(
		`UPDATE espacio_comercial.cfg_playlist_videos
		 SET youtube_url = ?, orden = ?
		 WHERE id = ?`,
		req.YoutubeURL, req.Orden, videoID,
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteVideo elimina un video de una playlist.
func DeleteVideo(db *gorm.DB, videoID int64) error {
	tx := db.Exec(`DELETE FROM espacio_comercial.cfg_playlist_videos WHERE id = ?`, videoID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ─── Ads ────────────────────────────────────────────────────────────────────

// ListAds devuelve los videos publicitarios de una empresa, en orden de rotación.
func ListAds(db *gorm.DB, empresaID int64) ([]AdResponse, error) {
	var results []AdResponse
	err := db.Raw(
		`SELECT id, nombre, archivo_path, orden, activo, tamano_bytes, mime_type
		 FROM espacio_comercial.cfg_ads
		 WHERE id_empresa = ?
		 ORDER BY orden ASC, id ASC`,
		empresaID,
	).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []AdResponse{}
	}
	return results, nil
}

// GetAdByID retorna un ad validando que pertenezca a la empresa.
func GetAdByID(db *gorm.DB, id, empresaID int64) (*AdResponse, error) {
	var result AdResponse
	err := db.Raw(
		`SELECT id, nombre, archivo_path, orden, activo, tamano_bytes, mime_type
		 FROM espacio_comercial.cfg_ads
		 WHERE id = ? AND id_empresa = ?`,
		id, empresaID,
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if result.ID == 0 {
		return nil, nil
	}
	return &result, nil
}

// CreateAd inserta un nuevo video publicitario y retorna su ID.
func CreateAd(
	db *gorm.DB,
	empresaID, userID int64,
	nombre, archivoPath string,
	orden int,
	tamanoBytes int64,
	mimeType string,
) (int64, error) {
	var id int64
	err := db.Raw(
		`INSERT INTO espacio_comercial.cfg_ads
			(id_empresa, nombre, archivo_path, orden, activo, tamano_bytes, mime_type, created_at, created_by)
		 VALUES (?, ?, ?, ?, TRUE, ?, ?, ?, ?)
		 RETURNING id`,
		empresaID, nombre, archivoPath, orden, tamanoBytes, mimeType, time.Now(), userID,
	).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateAd actualiza nombre, orden y estado activo de un ad.
func UpdateAd(db *gorm.DB, id, empresaID int64, req *ActualizarAdRequest) error {
	tx := db.Exec(
		`UPDATE espacio_comercial.cfg_ads
		 SET nombre = ?, orden = ?, activo = ?
		 WHERE id = ? AND id_empresa = ?`,
		req.Nombre, req.Orden, req.Activo, id, empresaID,
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteAd elimina un ad y retorna su archivo_path para borrar el archivo físico.
func DeleteAd(db *gorm.DB, id, empresaID int64) (string, error) {
	ad, err := GetAdByID(db, id, empresaID)
	if err != nil {
		return "", err
	}
	if ad == nil {
		return "", sql.ErrNoRows
	}

	tx := db.Exec(
		`DELETE FROM espacio_comercial.cfg_ads WHERE id = ? AND id_empresa = ?`,
		id, empresaID,
	)
	if tx.Error != nil {
		return "", tx.Error
	}
	return ad.ArchivoPath, nil
}
