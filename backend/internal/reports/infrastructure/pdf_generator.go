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

const (
	pdfAISectionTitle        = "SÃ­ntesis generada por IA:"
	pdfTasksSectionTitle     = "Tareas reportadas:"
	pdfMonitorAssistantLabel = "Monitor/Asistente"
	pdfReportedHoursLabel    = "Horas reportadas"
	pdfDescriptionPrefix     = "descripcion:"
	pdfObservationsPrefix    = "observaciones:"
)

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
		case pdfAISectionTitle:
			currentSection = "ai"
			continue
		case pdfTasksSectionTitle:
			currentSection = "tasks"
			continue
		}

		if key, value, ok := parsePDFInfoLine(line); ok {
			result.info[key] = value
			continue
		}

		switch currentSection {
		case "ai":
			appendAIText(&result.aiText, line)
		case "tasks":
			currentTask = parseTaskSectionLine(line, currentTask, &result.tasks)
		}
	}

	if currentTask != nil {
		result.tasks = append(result.tasks, *currentTask)
	}

	return result
}

func parsePDFInfoLine(line string) (string, string, bool) {
	infoPrefixes := []struct {
		key    string
		prefix string
	}{
		{key: "Workspace", prefix: "Workspace:"},
		{key: pdfMonitorAssistantLabel, prefix: pdfMonitorAssistantLabel + ":"},
		{key: "Rol", prefix: "Rol:"},
		{key: "Semana", prefix: "Semana:"},
		{key: pdfReportedHoursLabel, prefix: pdfReportedHoursLabel + ":"},
	}

	for _, item := range infoPrefixes {
		if strings.HasPrefix(line, item.prefix) {
			return item.key, strings.TrimSpace(strings.TrimPrefix(line, item.prefix)), true
		}
	}

	return "", "", false
}

func appendAIText(aiText *string, line string) {
	if line == "" {
		return
	}

	if *aiText != "" {
		*aiText += " "
	}

	*aiText += line
}

func parseTaskSectionLine(line string, currentTask *taskPDFRow, tasks *[]taskPDFRow) *taskPDFRow {
	if strings.HasPrefix(line, "- ") {
		if currentTask != nil {
			*tasks = append(*tasks, *currentTask)
		}
		return parseTaskHeader(line)
	}

	if currentTask == nil {
		return currentTask
	}

	switch {
	case strings.HasPrefix(line, pdfDescriptionPrefix):
		currentTask.Description = strings.TrimSpace(strings.TrimPrefix(line, pdfDescriptionPrefix))
	case strings.HasPrefix(line, pdfObservationsPrefix):
		currentTask.Observation = strings.TrimSpace(strings.TrimPrefix(line, pdfObservationsPrefix))
	}

	return currentTask
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
	writeSectionTitle(pdf, tr, "InformaciÃ³n general")

	rows := [][]string{
		{"Workspace", info["Workspace"]},
		{pdfMonitorAssistantLabel, info[pdfMonitorAssistantLabel]},
		{"Rol", info["Rol"]},
		{"Semana", info["Semana"]},
		{pdfReportedHoursLabel, info[pdfReportedHoursLabel]},
	}

	widths := []float64{48, 127}

	pdf.SetFont("Arial", "", 9)

	for _, row := range rows {
		writeWrappedRow(pdf, tr, row, widths, 5, []string{"B", ""}, []bool{true, false})
	}

	pdf.Ln(5)
}

func writeAISection(pdf *gofpdf.Fpdf, tr func(string) string, text string) {
	writeSectionTitle(pdf, tr, "SÃ­ntesis generada por IA")

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

	headers := []string{"Tarea", "Estado", "Horas", "ObservaciÃ³n / descripciÃ³n"}
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
		setRowCellFont(pdf, styles, i)
		wrappedCells[i] = splitCellText(pdf, tr(cell), widths[i]-2)

		if len(wrappedCells[i]) > maxLines {
			maxLines = len(wrappedCells[i])
		}
	}

	rowHeight := float64(maxLines)*lineHeight + 3

	ensureSpace(pdf, rowHeight)

	startX, startY := pdf.GetXY()

	for i, cellLines := range wrappedCells {
		x := calculateRowStartX(startX, widths, i)
		setRowCellFont(pdf, styles, i)
		fill := shouldFillCell(fills, i)
		drawWrappedCellBackground(pdf, x, startY, widths[i], rowHeight, fill)

		pdf.SetXY(x+1, startY+1.5)

		for _, cellLine := range cellLines {
			pdf.CellFormat(widths[i]-2, lineHeight, cellLine, "", 2, "L", false, 0, "")
		}
	}

	pdf.SetXY(startX, startY+rowHeight)
}

func setRowCellFont(pdf *gofpdf.Fpdf, styles []string, index int) {
	if index < len(styles) && styles[index] != "" {
		pdf.SetFont("Arial", styles[index], 9)
		return
	}

	pdf.SetFont("Arial", "", 8)
}

func calculateRowStartX(startX float64, widths []float64, index int) float64 {
	x := startX
	for i := 0; i < index; i++ {
		x += widths[i]
	}
	return x
}

func shouldFillCell(fills []bool, index int) bool {
	return index < len(fills) && fills[index]
}

func drawWrappedCellBackground(
	pdf *gofpdf.Fpdf,
	x float64,
	y float64,
	width float64,
	height float64,
	fill bool,
) {
	if fill {
		pdf.SetFillColor(238, 238, 238)
		pdf.Rect(x, y, width, height, "F")
		pdf.Rect(x, y, width, height, "D")
		return
	}

	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(x, y, width, height, "D")
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
