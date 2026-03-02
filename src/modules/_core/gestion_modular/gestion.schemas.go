package modular

import "encoding/json"

//** STRUCT'S FOR PRODUCTOS DE SOFTWARE **/

type ProductRequest struct {
	Codigo      string `json:"codigo" validate:"required,min=5,max=12"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	Logo        string `json:"logo" validate:"omitempty,url,max=300"`
	Verssion    string `json:"verssion" validate:"required,min=1"`
	IsActive    *bool  `json:"is_active" validate:"required"`
	UserID      int64  `json:"user_id" validate:"omitempty"`
	EmpresaID   int64  `json:"empresa_id" validate:"omitempty"`
	SedeID      int64  `json:"sede_id" validate:"omitempty"`
}

type ProductsResponse struct {
	ID          int    `json:"id"`
	Codigo      string `json:"codigo"`
	Logo        string `json:"logo" validate:"omitempty,url,max=300"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Verssion    string `json:"verssion"`
	IsActive    *bool  `json:"is_active"`
}

type ProductUpdateRequest struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Logo        string `json:"logo"`
	Verssion    string `json:"verssion"`
	IsActive    *bool  `json:"is_active"`
	UserID      int64  `json:"user_id,omitempty"`
	EmpresaID   int64  `json:"empresa_id" validate:"omitempty"`
	SedeID      int64  `json:"sede_id" validate:"omitempty"`
}

//** STRUCT'S FOR CATEGORIAS **/

type CategoryRequest struct {
	IDProduct   int    `json:"id_producto" validate:"required,gt=0"`
	Codigo      string `json:"codigo" validate:"required,min=5,max=12"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	IsActive    *bool  `json:"is_active" validate:"required"`
	UserID      int64  `json:"user_id" validate:"omitempty"`
	EmpresaID   int64  `json:"empresa_id" validate:"omitempty"`
	SedeID      int64  `json:"sede_id" validate:"omitempty"`
}

type CategoryResponse struct {
	ID          int    `json:"id"`
	IDProduct   int    `json:"id_producto"`
	Producto    string `json:"producto"`
	IDCategory  int    `json:"id_categoria"`
	Categoty    string `json:"categria"`
	Codigo      string `json:"codigo"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	IsActive    *bool  `json:"is_active" validate:"required"`
}

type CategoryUpdateRequest struct {
	IDProduct   int    `json:"id_producto" validate:"required,gt=0"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	IsActive    *bool  `json:"is_active"`
	UserID      int64  `json:"user_id" validate:"omitempty"`
	EmpresaID   int64  `json:"empresa_id" validate:"omitempty"`
	SedeID      int64  `json:"sede_id" validate:"omitempty"`
}

//** STRUCT'S FOR MODULES **/

type ModuleRequest struct {
	Codigo      string `json:"codigo" validate:"required,min=5,max=12"`
	IDProduct   int    `json:"id_producto" validate:"required,gt=0"`
	IDCategory  int    `json:"id_category" validate:"required,gt=0"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Ordenista   int    `json:"orden_lista" validate:"required,gt=0"`
	IdEstado    int    `json:"id_estado" validate:"required,gt=0"`
	Color       string `json:"color" validate:"required,min=3"`
	BgColor     string `json:"bg_color" validate:"required,min=3"`
	BorderColor string `json:"border_color" validate:"required,min=3"`
	Icono       string `json:"icono" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	EsInterno   *bool  `json:"es_interno" validate:"required"`
	IsActive    *bool  `json:"is_active" validate:"required"`
	UserID      int64  `json:"user_id" validate:"omitempty"`
	EmpresaID   int64  `json:"empresa_id" validate:"omitempty"`
	SedeID      int64  `json:"sede_id" validate:"omitempty"`
}

type ModuleResponse struct {
	ID          int    `json:"id"`
	IDProduct   int    `json:"id_producto"`
	Producto    string `json:"producto"`
	IDCategory  int    `json:"id_category" validate:"required,gt=0"`
	Categoria   string `json:"categoria"`
	Codigo      string `json:"codigo"`
	IdEstado    int    `json:"id_estado" validate:"required,gt=0"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	Color       string `json:"color" validate:"required,min=3"`
	BgColor     string `json:"bg_color" validate:"required,min=3"`
	BorderColor string `json:"border_color" validate:"required,min=3"`
	Icono       string `json:"icono" validate:"required,min=3"`
	IsActive    *bool  `json:"is_active" validate:"required"`
	EsInterno   *bool  `json:"es_interno" validate:"required"`
}

