package uploads

import (
	"fmt"

	"github.com/ecosistema/core/src/core"
	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func UploadImageController(c *fiber.Ctx) error {

	// Obtener archivo
	file, err := c.FormFile("image")
	if err != nil {
		logging.Error.Printf("❌ No se recibió la imagen: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "missing_file",
			"message": "No se encontró la imagen en el campo 'image'",
		})
	}

	// Validar imagen
	if err := utils.ValidteImage(file); err != nil {
		logging.Error.Printf("❌ Validación fallida: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "validation_error",
			"message": err.Error(),
		})
	}

	// Guardar imagen
	filename, err := utils.SaveImage(file)
	if err != nil {
		logging.Error.Printf("❌ Error guardando imagen: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "save_failed",
			"message": "Error al guardar la imagen",
		})
	}

	url := fmt.Sprintf("%s/uploads/%s", core.Cfg.Backend_public_url, filename)

	//guardar datos en bd

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success":  true,
		"message":  "Imagen subida exitosamente",
		"url":      url,
		"path":     "/uploads/" + filename,
		"fileName": filename,
		"fileSize": file.Size,
		"mimeType": file.Header.Get("Content-Type"),
	})

}

func DeleteImageController(c *fiber.Ctx) error {

	filenameOrUrl := c.Params("filename")

	logging.Info.Printf("🗑️  Argumento recibido: %v", filenameOrUrl)

	// Extraer el filename si viene una URL completa
	filename := utils.ExtractFilenameFromUrl(filenameOrUrl)

	logging.Info.Printf("📂 Filename extraído: %v", filename)

	// Eliminar archivo físico
	err := utils.DeleteImage(filename)
	if err != nil {
		logging.Error.Printf("❌ Error eliminando imagen: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Error eliminando imagen"})
	}

	logging.Info.Printf("✅ Imagen eliminada correctamente: %s", filename)

	return c.JSON(fiber.Map{
		"message":  "Imagen eliminada",
		"filename": filename,
	})

}
