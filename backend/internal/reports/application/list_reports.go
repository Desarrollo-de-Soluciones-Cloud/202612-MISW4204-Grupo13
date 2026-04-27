package application

import reportsDomain "backend/internal/reports/domain"

type ListReports struct {
	reportRepo reportsDomain.ReportRepository
}

type ListReportsInput struct {
	WorkspaceID *uint
	WeekID      *uint
}

type ListReportsOutput struct {
	Reports []ReportOutput `json:"reports"`
}

func NewListReports(reportRepo reportsDomain.ReportRepository) *ListReports {
	return &ListReports{reportRepo: reportRepo}
}

func (uc *ListReports) Execute(input ListReportsInput) (*ListReportsOutput, error) {
	reports, err := uc.reportRepo.FindAll(input.WorkspaceID, input.WeekID)
	if err != nil {
		return nil, err
	}

	outputs := make([]ReportOutput, 0, len(reports))
	for _, report := range reports {
		outputs = append(outputs, toReportOutput(report))
	}

	return &ListReportsOutput{Reports: outputs}, nil
}
