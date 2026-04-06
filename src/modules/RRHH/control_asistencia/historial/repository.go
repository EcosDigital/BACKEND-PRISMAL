package historial_asistencia

import (
	"gorm.io/gorm"
)

// GetHistorialAsistencia construye el reporte completo de asistencia de un tercero.
// El cálculo del estado (PUNTUAL / TOLERANCIA / TARDE / SIN_HORARIO) lo realiza
// la función PostgreSQL fn_historial_asistencia — el backend Go solo invoca
// la función y mapea el resultado.
func GetHistorialAsistencia(
	db *gorm.DB,
	req *HistorialAsistenciaRequest,
	empresaID int64,
) (*HistorialAsistenciaReporte, error) {

	// ── 1. Datos del tercero — ORM puro ──────────────────────────
	type TerceroInfo struct {
		PrimerNombre    string `gorm:"column:primer_nombre"`
		SegundoNombre   string `gorm:"column:segundo_nombre"`
		PrimerApellido  string `gorm:"column:primer_apellido"`
		SegundoApellido string `gorm:"column:segundo_apellido"`
		NumeroDocumento string `gorm:"column:numero_documento"`
	}
	var t TerceroInfo

	err := db.
		Table("configuracion.cfg_terceros").
		Select("primer_nombre, segundo_nombre, primer_apellido, segundo_apellido, numero_documento").
		Where("id = ?", req.IDTercero).
		Scan(&t).Error
	if err != nil {
		return nil, err
	}

	// ── 2. Datos de la empresa — ORM con Joins ───────────────────
	var empresa HistorialAsistenciaEmpresa

	err = db.
		Table("configuracion.cfg_empresas emp").
		Select(`
			emp.razon_social,
			emp.nit,
			emp.dv,
			emp.direccion,
			mun.nombre_municipio AS ciudad,
			dep.nombre           AS departamento,
			emp.telefono,
			COALESCE(emp.logo_url, '') AS logo_url`).
		Joins("INNER JOIN configuracion.ref_municipios    mun ON mun.id = emp.id_ciudad").
		Joins("INNER JOIN configuracion.ref_departamentos dep ON dep.id = emp.id_departamento").
		Where("emp.id = ?", empresaID).
		Scan(&empresa).Error
	if err != nil {
		return nil, err
	}

	// ── 3. Llamar a fn_historial_asistencia ──────────────────────
	// La función PostgreSQL encapsula todo el cruce de tablas y el
	// cálculo de estado. GORM invoca la función con parámetros
	// tipados — sin SQL plano expuesto en el código Go.
	type filaRaw struct {
		FilaOrden         int    `gorm:"column:fila_orden"`
		Fecha             string `gorm:"column:fecha"`
		Hora              string `gorm:"column:hora"`
		DiaSemana         string `gorm:"column:dia_semana"`
		TipoMarcacion     string `gorm:"column:tipo_marcacion"`
		PuntoMarcacion    string `gorm:"column:punto_marcacion"`
		Observaciones     string `gorm:"column:observaciones"`
		HorarioAsignado   string `gorm:"column:horario_asignado"`
		HoraProgramada    string `gorm:"column:hora_programada"`
		MinutosDiferencia int    `gorm:"column:minutos_diferencia"`
		Estado            string `gorm:"column:estado"`
	}

	var rawFilas []filaRaw

	err = db.
		Table("control_asistencia.fn_historial_asistencia(?, ?, ?)",
			req.IDTercero,
			req.FechaInicio,
			req.FechaFin,
		).
		Scan(&rawFilas).Error
	if err != nil {
		return nil, err
	}

	// ── 4. Mapear a HistorialAsistenciaFila ──────────────────────
	filas := make([]HistorialAsistenciaFila, len(rawFilas))
	for i, r := range rawFilas {
		filas[i] = HistorialAsistenciaFila{
			FilaOrden:         r.FilaOrden,
			Fecha:             r.Fecha,
			Hora:              r.Hora,
			DiaSemana:         r.DiaSemana,
			TipoMarcacion:     r.TipoMarcacion,
			PuntoMarcacion:    r.PuntoMarcacion,
			Estado:            r.Estado,
			HorarioAsignado:   r.HorarioAsignado,
			HoraProgramada:    r.HoraProgramada,
			MinutosDiferencia: r.MinutosDiferencia,
			Observaciones:     r.Observaciones,
		}
	}

	// ── 5. Construir reporte ──────────────────────────────────────
	return &HistorialAsistenciaReporte{
		Empresa: empresa,
		Cabecera: HistorialAsistenciaCabecera{
			IDTercero:       req.IDTercero,
			NombreEmpleado:  buildNombreCompleto(t.PrimerNombre, t.SegundoNombre, t.PrimerApellido, t.SegundoApellido),
			NumeroDocumento: t.NumeroDocumento,
			FechaInicio:     req.FechaInicio,
			FechaFin:        req.FechaFin,
			TotalRegistros:  len(filas),
		},
		Filas: filas,
	}, nil
}

// buildNombreCompleto construye el nombre completo en Go
// omitiendo partes vacías para evitar espacios dobles.
func buildNombreCompleto(primerNombre, segundoNombre, primerApellido, segundoApellido string) string {
	result := ""
	for _, p := range []string{primerNombre, segundoNombre, primerApellido, segundoApellido} {
		if p == "" {
			continue
		}
		if result != "" {
			result += " "
		}
		result += p
	}
	return result
}
