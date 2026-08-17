package asistente

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BuildContext arma el contexto que se envía al modelo: identidad.md
// (tono/reglas del asistente, siempre) + solo los manuales relevantes a
// la pregunta (según ResolveManualFiles). Mandar únicamente el manual
// del tema evita que el modelo mezcle contenido de temas distintos y
// hace que responda más rápido, al procesar menos texto. Si ningún tema
// coincide con la pregunta, se mandan todos los manuales como respaldo.
//
// Se lee del disco en cada solicitud para que editar un manual no
// requiera recompilar ni reiniciar el servidor.
func BuildContext(pregunta string) (string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error obteniendo directorio de trabajo: %v", err)
	}

	iaDir := filepath.Join(wd, "ia")

	identidad, err := os.ReadFile(filepath.Join(iaDir, "identidad.md"))
	if err != nil {
		return "", fmt.Errorf("error leyendo identidad.md: %v", err)
	}

	var sb strings.Builder
	sb.Write(identidad)
	sb.WriteString("\n\n")

	sistemaDir := filepath.Join(iaDir, "knowledge", "sistema")

	files, err := resolveManualPaths(sistemaDir, pregunta)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("error leyendo manual '%s': %v", filepath.Base(file), err)
		}
		sb.Write(content)
		sb.WriteString("\n\n")
	}

	return sb.String(), nil
}

// resolveManualPaths decide qué manuales cargar para esta pregunta.
func resolveManualPaths(sistemaDir, pregunta string) ([]string, error) {

	nombres := ResolveManualFiles(pregunta)

	if len(nombres) > 0 {
		rutas := make([]string, 0, len(nombres))
		for _, nombre := range nombres {
			rutas = append(rutas, filepath.Join(sistemaDir, nombre))
		}
		return rutas, nil
	}

	// Ningún tema coincidió: como respaldo, se mandan todos los manuales
	// para no dejar la pregunta sin ningún contexto.
	files, err := filepath.Glob(filepath.Join(sistemaDir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("error listando manuales: %v", err)
	}

	sort.Strings(files)

	return files, nil
}
