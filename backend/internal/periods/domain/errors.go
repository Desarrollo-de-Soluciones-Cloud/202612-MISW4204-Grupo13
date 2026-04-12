package domain

import "errors"

var (
	ErrPeriodNotFound                        = errors.New("period not found")
	ErrInvalidInput                          = errors.New("invalid input")
	ErrPeriodNameAlreadyExists               = errors.New("a period with the same name already exists")
	ErrPeriodNameRequired                    = errors.New("period name is required")
	ErrPeriodNameWrongFormat                 = errors.New("period name must have format YYYY-SS. For instance, 2026-10")
	ErrPeriodInitialDateRequired             = errors.New("period initial date is required")
	ErrPeriodInitialDateWrongFormat          = errors.New("period initial date must have format YYYY-MM-DD. For instance, 2026-10-01")
	ErrPeriodFinalDateRequired               = errors.New("period final date is required")
	ErrPeriodFinalDateWrongFormat            = errors.New("period final date must have format YYYY-MM-DD. For instance, 2026-10-01")
	ErrPeriodInscriptionFinalDateRequired    = errors.New("period inscription final date is required")
	ErrPeriodInscriptionFinalDateWrongFormat = errors.New("period inscription final date must have format YYYY-MM-DD. For instance, 2026-10-01")
	ErrPeriodDateSequenceInvalid             = errors.New("period date sequence is invalid: initial_date must be before final_date and inscription_final_date. Also, inscription_final_date must be before final_date.")
	ErrPeriodStateRequired                   = errors.New("period state is required")
	ErrPeriodStateInvalid                    = errors.New("period state is invalid. The following states are allowed: " + ValidPeriodStatesString())
)
