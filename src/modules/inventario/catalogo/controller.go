package catalogo

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateCatalogoController(c *fiber.Ctx) error {

	var req CatalogoRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	newID, err := RegisterCatalogo(db, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Catálogo creado exitosamente",
		"id":      newID,
	})
}

func FindCatalogosController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	results, err := FilterCatalogos(db, empresaID)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

func FindCatalogoByIDController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	result, err := FilterCatalogoByID(db, id)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result == nil || result.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Catálogo no encontrado"})
	}

	return c.JSON(result)
}

func ChangeCatalogoController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CatalogoUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	tenantSlug := middlewares.GetTenantSlug(c)

	uptID, syncWarning, err := EditCatalogo(db, tenantSlug, id, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Catálogo actualizado exitosamente",
		"id":           uptID,
		"sync_warning": syncWarning,
	})
}

// ─── Detalle: artículos publicados dentro del catálogo ──────────────────────

func CreateCatalogoArticuloController(c *fiber.Ctx) error {

	catalogoID, err := strconv.ParseInt(c.Params("id_catalogo"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de catálogo inválido"})
	}

	var req CatalogoArticuloRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	tenantSlug := middlewares.GetTenantSlug(c)

	newID, syncWarning, err := RegisterCatalogoArticulo(db, tenantSlug, catalogoID, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":      "Artículo publicado exitosamente",
		"id":           newID,
		"sync_warning": syncWarning,
	})
}

func FindCatalogoArticulosController(c *fiber.Ctx) error {

	catalogoID, err := strconv.ParseInt(c.Params("id_catalogo"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de catálogo inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterCatalogoArticulos(db, catalogoID)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

func ChangeCatalogoArticuloController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var req CatalogoArticuloUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	tenantSlug := middlewares.GetTenantSlug(c)

	uptID, syncWarning, err := EditCatalogoArticulo(db, tenantSlug, id, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Publicación actualizada exitosamente",
		"id":           uptID,
		"sync_warning": syncWarning,
	})
}

func DeleteCatalogoArticuloController(c *fiber.Ctx) error {

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	tenantSlug := middlewares.GetTenantSlug(c)

	syncWarning, err := RemoveCatalogoArticulo(db, tenantSlug, id)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Artículo despublicado exitosamente",
		"sync_warning": syncWarning,
	})
}
