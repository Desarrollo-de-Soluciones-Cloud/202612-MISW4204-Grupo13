package application

import (
	"bytes"
	"context"
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GenerateWeeklyReports struct {
	reportRepo        reportsDomain.ReportRepository
	workspaceReader   WorkspaceReader
	weekReader        WeekReader
	assignmentReader  AssignmentReader
	taskReader        TaskReader
	userReader        UserReader
	pdfGenerator      PDFGenerator
	aiReportGenerator AIReportGenerator
	reportFileStorage ReportFileStorage
	reportsStorageDir string
	reportsGCSPrefix  string
	now               func() time.Time
}

type GenerateWeeklyReportsDependencies struct {
	ReportRepo        reportsDomain.ReportRepository
	WorkspaceReader   WorkspaceReader
	WeekReader        WeekReader
	AssignmentReader  AssignmentReader
	TaskReader        TaskReader
	UserReader        UserReader
	PDFGenerator      PDFGenerator
	AIReportGenerator AIReportGenerator
	ReportFileStorage ReportFileStorage
}

type GenerateWeeklyReportsOptions struct {
	ReportsStorageDir string
	ReportsGCSPrefix  string
	Now               func() time.Time
}

type GenerateWeeklyReportsInput struct {
	WorkspaceID uint
	WeekID      uint
}

type GenerateWeeklyReportsOutput struct {
	Reports []ReportOutput `json:"reports"`
}

func NewGenerateWeeklyReports(
	deps GenerateWeeklyReportsDependencies,
	options *GenerateWeeklyReportsOptions,
) *GenerateWeeklyReports {
	reportsStorageDir := ""
	reportsGCSPrefix := "reports"
	var now func() time.Time

	if options != nil {
		reportsStorageDir = options.ReportsStorageDir
		reportsGCSPrefix = strings.Trim(options.ReportsGCSPrefix, "/")
		now = options.Now
	}

	if reportsStorageDir == "" {
		reportsStorageDir = filepath.Join("storage", "reports")
	}

	if reportsGCSPrefix == "" {
		reportsGCSPrefix = "reports"
	}

	if now == nil {
		now = time.Now
	}

	return &GenerateWeeklyReports{
		reportRepo:        deps.ReportRepo,
		workspaceReader:   deps.WorkspaceReader,
		weekReader:        deps.WeekReader,
		assignmentReader:  deps.AssignmentReader,
		taskReader:        deps.TaskReader,
		userReader:        deps.UserReader,
		pdfGenerator:      deps.PDFGenerator,
		aiReportGenerator: deps.AIReportGenerator,
		reportFileStorage: deps.ReportFileStorage,
		reportsStorageDir: reportsStorageDir,
		reportsGCSPrefix:  reportsGCSPrefix,
		now:               now,
	}
}

func (uc *GenerateWeeklyReports) Execute(input GenerateWeeklyReportsInput) (*GenerateWeeklyReportsOutput, error) {
	if err := validateGenerateWeeklyReportsInput(input); err != nil {
		return nil, err
	}

	workspace, week, err := uc.resolveGenerationReferences(input)
	if err != nil {
		return nil, err
	}

	outputs, err := uc.generateReportsForAssignments(input, workspace, week)
	if err != nil {
		return nil, err
	}

	return &GenerateWeeklyReportsOutput{Reports: outputs}, nil
}

func validateGenerateWeeklyReportsInput(input GenerateWeeklyReportsInput) error {
	if input.WorkspaceID == 0 {
		return reportsDomain.ErrReportWorkspaceIDRequired
	}

	if input.WeekID == 0 {
		return reportsDomain.ErrReportWeekIDRequired
	}

	return nil
}

func (uc *GenerateWeeklyReports) resolveGenerationReferences(
	input GenerateWeeklyReportsInput,
) (*workspacesDomain.Workspace, *weeksDomain.Week, error) {
	workspace, err := uc.workspaceReader.FindByID(input.WorkspaceID)
	if err != nil {
		return nil, nil, reportsDomain.ErrReportWorkspaceNotFound
	}

	week, err := uc.weekReader.FindByID(input.WeekID)
	if err != nil {
		return nil, nil, reportsDomain.ErrReportWeekNotFound
	}

	return workspace, week, nil
}

func (uc *GenerateWeeklyReports) generateReportsForAssignments(
	input GenerateWeeklyReportsInput,
	workspace *workspacesDomain.Workspace,
	week *weeksDomain.Week,
) ([]ReportOutput, error) {
	assignments, err := uc.assignmentReader.FindAllByWorkspaceID(input.WorkspaceID)
	if err != nil {
		return nil, err
	}

	reportableAssignments := filterReportableAssignments(assignments)
	if len(reportableAssignments) == 0 {
		return nil, reportsDomain.ErrReportNoAssignmentsFound
	}

	allTasks, err := uc.taskReader.FindAllByWorkspaceAndWeek(
		input.WorkspaceID,
		input.WeekID,
		week.InitialDate,
	)
	if err != nil {
		return nil, err
	}

	outputs := make([]ReportOutput, 0, len(reportableAssignments))

	for _, assignment := range reportableAssignments {
		filteredTasks, totalHours := collectTasksForAssignment(allTasks, assignment.ID)
		if len(filteredTasks) == 0 {
			continue
		}

		output, err := uc.generateReportForAssignment(
			input,
			workspace,
			week,
			assignment,
			filteredTasks,
			totalHours,
		)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, output)
	}

	if len(outputs) == 0 {
		return nil, reportsDomain.ErrReportNoTasksFoundForWeek
	}

	return outputs, nil
}

