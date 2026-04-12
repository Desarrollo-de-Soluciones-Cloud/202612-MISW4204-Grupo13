package domain

type PeriodState string

const (
	ActivePeriod PeriodState = "active"
	ClosedPeriod PeriodState = "closed"
)

func IsValidPeriodState(state PeriodState) bool {
	switch state {
	case ActivePeriod, ClosedPeriod:
		return true
	default:
		return false
	}
}