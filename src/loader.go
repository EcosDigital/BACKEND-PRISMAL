package loader

import (
	"github.com/ecosistema/core/src/modules/RRHH"
	core "github.com/ecosistema/core/src/modules/_core"
	internal "github.com/ecosistema/core/src/modules/_functions"
	reports "github.com/ecosistema/core/src/modules/_reports"
	"github.com/ecosistema/core/src/modules/contabilidad"
	"github.com/ecosistema/core/src/modules/inventario"
	"github.com/ecosistema/core/src/modules/ventas"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	//core
	core.RegisterCoreRoutes(v1)

	//contabilidad
	contabilidad.RegisterContabilidadRoutes(v1)

	//inventario
	inventario.RegisterInvetarioRoutes(v1)
	//ventas
	ventas.RegisterVentasRoutes(v1)

	//RRHH
	RRHH.RegisterRRHHRoutes(v1)

	//INTERNAL (Modulos de uso interno)
	internal.RegisterFunctionsRoutes(v1)

	//REPORTES
	reports.GestionReportsRoutes(v1)
}
