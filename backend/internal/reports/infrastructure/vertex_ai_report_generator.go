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
			Temperature:     genai.Ptr(float32(0.2)),
			MaxOutputTokens: 180,
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
Eres un asistente académico que genera reportes semanales breves para profesores.

Reglas:
- Usa únicamente la información del JSON.
- No inventes nombres, tareas, fechas, horas ni observaciones.
- No califiques el desempeño.
- Escribe en español.
- Máximo 120 palabras.
- Entrega solo el texto final del reporte.

Incluye:
1. Actividades realizadas.
2. Horas reportadas.
3. Observaciones relevantes.
4. Recomendación breve de seguimiento.

JSON:
%s
`, inputJSON)
}