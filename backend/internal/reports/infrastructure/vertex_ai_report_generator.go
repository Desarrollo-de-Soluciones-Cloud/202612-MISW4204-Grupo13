package infrastructure

import (
	reportsApplication "backend/internal/reports/application"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

type VertexAIReportGenerator struct {
	client *genai.Client
	model  string
}

func NewVertexAIReportGenerator(projectID, location, model string) (*VertexAIReportGenerator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("gcp project id is required")
	}

	if strings.TrimSpace(location) == "" {
		location = "us-central1"
	}

	if strings.TrimSpace(model) == "" {
		model = "gemini-2.5-flash-lite"
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, err
	}

	return &VertexAIReportGenerator{
		client: client,
		model:  model,
	}, nil
}

func (g *VertexAIReportGenerator) GenerateWeeklyReport(input reportsApplication.AIWeeklyReportInput) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	prompt := buildWeeklyReportPrompt(string(inputJSON))

	response, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature:     genai.Ptr(float32(0.25)),
			MaxOutputTokens: 220,
		},
	)
	if err != nil {
		return "", err
	}

	reportText := strings.TrimSpace(response.Text())
	if reportText == "" {
		return "", fmt.Errorf("empty ai report")
	}

	return reportText, nil
}

func buildWeeklyReportPrompt(inputJSON string) string {
	return fmt.Sprintf(`
Eres un asistente académico que redacta reportes semanales breves para profesores universitarios.

A partir del JSON, redacta una síntesis clara, natural y útil. No copies literalmente cada campo; transforma la información en un reporte ejecutivo breve.

Reglas:
- Usa únicamente la información del JSON.
- No inventes nombres, fechas, tareas, horas ni resultados.
- No califiques el desempeño con juicios exagerados.
- No uses Markdown, asteriscos, listas, encabezados ni títulos.
- Escribe en español formal y natural.
- Máximo 130 palabras.
- Entrega solo el cuerpo del reporte.

El texto debe mencionar de forma natural:
- El monitor o asistente.
- El workspace.
- El tipo de actividades realizadas.
- Las horas reportadas.
- Observaciones o pendientes relevantes.
- Una recomendación breve de seguimiento para el profesor.

JSON:
%s
`, inputJSON)
}