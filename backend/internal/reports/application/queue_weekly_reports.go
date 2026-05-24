package application

import (
	"context"

	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
)

type QueueWeeklyReports struct {
	workspaceReader  WorkspaceReader
	weekReader       WeekReader
	assignmentReader AssignmentReader
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
	jobPublisher ReportJobPublisher,
) *QueueWeeklyReports {
	return &QueueWeeklyReports{
		workspaceReader:  workspaceReader,
		weekReader:       weekReader,
		assignmentReader: assignmentReader,
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

	if err := uc.validateReferences(input); err != nil {
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

	queuedCount, err := uc.publishJobs(input, reportableAssignments)
	if err != nil {
		return nil, err
	}

	return &QueueWeeklyReportsOutput{QueuedCount: queuedCount}, nil
}

func (uc *QueueWeeklyReports) validateReferences(input QueueWeeklyReportsInput) error {
	_, err := uc.workspaceReader.FindByID(input.WorkspaceID)
	if err != nil {
		return reportsDomain.ErrReportWorkspaceNotFound
	}

	_, err = uc.weekReader.FindByID(input.WeekID)
	if err != nil {
		return reportsDomain.ErrReportWeekNotFound
	}

	return nil
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
