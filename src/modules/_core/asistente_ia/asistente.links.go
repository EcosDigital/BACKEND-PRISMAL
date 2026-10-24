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
	// Ventas e Inventario van primero porque el enlace se toma del primer
	// tema que coincide, y sus preguntas suelen mencionar también
	// "proveedor", "cliente", "bodega" o "artículo".
	{
		Keywords:   []string{"orden", "órdenes", "ordenes", "mesa", "mesas", "comanda", "cocina"},
		ManualFile: "ventas-ordenes.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Órdenes", Path: "/ventas/ordenes"},
	},
	{
		Keywords:   []string{"pedido", "pedidos", "domicilios", "mi llave"},
		ManualFile: "ventas-pedidos.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Pedidos", Path: "/ventas/search-pedidos"},
	},
	{
		Keywords:   []string{"fiado", "fiados", "fiar", "fío", "abono", "abonar", "cartera", "deuda", "me debe"},
		ManualFile: "ventas-fiados.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Fiados", Path: "/ventas/search-fiados"},
	},
	{
		Keywords:   []string{"kardex", "kárdex", "historial de movimientos"},
		ManualFile: "inventario-kardex.md",
		Enlace:     EnlaceSugerido{Label: "Ir al Kárdex", Path: "/inventario/kardex"},
	},
	{
		Keywords:   []string{"entrada", "entradas", "ingreso de mercancía", "ingresar mercancía", "meter al inventario"},
		ManualFile: "inventario-entradas.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Entradas", Path: "/inventario/search-entradas"},
	},
	{
		Keywords:   []string{" baja", "bajas", "dar de baja", "vencid", "dañad"},
		ManualFile: "inventario-bajas.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Bajas", Path: "/inventario/search-bajas"},
	},
	{
		Keywords:   []string{"traslado", "traslados", "trasladar"},
		ManualFile: "inventario-traslados.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Traslados", Path: "/inventario/search-traslados"},
	},
	{
		Keywords:   []string{"catálogo", "catalogo", "precio", "domicilio"},
		ManualFile: "inventario-catalogo.md",
		Enlace:     EnlaceSugerido{Label: "Ir al Catálogo", Path: "/inventario/search-catalogos"},
	},
	{
		Keywords:   []string{"bodega", "bodegas", "almacén", "almacen", "stock mínimo"},
		ManualFile: "inventario-bodegas.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Bodegas", Path: "/inventario/search-bodegas"},
	},
	{
		Keywords:   []string{"artículo", "articulo", "producto", "insumo"},
		ManualFile: "inventario-articulos.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Artículos", Path: "/inventario/search-articulos"},
	},
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
	{
		Keywords: []string{
			"usuario", "usuarios",
		},
		ManualFile: "gestionar-usuarios.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Usuarios", Path: "/config/search-user"},
	},
	{
		Keywords: []string{
			"roles", " rol ", "permiso", "permisos",
		},
		ManualFile: "gestionar-roles.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Roles", Path: "/config/search-roles"},
	},
	{
		Keywords: []string{
			"agente", "agentes",
		},
		ManualFile: "gestionar-agentes.md",
		Enlace:     EnlaceSugerido{Label: "Ir a Agentes", Path: "/agentes"},
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
