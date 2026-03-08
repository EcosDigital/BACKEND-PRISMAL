package modular

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func RegisterProductoSoftware(db *gorm.DB, req *ProductRequest) (int64, error) {

	// Verifica si existe registro
	exists, err := ListProductByCode(db, req.Codigo)
	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateProductSoftware(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func FilterLastProductsSoftware(db *gorm.DB) ([]ProductsResponse, error) {
	return ListProductLast(db)
}

func FilterProductsByID(db *gorm.DB, id int64) ([]ProductsResponse, error) {
	return ListProductByID(db, id)
}

func EditProductSofware(db *gorm.DB, req *ProductUpdateRequest, id int64) (int64, error) {

	//verificar si existe el registro por ID
	exists, err := ListProductByID(db, id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No se encontraron registros para actualizar")
	}

	//actualizar registro
	regID, err := UpdateProduct(db, req, id)
	if err != nil {
		return 0, err
	}

	return regID, err

}

// * CATEGORIAS FOR PRODUCT SOFTWARE * //
func RegisterCategory(db *gorm.DB, req *CategoryRequest) (int64, error) {

	//verificar existencia por codigo
	exists, err := ListCategoryByCode(db, req.Codigo)
	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateCategoryForProduct(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err
}

func FilterLastCategory(db *gorm.DB) ([]CategoryResponse, error) {
	return ListCategoryLast(db)
}

func FilterCategoryByIdProduct(db *gorm.DB, id int64) ([]CategoryResponse, error) {
	return ListCategoryByIdProduct(db, id)
}

func FilerCategoryByID(db *gorm.DB, id int64) ([]CategoryResponse, error) {
	return ListCategoryById(db, id)
}

func EditCategory(db *gorm.DB, req *CategoryUpdateRequest, id int64) (int64, error) {

	//verificar registro por ID
	exists, err := ListCategoryById(db, id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("No Se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := UpdateCategory(db, req, id)
	if err != nil {
		return 0, err
	}

	return regID, err

}

// * MODULOS * //
func RegisterModule(db *gorm.DB, req *ModuleRequest) (int64, error) {

	//verificar existencia por codigo
	exists, err := ListModuleByCode(db, req.Codigo)
	if err != nil {
		return 0, err
	}

	if len(exists) > 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//insertar registro
	regID, err := CreateModule(db, req)
	if err != nil {
		return 0, err
	}

	//guardar dependencias
	if len(req.Dependencias) > 0 {
		if err := AddDependencias(db, regID, req.Dependencias, req.UserID); err != nil {
			return 0, fmt.Errorf("error guardando dependencias: %v", err)
		}
	}

	return regID, nil

}

func FilterLastModule(db *gorm.DB) ([]ModuleResponse, error) {
	return ListModuleLast(db)
}

func FilterAllModule(db *gorm.DB) ([]ModuleResponse, error) {
	return ListModuleAll(db)
}

func FilterModuleByID(db *gorm.DB, id int64) ([]ModuleResponse, error) {
	return ListModuleById(db, id)
}

func FilterModuleByIDProduct(db *gorm.DB, id int64) ([]ModuleResponse, error) {
	return ListModuleByIdProduct(db, id)
}

func EditModule(db *gorm.DB, req *ModuleUpdateRequest, id int64) (int64, error) {

	//verificar registro por ID
	exists, err := ListModuleById(db, id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := UpdateModule(db, req, id)
	if err != nil {
		return 0, err
	}

	if err := SyncDependencias(db, id, req.Dependencias, req.UserID); err != nil {
		return 0, fmt.Errorf("error sincronizando dependencias: %v", err)
	}

	return regID, err

}

// ** Funciones **//
func RegisterFuncion(db *gorm.DB, req *FuncionRequest) (int64, error) {

	//insertar registro
	regID, err := CreateFuncion(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func FilterLastFuncion(db *gorm.DB) ([]FuncionResponse, error) {
	return ListFuncionLast(db)
}

func FilterFuncionByID(db *gorm.DB, id int64) ([]FuncionResponse, error) {
	return ListFuncionByID(db, id)
}

func FilterFuncionByIDModule(db *gorm.DB, id int64) ([]FuncionesResponse, error) {
	return ListFuncionByIDModule(db, id)
}

func EditFuncion(db *gorm.DB, req *FuncionRequest, id int64) (int64, error) {

	//verificar registro por ID
	exists, err := ListFuncionByID(db, id)

	if err != nil {
		return 0, err
	}

	if len(exists) <= 0 {
		// Existe registro con número de documento
		return 0, errors.New("Se encontraron resultados con este codigo de registro")
	}

	//actualizar registro
	regID, err := UpdateFuncion(db, req, id)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func RegisterSubFuncion(db *gorm.DB, req *SubFuncionesRequest) (int64, error) {

	//insertar registro
	regID, err := CreateSubFuncion(db, req)
	if err != nil {
		return 0, err
	}

	return regID, err

}

func FilterSubByIDFunc(db *gorm.DB, id int64) ([]SubFuncionesRequest, error) {
	return ListSubByIDFuncion(db, id)
}

// ** ADD **/
func FilterModulesByCodeLicence(codigo string) ([]ModuleResponse, error) {
	return ListModulesByCodeLicence(codigo)
}

func FilterModulesByRol(db *gorm.DB, id int64) (*AccessReponse, error) {
	return ListModuleByIdRol(db, id)
}
