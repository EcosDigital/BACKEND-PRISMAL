package usuarios

type UserRequest struct {
	Email        string `json:"email" validate:"required,email,max=100"`
	Password     string `json:"password" validate:"required,min=8"`
	IdTercero    int    `json:"id_tercero" validate:"required,gt=0"`
	IdRol        int    `json:"id_rol" validate:"required,gt=0"`
	ImageProfile string `json:"image_profile" validate:"omitempty,url"`
	Activo       *bool  `json:"activo" validate:"required"`
	UserID       int64  `json:"user_id" validate:"omitempty"`
	EmpresaID    int64  `json:"empresa_id" validate:"omitempty"`
	SedeID       int64  `json:"sede_id" validate:"omitempty"`
}

type UserResp struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Rol       string `json:"rol"`
	Tercero   string `json:"tercero"`
	Is_active bool   `json:"is_active"`
}

type UserResponseFull struct {
	ID        int    `json:"id"`
	IdTercero int    `json:"id_tercero"`
	Email     string `json:"email"`
	IdRol     int    `json:"id_rol"`
	Is_active bool   `json:"is_active"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IdRol    int    `json:"id_rol"`
	Rol      string `json:"rol"`
	Nombre   string `json:"nombre"`
	Activo   bool   `json:"activo"`
}

type UserUpdateRequest struct {
	IdTercero    int    `json:"id_tercero" validate:"required,gt=0"`
	IdRol        int    `json:"id_rol" validate:"required,gt=0"`
	ImageProfile string `json:"image_profile" validate:"omitempty,url"`
	Activo       *bool  `json:"activo" validate:"required"`
	UserID       int64  `json:"user_id" validate:"omitempty"`
	EmpresaID    int64  `json:"empresa_id" validate:"omitempty"`
	SedeID       int64  `json:"sede_id" validate:"omitempty"`
}

type NewPasswordUser struct {
	Password string `json:"password" validate:"required,min=8"`
}
