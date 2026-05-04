package scrapping

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

// Rutas_prospeccion registra todas las rutas del módulo bajo /prospeccion.
// Sigue exactamente el mismo patrón que Rutas_articulos.
func Rutas_prospeccion(r fiber.Router) {

	protected := r.Group("/prospeccion",
		middlewares.JWTProtected,
		middlewares.EnrichContext,
	)

	// Iniciar un job de scraping
	protected.Post("/run",
		middlewares.VallidateBody(&ScrapingJobRequest{}),
		RunScrapingJobController,
	)

	// Historial de jobs
	protected.Get("/jobs",
		FindJobsController,
	)

	// Listado paginado de leads con filtros
	protected.Get("/leads",
		FindLeadsController,
	)

	// Exportar leads a Excel
	protected.Post("/run-batch", RunBatchJobController)

	protected.Get("/export",
		ExportLeadsController,
	)
}
