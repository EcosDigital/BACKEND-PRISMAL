package roles

type RolesRequest struct {
	Nombre         string `json:"nombre" validate:"required,min=3"`
	IDTipo         int    `json:"id_tipo" validate:"omitempty,gt=0"`
	IntentosLogin  int    `json:"intentos_login" validate:"required"`
	IsMultisession *bool  `json:"is_multisession" validate:"required"`
	IsExterno      *bool  `json:"is_externo" validate:"required"`
	UserID         int64  `json:"user_id" validate:"omitempty"`
	EmpresaID      int64  `json:"empresa_id" validate:"omitempty"`
	SedeID         int64  `json:"sede_id" validate:"omitempty"`
}

type RolesResponse struct {
	ID             int    `json:"id"`
	Nombre         string `json:"nombre"`
	IDTipo         int    `json:"id_tipo"`
	Tipo           string `json:"tipo"`
	IntentosLogin  int    `json:"intentos_login"`
	IsMultisession *bool  `json:"is_multisession"`
	IsExterno      *bool  `json:"is_externo"`
}

type Item struct {
	ID       *int    `json:"id"`
	Title    *string `json:"title"`
	Icon     *string `json:"icon,omitempty"`
	Path     string  `json:"path,omitempty"`
	BgColor  *string `json:"bg_color,omitempty"`
	BrColor  *string `json:"br_color,omitempty"`
	Children []Item  `json:"children,omitempty"` // corregir typo cuando puedas
}

type ConfigRolRequest struct {
	IDRol     int    `json:"id_rol"`
	Sedes     []int  `json:"sedes"`
	Json      []Item `json:"json"`
	UserID    int64  `json:"user_id"`
	EmpresaID int64  `json:"empresa_id"`
	SedeID    int64  `json:"sede_id"`
}