type ModuleUpdateRequest struct {
	IDProduct   int    `json:"id_producto" validate:"required,gt=0"`
	IDCategory  int    `json:"id_category" validate:"required,gt=0"`
	IdEstado    int    `json:"id_estado" validate:"required,gt=0"`
	Nombre      string `json:"nombre" validate:"required,min=3"`
	Descripcion string `json:"descripcion" validate:"required,min=20"`
	Ordenista   int    `json:"orden_lista" validate:"required,gt=0"`
	EsInterno   *bool  `json:"es_interno" validate:"required"`
	Color       string `json:"color" validate:"required,min=3"`
	BgColor     string `json:"bg_color" validate:"required,min=3"`
	BorderColor string `json:"border_color" validate:"required,min=3"`
	Icono       string `json:"icono" validate:"required,min=3"`
	IsActive    *bool  `json:"is_active" validate:"required"`
	UserID      int64  `json:"user_id" validate:"omitempty"`
	EmpresaID   int64  `json:"empresa_id" validate:"omitempty"`
	SedeID      int64  `json:"sede_id" validate:"omitempty"`
}

type FuncionRequest struct {
	IDProduct  int    `json:"id_producto" validate:"required,gt=0"`
	IDModulo   int    `json:"id_modulo" validate:"required,gt=0"`
	Nombre     string `json:"nombre" validate:"required,min=3"`
	OrdenLista int    `json:"orden_lista" validate:"gt=0"`
	RutaAcceso string `json:"ruta_acceso" validate:"omitempty"`
	IsActive   *bool  `json:"is_active" validate:"required"`
	UserID     int64  `json:"user_id" validate:"omitempty"`
	EmpresaID  int64  `json:"empresa_id" validate:"omitempty"`
	SedeID     int64  `json:"sede_id" validate:"omitempty"`
}

type FuncionResponse struct {
	ID         int    `json:"id"`
	IDProduct  int    `json:"id_producto"`
	Producto   string `json:"producto"`
	IDModulo   int    `json:"id_modulo"`
	Modulo     string `json:"modulo"`
	Nombre     string `json:"nombre"`
	OrdenLista int    `json:"orden_lista" validate:"gt=0"`
	RutaAcceso string `json:"ruta_acceso" validate:"required,min=3"`
	IsActive   *bool  `json:"is_active"`
}

type FuncionesResponse struct {
	ID       int                 `json:"id"`
	IDModulo int                 `json:"id_modulo"`
	Title    string              `json:"title"`
	Path     string              `json:"path"`
	Children []FuncionesResponse `json:"children,omitempty" gorm:"-"`
}

type SubFuncionesRequest struct {
	ID         int    `json:"id"`
	IDFuncion  int    `json:"id_funcion" validate:"required,gt=0"`
	Nombre     string `json:"nombre"`
	RutaAcceso string `json:"ruta_acceso" validate:"required,min=3"`
	UserID     int64  `json:"user_id" validate:"omitempty"`
	EmpresaID  int64  `json:"empresa_id" validate:"omitempty"`
	SedeID     int64  `json:"sede_id" validate:"omitempty"`
}

type AccessReponse struct {
	ID         int             `json:"id"`
	IDRol      int             `json:"id_rol"`
	JsonAccess json.RawMessage `json:"json_access"`
}
