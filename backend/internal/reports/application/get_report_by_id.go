package application

import reportsDomain "backend/internal/reports/domain"

type GetReportByID struct {
	reportRepo reportsDomain.ReportRepository
}

func NewGetReportByID(reportRepo reportsDomain.ReportRepository) *GetReportByID {
	return &GetReportByID{reportRepo: reportRepo}
}

func (uc *GetReportByID) Execute(id uint) (*ReportOutput, error) {
	report, err := uc.reportRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	output := toReportOutput(*report)
	return &output, nil
}
