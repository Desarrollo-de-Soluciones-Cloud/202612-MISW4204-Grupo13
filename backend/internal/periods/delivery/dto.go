package delivery

type CreatePeriodRequest struct {
	Name                 string `json:"name" binding:"required"`
	InitialDate          string `json:"initial_date" binding:"required"`
	FinalDate            string `json:"final_date" binding:"required"`
	InscriptionFinalDate string `json:"inscription_final_date" binding:"required"`
	PeriodState          string `json:"period_state" binding:"required"`
}

type CreatePeriodResponse struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	InitialDate          string `json:"initial_date"`
	FinalDate            string `json:"final_date"`
	InscriptionFinalDate string `json:"inscription_final_date"`
	PeriodState          string `json:"period_state"`
}

type UpdatePeriodRequest struct {
	Name                 string `json:"name" binding:"required"`
	InitialDate          string `json:"initial_date" binding:"required"`
	FinalDate            string `json:"final_date" binding:"required"`
	InscriptionFinalDate string `json:"inscription_final_date" binding:"required"`
	PeriodState          string `json:"period_state" binding:"required"`
}

type ListPeriodsResponse struct {
	Periods []PeriodResponse `json:"periods"`
}

type PeriodResponse struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	InitialDate          string `json:"initial_date"`
	FinalDate            string `json:"final_date"`
	InscriptionFinalDate string `json:"inscription_final_date"`
	PeriodState          string `json:"period_state"`
}
