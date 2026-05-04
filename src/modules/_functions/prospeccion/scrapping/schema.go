package scrapping

// ─── Requests ─────────────────────────────────────────────────────────────────

// ScrapingJobRequest es el body que recibe el endpoint POST /prospeccion/run.
// keyword: término de búsqueda (ej: "restaurantes", "ferreterías")
// ciudad:  ciudad donde buscar (ej: "Montería", "Bogotá")
// limite:  máximo de leads a extraer (default 50, max 200)
type ScrapingJobRequest struct {
	Keyword string `json:"keyword" validate:"required,min=2,max=250"`
	Ciudad  string `json:"ciudad"  validate:"required,min=2,max=150"`
	Limite  int    `json:"limite"  validate:"omitempty,min=1,max=200"`
}

// LeadExportRequest filtra qué leads exportar al Excel.
// Todos los campos son opcionales; si se omiten se exporta todo.
type LeadExportRequest struct {
	IDJob    *int64 `json:"id_job"       query:"id_job"`
	Ciudad   string `json:"ciudad"       query:"ciudad"`
	Keyword  string `json:"keyword"      query:"keyword"`
	IDEstado *int   `json:"id_estado"    query:"id_estado"`
}

// LeadFilterRequest filtra el listado paginado de leads.
type LeadFilterRequest struct {
	IDJob    *int64 `query:"id_job"`
	Ciudad   string `query:"ciudad"`
	Keyword  string `query:"keyword"`
	IDEstado *int   `query:"id_estado"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

// ─── Responses ────────────────────────────────────────────────────────────────

// ScrapingJobResponse es lo que retorna POST /prospeccion/run al terminar.
type ScrapingJobResponse struct {
	IDJob            int64  `json:"id_job"`
	Keyword          string `json:"keyword"`
	Ciudad           string `json:"ciudad"`
	TotalEncontrados int    `json:"total_encontrados"`
	TotalGuardados   int    `json:"total_guardados"`
	TotalDuplicados  int    `json:"total_duplicados"`
	TotalErrores     int    `json:"total_errores"`
	Estado           string `json:"estado"`
}

// JobListResponse representa un job en el listado histórico.
type JobListResponse struct {
	ID               int64   `json:"id"`
	Keyword          string  `json:"keyword"`
	Ciudad           string  `json:"ciudad"`
	Limite           int     `json:"limite"`
	Estado           string  `json:"estado"`
	TotalEncontrados int     `json:"total_encontrados"`
	TotalGuardados   int     `json:"total_guardados"`
	TotalDuplicados  int     `json:"total_duplicados"`
	TotalErrores     int     `json:"total_errores"`
	CreatedAt        string  `json:"created_at"`
	FinishedAt       *string `json:"finished_at"`
}

// LeadResponse es la representación de un lead en el listado HTTP.
type LeadResponse struct {
	ID           int64   `json:"id"`
	IDJob        int64   `json:"id_job"`
	Nombre       string  `json:"nombre"`
	TipoNegocio  string  `json:"tipo_negocio"`
	Ciudad       string  `json:"ciudad"`
	Direccion    string  `json:"direccion"`
	Telefono     string  `json:"telefono"`
	SitioWeb     string  `json:"sitio_web"`
	Rating       float64 `json:"rating"`
	TotalReviews int     `json:"total_reviews"`
	MapsURL      string  `json:"maps_url"`
	EstadoLead   string  `json:"estado_lead"`
	CreatedAt    string  `json:"created_at"`
}

// LeadExportRow es la fila que se escribe en el Excel.
// Los campos son planos y listos para consumo comercial.
type LeadExportRow struct {
	Nombre       string
	TipoNegocio  string
	Ciudad       string
	Direccion    string
	Telefono     string
	SitioWeb     string
	Rating       float64
	TotalReviews int
	MapsURL      string
	Estado       string
}

// ─── Tipos internos (engine → service) ────────────────────────────────────────

// LeadRaw es el dato crudo que retorna el scraping engine.
// No toca la DB ni los middlewares; solo viaja del engine al service.
type LeadRaw struct {
	Nombre       string
	TipoNegocio  string
	Ciudad       string
	Direccion    string
	Telefono     string
	SitioWeb     string
	Rating       float64
	TotalReviews int
	MapsURL      string
}

// DuplicadoCheck es la respuesta de fn_check_lead_duplicado.
type DuplicadoCheck struct {
	EsDuplicado     bool
	Motivo          string // "TELEFONO" | "NOMBRE_CIUDAD" | ""
	IDLeadExistente *int64
}

// JobCounters acumula los contadores durante la ejecución del job.
// El service lo actualiza en memoria y lo persiste al finalizar.
type JobCounters struct {
	Encontrados int
	Guardados   int
	Duplicados  int
	Errores     int
}

// ─── Batch ────────────────────────────────────────────────────────────────────

// municipiosCórdoba es la lista fija de municipios del departamento.
// Se usa cuando el usuario no especifica ciudades.
var MunicipiosCordoba = []string{
	"Montería", "Cereté", "Sahagún", "Lorica", "Montelíbano",
	"Tierralta", "Planeta Rica", "Ciénaga de Oro", "San Pelayo",
	"Chinú", "Ayapel", "Buenavista", "Chima", "Cotorra",
	"La Apartada", "Moñitos", "Purísima", "San Andrés de Sotavento",
	"San Bernardo del Viento", "San Carlos", "Tuchín", "Valencia",
}

// KeywordsComercio agrupa los sectores con mayor potencial para software de gestión.
var KeywordsComercio = []string{
	"ferretería", "droguería", "farmacia", "distribuidora",
	"supermercado", "minimercado", "repuestos", "papelería",
	"veterinaria", "taller mecánico", "óptica", "materiales construcción",
}

// BatchJobRequest permite lanzar múltiples combinaciones keyword×ciudad en un solo request.
// Si Keywords está vacío usa KeywordsComercio.
// Si Ciudades está vacío usa todos los municipios de Córdoba.
type BatchJobRequest struct {
	Keywords []string `json:"keywords" validate:"omitempty"`
	Ciudades []string `json:"ciudades" validate:"omitempty"`
	Limite   int      `json:"limite"   validate:"omitempty,min=1,max=50"`
}

// BatchJobResponse resume el resultado de todas las combinaciones.
type BatchJobResponse struct {
	TotalCombinaciones int                   `json:"total_combinaciones"`
	TotalEncontrados   int                   `json:"total_encontrados"`
	TotalGuardados     int                   `json:"total_guardados"`
	TotalDuplicados    int                   `json:"total_duplicados"`
	TotalErrores       int                   `json:"total_errores"`
	Jobs               []ScrapingJobResponse `json:"jobs"`
}
