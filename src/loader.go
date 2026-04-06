package loader

import (
	"github.com/ecosistema/core/src/modules/RRHH"
	core "github.com/ecosistema/core/src/modules/_core"
	reports "github.com/ecosistema/core/src/modules/_reports"
	"github.com/ecosistema/core/src/modules/inventario"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	//core
	core.RegisterCoreRoutes(v1)

	//contabilidad

	//inventario
	inventario.RegisterInvetarioRoutes(v1)

	//RRHH
	RRHH.RegisterRRHHRoutes(v1)

	//REPORTES
	reports.GestionReportsRoutes(v1)
}
