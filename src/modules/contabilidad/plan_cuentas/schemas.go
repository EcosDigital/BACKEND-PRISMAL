package plan_cuentas

// ─── Requests ─────────────────────────────────────────────────────────────────

type CuentaContableRequest struct {
	CodigoCuenta        string `json:"codigo_cuenta"          validate:"required,min=1,max=20"`
	NombreCuenta        string `json:"nombre_cuenta"          validate:"required,min=3"`
	IDCuentaPadre       *int   `json:"id_cuenta_padre"        validate:"omitempty"`
	IDNaturaleza        int    `json:"id_naturaleza"          validate:"gt=0"`
	IDTipoCuenta        *int   `json:"id_tipo_cuenta"         validate:"omitempty"`
	IDNivelCuenta       int    `json:"id_nivel_cuenta"        validate:"gt=0"`
	PermiteMovimientos  *bool  `json:"permite_movimientos"    validate:"required"`
	RequiereTercero     *bool  `json:"requiere_tercero"       validate:"required"`
	RequiereCentroCosto *bool  `json:"requiere_centro_costo"  validate:"required"`
	IsActive            *bool  `json:"is_active"              validate:"required"`
	// Inyectados desde el contexto JWT — no vienen del body
	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
	SedeID    int64 `json:"sede_id"    validate:"omitempty"`
}

type CuentaContableUpdateRequest struct {
	NombreCuenta        string `json:"nombre_cuenta"          validate:"required,min=3"`
	IDCuentaPadre       *int   `json:"id_cuenta_padre"        validate:"omitempty"`
	IDNaturaleza        int    `json:"id_naturaleza"          validate:"gt=0"`
	IDTipoCuenta        *int   `json:"id_tipo_cuenta"         validate:"omitempty"`
	IDNivelCuenta       int    `json:"id_nivel_cuenta"        validate:"gt=0"`
	PermiteMovimientos  *bool  `json:"permite_movimientos"    validate:"required"`
	RequiereTercero     *bool  `json:"requiere_tercero"       validate:"required"`
	RequiereCentroCosto *bool  `json:"requiere_centro_costo"  validate:"required"`
	IsActive            *bool  `json:"is_active"              validate:"required"`
	UserID              int64  `json:"user_id"                validate:"omitempty"`
	EmpresaID           int64  `json:"empresa_id"             validate:"omitempty"`
	SedeID              int64  `json:"sede_id"                validate:"omitempty"`
}

// ─── Responses ────────────────────────────────────────────────────────────────

type CuentaContableResponse struct {
	ID                  int    `json:"id"                    gorm:"column:id"`
	CodigoCuenta        string `json:"codigo_cuenta"         gorm:"column:codigo_cuenta"`
	NombreCuenta        string `json:"nombre_cuenta"         gorm:"column:nombre_cuenta"`
	CuentaPadre         string `json:"cuenta_padre"          gorm:"column:cuenta_padre"`
	Naturaleza          string `json:"naturaleza"            gorm:"column:naturaleza"`
	TipoCuenta          string `json:"tipo_cuenta"           gorm:"column:tipo_cuenta"`
	NivelCuenta         string `json:"nivel_cuenta"          gorm:"column:nivel_cuenta"`
	PermiteMovimientos  *bool  `json:"permite_movimientos"   gorm:"column:permite_movimientos"`
	RequiereTercero     *bool  `json:"requiere_tercero"      gorm:"column:requiere_tercero"`
	RequiereCentroCosto *bool  `json:"requiere_centro_costo" gorm:"column:requiere_centro_costo"`
	IsActive            *bool  `json:"is_active"             gorm:"column:is_active"`
}

// CuentaContableResponseFull se usa en el detalle (edición)
type CuentaContableResponseFull struct {
	ID                  int    `json:"id"                    gorm:"column:id"`
	CodigoCuenta        string `json:"codigo_cuenta"         gorm:"column:codigo_cuenta"`
	NombreCuenta        string `json:"nombre_cuenta"         gorm:"column:nombre_cuenta"`
	IDCuentaPadre       *int   `json:"id_cuenta_padre"       gorm:"column:id_cuenta_padre"`
	CuentaPadre         string `json:"cuenta_padre"          gorm:"column:cuenta_padre"`
	IDNaturaleza        int    `json:"id_naturaleza"         gorm:"column:id_naturaleza"`
	Naturaleza          string `json:"naturaleza"            gorm:"column:naturaleza"`
	IDTipoCuenta        *int   `json:"id_tipo_cuenta"        gorm:"column:id_tipo_cuenta"`
	TipoCuenta          string `json:"tipo_cuenta"           gorm:"column:tipo_cuenta"`
	IDNivelCuenta       int    `json:"id_nivel_cuenta"       gorm:"column:id_nivel_cuenta"`
	NivelCuenta         string `json:"nivel_cuenta"          gorm:"column:nivel_cuenta"`
	PermiteMovimientos  *bool  `json:"permite_movimientos"   gorm:"column:permite_movimientos"`
	RequiereTercero     *bool  `json:"requiere_tercero"      gorm:"column:requiere_tercero"`
	RequiereCentroCosto *bool  `json:"requiere_centro_costo" gorm:"column:requiere_centro_costo"`
	IsActive            *bool  `json:"is_active"             gorm:"column:is_active"`
}

// ─── Referencias ──────────────────────────────────────────────────────────────

