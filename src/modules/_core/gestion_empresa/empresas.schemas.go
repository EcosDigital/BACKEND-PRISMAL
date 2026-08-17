package empresa

type EmpresaRequest struct {
	IdTipoEmpresa          int    `json:"id_tipo_empresa" validate:"required,gt=0"`
	Nit                    string `json:"nit" validate:"required,min=5,max=20"`
	Dv                     string `json:"dv" validate:"required,len=1"`
	RazonSocial            string `json:"razon_social" validate:"required,min=2,max=150"`
	Descripcion            string `json:"descripcion" validate:"omitempty,min=2,max=150"`
	Direccion              string `json:"direccion" validate:"required,min=5,max=150"`
	IdPais                 int    `json:"id_pais" valdate:"required,gt=0"`
	IdDepartamento         int    `json:"id_departamento" validate:"required,gt=0"`
	Id_ciudad              int    `json:"id_ciudad" validate:"required,gt=0"`
	IdZona                 int    `json:"id_zona" valdate:"required,gt=0"`
	Telefono               string `json:"telefono" validate:"required,numeric,min=7,max=15"`
	Telefono_2             string `json:"telefono_2" validate:"omitempty,numeric,min=7,max=15"`
	Email                  string `json:"email" validate:"required,email,max=100"`
	Fax                    string `json:"fax" validate:"omitempty,max=70"`
	PaginaWeb              string `json:"pagina_web" validate:"omitempty,url,max=300"`
	LogoUrl                string `json:"logo_url" validate:"omitempty,url,max=300"`
	UsaSedes               *bool  `json:"usa_sedes" validate:"required"`
	NumeroSedes            int    `json:"numero_sedes" valdate:"required"`
	IdNaturaleza           int    `json:"id_naturaleza" validate:"required,gt=0"`
	IdActividadEconomica   int    `json:"id_actividad_economica" validate:"required,gt=0"`
	IdTipoContribuyente    int    `json:"id_tipo_contribuyente" validate:"required,gt=0"`
	IdRegimenTributario    int    `json:"id_regimen_tributario" validate:"required,gt=0"`
	IdTipoDocRepresentante int    `json:"id_tipo_doc_representante" validate:"required,gt=0"`
	NumeroDocumento        string `json:"numero_documento" validate:"required,min=5,max=20"`
	RepresentanteLegal     string `json:"representante_legal" validate:"omitempty,min=2,max=200"`
	CodigoLicencia         string `json:"codigo_licencia" validate:"omitempty,min=2,max=200"`
	Estado                 *bool  `json:"is_active" validate:"required"`
	IaEndpoint             string `json:"ia_endpoint" validate:"omitempty,url,max=300"`
	IaToken                string `json:"ia_token" validate:"omitempty,max=1000"`
	UserID                 int64  `json:"user_id" validate:"omitempty"`
	EmpresaID              int64  `json:"empresa_id" validate:"omitempty"`
	SedeID                 int64  `json:"sede_id" validate:"omitempty"`
}

type EmpresaResponse struct {
	ID          int    `json:"id"`
	Nit         string `json:"nit"`
	Dv          string `json:"dv"`
	RazonSocial string `json:"razon_social"`
	UsaSedes    string `json:"usa_sedes"`
	Tipo        string `json:"tipo"`
	IsActive    bool   `json:"is_active"`
}

