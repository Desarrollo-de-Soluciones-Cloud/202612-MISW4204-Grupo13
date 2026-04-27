package application

import weeksDomain "backend/internal/weeks/domain"

type TaskWeekRepository interface {
	FindByPeriodIDAndStartDate(periodID uint, startDate string) (*weeksDomain.Week, error)
}