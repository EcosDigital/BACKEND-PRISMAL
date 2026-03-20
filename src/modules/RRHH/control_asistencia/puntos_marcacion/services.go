package control_asistencia

import (
	"errors"
	"fmt"
	"math"

	"gorm.io/gorm"
)

func RegisterPuntoMarcacion(db *gorm.DB, req *PuntoMarcacionRequest) (int64, error) {

	existe, err := FindPuntoMarcacionByNombre(db, req.Nombre, 0)
	if err != nil {
		return 0, err
	}
	if existe {
		return 0, errors.New("ya existe un punto de marcación con ese nombre")
	}

	return CreatePuntoMarcacion(db, req)
}

func FilterPuntosMarcacion(db *gorm.DB, empresaID int64) ([]PuntoMarcacionListResponse, error) {
	return FindPuntosMarcacion(db, empresaID)
}

func FilterPuntoMarcacionByID(db *gorm.DB, id int64) (*PuntoMarcacionResponse, error) {

	result, err := FindPuntoMarcacionByID(db, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("punto de marcación con id %d no encontrado", id)
	}

	return result, nil
}

func UpdatePuntoMarcacion(db *gorm.DB, id int64, req *UpdatePuntoMarcacionRequest) (*PuntoMarcacionResponse, error) {

	// Reutiliza FilterPuntoMarcacionByID que ya maneja el caso not found con error
	existing, err := FindPuntoMarcacionByID(db, id)
	if err != nil {
		return nil, err
	}

	// Validar unicidad del nombre solo si está cambiando
	if req.Nombre != "" && req.Nombre != existing.Nombre {
		existe, err := FindPuntoMarcacionByNombre(db, req.Nombre, id)
		if err != nil {
			return nil, err
		}
		if existe {
			return nil, errors.New("ya existe un punto de marcación con ese nombre")
		}
	}

	return EditPuntoMarcacion(db, id, req)
}

func ProcessCambioEstado(db *gorm.DB, id int64, req *CambioEstadoRequest) error {
	return ChangePuntoMarcacionEstado(db, id, req.Activo)
}

func FilterTiposPuntoMarcacion(db *gorm.DB) ([]RefTipoPuntoMarcacion, error) {
	return FindTiposPuntoMarcacion(db)
}

// ===== MARCACIONES (MOV) ==== ///
const (
	// segundosEntreMaraciones define el tiempo mínimo obligatorio entre dos
	// marcaciones del mismo empleado en el mismo punto.
	segundosEntreMarcaciones = 60

	// radioTierraKm es el radio medio de la Tierra en kilómetros.
	radioTierraKm = 6371.0
)

// ─── Errores ESPERRADOS  ───────────────────────────────────────────────────────

var (
	ErrTokenInvalido        = errors.New("token de marcación inválido o inactivo")
	ErrFueraDeRadio         = errors.New("ubicación fuera del radio permitido para este punto de marcación")
	ErrMarcacionDuplicada   = errors.New("ya registraste una marcación en este punto recientemente, espera un momento")
	ErrOrigenNoConfigurado  = errors.New("el origen QR no está configurado en el sistema")
	ErrCoordenadasInvalidas = errors.New("las coordenadas enviadas no son válidas")
)

// ─── funciones de calculos  ────────────────────────────────────────────────────────────────

func calcularDistanciaKm(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return radioTierraKm * c

}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}

//-------- validaciones Haversine -- //

// coordenadasValidas verifica que latitud y longitud estén dentro de rangos
// geográficos posibles. No confía ciegamente en el frontend.
func coordenadasValidas(lat, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

// ─── Service principal ────────────────────────────────────────────────────────

// ProcessMarcacion ejecuta el flujo completo de registro de una marcación QR:
// validación de token → geofencing → anti-duplicado → inserción.

func ProcessMarcacion(db *gorm.DB, req *MarcacionRequest) (*MarcacionResponse, error) {

	// 1. Validar coordenadas del request (no confiar en el frontend)
	if !coordenadasValidas(req.Latitud, req.Longitud) {
		return nil, ErrCoordenadasInvalidas
	}

	// 2. Validar token — buscar punto activo en BD
	punto, err := FindPuntoByToken(db, req.Token)
	if err != nil {
		return nil, err
	}
	if punto == nil {
		return nil, ErrTokenInvalido
	}

	// 3. Validar geolocalización con Haversine
	distanciaKm := calcularDistanciaKm(
		req.Latitud, req.Longitud,
		punto.Latitud, punto.Longitud,
	)
	distanciaMetros := distanciaKm * 1000

	const maxPrecisionTolerable = 200
	precisionEfectiva := req.Precision
	if precisionEfectiva > maxPrecisionTolerable {
		precisionEfectiva = maxPrecisionTolerable
	}
	radioEfectivo := float64(punto.RadioMetros) + float64(precisionEfectiva)

	if distanciaMetros > radioEfectivo {
		return nil, ErrFueraDeRadio
	}

	// 4. Validar duplicidad — evitar marcaciones en menos de N segundos
	duplicada, err := ExisteMarcacionReciente(db, req.IDEmpleado, punto.ID, segundosEntreMarcaciones)
	if err != nil {
		return nil, err
	}
	if duplicada {
		return nil, ErrMarcacionDuplicada
	}

	// 5. Obtener id del origen QR desde ref_origen_marcacion
	idOrigen, err := FindOrigenQR(db)
	if err != nil {
		return nil, err
	}
	if idOrigen == 0 {
		return nil, ErrOrigenNoConfigurado
	}

	// 6. Insertar marcación
	id, fechaHora, err := InsertMarcacion(db, req, punto.ID, idOrigen)
	if err != nil {
		return nil, err
	}

	return &MarcacionResponse{
		ID:               id,
		IDEmpleado:       req.IDEmpleado,
		IDPuntoMarcacion: punto.ID,
		FechaHora:        fechaHora.Format("2006-01-02T15:04:05"),
		TipoEvento:       "ASISTENCIA", // cambiar por el ID ref tipo marcacion
		Mensaje:          "Marcación registrada exitosamente",
	}, nil

}
