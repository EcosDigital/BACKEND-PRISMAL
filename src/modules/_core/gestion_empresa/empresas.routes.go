package empresa

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_empresa(r fiber.Router) {

	protected := r.Group("/company",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	protected.Post("/",
		middlewares.VallidateBody(&EmpresaRequest{}), CreateEmpresaController)

	protected.Get("/", FindLastEmpresaController)

	protected.Get("/:id",
		FindEmpresaByIDController)

	protected.Put("/:id",
		middlewares.VallidateBody(&EmpresaUpdateRequest{}), ChangeEmpresaController)

}
