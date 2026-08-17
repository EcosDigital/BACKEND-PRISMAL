package asistente

import "strings"

// EnlaceSugerido es un acceso directo que se ofrece junto con la
// respuesta del asistente cuando la pregunta corresponde a un tema del
// sistema con una pantalla asociada (ej. Terceros).
type EnlaceSugerido struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

// temaSistema agrupa, para un mismo tema del sistema, las palabras clave
// que lo identifican, el manual que lo documenta y el acceso directo a
// su pantalla. Se usa tanto para decidir qué manual(es) mandarle al
// modelo (para que no reciba temas que no vienen al caso, lo que además
// lo hace responder más rápido) como para enganchar el botón de acceso
// directo en el chat.
type temaSistema struct {
	Keywords   []string
	ManualFile string
	Enlace     EnlaceSugerido
}

var temasSistema = []temaSistema{
	{
		Keywords: []string{
			"tercero", "terceros",
			"cliente", "clientes",
			"proveedor", "proveedores",
			"colaborador", "colaboradores",
			"empleado", "empleados",
			"contratista", "contratistas",
			"distribuidor", "distribuidores",
		},
		ManualFile: "gestionar-terceros.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Terceros", Path: "/config/search-tercero"},
	},
	{
		Keywords: []string{
			"sede", "sedes",
			"sucursal", "sucursales",
		},
		ManualFile: "gestionar-sedes.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Sedes", Path: "/config/search-sedes"},
	},
	{
		Keywords: []string{
			"empresa", "empresas",
			"razon social", "razón social",
		},
		ManualFile: "gestionar-empresas.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Empresas", Path: "/config/search-company"},
	},
}

// ResolveEnlace busca si la pregunta del usuario corresponde a alguno de
// los temas con pantalla asociada. Devuelve nil si no aplica ninguno.
func ResolveEnlace(pregunta string) *EnlaceSugerido {
	preguntaLower := strings.ToLower(pregunta)

	for _, tema := range temasSistema {
		for _, keyword := range tema.Keywords {
			if strings.Contains(preguntaLower, keyword) {
				enlace := tema.Enlace
				return &enlace
			}
		}
	}

	return nil
}

// ResolveManualFiles busca qué manuales de knowledge/sistema son
// relevantes para la pregunta, según las mismas palabras clave que
// ResolveEnlace. Si ninguno coincide, devuelve nil — quien la llama
// decide el respaldo (hoy: mandar todos los manuales, para no dejar la
// pregunta sin contexto).
func ResolveManualFiles(pregunta string) []string {
	preguntaLower := strings.ToLower(pregunta)

	seen := map[string]bool{}
	var archivos []string

	for _, tema := range temasSistema {
		for _, keyword := range tema.Keywords {
			if strings.Contains(preguntaLower, keyword) {
				if !seen[tema.ManualFile] {
					seen[tema.ManualFile] = true
					archivos = append(archivos, tema.ManualFile)
				}
				break
			}
		}
	}

	return archivos
}
