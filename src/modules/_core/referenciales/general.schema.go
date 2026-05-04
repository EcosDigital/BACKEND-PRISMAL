package referenciales

type TipoPersonaDB struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type TipoDocumentoDB struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Codigo string `json:"codigo"`
	Tipo   int    `json:"tipo"`
}

type GeneroDB struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type DepartamentoDB struct {
	ID     int    `json:"id"`
	Codigo int    `json:"codigo"`
	Nombre string `json:"nombre"`
}

type MunicipiosDB struct {
	ID                 int    `json:"id"`
	CodigoDepartamento int    `json:"codigo_departamento"`
	CodigoMunicipio    int    `json:"codigo_municipio"`
	NombreMunicipio    string `json:"nombre_municipio"`
}

type ZonaRuralesDB struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

type ActividadesEconomicas struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

type AmbitosTerceros struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Centralizaciones struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type ResponsabilidadDian struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type TipoContribuyente struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type RegimenIva struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type RegimenDian struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type ClasePersonas struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Paises struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type TipoEmpresa struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type NaturalezaEmpresa struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type TotalClientes struct {
	Total  int    `json:"total"`
	Codigo int    `json:"codigo"`
	Nombre string `json:"nombre"`
}

type FiltrosGenerales struct {
	IdTipo       *int    `json:"id_tipo,omitempty"`
	Estado       *bool   `json:"estado,omitempty"`
	NombreCodigo *string `json:"nombre_codigo,omitempty"`
	FechaInicio  *string `json:"fecha_inicio,omitempty"`
	FechaFin     *string `json:"fecha_final,omitempty"`
}

type EstadoModulo struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type TipoSede struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type TipoRol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// === COMPROBANTES ==//
type ModuleComprobantes struct {
	ID     int    `json:"id"`
	Codigo string `json:"codigo"`
	IDRef  int    `json:"id_ref"`
	Nombre string `json:"nombre"`
}

type TipoOperacion struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// ==== GESTIONES == //

type PrioridadTicket struct {
	ID          int    `json:"id"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	ColorHex    string `json:"color_hex"`
	Orden       string `json:"orden"`
	Nombre      string `json:"nombre"`
}

type EstadoTicket struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	ColorHex string `json:"color_hex"`
	Orden    string `json:"orden"`
}

type Colaboradores struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// == INVENTARIO === ///

type TipoBodega struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type UnidadeMedida struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type GrupoArticulos struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Presentacionrticulos struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type ResultGeneral struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}