type EmpresaResponseFull struct {
	ID                     int    `json:"id"`
	IdTipoEmpresa          int    `json:"id_tipo_empresa"`
	Nit                    string `json:"nit"`
	Dv                     string `json:"dv"`
	RazonSocial            string `json:"razon_social"`
	Descripcion            string `json:"descripcion"`
	Direccion              string `json:"direccion"`
	IdPais                 int    `json:"id_pais"`
	IdDepartamento         int    `json:"id_departamento"`
	Id_ciudad              int    `json:"id_ciudad"`
	IdZona                 int    `json:"id_zona"`
	Telefono               string `json:"telefono"`
	Telefono_2             string `json:"telefono_2"`
	Email                  string `json:"email"`
	Fax                    string `json:"fax"`
	PaginaWeb              string `json:"pagina_web"`
	LogoUrl                string `json:"logo_url"`
	UsaSedes               *bool  `json:"usa_sedes"`
	NumeroSedes            int    `json:"numero_sedes"`
	IdNaturaleza           int    `json:"id_naturaleza"`
	IdActividadEconomica   int    `json:"id_actividad_economica"`
	IdTipoContribuyente    int    `json:"id_tipo_contribuyente"`
	IdRegimenTributario    int    `json:"id_regimen_tributario"`
	IdTipoDocRepresentante int    `json:"id_tipo_doc_representante"`
	NumeroDocumento        string `json:"numero_documento"`
	RepresentanteLegal     string `json:"representante_legal"`
	CodigoLicencia         string `json:"codigo_licencia"`
	IsActive               *bool  `json:"is_active"`
	IaEndpoint             string `json:"ia_endpoint"`
	IaToken                string `json:"ia_token"`
}

type EmpresaUpdateRequest struct {
	IdTipoEmpresa          int    `json:"id_tipo_empresa" validate:"required,gt=0"`
	Nit                    string `json:"nit" validate:"required,min=5,max=20"`
	Dv                     string `json:"dv" validate:"required,len=1"`
	RazonSocial            string `json:"razon_social" validate:"required,min=2,max=150"`
	Descripcion            string `json:"descripcion" validate:"omitempty,min=2,max=150"`
	Direccion              string `json:"direccion" validate:"required,min=5,max=150"`
	IdPais                 int    `json:"id_pais" valdate:"required,gt=0"`
	IdDepartamento         int    `json:"id_departamento" validate:"required,gt=0"`
	Id_ciudad              int    `json:"id_ciudad" validate:"required,gt=0"`
	IdZona                 int    `json:"id_zona" valdate:"required,gt=0"`
	Telefono               string `json:"telefono" validate:"required,numeric,min=7,max=15"`
	Telefono_2             string `json:"telefono_2" validate:"omitempty,numeric,min=7,max=15"`
	Email                  string `json:"email" validate:"required,email,max=100"`
	Fax                    string `json:"fax" validate:"omitempty,max=70"`
	PaginaWeb              string `json:"pagina_web" validate:"omitempty,url,max=300"`
	LogoUrl                string `json:"logo_url" validate:"omitempty,url,max=300"`
	UsaSedes               *bool  `json:"usa_sedes" validate:"required"`
	NumeroSedes            int    `json:"numero_sedes" valdate:"required"`
	IdNaturaleza           int    `json:"id_naturaleza" validate:"required,gt=0"`
	IdActividadEconomica   int    `json:"id_actividad_economica" validate:"required,gt=0"`
	IdTipoContribuyente    int    `json:"id_tipo_contribuyente" validate:"required,gt=0"`
	IdRegimenTributario    int    `json:"id_regimen_tributario" validate:"required,gt=0"`
	IdTipoDocRepresentante int    `json:"id_tipo_doc_representante" validate:"required,gt=0"`
	NumeroDocumento        string `json:"numero_documento" validate:"required,min=5,max=20"`
	RepresentanteLegal     string `json:"representante_legal" validate:"omitempty,min=2,max=200"`
	CodigoLicencia         string `json:"codigo_licencia" validate:"omitempty,min=2,max=200"`
	Estado                 *bool  `json:"is_active" validate:"required"`
	IaEndpoint             string `json:"ia_endpoint" validate:"omitempty,url,max=300"`
	IaToken                string `json:"ia_token" validate:"omitempty,max=1000"`
	UserID                 int64  `json:"user_id" validate:"omitempty"`
	EmpresaID              int64  `json:"empresa_id" validate:"omitempty"`
	SedeID                 int64  `json:"sede_id" validate:"omitempty"`
}

