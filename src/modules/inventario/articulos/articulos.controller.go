package articulos

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateArticuloController(c *fiber.Ctx) error {

	var req ArticuloRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	newID, err := RegisterArticulo(db, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Artículo registrado exitosamente",
		"id":      newID,
	})
}

func FindArticulosController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := FilterArticulos(db, empresaID)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

func FindArticuloByIDController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	result, err := FilterArticuloByID(db, id)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result == nil || result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Artículo no encontrado"})
	}

	return c.JSON(result)
}

func ChangeArticuloController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req ArticuloUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	uptID, err := EditArticulo(db, id, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Artículo actualizado exitosamente",
		"id":      uptID,
	})
}

func SearchArticulosController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))
	nombreCodigo := c.Query("q", "")

	results, err := FilterSearchArticulos(db, nombreCodigo, empresaID)
	if err != nil {
		logging.Error.Printf("Error buscando artículos: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": results})
}

// ─── Imagen de producto ─────────────────────────────────────────────────────

// SetArticuloImagenController sube/reemplaza la imagen de un artículo.
//
// POST /inventario/articulos/:id/imagen (multipart, campo "image")
func SetArticuloImagenController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No se encontró la imagen en el campo 'image'"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	userID := int64(middlewares.GetUserID(c))

	imagenURL, err := SetArticuloImagen(db, id, file, userID)
	if err != nil {
		logging.Error.Printf("Error subiendo imagen de artículo: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":    "Imagen actualizada exitosamente",
		"id":         id,
		"imagen_url": imagenURL,
	})
}

// ─── Carga masiva ─────────────────────────────────────────────────────────────

// ImportArticulosController recibe el JSON con las filas del plano Excel,
// resuelve las referencias por nombre, y hace upsert de cada artículo.
//
// POST /inventario/articulos/import
// Body: { "filas": [ { "CODIGO": "...", "NOMBRE": "...", ... } ] }
func ImportArticulosController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var req ImportArticulosRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	if len(req.Filas) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "El plano no contiene filas"})
	}

	userID := int64(middlewares.GetUserID(c))
	empresaID := int64(middlewares.GetEmpresaID(c))
	sedeID := int64(middlewares.GetSedeID(c))

	result, err := ImportArticulos(db, &req, userID, empresaID, sedeID)
	if err != nil {
		logging.Error.Printf("Error en carga masiva de artículos: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(result)
}