func filterReportableAssignments(assignments []assignmentsDomain.Assignment) []assignmentsDomain.Assignment {
	reportableAssignments := make([]assignmentsDomain.Assignment, 0, len(assignments))

	for _, assignment := range assignments {
		if isReportableAssignmentRole(assignment.Role) {
			reportableAssignments = append(reportableAssignments, assignment)
		}
	}

	return reportableAssignments
}

func isReportableAssignmentRole(role assignmentsDomain.AssignmentRole) bool {
	return role == assignmentsDomain.RoleMonitor || role == assignmentsDomain.RoleAssistant
}

func (uc *GenerateWeeklyReports) generateReportForAssignment(
	input GenerateWeeklyReportsInput,
	workspace *workspacesDomain.Workspace,
	week *weeksDomain.Week,
	assignment assignmentsDomain.Assignment,
	filteredTasks []tasksDomain.Task,
	totalHours int,
) (ReportOutput, error) {
	user, err := uc.userReader.FindByID(assignment.UserID)
	if err != nil {
		return ReportOutput{}, reportsDomain.ErrReportUserNotFound
	}

	aiReport, err := uc.aiReportGenerator.GenerateWeeklyReport(
		buildAIWeeklyReportInput(input, workspace, week, assignment, user, filteredTasks, totalHours),
	)
	if err != nil || strings.TrimSpace(aiReport) == "" {
		return ReportOutput{}, reportsDomain.ErrReportAIGenerationFailed
	}

	fileName := uc.buildFileName(input.WorkspaceID, input.WeekID, assignment.ID)
	localFilePath := filepath.Join(uc.reportsStorageDir, fileName)
	objectName := uc.buildObjectName(fileName)

	title := "Reporte generado por IA"
	lines := buildPDFLines(workspace, user, assignment, week, filteredTasks, totalHours, aiReport)

	if err := uc.pdfGenerator.Generate(localFilePath, title, lines); err != nil {
		return ReportOutput{}, reportsDomain.ErrReportPDFGenerationFailed
	}

	pdfBytes, err := os.ReadFile(localFilePath)
	if err != nil {
		return ReportOutput{}, reportsDomain.ErrReportPDFGenerationFailed
	}

	if uc.reportFileStorage != nil {
		if err := uc.reportFileStorage.Upload(
			context.Background(),
			objectName,
			bytes.NewReader(pdfBytes),
			"application/pdf",
		); err != nil {
			return ReportOutput{}, reportsDomain.ErrReportPDFGenerationFailed
		}
	}

	report, err := reportsDomain.NewWeeklyReport(
		input.WorkspaceID,
		input.WeekID,
		assignment.ID,
		assignment.UserID,
		objectName,
	)
	if err != nil {
		return ReportOutput{}, reportsDomain.ErrReportPDFGenerationFailed
	}

	if err := uc.reportRepo.Save(report); err != nil {
		return ReportOutput{}, reportsDomain.ErrReportPDFGenerationFailed
	}

	return toReportOutput(*report), nil
}

