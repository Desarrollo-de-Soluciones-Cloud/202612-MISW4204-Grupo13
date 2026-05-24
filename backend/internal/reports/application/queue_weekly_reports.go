package application

import (
	"context"

	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type QueueWeeklyReports struct {
	workspaceReader  WorkspaceReader
	weekReader       WeekReader
	assignmentReader AssignmentReader
	taskReader       TaskReader
	jobPublisher     ReportJobPublisher
}

type QueueWeeklyReportsInput struct {
	WorkspaceID uint
	WeekID      uint
}

type QueueWeeklyReportsOutput struct {
	QueuedCount int `json:"queued_count"`
}

func NewQueueWeeklyReports(
	workspaceReader WorkspaceReader,
	weekReader WeekReader,
	assignmentReader AssignmentReader,
	taskReader TaskReader,
	jobPublisher ReportJobPublisher,
) *QueueWeeklyReports {
	return &QueueWeeklyReports{
		workspaceReader:  workspaceReader,
		weekReader:       weekReader,
		assignmentReader: assignmentReader,
		taskReader:       taskReader,
		jobPublisher:     jobPublisher,
	}
}

func (uc *QueueWeeklyReports) Execute(input QueueWeeklyReportsInput) (*QueueWeeklyReportsOutput, error) {
	if err := validateGenerateWeeklyReportsInput(GenerateWeeklyReportsInput{
		WorkspaceID: input.WorkspaceID,
		WeekID:      input.WeekID,
	}); err != nil {
		return nil, err
	}

	_, week, err := uc.resolveReferences(input)
	if err != nil {
		return nil, err
	}

	assignments, err := uc.assignmentReader.FindAllByWorkspaceID(input.WorkspaceID)
	if err != nil {
		return nil, err
	}

	reportableAssignments := filterReportableAssignments(assignments)
	if len(reportableAssignments) == 0 {
		return nil, reportsDomain.ErrReportNoAssignmentsFound
	}

	assignmentsWithTasks, err := uc.filterAssignmentsWithTasks(input, week.InitialDate, reportableAssignments)
	if err != nil {
		return nil, err
	}
	if len(assignmentsWithTasks) == 0 {
		return nil, reportsDomain.ErrReportNoTasksFoundForWeek
	}

	queuedCount, err := uc.publishJobs(input, assignmentsWithTasks)
	if err != nil {
		return nil, err
	}

	return &QueueWeeklyReportsOutput{QueuedCount: queuedCount}, nil
}

func (uc *QueueWeeklyReports) resolveReferences(
	input QueueWeeklyReportsInput,
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

func (uc *QueueWeeklyReports) filterAssignmentsWithTasks(
	input QueueWeeklyReportsInput,
	weekInitialDate string,
	assignments []assignmentsDomain.Assignment,
) ([]assignmentsDomain.Assignment, error) {
	allTasks, err := uc.taskReader.FindAllByWorkspaceAndWeek(
		input.WorkspaceID,
		input.WeekID,
		weekInitialDate,
	)
	if err != nil {
		return nil, err
	}

	assignmentsWithTasks := make([]assignmentsDomain.Assignment, 0, len(assignments))
	for _, assignment := range assignments {
		if hasTasksForAssignment(allTasks, assignment.ID) {
			assignmentsWithTasks = append(assignmentsWithTasks, assignment)
		}
	}

	return assignmentsWithTasks, nil
}

func hasTasksForAssignment(allTasks []tasksDomain.Task, assignmentID uint) bool {
	for _, task := range allTasks {
		if task.AssignmentID == assignmentID {
			return true
		}
	}

	return false
}

func (uc *QueueWeeklyReports) publishJobs(
	input QueueWeeklyReportsInput,
	assignments []assignmentsDomain.Assignment,
) (int, error) {
	queuedCount := 0

	for _, assignment := range assignments {
		err := uc.jobPublisher.PublishWeeklyReportJob(context.Background(), WeeklyReportJobMessage{
			WorkspaceID:  input.WorkspaceID,
			WeekID:       input.WeekID,
			AssignmentID: assignment.ID,
			UserID:       assignment.UserID,
		})
		if err != nil {
			return 0, err
		}

		queuedCount++
	}

	return queuedCount, nil
}
