package scrapping

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// headers define el orden y nombre de las columnas del Excel.
// Se declara aquí para que cualquier cambio futuro sea en un solo lugar.
var headers = []string{
	"Nombre",
	"Tipo de Negocio",
	"Ciudad",
	"Dirección",
	"Teléfono",
	"Sitio Web",
	"Rating",
	"Total Reseñas",
	"Google Maps",
	"Estado",
}

// GenerarExcelLeads construye el archivo .xlsx en memoria y retorna el buffer.
// El controller lo sirve directamente como descarga sin tocar el disco.
func GenerarExcelLeads(leads []LeadExportRow, keyword, ciudad string) (*bytes.Buffer, error) {

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Leads"
	f.SetSheetName("Sheet1", sheet)

	// ── Estilo encabezado ────────────────────────────────────────────────────
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
			Size:  11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"1F4E79"}, // azul corporativo oscuro
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "FFFFFF", Style: 2},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creando estilo encabezado: %w", err)
	}

	// ── Estilo fila par ───────────────────────────────────────────────────────
	evenStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"EBF3FB"},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creando estilo par: %w", err)
	}

	// ── Estilo fila impar ─────────────────────────────────────────────────────
	oddStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"FFFFFF"},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creando estilo impar: %w", err)
	}

	// ── Metadatos en fila 1 ───────────────────────────────────────────────────
	meta := fmt.Sprintf("Búsqueda: %s  |  Ciudad: %s  |  Total leads: %d", keyword, ciudad, len(leads))
	f.SetCellValue(sheet, "A1", meta)
	f.MergeCell(sheet, "A1", colName(len(headers))+"1")

	// ── Encabezados en fila 2 ─────────────────────────────────────────────────
	for i, h := range headers {
		cell := colName(i+1) + "2"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 2, 30)

	// ── Datos desde fila 3 ────────────────────────────────────────────────────
	for rowIdx, lead := range leads {
		excelRow := rowIdx + 3
		style := oddStyle
		if rowIdx%2 == 0 {
			style = evenStyle
		}

		row := []interface{}{
			lead.Nombre,
			lead.TipoNegocio,
			lead.Ciudad,
			lead.Direccion,
			lead.Telefono,
			lead.SitioWeb,
			formatRating(lead.Rating),
			lead.TotalReviews,
			lead.MapsURL,
			lead.Estado,
		}

		for colIdx, val := range row {
			cell := colName(colIdx+1) + strconv.Itoa(excelRow)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, style)
		}

		f.SetRowHeight(sheet, excelRow, 20)
	}

	// ── Ancho de columnas ────────────────────────────────────────────────────
	widths := map[int]float64{
		1:  35, // Nombre
		2:  25, // Tipo negocio
		3:  18, // Ciudad
		4:  40, // Dirección
		5:  18, // Teléfono
		6:  35, // Sitio web
		7:  10, // Rating
		8:  14, // Total reseñas
		9:  50, // Google Maps URL
		10: 15, // Estado
	}
	for col, width := range widths {
		f.SetColWidth(sheet, colName(col), colName(col), width)
	}

	// ── Auto-filtro ───────────────────────────────────────────────────────────
	lastCol := colName(len(headers))
	lastRow := len(leads) + 2
	f.AutoFilter(sheet,
		fmt.Sprintf("A2:%s%d", lastCol, lastRow),
		[]excelize.AutoFilterOptions{},
	)

	// ── Freeze primera fila de datos ─────────────────────────────────────────
	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      2,
		TopLeftCell: "A3",
		ActivePane:  "bottomLeft",
	})

	// ── Serializar a buffer ───────────────────────────────────────────────────
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error serializando Excel: %w", err)
	}

	return buf, nil
}

// ─── Helpers privados ─────────────────────────────────────────────────────────

// colName convierte un índice base-1 al nombre de columna Excel (1→A, 27→AA…).
func colName(col int) string {
	name := ""
	for col > 0 {
		col--
		name = string(rune('A'+col%26)) + name
		col /= 26
	}
	return name
}

// formatRating muestra "—" si el rating es 0.
func formatRating(r float64) interface{} {
	if r == 0 {
		return "—"
	}
	return r
}
