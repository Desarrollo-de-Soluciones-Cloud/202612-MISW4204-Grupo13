package delivery

import "time"

type CreateWeekRequest struct {
	Number      int       `json:"number" binding:"required"`
	InitialDate time.Time `json:"initial_date" binding:"required"`
	FinalDate   time.Time `json:"final_date" binding:"required"`
	PeriodID    uint      `json:"period_id" binding:"required"`
}

type CreateWeekResponse struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type UpdateWeekRequest struct {
	Number      int       `json:"number" binding:"required"`
	InitialDate time.Time `json:"initial_date" binding:"required"`
	FinalDate   time.Time `json:"final_date" binding:"required"`
	PeriodID    uint      `json:"period_id" binding:"required"`
}

type ListWeeksResponse struct {
	Weeks []WeekResponse `json:"weeks"`
}

type WeekResponse struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}
