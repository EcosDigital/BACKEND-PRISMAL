package referenciales

import (
	"github.com/ecosistema/core/src/shared/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(r fiber.Router) {

	api := r.Group("/ref") // rutas bajo (/api)

	//rutas generales
	api.Get("/TipoPersona", middlewares.JWTProtected, GetTipoPersonaController)
	api.Get("/TipoDocument", middlewares.JWTProtected, GetTipoDocumentController)
	api.Get("/Generos", middlewares.JWTProtected, GetGeneroController)
	api.Get("/Departamentos", middlewares.JWTProtected, GetDepartamentController)
	api.Get("/Municipios/:id", middlewares.JWTProtected, GetMunicipioController)
	api.Get("/ZonasResidenciales", middlewares.JWTProtected, GetZonasRuralesController)
	api.Get("/ActividadesEconomic", middlewares.JWTProtected, GetActividadesController)
	api.Get("/Ambitos", middlewares.JWTProtected, GetAmbitosTerceroController)
	api.Get("/Centralizacion", middlewares.JWTProtected, GetCentralizacionesController)
	api.Get("/ResponsabilidadDian", middlewares.JWTProtected, GetResponsabilidadesDianController)
	api.Get("/TipoContribuyente", middlewares.JWTProtected, GetTipoContribuyenteController)
	api.Get("/RegimenIva", middlewares.JWTProtected, GetRegimenIvaController)
	api.Get("/RegimenDian", middlewares.JWTProtected, GetRegimenDianController)
	api.Get("/ClaseTerceros", middlewares.JWTProtected, GetClaseTercerosController)
	api.Get("/EstadoModulo", middlewares.JWTProtected, GetEstadoModuloController)
	api.Get("/TipoEmpresa", middlewares.JWTProtected, GetTipoEmpresaController)
	api.Get("/Paises", middlewares.JWTProtected, GetPaisesController)
	api.Get("/NatEmpresa", middlewares.JWTProtected, GetNaturalezaEmpresaController)
	api.Get("/TipoSedes", middlewares.JWTProtected, GetTipoSedeController)
	api.Get("/TipoRoles", middlewares.JWTProtected, GetTipoRolesController)

	api.Get("/comprobantes/modulos", middlewares.JWTProtected, GetModulesComprobanteController)
	api.Get("/comprobantes/operaciones/:id", middlewares.JWTProtected, GetTipoOperacionController)

	api.Get("/gestiones/niveles", middlewares.JWTProtected, GetNivelesTicketController)
	api.Get("/gestiones/estados", middlewares.JWTProtected, GetEstadoTicketController)
	api.Get("/gestiones/colaboradores", middlewares.JWTProtected, GetColaboradoresController)

	api.Get("/inventario/tipo-bodega", middlewares.JWTProtected, GetTipoBodegaController)
	api.Get("/inventario/unidades-medida", middlewares.JWTProtected, GetUnidadesMedidaController)
	api.Get("/inventario/grupos-articulo", middlewares.JWTProtected, GetGrupoArticulosController)
	api.Get("/articulos/presentaciones", middlewares.JWTProtected, GetPresentacionArticulosController)

}
