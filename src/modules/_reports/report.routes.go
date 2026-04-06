package reports

import (
	reports "github.com/ecosistema/core/src/modules/_reports/admin"
	"github.com/gofiber/fiber/v2"
)

func GestionReportsRoutes(r fiber.Router) {

	rpt := r.Group("/report")

	reports.RutasReport(rpt)

}
