package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type ProcessWeeklyReportJob struct {
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
}

func NewProcessWeeklyReportJob(
	reportRepo reportsDomain.ReportRepository,
	workspaceReader WorkspaceReader,
	weekReader WeekReader,
	assignmentReader AssignmentReader,
	taskReader TaskReader,
	userReader UserReader,
	pdfGenerator PDFGenerator,
	aiReportGenerator AIReportGenerator,
	reportFileStorage ReportFileStorage,
	options *GenerateWeeklyReportsOptions,
) *ProcessWeeklyReportJob {
	generator := NewGenerateWeeklyReports(
		reportRepo,
		workspaceReader,
		weekReader,
		assignmentReader,
		taskReader,
		userReader,
		pdfGenerator,
		aiReportGenerator,
		reportFileStorage,
		options,
	)

	return &ProcessWeeklyReportJob{
		reportRepo:        reportRepo,
		workspaceReader:   workspaceReader,
		weekReader:        weekReader,
		assignmentReader:  assignmentReader,
		taskReader:        taskReader,
		userReader:        userReader,
		pdfGenerator:      pdfGenerator,
		aiReportGenerator: aiReportGenerator,
		reportFileStorage: reportFileStorage,
		reportsStorageDir: generator.reportsStorageDir,
		reportsGCSPrefix:  generator.reportsGCSPrefix,
	}
}

func (uc *ProcessWeeklyReportJob) Execute(input WeeklyReportJobMessage) (*ReportOutput, error) {
	if err := validateWeeklyReportJobMessage(input); err != nil {
		return nil, err
	}

	workspace, week, err := uc.resolveReferences(input)
	if err != nil {
		return nil, err
	}

	assignment, err := uc.resolveAssignment(input)
	if err != nil {
		return nil, err
	}

	allTasks, err := uc.taskReader.FindAllByWorkspaceAndWeek(input.WorkspaceID, input.WeekID, week.InitialDate)
	if err != nil {
		return nil, err
	}

	filteredTasks, totalHours := collectTasksForAssignment(allTasks, assignment.ID)
	if len(filteredTasks) == 0 {
		return nil, reportsDomain.ErrReportNoTasksFoundForWeek
	}

	output, err := uc.generateReport(input, workspace, week, *assignment, filteredTasks, totalHours)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func validateWeeklyReportJobMessage(input WeeklyReportJobMessage) error {
	if input.WorkspaceID == 0 {
		return reportsDomain.ErrReportWorkspaceIDRequired
	}

	if input.WeekID == 0 {
		return reportsDomain.ErrReportWeekIDRequired
	}

	if input.AssignmentID == 0 {
		return reportsDomain.ErrReportAssignmentIDRequired
	}

	if input.UserID == 0 {
		return reportsDomain.ErrReportUserIDRequired
	}

	return nil
}

func (uc *ProcessWeeklyReportJob) resolveReferences(
	input WeeklyReportJobMessage,
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

func (uc *ProcessWeeklyReportJob) resolveAssignment(
	input WeeklyReportJobMessage,
) (*assignmentsDomain.Assignment, error) {
	assignments, err := uc.assignmentReader.FindAllByWorkspaceID(input.WorkspaceID)
	if err != nil {
		return nil, err
	}

	for _, assignment := range assignments {
		if assignment.ID != input.AssignmentID || assignment.UserID != input.UserID {
			continue
		}

		if !isReportableAssignmentRole(assignment.Role) {
			return nil, reportsDomain.ErrReportNoAssignmentsFound
		}

		copyAssignment := assignment
		return &copyAssignment, nil
	}

	return nil, reportsDomain.ErrReportNoAssignmentsFound
}

func (uc *ProcessWeeklyReportJob) generateReport(
	input WeeklyReportJobMessage,
	workspace *workspacesDomain.Workspace,
	week *weeksDomain.Week,
	assignment assignmentsDomain.Assignment,
	filteredTasks []tasksDomain.Task,
	totalHours int,
) (ReportOutput, error) {
	generator := &GenerateWeeklyReports{
		reportRepo:        uc.reportRepo,
		workspaceReader:   uc.workspaceReader,
		weekReader:        uc.weekReader,
		assignmentReader:  uc.assignmentReader,
		taskReader:        uc.taskReader,
		userReader:        uc.userReader,
		pdfGenerator:      uc.pdfGenerator,
		aiReportGenerator: uc.aiReportGenerator,
		reportFileStorage: uc.reportFileStorage,
		reportsStorageDir: uc.reportsStorageDir,
		reportsGCSPrefix:  uc.reportsGCSPrefix,
	}

	return generator.generateReportForAssignment(
		GenerateWeeklyReportsInput{
			WorkspaceID: input.WorkspaceID,
			WeekID:      input.WeekID,
		},
		workspace,
		week,
		assignment,
		filteredTasks,
		totalHours,
	)
}
