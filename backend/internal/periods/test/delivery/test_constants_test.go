package delivery

const (
	testPeriodsPath             = "/periods"
	testPeriodByIDPath          = "/periods/1"
	testPeriod202610Name        = "2026-10"
	testPeriod202611Name        = "2026-11"
	testPeriodInitialDate1005   = "2026-10-05"
	testPeriodInitialDate1012   = "2026-10-12"
	testJSONContentType         = "application/json"
	testHeaderContentType       = "Content-Type"
	errCreatePeriodStatus       = "CreatePeriod() status = %d, want %d"
	errListPeriodsStatus        = "ListPeriods() status = %d, want %d"
)
