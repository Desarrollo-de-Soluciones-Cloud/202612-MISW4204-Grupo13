package infrastructure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jung-kurt/gofpdf"
)

func TestGenerateRejectsMissingFilePath(t *testing.T) {
	generator := NewPDFGenerator()

	if err := generator.Generate("", "Weekly report", nil); err == nil {
		t.Fatalf("expected file path required error")
	}
}

func TestGenerateCreatesPDFFile(t *testing.T) {
	generator := NewPDFGenerator()
	filePath := filepath.Join(t.TempDir(), "reports", "weekly.pdf")
	lines := []string{
		"Workspace: Algorithms",
		"Monitor/Asistente: Ana Gomez",
		"Rol: assistant",
		"Semana: 3",
		"Horas reportadas: 6",
		pdfAISectionTitle,
		"Resumen generado",
		pdfTasksSectionTitle,
		"- Preparar lab | estado=finalizado | horas=2",
		"descripcion: Slides",
		"observaciones: Todo bien",
	}

	if err := generator.Generate(filePath, "Weekly report", lines); err != nil {
		t.Fatalf("expected no error generating pdf, got %v", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("expected generated file, got %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("expected generated pdf to have content")
	}
}

func TestParsePDFSections(t *testing.T) {
	sections := parsePDFSections([]string{
		"Workspace: Algorithms",
		"Monitor/Asistente: Ana Gomez",
		"Rol: assistant",
		"Semana: 3",
		"Horas reportadas: 6",
		pdfAISectionTitle,
		"Linea uno",
		"Linea dos",
		pdfTasksSectionTitle,
		"- Preparar lab | estado=finalizado | horas=2",
		"descripcion: Slides",
		"observaciones: Todo bien",
	})

	if sections.info["Workspace"] != "Algorithms" {
		t.Fatalf("expected workspace info, got %#v", sections.info)
	}
	if sections.aiText != "Linea uno Linea dos" {
		t.Fatalf("expected combined ai text, got %q", sections.aiText)
	}
	if len(sections.tasks) != 1 {
		t.Fatalf("expected 1 parsed task, got %d", len(sections.tasks))
	}
	if sections.tasks[0].Observation != "Todo bien" {
		t.Fatalf("expected parsed observation, got %#v", sections.tasks[0])
	}
}

func TestParseTaskHelpers(t *testing.T) {
	key, value, ok := parsePDFInfoLine("Workspace: Algorithms")
	if !ok || key != "Workspace" || value != "Algorithms" {
		t.Fatalf("unexpected parsed info line: %q %q %v", key, value, ok)
	}

	row := parseTaskHeader("- Preparar lab | estado=finalizado | horas=2")
	if row.Title != "Preparar lab" || row.Status != "finalizado" || row.Hours != "2" {
		t.Fatalf("unexpected parsed header: %#v", row)
	}

	var aiText string
	appendAIText(&aiText, "Linea uno")
	appendAIText(&aiText, "Linea dos")
	if aiText != "Linea uno Linea dos" {
		t.Fatalf("unexpected ai text %q", aiText)
	}
}

func TestTaskSectionAndLayoutHelpers(t *testing.T) {
	tasks := make([]taskPDFRow, 0)
	current := parseTaskSectionLine("- Preparar lab | estado=finalizado | horas=2", nil, &tasks)
	current = parseTaskSectionLine("descripcion: Slides", current, &tasks)
	current = parseTaskSectionLine("observaciones: Todo bien", current, &tasks)
	if current == nil || current.Description != "Slides" || current.Observation != "Todo bien" {
		t.Fatalf("unexpected current task %#v", current)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 8)

	lines := splitCellText(pdf, "Texto de prueba", 40)
	if len(lines) == 0 {
		t.Fatalf("expected wrapped lines")
	}

	if x := calculateRowStartX(10, []float64{20, 30, 40}, 2); x != 60 {
		t.Fatalf("expected row start x 60, got %v", x)
	}

	if !shouldFillCell([]bool{false, true}, 1) {
		t.Fatalf("expected fill for cell 1")
	}

	if !hasSpaceForRow(pdf, func(s string) string { return s }, []string{"a", "b"}, []float64{20, 20}, 4) {
		t.Fatalf("expected enough space for small row")
	}
}
