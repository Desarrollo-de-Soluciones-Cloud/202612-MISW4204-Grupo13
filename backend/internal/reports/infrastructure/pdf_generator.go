package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) Generate(filePath string, title string, lines []string) error {
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("file path is required")
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(16, 16, 16)
	pdf.SetAutoPageBreak(true, 18)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "B", 17)
	pdf.CellFormat(0, 10, tr(title), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	sections := parsePDFSections(lines)

	writeInfoTable(pdf, tr, sections.info)
	writeAISection(pdf, tr, sections.aiText)
	writeTasksTable(pdf, tr, sections.tasks)

	return pdf.OutputFileAndClose(filePath)
}

type parsedPDFSections struct {
	info   map[string]string
	aiText string
	tasks  []taskPDFRow
}

type taskPDFRow struct {
	Title       string
	Status      string
	Hours       string
	Description string
	Observation string
}

func parsePDFSections(lines []string) parsedPDFSections {
	result := parsedPDFSections{
		info:  make(map[string]string),
		tasks: make([]taskPDFRow, 0),
	}

	currentSection := ""
	var currentTask *taskPDFRow

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)

		switch line {
		case "Síntesis generada por IA:":
			currentSection = "ai"
			continue
		case "Tareas reportadas:":
			currentSection = "tasks"
			continue
		}

		if strings.HasPrefix(line, "Workspace:") {
			result.info["Workspace"] = strings.TrimSpace(strings.TrimPrefix(line, "Workspace:"))
			continue
		}

		if strings.HasPrefix(line, "Monitor/Asistente:") {
			result.info["Monitor/Asistente"] = strings.TrimSpace(strings.TrimPrefix(line, "Monitor/Asistente:"))
			continue
		}

		if strings.HasPrefix(line, "Rol:") {
			result.info["Rol"] = strings.TrimSpace(strings.TrimPrefix(line, "Rol:"))
			continue
		}

		if strings.HasPrefix(line, "Semana:") {
			result.info["Semana"] = strings.TrimSpace(strings.TrimPrefix(line, "Semana:"))
			continue
		}

		if strings.HasPrefix(line, "Horas reportadas:") {
			result.info["Horas reportadas"] = strings.TrimSpace(strings.TrimPrefix(line, "Horas reportadas:"))
			continue
		}

		if currentSection == "ai" {
			if line != "" {
				if result.aiText != "" {
					result.aiText += " "
				}
				result.aiText += line
			}
			continue
		}

		if currentSection == "tasks" {
			if strings.HasPrefix(line, "- ") {
				if currentTask != nil {
					result.tasks = append(result.tasks, *currentTask)
				}

				currentTask = parseTaskHeader(line)
				continue
			}

			if currentTask != nil && strings.HasPrefix(line, "descripcion:") {
				currentTask.Description = strings.TrimSpace(strings.TrimPrefix(line, "descripcion:"))
				continue
			}

			if currentTask != nil && strings.HasPrefix(line, "observaciones:") {
				currentTask.Observation = strings.TrimSpace(strings.TrimPrefix(line, "observaciones:"))
				continue
			}
		}
	}

	if currentTask != nil {
		result.tasks = append(result.tasks, *currentTask)
	}

	return result
}

func parseTaskHeader(line string) *taskPDFRow {
	line = strings.TrimPrefix(line, "- ")
	parts := strings.Split(line, "|")

	row := &taskPDFRow{
		Title: strings.TrimSpace(parts[0]),
	}

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "estado=") {
			row.Status = strings.TrimSpace(strings.TrimPrefix(part, "estado="))
		}

		if strings.HasPrefix(part, "horas=") {
			row.Hours = strings.TrimSpace(strings.TrimPrefix(part, "horas="))
		}
	}

	return row
}

func writeInfoTable(pdf *gofpdf.Fpdf, tr func(string) string, info map[string]string) {
	writeSectionTitle(pdf, tr, "Información general")

	rows := [][]string{
		{"Workspace", info["Workspace"]},
		{"Monitor/Asistente", info["Monitor/Asistente"]},
		{"Rol", info["Rol"]},
		{"Semana", info["Semana"]},
		{"Horas reportadas", info["Horas reportadas"]},
	}

	widths := []float64{48, 127}

	pdf.SetFont("Arial", "", 9)

	for _, row := range rows {
		writeWrappedRow(pdf, tr, row, widths, 5, []string{"B", ""}, []bool{true, false})
	}

	pdf.Ln(5)
}

func writeAISection(pdf *gofpdf.Fpdf, tr func(string) string, text string) {
	writeSectionTitle(pdf, tr, "Síntesis generada por IA")

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 6, tr(text), "", "L", false)
	pdf.Ln(5)
}

