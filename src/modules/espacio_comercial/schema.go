package espacio_comercial

// ─── Config ─────────────────────────────────────────────────────────────────

// MaxAdFileSize límite de peso para un video publicitario (350MB).
const MaxAdFileSize = 350 * 1024 * 1024

// AllowedAdExtensions extensiones de video aceptadas para los anuncios.
var AllowedAdExtensions = map[string]string{
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".mov":  "video/quicktime",
}

// ConfigResponse fila de espacio_comercial.cfg_espacio_comercial
type ConfigResponse struct {
	ID            int64 `json:"id"             gorm:"column:id"`
	EmpresaID     int64 `json:"-"              gorm:"column:id_empresa"`
	FrecuenciaAds int   `json:"frecuencia_ads" gorm:"column:frecuencia_ads"`
	Activo        bool  `json:"activo"         gorm:"column:activo"`
}

// ActualizarConfigRequest body del PUT /espacio-comercial/config
type ActualizarConfigRequest struct {
	FrecuenciaAds int   `json:"frecuencia_ads" validate:"gt=0"`
	Activo        *bool `json:"activo"         validate:"required"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ─── Playlists ──────────────────────────────────────────────────────────────

// CrearPlaylistRequest body del POST /espacio-comercial/playlists
type CrearPlaylistRequest struct {
	Nombre string `json:"nombre" validate:"required,min=1"`
	Orden  int    `json:"orden"  validate:"omitempty"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ActualizarPlaylistRequest body del PUT /espacio-comercial/playlists/:id
type ActualizarPlaylistRequest struct {
	Nombre string `json:"nombre" validate:"required,min=1"`
	Orden  int    `json:"orden"  validate:"omitempty"`
	Activo *bool  `json:"activo" validate:"required"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// PlaylistResponse fila de espacio_comercial.cfg_playlists
type PlaylistResponse struct {
	ID     int64                   `json:"id"     gorm:"column:id"`
	Nombre string                  `json:"nombre" gorm:"column:nombre"`
	Orden  int                     `json:"orden"  gorm:"column:orden"`
	Activo bool                    `json:"activo" gorm:"column:activo"`
	Videos []PlaylistVideoResponse `json:"videos" gorm:"-"`
}

// PlaylistVideoResponse fila de espacio_comercial.cfg_playlist_videos
type PlaylistVideoResponse struct {
	ID         int64  `json:"id"          gorm:"column:id"`
	PlaylistID int64  `json:"playlist_id" gorm:"column:playlist_id"`
	YoutubeURL string `json:"youtube_url" gorm:"column:youtube_url"`
	Orden      int    `json:"orden"       gorm:"column:orden"`
}

// AgregarVideoRequest body del POST /espacio-comercial/playlists/:id/videos
type AgregarVideoRequest struct {
	YoutubeURL string `json:"youtube_url" validate:"required,url"`
	Orden      int    `json:"orden"       validate:"omitempty"`

	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ActualizarVideoRequest body del PUT /espacio-comercial/videos/:id
type ActualizarVideoRequest struct {
	YoutubeURL string `json:"youtube_url" validate:"required,url"`
	Orden      int    `json:"orden"       validate:"omitempty"`

	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ─── Ads ────────────────────────────────────────────────────────────────────

// AdResponse fila de espacio_comercial.cfg_ads
type AdResponse struct {
	ID          int64  `json:"id"           gorm:"column:id"`
	Nombre      string `json:"nombre"       gorm:"column:nombre"`
	ArchivoPath string `json:"-"            gorm:"column:archivo_path"`
	URL         string `json:"url"          gorm:"-"`
	Orden       int    `json:"orden"        gorm:"column:orden"`
	Activo      bool   `json:"activo"       gorm:"column:activo"`
	TamanoBytes int64  `json:"tamano_bytes" gorm:"column:tamano_bytes"`
	MimeType    string `json:"mime_type"    gorm:"column:mime_type"`
}

// ActualizarAdRequest body del PUT /espacio-comercial/ads/:id
type ActualizarAdRequest struct {
	Nombre string `json:"nombre" validate:"required,min=1"`
	Orden  int    `json:"orden"  validate:"omitempty"`
	Activo *bool  `json:"activo" validate:"required"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ─── Público (pantalla de reproducción) ──────────────────────────────────────

// PublicPlaybackResponse payload consumido por la pantalla pública sin login.
type PublicPlaybackResponse struct {
	FrecuenciaAds int                `json:"frecuencia_ads"`
	Playlists     []PlaylistResponse `json:"playlists"`
	Ads           []AdResponse       `json:"ads"`
}