type RefNaturaleza struct {
	ID     int    `json:"id"     gorm:"column:id"`
	Nombre string `json:"nombre" gorm:"column:nombre"`
}

type RefTipoCuenta struct {
	ID     int    `json:"id"     gorm:"column:id"`
	Nombre string `json:"nombre" gorm:"column:nombre"`
}

type RefNivelCuenta struct {
	ID     int    `json:"id"     gorm:"column:id"`
	Nombre string `json:"nombre" gorm:"column:nombre"`
}

// RefCuentaPadre se usa en el select de cuenta padre del formulario
type RefCuentaPadre struct {
	ID           int    `json:"id"            gorm:"column:id"`
	CodigoCuenta string `json:"codigo_cuenta" gorm:"column:codigo_cuenta"`
	NombreCuenta string `json:"nombre_cuenta" gorm:"column:nombre_cuenta"`
	NivelCuenta  string `json:"nivel_cuenta"  gorm:"column:nivel_cuenta"`
}

// ─── Carga masiva ─────────────────────────────────────────────────────────────
//
// Columnas del plano plantillaPUC.xlsx:
//   A  CODIGO_CUENTA           obligatorio — código único por empresa
//   B  NOMBRE                  obligatorio — nombre de la cuenta
//   C  CUENTA_PADRE            opcional  — CÓDIGO de la cuenta padre (no ID)
//   D  ID_NATURALEZA           obligatorio — nombre: Activo|Pasivo|Patrimonio|Ingreso|Gasto|Costo
//   E  ID_TIPO_CUENTA          opcional  — nombre: Inventario|Caja|Bancos|...
//   F  NIVEL_CUENTA            obligatorio — número 1-4 o texto "Nivel N"
//   G  PERMITE MOVIMIENTOS     obligatorio — Si|No
//   H  REQUIERE_TERCEROS       obligatorio — Si|No
//   I  REQUIERE_CENTROS_COSTO  obligatorio — Si|No
//   J  ESTADO                  obligatorio — Activo|Inactivo

type PlanoCuentaRow struct {
	CodigoCuenta         string `json:"CODIGO_CUENTA"`
	Nombre               string `json:"NOMBRE"`
	CuentaPadre          string `json:"CUENTA_PADRE"`
	IDNaturaleza         string `json:"ID_NATURALEZA"`
	IDTipoCuenta         string `json:"ID_TIPO_CUENTA"`
	NivelCuenta          string `json:"NIVEL_CUENTA"`
	PermiteMovimientos   string `json:"PERMITE MOVIMIENTOS"`
	RequiereTerceros     string `json:"REQUIERE_TERCEROS"`
	RequiereCentrosCosto string `json:"REQUIERE_CENTROS_COSTO"`
	Estado               string `json:"ESTADO"`
}

type ImportCuentasRequest struct {
	Filas []PlanoCuentaRow `json:"filas" validate:"required,min=1"`
}

type ImportCuentaRowResult struct {
	Fila   int    `json:"fila"`
	Codigo string `json:"codigo"`
	Accion string `json:"accion"` // "creada" | "actualizada" | "error"
	Error  string `json:"error,omitempty"`
}

type ImportCuentasResponse struct {
	Creadas      int                     `json:"creadas"`
	Actualizadas int                     `json:"actualizadas"`
	Errores      int                     `json:"errores"`
	Detalle      []ImportCuentaRowResult `json:"detalle"`
}

// cuentasRefCache centraliza todos los lookups nombre→id para el import
type cuentasRefCache struct {
	Naturalezas       map[string]int // "activo" → id
	TiposCuenta       map[string]int // "inventario" → id
	NivelesCuenta     map[string]int // "nivel 1" → id  AND  "1" → id
	CodigosExistentes map[string]int // codigo_lower → id (saber si crear o actualizar)
}

// ─── Exportación Excel ────────────────────────────────────────────────────────
//
// CuentaExportRow es la fila que se serializa en la respuesta JSON del endpoint
// GET /contabilidad/cuentas/export y luego el frontend la convierte a XLSX.
// Los booleanos viajan como bool (true/false); la conversión a "Sí"/"No"
// se realiza en el frontend con SheetJS para no acoplar lógica de presentación
// al backend.
type CuentaExportRow struct {
	CodigoCuenta        string `json:"codigo_cuenta"         gorm:"column:codigo_cuenta"`
	NombreCuenta        string `json:"nombre_cuenta"         gorm:"column:nombre_cuenta"`
	CuentaPadre         string `json:"cuenta_padre"          gorm:"column:cuenta_padre"`
	Naturaleza          string `json:"naturaleza"            gorm:"column:naturaleza"`
	TipoCuenta          string `json:"tipo_cuenta"           gorm:"column:tipo_cuenta"`
	NivelCuenta         string `json:"nivel_cuenta"          gorm:"column:nivel_cuenta"`
	PermiteMovimientos  bool   `json:"permite_movimientos"   gorm:"column:permite_movimientos"`
	RequiereTercero     bool   `json:"requiere_tercero"      gorm:"column:requiere_tercero"`
	RequiereCentroCosto bool   `json:"requiere_centro_costo" gorm:"column:requiere_centro_costo"`
	IsActive            bool   `json:"is_active"             gorm:"column:is_active"`
	CreatedAt           string `json:"created_at"            gorm:"column:created_at"`
}
