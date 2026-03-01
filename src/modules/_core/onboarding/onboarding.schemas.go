package onboarding

type EmpresaPayload struct {
	RazonSocial        string `json:"razon_social" validate:"required,min=3"`
	RepresentanteLegal string `json:"representante_legal" validate:"omitempty,min=2,max=200"`
	Nit                string `json:"nit" validate:"required"`
	Dv                 string `json:"dv" validate:"omitempty,len=1"`
	Telefono           string `json:"telefono"`
	EmailContacto      string `json:"email_contacto" validate:"required,email"`
	Subdominio         string `json:"subdominio" validate:"required,min=1"`
	IdDepartamento     int    `json:"id_departamento" validate:"required,gt=0"`
	IdCiudad           int    `json:"id_ciudad"`
}

type AdminPayload struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	NumeroDocumento string `json:"numero_documento" validate:"required,min=5,max=20"`
	PrimerNombre    string `json:"primer_nombre" validate:"omitempty,alphaunicode,min=2,max=30"`
	SegundoNombre   string `json:"segundo_nombre" validate:"omitempty,alphaunicode,max=30"`
	PrimerApellido  string `json:"primer_apellido" validate:"omitempty,alphaunicode,min=2,max=30"`
	SegundoApellido string `json:"segundo_apellido" validate:"omitempty,alphaunicode,min=2,max=30"`
}

type TenantRequest struct {
	Empresa EmpresaPayload `json:"empresa" validate:"required"`
	Admin   AdminPayload   `json:"admin" validate:"required"`
	Modulos []string       `json:"modulos" validate:"required,min=1"`
}

type ModuleInfo struct {
	ID   int
	Code string
	Name string
}

type ModuleCatalogItem struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	Icono       string `json:"icono"`
	Color       string `json:"color"`
	BgColor     string `json:"bg_color"`
	BorderColor string `json:"border_color"`
}

type ModuleCatalogCategory struct {
	CategoryID          int64               `json:"id_categoria" gorm:"column:id_categoria"`
	CategoryName        string              `json:"nombre_categoria" gorm:"column:nombre_categoria"`
	CategoryDescripcion string              `json:"descripcion_categoria" gorm:"column:descripcion_categoria"`
	Modules             []ModuleCatalogItem `json:"modulos" gorm:"-"`
}
