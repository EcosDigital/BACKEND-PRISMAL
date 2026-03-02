package modular

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Rutas_Modular(r fiber.Router) {

	protected := r.Group("/",
		middlewares.JWTProtected,
		middlewares.EnrichContext)

	api := r.Group("/licence")

	api.Get("/modules/:codigo", FindModulesByCodeLicense)
	api.Get("/:codigo/modules/:id_modulo/functions", FindFuncionesByCodeLicence)
	api.Get("/:codigo/modules/functions/:id_funcion/subfunctions", FindSubFuncionesByCodeLicence)

	/*PRODUCT SOFTWARE*/
	protected.Post("/product",
		middlewares.VallidateBody(&ProductRequest{}), CreateProductSoftController)

	protected.Get("/product", FindLastProductsController)

	protected.Get("/product/:id", FindProductByIDController)

	protected.Put("/product/:id",
		middlewares.VallidateBody(&ProductUpdateRequest{}), ChangeProductSoftController)

	/*CATEGORY FOR PRODUCT SOFTWARE*/

	protected.Post("/category",
		middlewares.VallidateBody(&CategoryRequest{}), CreateCategoryController)

	protected.Get("/category", FindLastCategoryController)

	protected.Get("/category/:id", FindCategoryByIDController)

	protected.Get("/category/product/:id", FindCategoryByIDProductController)

	protected.Put("/category/:id",
		middlewares.VallidateBody(&CategoryUpdateRequest{}), ChangeCategoryController)

	/*MODULES FOR PRODUCT SOFTWARE*/
	protected.Post("/module",
		middlewares.VallidateBody(&ModuleRequest{}), CreateModuleController)

	protected.Get("/module", FindLastModuleController)

	protected.Get("/module/:id", FindModuleByIDController)

	protected.Get("/module/product/:id", FindModuleByIDProductController)

	protected.Put("/module/:id",
		middlewares.VallidateBody(&ModuleUpdateRequest{}), ChangeModuleController)

	protected.Get("/module/user/all", FindAllModuleController)

	//ADD
	protected.Get("/modules/rol", FindModulesByRolController)

	/*FUNCION FOR FUNCION*/
	protected.Post("/mod/funcion",
		middlewares.VallidateBody(&FuncionRequest{}), CreateFuncionController)

	protected.Get("/mod/funcions", FindLastFuncionController)

	protected.Get("/mod/funcion/:id", FindFuncionByIDController)

	protected.Get("/mod/funcions/:id", FindFuncionByIDModuleController)

	protected.Put("/mod/funcion/:id",
		middlewares.VallidateBody(&FuncionRequest{}), ChangeFuncionController)

	/*SUBFUNCION FOR FUNCION*/
	protected.Post("/mod/funcion/sub",
		middlewares.VallidateBody(&SubFuncionesRequest{}), CreateSubFuncionController)

	protected.Get("/mod/funcion/sub/:id", FindSubByIDFunController)

	//Add

}