func buildAIWeeklyReportInput(
	input GenerateWeeklyReportsInput,
	workspace *workspacesDomain.Workspace,
	week *weeksDomain.Week,
	assignment assignmentsDomain.Assignment,
	user *usersDomain.User,
	tasks []tasksDomain.Task,
	totalHours int,
) AIWeeklyReportInput {
	aiTasks := make([]AIWeeklyReportTask, 0, len(tasks))

	for _, task := range tasks {
		aiTasks = append(aiTasks, AIWeeklyReportTask{
			Title:        task.Title,
			Description:  task.Description,
			Status:       string(task.Status),
			SpentHours:   task.SpentHours,
			Observations: task.Observations,
			Late:         task.Late,
		})
	}

	return AIWeeklyReportInput{
		WorkspaceID:   input.WorkspaceID,
		WorkspaceName: workspace.Name,
		WeekID:        input.WeekID,
		WeekNumber:    week.Number,
		InitialDate:   week.InitialDate,
		FinalDate:     week.FinalDate,
		AssignmentID:  assignment.ID,
		UserID:        assignment.UserID,
		UserName:      user.Name,
		Role:          string(assignment.Role),
		TotalHours:    totalHours,
		Tasks:         aiTasks,
	}
}

func (uc *GenerateWeeklyReports) buildFileName(workspaceID, weekID, assignmentID uint) string {
	return fmt.Sprintf(
		"workspace_%d_week_%d_assignment_%d.pdf",
		workspaceID,
		weekID,
		assignmentID,
	)
}

func (uc *GenerateWeeklyReports) buildObjectName(fileName string) string {
	prefix := strings.Trim(uc.reportsGCSPrefix, "/")
	if prefix == "" {
		return fileName
	}

	return prefix + "/" + fileName
}

func collectTasksForAssignment(
	allTasks []tasksDomain.Task,
	assignmentID uint,
) ([]tasksDomain.Task, int) {
	result := make([]tasksDomain.Task, 0)
	totalHours := 0

	for _, task := range allTasks {
		if task.AssignmentID != assignmentID {
			continue
		}

		result = append(result, task)
		totalHours += task.SpentHours
	}

	return result, totalHours
}

func buildPDFLines(
	workspace *workspacesDomain.Workspace,
	user *usersDomain.User,
	assignment assignmentsDomain.Assignment,
	week *weeksDomain.Week,
	tasks []tasksDomain.Task,
	totalHours int,
	aiReport string,
) []string {
	lines := []string{
		fmt.Sprintf("Workspace: %s", workspace.Name),
		fmt.Sprintf("Monitor/Asistente: %s", user.Name),
		fmt.Sprintf("Rol: %s", assignment.Role),
		fmt.Sprintf("Semana: %d (%s a %s)", week.Number, week.InitialDate, week.FinalDate),
		fmt.Sprintf("Horas reportadas: %d", totalHours),
		"",
		"Síntesis generada por IA:",
	}

	lines = append(lines, splitTextIntoPDFLines(aiReport)...)

	lines = append(lines,
		"",
		"Tareas reportadas:",
	)

	for _, task := range tasks {
		lines = append(lines,
			fmt.Sprintf("- %s | estado=%s | horas=%d", task.Title, task.Status, task.SpentHours),
			fmt.Sprintf("descripcion: %s", task.Description),
		)

		if strings.TrimSpace(task.Observations) != "" {
			lines = append(lines, fmt.Sprintf("observaciones: %s", task.Observations))
		}
	}

	return lines
}

func splitTextIntoPDFLines(text string) []string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	paragraphs := strings.Split(normalized, "\n")

	lines := make([]string, 0)

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}

		lines = append(lines, wrapText(paragraph, 95)...)
	}

	return lines
}

func wrapText(text string, maxLen int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0)
	current := words[0]

	for _, word := range words[1:] {
		if len(current)+1+len(word) > maxLen {
			lines = append(lines, current)
			current = word
			continue
		}

		current += " " + word
	}

	lines = append(lines, current)
	return lines
}
