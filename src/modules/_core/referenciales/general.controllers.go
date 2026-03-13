package referenciales

import (
	"strconv"

	"github.com/ecosistema/core/src/shared/logging"
	"github.com/ecosistema/core/src/shared/utils"
	"github.com/gofiber/fiber/v2"
)

func GetTipoPersonaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetTipoPersona(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetTipoDocumentController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetTipoDocumento(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetGeneroController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetGeneros(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetDepartamentController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetDepartamentos(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetMunicipioController(c *fiber.Ctx) error {

	idParam := c.Params("id")
	id_dep, err := strconv.Atoi(idParam)
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

	result, err := GetMunicipios(db, id_dep)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetZonasRuralesController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetZonaRurales(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetActividadesController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetActividadesEco(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetAmbitosTerceroController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetAmbitosTerceros(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetCentralizacionesController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetCentralizaciones(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetResponsabilidadesDianController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetResponsabilidadDian(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetTipoContribuyenteController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetTipoContribuyente(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetRegimenIvaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetRegimenIva(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetRegimenDianController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetRegimenDian(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetClaseTercerosController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetClaseTercero(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetEstadoModuloController(c *fiber.Ctx) error {
	results, err := GetEstadoModule()
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(results)
}

func GetTipoEmpresaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := GetTipoEmpresa(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(results)

}

func GetPaisesController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := GetPaises(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(results)

}

func GetNaturalezaEmpresaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := GetNaturalezaEmpresa(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(results)

}

func GetTipoSedeController(c *fiber.Ctx) error {
	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	results, err := GetTipoSede(db)
	if err != nil {
		logging.Error.Printf("Hubo un error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(results)
}

func GetTipoRolesController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetTipoRol(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetModulesComprobanteController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetModulosComprobantes(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

//==== GESTIONES ========//

func GetNivelesTicketController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetNivelesCaso(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetTipoOperacionController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})
	}

	result, err := GetTipoOperacion(db, id)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetEstadoTicketController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetEstadosCaso(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetColaboradoresController(c *fiber.Ctx) error {
	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetColaboradores(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

//===== INVENTARIO ====== //

func GetTipoBodegaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetTipoBodega(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetUnidadesMedidaController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetUnidadesMedida(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetGrupoArticulosController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetGrupoArticulos(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}

func GetPresentacionArticulosController(c *fiber.Ctx) error {

	//obtiene base de datos
	db, err := utils.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := GetPresentacionArticulos(db)
	if err != nil {
		logging.Info.Printf("Hubo un error %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(result)
}
