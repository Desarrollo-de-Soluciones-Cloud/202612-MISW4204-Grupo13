package domain

import "errors"

var (
	ErrPeriodNotFound                     = errors.New("period not found")
	ErrInvalidInput                       = errors.New("invalid input")
	ErrPeriodNameRequired                 = errors.New("period name is required")
	ErrPeriodNameWrongFormat              = errors.New("period name must have format YYYY-SS")
	ErrPeriodInitialDateRequired          = errors.New("period initial date is required")
	ErrPeriodInitialDateInvalid           = errors.New("period initial date is invalid")
	ErrPeriodFinalDateRequired            = errors.New("period final date is required")
	ErrPeriodFinalDateInvalid             = errors.New("period final date is invalid")
	ErrPeriodInscriptionFinalDateRequired = errors.New("period inscription final date is required")
	ErrPeriodInscriptionFinalDateInvalid  = errors.New("period inscription final date is invalid")
	ErrPeriodDateSequenceInvalid          = errors.New("period date sequence is invalid: initial_date must be before final_date and inscription_final_date")
	ErrPeriodStateRequired                = errors.New("period state is required")
	ErrPeriodStateInvalid                 = errors.New("period state is invalid")
	ErrPeriodCannotBeDeleted              = errors.New("period cannot be deleted while it is active")
)
