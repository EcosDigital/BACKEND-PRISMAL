package reports

type EjecutarReportesRequest struct {
	IdInforme int                    `json:"id_informe" validate:"required"`
	Filtros   map[string]interface{} `json:"filtros"`
}

type InformeDB struct {
	ID         int64
	FuncionSQL string
	IsActive   bool
}

type FiltroInformeDB struct {
	Codigo    string
	Requerido bool
}