type SedeRequest struct {
	IdEmpresa        int    `json:"id_empresa" validate:"required,gt=0"`
	Nombre           string `json:"nombre" validate:"required,min=2,max=150"`
	Codigo           string `json:"codigo"`
	IdTipoSede       int    `json:"id_tipo_sede" validate:"required,gt=0"`
	Direccion        string `json:"direccion" validate:"required,min=5,max=150"`
	IdPais           int    `json:"id_pais" valdate:"required,gt=0"`
	IdDepartamento   int    `json:"id_departamento" validate:"required,gt=0"`
	Id_ciudad        int    `json:"id_ciudad" validate:"required,gt=0"`
	IdZona           int    `json:"id_zona" valdate:"required,gt=0"`
	Telefono         string `json:"telefono" validate:"required,numeric,min=7,max=15"`
	Telefono_2       string `json:"telefono_2" validate:"omitempty,numeric,min=7,max=15"`
	Email            string `json:"email" validate:"required,email,max=100"`
	Fax              string `json:"fax" validate:"omitempty,max=70"`
	PaginaWeb        string `json:"pagina_web" validate:"omitempty,url,max=300"`
	ResponsableSede  string `json:"responsable" validate:"required,min=2,max=150"`
	CargoResponsable string `json:"cargo_responsable" validate:"required,min=3,max=150"`
	Estado           *bool  `json:"estado" validate:"required"`
	UserID           int64  `json:"user_id" validate:"omitempty"`
	EmpresaID        int64  `json:"empresa_id" validate:"omitempty"`
	SedeID           int64  `json:"sede_id" validate:"omitempty"`
}

type SedeUpdateRequest struct {
	IdEmpresa        int    `json:"id_empresa" validate:"required,gt=0"`
	Nombre           string `json:"nombre" validate:"required,min=2,max=150"`
	Codigo           string `json:"codigo"`
	IdTipoSede       int    `json:"id_tipo_sede" validate:"required,gt=0"`
	Direccion        string `json:"direccion" validate:"required,min=5,max=150"`
	IdPais           int    `json:"id_pais" valdate:"required,gt=0"`
	IdDepartamento   int    `json:"id_departamento" validate:"required,gt=0"`
	Id_ciudad        int    `json:"id_ciudad" validate:"required,gt=0"`
	IdZona           int    `json:"id_zona" valdate:"required,gt=0"`
	Telefono         string `json:"telefono" validate:"required,numeric,min=7,max=15"`
	Telefono_2       string `json:"telefono_2" validate:"omitempty,numeric,min=7,max=15"`
	Email            string `json:"email" validate:"required,email,max=100"`
	Fax              string `json:"fax" validate:"omitempty,max=70"`
	PaginaWeb        string `json:"pagina_web" validate:"omitempty,url,max=300"`
	ResponsableSede  string `json:"responsable" validate:"required,min=2,max=150"`
	CargoResponsable string `json:"cargo_responsable" validate:"required,min=3,max=150"`
	Estado           *bool  `json:"estado" validate:"required"`
	UserID           int64  `json:"user_id" validate:"omitempty"`
	EmpresaID        int64  `json:"empresa_id" validate:"omitempty"`
	SedeID           int64  `json:"sede_id" validate:"omitempty"`
}

type SedeResponse struct {
	ID       int    `json:"id"`
	Codigo   string `json:"codigo"`
	Nombre   string `json:"nombre"`
	Empresa  string `json:"empresa"`
	Tipo     string `json:"tipo"`
	IsActive bool   `json:"is_active"`
}

type SedeResponseFull struct {
	ID               int    `json:"id"`
	IdEmpresa        int    `json:"id_empresa" `
	Nombre           string `json:"nombre"`
	Codigo           string `json:"codigo"`
	IdTipoSede       int    `json:"id_tipo_sede"`
	Direccion        string `json:"direccion" `
	IdPais           int    `json:"id_pais" `
	IdDepartamento   int    `json:"id_departamento" `
	Id_ciudad        int    `json:"id_ciudad" `
	IdZona           int    `json:"id_zona" `
	Telefono         string `json:"telefono" `
	Telefono_2       string `json:"telefono_2" `
	Email            string `json:"email"`
	Fax              string `json:"fax"`
	PaginaWeb        string `json:"pagina_web"`
	ResponsableSede  string `json:"responsable"`
	CargoResponsable string `json:"cargo_responsable"`
	Estado           *bool  `json:"estado"`
}
