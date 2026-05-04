package fiados

// ─── Requests ────────────────────────────────────────────────────────────────

// RegistrarFiadoRequest body del POST /fiados/fiado
// El tendero puede pasar un id_cuenta existente (cliente ya registrado)
// o los datos del cliente nuevo para crearlo en cfg_terceros automáticamente.

type RegistrarFiadoRequest struct {
	IDCuenta        int64  `json:"id_cuenta" validate:"omitempty"`
	NombreCliente   string `json:"nombre_cliente" validate:"omitempty"`
	ApellidoCliente string `json:"apellido_cliente" validate:"omitempty"`
	Cedula          string `json:"cedula"         validate:"omitempty"`
	Telefono        string `json:"telefono"       validate:"omitempty"`
	Direccion       string `json:"direccion"      validate:"omitempty"`
	// Datos del movimiento (siempre requeridos)
	IDComprobante int64   `json:"id_comprobante" validate:"gt=0"`
	Valor         float64 `json:"valor"          validate:"gt=0"`
	Observaciones string  `json:"observaciones"  validate:"omitempty"` // lista de artículos
	Fecha         string  `json:"fecha"          validate:"omitempty"` // YYYY-MM-DD, opcional
	UserID        int64   `json:"user_id"    validate:"omitempty"`
	EmpresaID     int64   `json:"empresa_id" validate:"omitempty"`
	SedeID        int64   `json:"sede_id"    validate:"omitempty"`
}

// RegistrarAbonoRequest body del POST /fiados/abono
type RegistrarAbonoRequest struct {
	IDCuenta      int64   `json:"id_cuenta"      validate:"gt=0"`
	IDComprobante int64   `json:"id_comprobante" validate:"gt=0"`
	Valor         float64 `json:"valor"          validate:"gt=0"`
	Observaciones string  `json:"observaciones"  validate:"omitempty"`
	Fecha         string  `json:"fecha"          validate:"omitempty"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
	SedeID    int64 `json:"sede_id"    validate:"omitempty"`
}

// AnularMovimientoRequest body del PUT /fiados/movimientos/:id/anular
type AnularMovimientoRequest struct {
	Motivo string `json:"motivo" validate:"required,min=5"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ActualizarCuentaRequest body del PUT /fiados/cuentas/:id
// Permite editar límite de crédito, observaciones y QR de cobro.
type ActualizarCuentaRequest struct {
	LimiteCredito *float64 `json:"limite_credito" validate:"omitempty"`
	Observaciones string   `json:"observaciones"  validate:"omitempty"`
	QrCobroURL    string   `json:"qr_cobro_url"   validate:"omitempty"`

	UserID    int64 `json:"user_id"    validate:"omitempty"`
	EmpresaID int64 `json:"empresa_id" validate:"omitempty"`
}

// ─── Responses ───────────────────────────────────────────────────────────────

// FiadoMovResult respuesta de fn_registrar_fiado / fn_registrar_abono
type FiadoMovResult struct {
	IDMovimiento     int64   `json:"id_movimiento"       gorm:"column:id_movimiento"`
	IDCuenta         int64   `json:"id_cuenta"           gorm:"column:id_cuenta"`
	IDMovComprobante int64   `json:"id_mov_comprobante"  gorm:"column:id_mov_comprobante"`
	Consecutivo      int     `json:"consecutivo"         gorm:"column:consecutivo"`
	Prefijo          string  `json:"prefijo"             gorm:"column:prefijo"`
	SaldoResultante  float64 `json:"saldo_resultante"    gorm:"column:saldo_resultante"`
}

// AbonoMovResult extiende FiadoMovResult con el campo cuenta_saldada
type AbonoMovResult struct {
	FiadoMovResult
	CuentaSaldada bool `json:"cuenta_saldada" gorm:"column:cuenta_saldada"`
}

// CuentaFiadoResponse fila de la vista fiados.v_cuentas
type CuentaFiadoResponse struct {
	IDCuenta              int64    `json:"id_cuenta"                gorm:"column:id_cuenta"`
	IDTercero             int64    `json:"id_tercero"               gorm:"column:id_tercero"`
	NombreCliente         string   `json:"nombre_cliente"           gorm:"column:nombre_cliente"`
	Cedula                string   `json:"cedula"                   gorm:"column:cedula"`
	Telefono              string   `json:"telefono"                 gorm:"column:telefono"`
	Direccion             string   `json:"direccion"                gorm:"column:direccion"`
	SaldoActual           float64  `json:"saldo_actual"             gorm:"column:saldo_actual"`
	EstadoCodigo          string   `json:"estado_codigo"            gorm:"column:estado_codigo"`
	EstadoNombre          string   `json:"estado_nombre"            gorm:"column:estado_nombre"`
	LimiteCredito         *float64 `json:"limite_credito"           gorm:"column:limite_credito"`
	FechaApertura         string   `json:"fecha_apertura"           gorm:"column:fecha_apertura"`
	FechaCierre           *string  `json:"fecha_cierre"             gorm:"column:fecha_cierre"`
	Observaciones         *string  `json:"observaciones"            gorm:"column:observaciones"`
	FechaUltimoMovimiento *string  `json:"fecha_ultimo_movimiento"  gorm:"column:fecha_ultimo_movimiento"`
	UltimoTipoMovimiento  *string  `json:"ultimo_tipo_movimiento"   gorm:"column:ultimo_tipo_movimiento"`
}

// MovimientoFiadoResponse fila de fn_historial_cuenta
type MovimientoFiadoResponse struct {
	IDMovimiento      int64   `json:"id_movimiento"       gorm:"column:id_movimiento"`
	Fecha             string  `json:"fecha"               gorm:"column:fecha"`
	TipoCodigo        string  `json:"tipo_codigo"         gorm:"column:tipo_codigo"`
	TipoNombre        string  `json:"tipo_nombre"         gorm:"column:tipo_nombre"`
	Valor             float64 `json:"valor"               gorm:"column:valor"`
	Observaciones     string  `json:"observaciones"       gorm:"column:observaciones"`
	SaldoAnterior     float64 `json:"saldo_anterior"      gorm:"column:saldo_anterior"`
	SaldoMovimiento   float64 `json:"saldo_movimiento"    gorm:"column:saldo_movimiento"`
	SaldoDespues      float64 `json:"saldo_despues"       gorm:"column:saldo_despues"`
	Documento         string  `json:"documento"           gorm:"column:documento"`
	EstadoComprobante string  `json:"estado_comprobante"  gorm:"column:estado_comprobante"`
	EsAnulado         bool    `json:"es_anulado"          gorm:"column:es_anulado"`
}

// ResumenCarteraResponse fila de fn_resumen_cartera
type ResumenCarteraResponse struct {
	TotalCuentas      int     `json:"total_cuentas"       gorm:"column:total_cuentas"`
	CuentasAlDia      int     `json:"cuentas_al_dia"      gorm:"column:cuentas_al_dia"`
	CuentasPendientes int     `json:"cuentas_pendientes"  gorm:"column:cuentas_pendientes"`
	CuentasAtrasadas  int     `json:"cuentas_atrasadas"   gorm:"column:cuentas_atrasadas"`
	TotalCartera      float64 `json:"total_cartera"       gorm:"column:total_cartera"`
	TotalFiadosMes    float64 `json:"total_fiados_mes"    gorm:"column:total_fiados_mes"`
	TotalAbonosMes    float64 `json:"total_abonos_mes"    gorm:"column:total_abonos_mes"`
}
