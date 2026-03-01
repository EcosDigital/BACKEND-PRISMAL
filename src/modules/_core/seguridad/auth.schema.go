package seguridad

import "time"

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IdRol    int    `json:"id_rol"`
	Rol      string `json:"rol"`
	Nombre   string `json:"nombre"`
	Activo   bool   `json:"activo"`
}

type SignInResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VerifyTokenResponse struct {
	User struct {
		ID     int    `json:"id"`
		Email  string `json:"email"`
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
	} `json:"user"`
}

type CompanyAccess struct {
	ID          int    `json:"id"`
	IdEmpresa   int    `json:"id_empresa"`
	RazonSocial string `json:"razon_social"`
	Descripcion string `json:"descripcion"`
	UsaSedes    bool   `json:"usa_sedes"`
}

type SedesAccess struct {
	ID        int    `json:"id"`
	IdEmpresa int    `json:"id_empresa"`
	IdSede    int    `json:"id_sede"`
	Nombre    string `json:"nombre"`
}

type AccessAutorized struct {
	EmpresaID   int           `json:"empresa_id"`
	RazonSocial string        `json:"razon_social"`
	Descripcion string        `json:"descripcion"`
	UsaSedes    bool          `json:"usa_sedes"`
	Sedes       []SedesAccess `json:"sedes"`
}
