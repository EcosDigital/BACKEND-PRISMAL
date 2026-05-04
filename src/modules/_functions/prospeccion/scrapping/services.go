package scrapping

 
import (
	"fmt"
 
	"github.com/ecosistema/core/src/shared/logging"
	"gorm.io/gorm"
)
 
func RunJob(db *gorm.DB, req *ScrapingJobRequest, userID, empresaID, sedeID int64) (*ScrapingJobResponse, error) {
 
	// 1. Crear job en estado Pendiente
	jobID, err := CreateScrapingJob(db, req, userID, empresaID, sedeID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el job: %w", err)
	}
 
	// 2. Marcar Ejecutando
	if err := UpdateJobEstado(db, jobID, EstadoJobEjecutando); err != nil {
		logging.Error.Printf("[prospeccion] job %d: error marcando Ejecutando: %v", jobID, err)
	}
 
	// 3. Ejecutar engine
	cfg := DefaultEngineConfig(req.Limite)
	rawLeads, engineErr := RunScrapingEngine(req.Keyword, req.Ciudad, cfg)
 
	if engineErr != nil {
		logging.Error.Printf("[prospeccion] job %d: engine falló: %v", jobID, engineErr)
		UpdateJobError(db, jobID, engineErr.Error()) //nolint
		return nil, fmt.Errorf("el scraping falló: %w", engineErr)
	}
 
	// 4. Procesar leads
	counters := JobCounters{Encontrados: len(rawLeads)}
 
	for _, raw := range rawLeads {
		dup, err := CheckLeadDuplicado(db, raw.Telefono, raw.Nombre, raw.Ciudad, empresaID)
		if err != nil {
			logging.Error.Printf("[prospeccion] job %d: error chequeando duplicado (%s): %v", jobID, raw.Nombre, err)
			counters.Errores++
			continue
		}
 
		if dup.EsDuplicado {
			logging.Error.Printf("[prospeccion] job %d: duplicado (%s) motivo=%s id=%v",
				jobID, raw.Nombre, dup.Motivo, dup.IDLeadExistente)
			counters.Duplicados++
			continue
		}
 
		if err := CreateLead(db, &raw, jobID, userID, empresaID, sedeID); err != nil {
			logging.Error.Printf("[prospeccion] job %d: error guardando lead (%s): %v", jobID, raw.Nombre, err)
			counters.Errores++
			continue
		}
 
		counters.Guardados++
	}
 
	// 5. Estado final
	estadoFinal := EstadoJobCompletado
	if counters.Guardados == 0 && counters.Errores > 0 {
		estadoFinal = EstadoJobFallido
	}
 
	if err := UpdateJobCounters(db, jobID, counters, estadoFinal); err != nil {
		logging.Error.Printf("[prospeccion] job %d: error actualizando contadores: %v", jobID, err)
	}
 
	return &ScrapingJobResponse{
		IDJob:            jobID,
		Keyword:          req.Keyword,
		Ciudad:           req.Ciudad,
		TotalEncontrados: counters.Encontrados,
		TotalGuardados:   counters.Guardados,
		TotalDuplicados:  counters.Duplicados,
		TotalErrores:     counters.Errores,
		Estado:           estadoFinal,
	}, nil
}
 
func GetJobs(db *gorm.DB, empresaID int64) ([]JobListResponse, error) {
	return ListScrapingJobs(db, empresaID)
}
 
func GetLeads(db *gorm.DB, f LeadFilterRequest, empresaID int64) ([]LeadResponse, error) {
	return ListLeads(db, f, empresaID)
}
 
func GetLeadsParaExportar(db *gorm.DB, f LeadExportRequest, empresaID int64) ([]LeadExportRow, error) {
	return ListLeadsForExport(db, f, empresaID)
}
 
// RunBatchJob ejecuta todas las combinaciones keyword×ciudad secuencialmente.
// Si el request no trae keywords o ciudades usa los valores por defecto
// (KeywordsComercio y MunicipiosCordoba).
func RunBatchJob(db *gorm.DB, req *BatchJobRequest, userID, empresaID, sedeID int64) (*BatchJobResponse, error) {
 
	keywords := req.Keywords
	if len(keywords) == 0 {
		keywords = KeywordsComercio
	}
 
	ciudades := req.Ciudades
	if len(ciudades) == 0 {
		ciudades = MunicipiosCordoba
	}
 
	limite := req.Limite
	if limite == 0 {
		limite = 20
	}
 
	resp := &BatchJobResponse{}
 
	for _, kw := range keywords {
		for _, ciudad := range ciudades {
			jobReq := &ScrapingJobRequest{
				Keyword: kw,
				Ciudad:  ciudad,
				Limite:  limite,
			}
 
			jobResp, err := RunJob(db, jobReq, userID, empresaID, sedeID)
			if err != nil {
				logging.Error.Printf("[prospeccion] batch: error en %s/%s: %v", kw, ciudad, err)
				resp.TotalErrores++
				continue
			}
 
			resp.TotalCombinaciones++
			resp.TotalEncontrados += jobResp.TotalEncontrados
			resp.TotalGuardados   += jobResp.TotalGuardados
			resp.TotalDuplicados  += jobResp.TotalDuplicados
			resp.TotalErrores     += jobResp.TotalErrores
			resp.Jobs              = append(resp.Jobs, *jobResp)
		}
	}
 
	return resp, nil
}