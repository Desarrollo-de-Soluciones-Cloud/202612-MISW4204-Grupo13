package delivery

import "time"

type CreatePeriodRequest struct {
	Name                 string    `json:"name" binding:"required"`
	InitialDate          time.Time `json:"initial_date" binding:"required"`
	FinalDate            time.Time `json:"final_date" binding:"required"`
	InscriptionFinalDate time.Time `json:"inscription_final_date" binding:"required"`
	PeriodState          string    `json:"period_state" binding:"required"`
}

type CreatePeriodResponse struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	InitialDate          time.Time `json:"initial_date"`
	FinalDate            time.Time `json:"final_date"`
	InscriptionFinalDate time.Time `json:"inscription_final_date"`
	PeriodState          string    `json:"period_state"`
}

type UpdatePeriodRequest struct {
	Name                 string    `json:"name" binding:"required"`
	InitialDate          time.Time `json:"initial_date" binding:"required"`
	FinalDate            time.Time `json:"final_date" binding:"required"`
	InscriptionFinalDate time.Time `json:"inscription_final_date" binding:"required"`
	PeriodState          string    `json:"period_state" binding:"required"`
}

type ListPeriodsResponse struct {
	Periods []PeriodResponse `json:"periods"`
}

type PeriodResponse struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	InitialDate          time.Time `json:"initial_date"`
	FinalDate            time.Time `json:"final_date"`
	InscriptionFinalDate time.Time `json:"inscription_final_date"`
	PeriodState          string    `json:"period_state"`
}
