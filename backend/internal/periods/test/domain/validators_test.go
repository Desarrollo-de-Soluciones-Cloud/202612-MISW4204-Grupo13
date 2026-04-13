package domain

import (
	"backend/internal/periods/domain"
	"testing"
)

func TestValidatePeriodName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid name",
			input:   "2026-10",
			wantErr: false,
		},
		{
			name:    "valid name with leading/trailing spaces",
			input:   "  2026-10  ",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errType: domain.ErrPeriodNameWrongFormat,
		},
		{
			name:    "wrong format - missing dash",
			input:   "202610",
			wantErr: true,
			errType: domain.ErrPeriodNameWrongFormat,
		},
		{
			name:    "wrong format - wrong year length",
			input:   "26-10",
			wantErr: true,
			errType: domain.ErrPeriodNameWrongFormat,
		},
		{
			name:    "wrong format - wrong semester length",
			input:   "2026-1",
			wantErr: true,
			errType: domain.ErrPeriodNameWrongFormat,
		},
		{
			name:    "wrong format - non-numeric values",
			input:   "abcd-ef",
			wantErr: true,
			errType: domain.ErrPeriodNameWrongFormat,
		},
		{
			name:    "wrong format - too long",
			input:   "2026-101",
			wantErr: true,
			errType: domain.ErrPeriodNameWrongFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("ValidatePeriodName() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidatePeriodInitialDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid date",
			input:   "2026-10-01",
			wantErr: false,
		},
		{
			name:    "valid date",
			input:   "2026-10-01",
			wantErr: false,
		},
		{
			name:    "empty date",
			input:   "",
			wantErr: true,
			errType: domain.ErrPeriodInitialDateWrongFormat,
		},
		{
			name:    "wrong format - no dashes",
			input:   "20261001",
			wantErr: true,
			errType: domain.ErrPeriodInitialDateWrongFormat,
		},
		{
			name:    "wrong format - incomplete date",
			input:   "2026-10",
			wantErr: true,
			errType: domain.ErrPeriodInitialDateWrongFormat,
		},
		{
			name:    "wrong format - invalid month",
			input:   "2026-13-01",
			wantErr: true,
			errType: domain.ErrPeriodInitialDateWrongFormat,
		},
		{
			name:    "wrong format - invalid day",
			input:   "2026-10-32",
			wantErr: true,
			errType: domain.ErrPeriodInitialDateWrongFormat,
		},
		{
			name:    "wrong format - inverse order",
			input:   "01-10-2026",
			wantErr: true,
			errType: domain.ErrPeriodInitialDateWrongFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodInitialDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodInitialDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("ValidatePeriodInitialDate() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidatePeriodFinalDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid date",
			input:   "2026-12-31",
			wantErr: false,
		},
		{
			name:    "empty date",
			input:   "",
			wantErr: true,
			errType: domain.ErrPeriodFinalDateWrongFormat,
		},
		{
			name:    "wrong format - no dashes",
			input:   "20261231",
			wantErr: true,
			errType: domain.ErrPeriodFinalDateWrongFormat,
		},
		{
			name:    "wrong format - invalid month",
			input:   "2026-13-01",
			wantErr: true,
			errType: domain.ErrPeriodFinalDateWrongFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodFinalDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodFinalDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("ValidatePeriodFinalDate() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidatePeriodInscriptionFinalDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid date",
			input:   "2026-11-15",
			wantErr: false,
		},
		{
			name:    "empty date",
			input:   "",
			wantErr: true,
			errType: domain.ErrPeriodInscriptionFinalDateWrongFormat,
		},
		{
			name:    "wrong format - no dashes",
			input:   "20261115",
			wantErr: true,
			errType: domain.ErrPeriodInscriptionFinalDateWrongFormat,
		},
		{
			name:    "wrong format - invalid day",
			input:   "2026-11-32",
			wantErr: true,
			errType: domain.ErrPeriodInscriptionFinalDateWrongFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodInscriptionFinalDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodInscriptionFinalDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("ValidatePeriodInscriptionFinalDate() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidatorsValidatePeriodDateSequence(t *testing.T) {
	tests := []struct {
		name                  string
		initialDate           string
		finalDate             string
		inscriptionFinalDate  string
		wantErr               bool
	}{
		{
			name:                 "valid sequence",
			initialDate:           "2026-10-01",
			finalDate:             "2026-12-31",
			inscriptionFinalDate:  "2026-11-15",
			wantErr:               false,
		},
		{
			name:                 "valid sequence - same dates",
			initialDate:           "2026-10-01",
			finalDate:             "2026-10-01",
			inscriptionFinalDate:  "2026-10-01",
			wantErr:               false,
		},
		{
			name:                 "invalid - initialDate after finalDate",
			initialDate:           "2026-12-31",
			finalDate:             "2026-10-01",
			inscriptionFinalDate:  "2026-11-15",
			wantErr:               true,
		},
		{
			name:                 "invalid - initialDate after inscriptionFinalDate",
			initialDate:           "2026-12-01",
			finalDate:             "2026-12-31",
			inscriptionFinalDate:  "2026-10-01",
			wantErr:               true,
		},
		{
			name:                 "invalid - inscriptionFinalDate after finalDate",
			initialDate:           "2026-10-01",
			finalDate:             "2026-11-30",
			inscriptionFinalDate:  "2026-12-31",
			wantErr:               true,
		},
		{
			name:                 "invalid - all out of sequence",
			initialDate:           "2026-12-31",
			finalDate:             "2026-10-01",
			inscriptionFinalDate:  "2026-11-01",
			wantErr:               true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodDateSequence(tt.initialDate, tt.finalDate, tt.inscriptionFinalDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodDateSequence() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != domain.ErrPeriodDateSequenceInvalid {
				t.Errorf("ValidatePeriodDateSequence() error = %v, want %v", err, domain.ErrPeriodDateSequenceInvalid)
			}
		})
	}
}

func TestValidatePeriodState(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.PeriodState
		wantErr bool
		errType error
	}{
		{
			name:    "valid state - active",
			input:   domain.ActivePeriod,
			wantErr: false,
		},
		{
			name:    "valid state - closed",
			input:   domain.ClosedPeriod,
			wantErr: false,
		},
		{
			name:    "invalid state - empty",
			input:   domain.PeriodState(""),
			wantErr: true,
			errType: domain.ErrPeriodStateRequired,
		},
		{
			name:    "invalid state - unknown",
			input:   domain.PeriodState("invalid"),
			wantErr: true,
			errType: domain.ErrPeriodStateInvalid,
		},
		{
			name:    "invalid state - misspelled",
			input:   domain.PeriodState("Active"),
			wantErr: true,
			errType: domain.ErrPeriodStateInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodState(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodState() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("ValidatePeriodState() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestValidatorsNormalizePeriodName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no spaces",
			input:    "2026-10",
			expected: "2026-10",
		},
		{
			name:     "leading spaces",
			input:    "  2026-10",
			expected: "2026-10",
		},
		{
			name:     "trailing spaces",
			input:    "2026-10  ",
			expected: "2026-10",
		},
		{
			name:     "leading and trailing spaces",
			input:    "  2026-10  ",
			expected: "2026-10",
		},
		{
			name:     "multiple spaces",
			input:    "   2026-10   ",
			expected: "2026-10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domain.NormalizePeriodName(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePeriodName() = %v, want %v", result, tt.expected)
			}
		})
	}
}