func writeTasksTable(pdf *gofpdf.Fpdf, tr func(string) string, tasks []taskPDFRow) {
	writeSectionTitle(pdf, tr, "Tareas reportadas")

	if len(tasks) == 0 {
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, tr("No se registraron tareas."), "", "L", false)
		return
	}

	headers := []string{"Tarea", "Estado", "Horas", "Observación / descripción"}
	widths := []float64{54, 27, 18, 76}

	writeTableHeader(pdf, tr, headers, widths)

	for _, task := range tasks {
		observation := task.Observation
		if strings.TrimSpace(observation) == "" {
			observation = task.Description
		}

		row := []string{
			task.Title,
			task.Status,
			task.Hours,
			observation,
		}

		if !hasSpaceForRow(pdf, tr, row, widths, 4.5) {
			pdf.AddPage()
			writeTableHeader(pdf, tr, headers, widths)
		}

		writeWrappedRow(pdf, tr, row, widths, 4.5, []string{"", "", "", ""}, []bool{false, false, false, false})
	}
}

func writeTableHeader(pdf *gofpdf.Fpdf, tr func(string) string, headers []string, widths []float64) {
	ensureSpace(pdf, 12)

	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(230, 230, 230)

	for i, header := range headers {
		pdf.CellFormat(widths[i], 8, tr(header), "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)
}

func writeWrappedRow(
	pdf *gofpdf.Fpdf,
	tr func(string) string,
	cells []string,
	widths []float64,
	lineHeight float64,
	styles []string,
	fills []bool,
) {
	wrappedCells := make([][]string, len(cells))
	maxLines := 1

	for i, cell := range cells {
		if i < len(styles) && styles[i] != "" {
			pdf.SetFont("Arial", styles[i], 9)
		} else {
			pdf.SetFont("Arial", "", 8)
		}

		wrappedCells[i] = splitCellText(pdf, tr(cell), widths[i]-2)

		if len(wrappedCells[i]) > maxLines {
			maxLines = len(wrappedCells[i])
		}
	}

	rowHeight := float64(maxLines)*lineHeight + 3

	ensureSpace(pdf, rowHeight)

	startX, startY := pdf.GetXY()

	for i, cellLines := range wrappedCells {
		x := startX
		for j := 0; j < i; j++ {
			x += widths[j]
		}

		if i < len(styles) && styles[i] != "" {
			pdf.SetFont("Arial", styles[i], 9)
		} else {
			pdf.SetFont("Arial", "", 8)
		}

		fill := false
		if i < len(fills) {
			fill = fills[i]
		}

		if fill {
			pdf.SetFillColor(238, 238, 238)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.Rect(x, startY, widths[i], rowHeight, "D")

		if fill {
			pdf.Rect(x, startY, widths[i], rowHeight, "F")
			pdf.Rect(x, startY, widths[i], rowHeight, "D")
		}

		pdf.SetXY(x+1, startY+1.5)

		for _, cellLine := range cellLines {
			pdf.CellFormat(widths[i]-2, lineHeight, cellLine, "", 2, "L", false, 0, "")
		}
	}

	pdf.SetXY(startX, startY+rowHeight)
}

func splitCellText(pdf *gofpdf.Fpdf, text string, width float64) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{""}
	}

	lines := pdf.SplitLines([]byte(text), width)
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		result = append(result, string(line))
	}

	if len(result) == 0 {
		return []string{""}
	}

	return result
}

func hasSpaceForRow(pdf *gofpdf.Fpdf, tr func(string) string, cells []string, widths []float64, lineHeight float64) bool {
	maxLines := 1

	pdf.SetFont("Arial", "", 8)

	for i, cell := range cells {
		lines := splitCellText(pdf, tr(cell), widths[i]-2)
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	rowHeight := float64(maxLines)*lineHeight + 3

	_, y := pdf.GetXY()
	_, pageHeight := pdf.GetPageSize()
	_, _, _, bottomMargin := pdf.GetMargins()

	return y+rowHeight <= pageHeight-bottomMargin
}

func writeSectionTitle(pdf *gofpdf.Fpdf, tr func(string) string, title string) {
	ensureSpace(pdf, 14)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, tr(title), "", 1, "L", false, 0, "")
}

func ensureSpace(pdf *gofpdf.Fpdf, requiredHeight float64) {
	_, y := pdf.GetXY()
	_, pageHeight := pdf.GetPageSize()
	_, _, _, bottomMargin := pdf.GetMargins()

	if y+requiredHeight > pageHeight-bottomMargin {
		pdf.AddPage()
	}
}