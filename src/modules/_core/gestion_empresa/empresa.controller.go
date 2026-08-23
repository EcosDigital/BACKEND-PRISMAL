package empresa

import (
	"net/http"
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateEmpresaController(c *fiber.Ctx) error {

	var req EmpresaRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//invocar services
	newID, err := RegisterEmpresa(db, &req)
	if err != nil {

		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Registro exitoso...",
		"id":      newID,
	})
}

func FindLastEmpresaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := FilterLastEmpresa(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func FindEmpresaByIDController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})
	}

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := FilterEmpresaByID(db, id)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tercero no encontrado"})
	}

	return c.JSON(result)

}

func ChangeEmpresaController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})
	}

	var req EmpresaUpdateRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//services (logica negocio)
	uptID, err := EditEmpresa(db, id, &req)
	if err != nil {

		logging.Error.Printf("Hubo un error: %v", err)

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Registro actualizado...",
		"id":      uptID,
	})

}

//sedes

func CreateSedeController(c *fiber.Ctx) error {

	var req SedeRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//invocar services
	newID, err := RegisterSede(db, &req)
	if err != nil {

		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Registro exitoso...",
		"id":      newID,
	})
}

func FindLastSedeController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := FilterLsatSede(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func FindSedeByIDController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})
	}

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := FilterSedeByID(db, id)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tercero no encontrado"})
	}

	return c.JSON(result)

}

func ChangeSedeController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})
	}

	var req SedeUpdateRequest

	//parsear JSON de la request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": "Json invalido",
		})
	}

	req.UserID = int64(middlewares.GetUserID(c))
	req.EmpresaID = int64(middlewares.GetEmpresaID(c))
	req.SedeID = int64(middlewares.GetSedeID(c))

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	tenantSlug := middlewares.GetTenantSlug(c)

	//services (logica negocio)
	uptID, syncWarning, err := EditSede(db, tenantSlug, id, &req)
	if err != nil {

		logging.Error.Printf("Hubo un error: %v", err)

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":      "Registro actualizado...",
		"id":           uptID,
		"sync_warning": syncWarning,
	})

}

// ─── Imagen de sede ──────────────────────────────────────────────────────────

// SetSedeImagenController sube/reemplaza la imagen de una sede.
//
// POST /core/sedes/:id/imagen (multipart, campo "image")
func SetSedeImagenController(c *fiber.Ctx) error {

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

	tenantSlug := middlewares.GetTenantSlug(c)

	imagenURL, syncWarning, err := SetSedeImagen(db, tenantSlug, id, file)
	if err != nil {
		logging.Error.Printf("Error subiendo imagen de sede: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Imagen actualizada exitosamente",
		"id":           id,
		"imagen_url":   imagenURL,
		"sync_warning": syncWarning,
	})
}

// ─── Horario de atención por sede ───────────────────────────────────────────

// ListSedeHorariosController lista los días configurados de una sede.
//
// GET /core/sedes/:id/horarios
func ListSedeHorariosController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	results, err := FilterSedeHorarios(db, id)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(results)
}

// SaveSedeHorarioController crea o actualiza el horario de un día puntual.
//
// PUT /core/sedes/:id/horarios/:dia (dia = 0..6, 0=domingo)
func SaveSedeHorarioController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de sede inválido"})
	}

	diaParam := c.Params("dia")
	dia, err := strconv.Atoi(diaParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Día inválido"})
	}

	var req SedeHorarioRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	req.UserID = int64(middlewares.GetUserID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	tenantSlug := middlewares.GetTenantSlug(c)

	result, syncWarning, err := SaveSedeHorario(db, tenantSlug, id, dia, &req)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":      "Horario actualizado exitosamente",
		"data":         result,
		"sync_warning": syncWarning,
	})
}
