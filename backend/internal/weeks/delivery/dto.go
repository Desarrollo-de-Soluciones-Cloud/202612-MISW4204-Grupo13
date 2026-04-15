package delivery

type ListWeeksResponse struct {
	Weeks []WeekResponse `json:"weeks"`
}

type WeekResponse struct {
	ID          uint   `json:"id"`
	PeriodID    uint   `json:"period_id"`
	Number      int    `json:"number"`
	InitialDate string `json:"initial_date"`
	FinalDate   string `json:"final_date"`
}
