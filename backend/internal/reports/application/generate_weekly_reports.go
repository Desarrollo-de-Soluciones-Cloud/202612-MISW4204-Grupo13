package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
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
	pdfGenerator      PDFGenerator
	reportsStorageDir string
	now               func() time.Time
}

type GenerateWeeklyReportsOptions struct {
	ReportsStorageDir string
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
	reportRepo reportsDomain.ReportRepository,
	workspaceReader WorkspaceReader,
	weekReader WeekReader,
	assignmentReader AssignmentReader,
	taskReader TaskReader,
	pdfGenerator PDFGenerator,
	options *GenerateWeeklyReportsOptions,
) *GenerateWeeklyReports {
	reportsStorageDir := ""
	var now func() time.Time
	if options != nil {
		reportsStorageDir = options.ReportsStorageDir
		now = options.Now
	}

	if reportsStorageDir == "" {
		reportsStorageDir = filepath.Join("storage", "reports")
	}
	if now == nil {
		now = time.Now
	}

	return &GenerateWeeklyReports{
		reportRepo:        reportRepo,
		workspaceReader:   workspaceReader,
		weekReader:        weekReader,
		assignmentReader:  assignmentReader,
		taskReader:        taskReader,
		pdfGenerator:      pdfGenerator,
		reportsStorageDir: reportsStorageDir,
		now:               now,
	}
}

func (uc *GenerateWeeklyReports) Execute(input GenerateWeeklyReportsInput) (*GenerateWeeklyReportsOutput, error) {
	if err := validateGenerateWeeklyReportsInput(input); err != nil {
		return nil, err
	}

	week, err := uc.resolveGenerationReferences(input)
	if err != nil {
		return nil, err
	}

	outputs, err := uc.generateReportsForAssignments(input, week)
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

func (uc *GenerateWeeklyReports) resolveGenerationReferences(input GenerateWeeklyReportsInput) (*weeksDomain.Week, error) {
	if _, err := uc.workspaceReader.FindByID(input.WorkspaceID); err != nil {
		return nil, reportsDomain.ErrReportWorkspaceNotFound
	}

	week, err := uc.weekReader.FindByID(input.WeekID)
	if err != nil {
		return nil, reportsDomain.ErrReportWeekNotFound
	}

	return week, nil
}

func (uc *GenerateWeeklyReports) generateReportsForAssignments(input GenerateWeeklyReportsInput, week *weeksDomain.Week) ([]ReportOutput, error) {
	assignments, err := uc.assignmentReader.FindAllByWorkspaceID(input.WorkspaceID)
	if err != nil {
		return nil, err
	}

	allTasks, err := uc.taskReader.FindAll()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(uc.reportsStorageDir, 0o755); err != nil {
		return nil, err
	}

	outputs := make([]ReportOutput, 0, len(assignments))
	for _, assignment := range assignments {
		if !isReportableAssignmentRole(assignment.Role) {
			continue
		}

		output, err := uc.generateReportForAssignment(input, week, assignment, allTasks)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

func isReportableAssignmentRole(role assignmentsDomain.AssignmentRole) bool {
	return role == assignmentsDomain.RoleMonitor || role == assignmentsDomain.RoleAssistant
}

func (uc *GenerateWeeklyReports) generateReportForAssignment(
	input GenerateWeeklyReportsInput,
	week *weeksDomain.Week,
	assignment assignmentsDomain.Assignment,
	allTasks []tasksDomain.Task,
) (ReportOutput, error) {
	filteredTasks, totalHours := collectTasksForAssignmentWeek(allTasks, assignment.ID, input.WeekID, week)
	summary := buildSummary(assignment.UserID, week, totalHours, filteredTasks)
	filePath := uc.buildFilePath(input.WorkspaceID, input.WeekID, assignment.ID, assignment.UserID)

	title := fmt.Sprintf("Reporte semanal - Workspace %d - Semana %d", input.WorkspaceID, week.Number)
	lines := buildPDFLines(assignment, week, filteredTasks, totalHours)
	if err := uc.pdfGenerator.Generate(filePath, title, lines); err != nil {
		return ReportOutput{}, err
	}

	report, err := reportsDomain.NewWeeklyReport(input.WorkspaceID, input.WeekID, assignment.ID, assignment.UserID, summary, filePath)
	if err != nil {
		return ReportOutput{}, err
	}

	if err := uc.reportRepo.Create(report); err != nil {
		return ReportOutput{}, err
	}

	return toReportOutput(*report), nil
}

func (uc *GenerateWeeklyReports) buildFilePath(workspaceID, weekID, assignmentID, userID uint) string {
	fileName := fmt.Sprintf(
		"workspace_%d_week_%d_assignment_%d_user_%d_%d.pdf",
		workspaceID,
		weekID,
		assignmentID,
		userID,
		uc.now().UnixNano(),
	)
	return filepath.Join(uc.reportsStorageDir, fileName)
}

func collectTasksForAssignmentWeek(
	allTasks []tasksDomain.Task,
	assignmentID uint,
	weekID uint,
	week *weeksDomain.Week,
) ([]tasksDomain.Task, int) {
	result := make([]tasksDomain.Task, 0)
	totalHours := 0

	for _, task := range allTasks {
		if task.AssignmentID != assignmentID {
			continue
		}
		if !taskBelongsToWeek(task, weekID, week) {
			continue
		}

		result = append(result, task)
		totalHours += task.SpentHours
	}

	return result, totalHours
}

func taskBelongsToWeek(task tasksDomain.Task, weekID uint, week *weeksDomain.Week) bool {
	if task.WeekID != nil && *task.WeekID == weekID {
		return true
	}

	return task.WeekStartDate.Format("2006-01-02") == week.InitialDate
}

func buildSummary(userID uint, week *weeksDomain.Week, totalHours int, tasks []tasksDomain.Task) string {
	if len(tasks) == 0 {
		return fmt.Sprintf("Usuario %d, semana %d: sin tareas reportadas.", userID, week.Number)
	}

	titles := make([]string, 0, len(tasks))
	for _, task := range tasks {
		titles = append(titles, task.Title)
	}

	return fmt.Sprintf(
		"Usuario %d, semana %d, horas totales %d. Tareas: %s",
		userID,
		week.Number,
		totalHours,
		strings.Join(titles, "; "),
	)
}

func buildPDFLines(
	assignment assignmentsDomain.Assignment,
	week *weeksDomain.Week,
	tasks []tasksDomain.Task,
	totalHours int,
) []string {
	lines := []string{
		fmt.Sprintf("Usuario asignado: %d", assignment.UserID),
		fmt.Sprintf("Rol: %s", assignment.Role),
		fmt.Sprintf("Semana: %d (%s a %s)", week.Number, week.InitialDate, week.FinalDate),
		fmt.Sprintf("Horas reportadas: %d", totalHours),
		"",
		"Detalle de tareas:",
	}

	if len(tasks) == 0 {
		return append(lines, "- Sin tareas reportadas")
	}

	for _, task := range tasks {
		lines = append(lines,
			fmt.Sprintf("- %s | estado=%s | horas=%d", task.Title, task.Status, task.SpentHours),
			fmt.Sprintf("  descripcion: %s", task.Description),
		)
		if strings.TrimSpace(task.Observations) != "" {
			lines = append(lines, fmt.Sprintf("  observaciones: %s", task.Observations))
		}
	}

	return lines
}
