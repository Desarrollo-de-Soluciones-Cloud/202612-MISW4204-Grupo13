package application

import reportsDomain "backend/internal/reports/domain"

type ListReports struct {
	reportRepo       reportsDomain.ReportRepository
	workspaceReader WorkspaceReader
	weekReader      WeekReader
	userReader      UserReader
}

type ListReportsInput struct {
	WorkspaceID uint
	WeekID      *uint
	UserID      *uint
}

type ListReportsOutput struct {
	Reports []ReportOutput `json:"reports"`
}

func NewListReports(
	reportRepo reportsDomain.ReportRepository,
	workspaceReader WorkspaceReader,
	weekReader WeekReader,
	userReader UserReader,
) *ListReports {
	return &ListReports{
		reportRepo:       reportRepo,
		workspaceReader: workspaceReader,
		weekReader:      weekReader,
		userReader:      userReader,
	}
}

func (uc *ListReports) Execute(input ListReportsInput) (*ListReportsOutput, error) {
	if input.WorkspaceID == 0 {
		return nil, reportsDomain.ErrReportWorkspaceFilterRequired
	}

	if _, err := uc.workspaceReader.FindByID(input.WorkspaceID); err != nil {
		return nil, reportsDomain.ErrReportWorkspaceNotFound
	}

	if input.WeekID != nil {
		if _, err := uc.weekReader.FindByID(*input.WeekID); err != nil {
			return nil, reportsDomain.ErrReportWeekNotFound
		}
	}

	if input.UserID != nil {
		if _, err := uc.userReader.FindByID(*input.UserID); err != nil {
			return nil, reportsDomain.ErrReportUserNotFound
		}
	}

	reports, err := uc.reportRepo.FindAll(input.WorkspaceID, input.WeekID, input.UserID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ReportOutput, 0, len(reports))
	for _, report := range reports {
		outputs = append(outputs, toReportOutput(report))
	}

	return &ListReportsOutput{Reports: outputs}, nil
}