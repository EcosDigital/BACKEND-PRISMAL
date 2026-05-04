package scrapping

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

// RunScrapingJobController godoc
// POST /prospeccion/run
// Inicia un job de scraping sincrónico. Retorna el resumen al finalizar.
// Nota: el timeout de Fiber debe configurarse >= 15 min para jobs grandes.
func RunScrapingJobController(c *fiber.Ctx) error {

	var req ScrapingJobRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	if req.Limite == 0 {
		req.Limite = 50
	}

	userID := int64(middlewares.GetUserID(c))
	empresaID := int64(middlewares.GetEmpresaID(c))
	sedeID := int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := RunJob(db, &req, userID, empresaID, sedeID)
	if err != nil {
		logging.Error.Printf("[prospeccion] RunScrapingJobController: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(result)
}

// FindJobsController godoc
// GET /prospeccion/jobs
// Retorna el historial de jobs de la empresa autenticada.
func FindJobsController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	jobs, err := GetJobs(db, empresaID)
	if err != nil {
		logging.Error.Printf("[prospeccion] FindJobsController: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": jobs})
}

// FindLeadsController godoc
// GET /prospeccion/leads
// Retorna leads paginados con filtros opcionales por query params:
//
//	?id_job=1&ciudad=Montería&id_estado=1&page=1&limit=50
func FindLeadsController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	var filter LeadFilterRequest
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Parámetros inválidos"})
	}

	leads, err := GetLeads(db, filter, empresaID)
	if err != nil {
		logging.Error.Printf("[prospeccion] FindLeadsController: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": leads})
}

// ExportLeadsController godoc
// GET /prospeccion/export
// Genera y descarga un archivo Excel con los leads filtrados.
// Acepta los mismos query params que FindLeadsController excepto page/limit.
func ExportLeadsController(c *fiber.Ctx) error {

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	empresaID := int64(middlewares.GetEmpresaID(c))

	var filter LeadExportRequest
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Parámetros inválidos"})
	}

	rows, err := GetLeadsParaExportar(db, filter, empresaID)
	if err != nil {
		logging.Error.Printf("[prospeccion] ExportLeadsController: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(rows) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No hay leads para exportar con los filtros indicados"})
	}

	buf, err := GenerarExcelLeads(rows, filter.Keyword, filter.Ciudad)
	if err != nil {
		logging.Error.Printf("[prospeccion] ExportLeadsController - Excel: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generando el archivo Excel"})
	}

	filename := fmt.Sprintf("leads_%s.xlsx", time.Now().Format("20060102_150405"))

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return c.Send(buf.Bytes())
}

// RunBatchJobController godoc
// POST /prospeccion/run-batch
// Lanza scraping para todas las combinaciones keyword×ciudad automáticamente.
// Body vacío = KeywordsComercio × todos los municipios de Córdoba.
func RunBatchJobController(c *fiber.Ctx) error {

	var req BatchJobRequest
	// Body es opcional — si viene vacío usa los defaults
	c.BodyParser(&req) //nolint — fallo intencional si body vacío

	userID := int64(middlewares.GetUserID(c))
	empresaID := int64(middlewares.GetEmpresaID(c))
	sedeID := int64(middlewares.GetSedeID(c))

	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := RunBatchJob(db, &req, userID, empresaID, sedeID)
	if err != nil {
		logging.Error.Printf("[prospeccion] RunBatchJobController: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(result)
}
