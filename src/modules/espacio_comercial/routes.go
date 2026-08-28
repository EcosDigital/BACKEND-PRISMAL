package espacio_comercial

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// RegisterEspacioComercialRoutes registra todas las rutas del módulo Espacio
// Comercial. Se monta sobre /api/v1 desde el loader principal.
func RegisterEspacioComercialRoutes(r fiber.Router) {

	// ── Ruta pública (sin login) ───────────────────────────────
	// GET /espacio-comercial/public
	// Consumida por la pantalla de reproducción (TV/monitor). No lleva
	// JWTProtected, pero sí pasa por el middleware global de tenancy
	// (resuelve la DB por el header X-Tenant-ID del subdominio) — eso ya
	// es suficiente para identificar a qué tenant pertenece, sin token.
	public := r.Group("/espacio-comercial")
	public.Get("/public", FindPlaybackController)

	// ── Rutas protegidas (panel admin del tenant) ──────────────
	protected := r.Group("/espacio-comercial",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// Configuración
	protected.Get("/config", FindConfigController)
	protected.Put("/config",
		middlewares.VallidateBody(&ActualizarConfigRequest{}),
		ChangeConfigController,
	)

	// Playlists
	protected.Get("/playlists", FindPlaylistsController)
	protected.Post("/playlists",
		middlewares.VallidateBody(&CrearPlaylistRequest{}),
		CreatePlaylistController,
	)
	protected.Put("/playlists/:id",
		middlewares.VallidateBody(&ActualizarPlaylistRequest{}),
		ChangePlaylistController,
	)
	protected.Delete("/playlists/:id", DeletePlaylistController)

	// Videos de playlist
	protected.Post("/playlists/:id/videos",
		middlewares.VallidateBody(&AgregarVideoRequest{}),
		CreateVideoController,
	)
	protected.Put("/videos/:id",
		middlewares.VallidateBody(&ActualizarVideoRequest{}),
		ChangeVideoController,
	)
	protected.Delete("/videos/:id", DeleteVideoController)

	// Ads (videos publicitarios)
	protected.Get("/ads", FindAdsController)
	protected.Post("/ads", CreateAdController)
	protected.Put("/ads/:id",
		middlewares.VallidateBody(&ActualizarAdRequest{}),
		ChangeAdController,
	)
	protected.Delete("/ads/:id", DeleteAdController)
}
