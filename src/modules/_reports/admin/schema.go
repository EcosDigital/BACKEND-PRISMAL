package reports

type ReportServerConnectionRequest struct {
	Name      string `json:"name"     validate:"required,min=3,max=100"`
	BaseURL   string `json:"base_url" validate:"required,url"`
	APIKey    string `json:"api_key"  validate:"required,min=8"`
	UserID    int64  `json:"user_id"              validate:"omitempty"`
	EmpresaID int64  `json:"empresa_id"           validate:"omitempty"`
	SedeID    int64  `json:"sede_id"              validate:"omitempty"`
}

type ReportServerConnectionUpdateRequest struct {
	Name    string `json:"name"     validate:"required,min=3,max=100"`
	BaseURL string `json:"base_url" validate:"required,url"`
	APIKey  string `json:"api_key"  validate:"required,min=8"`
	// UserID is injected from JWT context, not from body
	UserID int64 `json:"user_id" validate:"omitempty"`
}

type ReportServerConnectionResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BaseURL   string `json:"base_url"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type ReportServerConnectionResponseFull struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BaseURL   string `json:"base_url"`
	APIKey    string `json:"api_key"`
	IsActive  bool   `json:"is_active"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SetActiveRequest struct {
	// No body needed — the ID comes from the URL param
}
