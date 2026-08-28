package espacio_comercial

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// ─── Config ─────────────────────────────────────────────────────────────────

// FindConfigController GET /espacio-comercial/config
func FindConfigController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	userID := int64(middlewares.GetUserID(c))

	result, err := GetConfig(db, empresaID, userID)
	if err != nil {
		logging.Error.Printf("FindConfig error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}

// ChangeConfigController PUT /espacio-comercial/config
func ChangeConfigController(c *fiber.Ctx) error {
	var req ActualizarConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := EditConfig(db, &req)
	if err != nil {
		logging.Error.Printf("ChangeConfig error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Configuración actualizada exitosamente",
		"data":    result,
	})
}

// ─── Playlists ──────────────────────────────────────────────────────────────

// FindPlaylistsController GET /espacio-comercial/playlists
func FindPlaylistsController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := GetPlaylists(db, empresaID)
	if err != nil {
		logging.Error.Printf("FindPlaylists error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// CreatePlaylistController POST /espacio-comercial/playlists
func CreatePlaylistController(c *fiber.Ctx) error {
	var req CrearPlaylistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := NuevaPlaylist(db, &req)
	if err != nil {
		logging.Error.Printf("CreatePlaylist error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Playlist creada exitosamente",
		"id":      id,
	})
}

// ChangePlaylistController PUT /espacio-comercial/playlists/:id
func ChangePlaylistController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req ActualizarPlaylistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := EditPlaylist(db, id, &req); err != nil {
		logging.Error.Printf("ChangePlaylist error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Playlist actualizada exitosamente", "id": id})
}

// DeletePlaylistController DELETE /espacio-comercial/playlists/:id
func DeletePlaylistController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	if err := RemovePlaylist(db, id, empresaID); err != nil {
		logging.Error.Printf("DeletePlaylist error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Playlist eliminada exitosamente", "id": id})
}

// ─── Videos de playlist ───────────────────────────────────────────────────────

// CreateVideoController POST /espacio-comercial/playlists/:id/videos
func CreateVideoController(c *fiber.Ctx) error {
	playlistID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req AgregarVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := NuevoVideo(db, playlistID, &req)
	if err != nil {
		logging.Error.Printf("CreateVideo error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Video agregado exitosamente",
		"id":      id,
	})
}

// ChangeVideoController PUT /espacio-comercial/videos/:id
func ChangeVideoController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req ActualizarVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := EditVideo(db, id, &req); err != nil {
		logging.Error.Printf("ChangeVideo error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Video actualizado exitosamente", "id": id})
}

// DeleteVideoController DELETE /espacio-comercial/videos/:id
func DeleteVideoController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	if err := RemoveVideo(db, id, empresaID); err != nil {
		logging.Error.Printf("DeleteVideo error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Video eliminado exitosamente", "id": id})
}

// ─── Ads ────────────────────────────────────────────────────────────────────

// FindAdsController GET /espacio-comercial/ads
func FindAdsController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := GetAds(db, empresaID)
	if err != nil {
		logging.Error.Printf("FindAds error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// CreateAdController POST /espacio-comercial/ads (multipart: campo "video", "nombre", "orden")
func CreateAdController(c *fiber.Ctx) error {
	file, err := c.FormFile("video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No se encontró el video en el campo 'video'"})
	}

	nombre := c.FormValue("nombre")
	if nombre == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "El nombre es obligatorio"})
	}

	orden, _ := strconv.Atoi(c.FormValue("orden", "0"))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	userID := int64(middlewares.GetUserID(c))

	result, err := SubirAd(db, empresaID, userID, nombre, orden, file)
	if err != nil {
		logging.Error.Printf("CreateAd error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Video publicitario subido exitosamente",
		"data":    result,
	})
}

// ChangeAdController PUT /espacio-comercial/ads/:id
func ChangeAdController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req ActualizarAdRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if err := EditAd(db, id, &req); err != nil {
		logging.Error.Printf("ChangeAd error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Video publicitario actualizado exitosamente", "id": id})
}

// DeleteAdController DELETE /espacio-comercial/ads/:id
func DeleteAdController(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	if err := RemoveAd(db, id, empresaID); err != nil {
		logging.Error.Printf("DeleteAd error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Video publicitario eliminado exitosamente", "id": id})
}

// ─── Público (pantalla de reproducción) ──────────────────────────────────────

// FindPlaybackController GET /espacio-comercial/public
// Sin autenticación — consumido por la pantalla de reproducción (TV/monitor).
// El tenant se resuelve por el header X-Tenant-ID (subdominio), igual que
// cualquier otra ruta; no requiere ningún token en la URL.
func FindPlaybackController(c *fiber.Ctx) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := GetPlayback(db)
	if err != nil {
		logging.Error.Printf("FindPlayback error: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(result)
}
