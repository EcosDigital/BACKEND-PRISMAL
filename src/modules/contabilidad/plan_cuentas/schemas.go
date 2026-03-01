package plan_cuentas

type CuentaContableRequest struct {
	CodigoCuenta       string `json:"codigo_cuenta" validate:"required,min=2,max=12"`
	NombreCuenta       string `json:"nombre_cuenta" validate:"required,min=3"`
	IDCuentaPadre      *int   `json:"id_cuenta_padre" validate:"omitempty,gt=0"`
	IDNaturaleza       int    `json:"id_naturaleza" validate:"gt=0"`
	IDTipoCuenta       *int   `json:"id_tipo_cuenta" validate:"omitempty,gt=0"`
	IDNivelCuenta      *int   `json:"id_nivel_cuenta" validate:"omitempty,gt=0"`
	PermiteMovimientos *bool  `json:"permite_movimientos" validate:"required"`
	ReqTercero         *bool  `json:"requiere_tercero" validate:"required"`
	ReqCentroCosto     *bool  `json:"requiere_centro_costo" validate:"required"`
	IsActive           *bool  `json:"is_active" validate:"required"`
	UserID             int64  `json:"user_id" validate:"omitempty"`
	EmpresaID          int64  `json:"empresa_id" validate:"omitempty"`
	SedeID             int64  `json:"sede_id" validate:"omitempty"`
}

type CuentaContableResponse struct {
	ID                 int    `json:"id"`
	CodigoCuenta       string `json:"codigo_cuenta"`
	NombreCuenta       string `json:"nombre_cuenta"`
	IDCuentaPadre      *int   `json:"id_cuenta_padre"`
	IDNaturaleza       int    `json:"id_naturaleza"`
	Naturaleza         string `json:"naturaleza"`
	IDTipoCuenta       *int   `json:"id_tipo_cuenta"`
	TipoCuenta         string `json:"tipo_cuenta"`
	IDNivelCuenta      *int   `json:"id_nivel_cuenta"`
	NivelCuenta        string `json:"nivel_cuenta"`
	PermiteMovimientos *bool  `json:"permite_movimientos" validate:"required"`
	ReqTercero         *bool  `json:"requiere_tercero" validate:"required"`
	IsActive           *bool  `json:"is_active" validate:"required"`
}
