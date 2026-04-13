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

func TestValidatePeriodInitialDateIsMonday(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid Monday date",
			input:   "2026-10-05",
			wantErr: false,
		},
		{
			name:    "invalid Tuesday date",
			input:   "2026-10-06",
			wantErr: true,
		},
		{
			name:    "invalid Wednesday date",
			input:   "2026-10-07",
			wantErr: true,
		},
		{
			name:    "invalid Thursday date",
			input:   "2026-10-08",
			wantErr: true,
		},
		{
			name:    "invalid Friday date",
			input:   "2026-10-09",
			wantErr: true,
		},
		{
			name:    "invalid Saturday date",
			input:   "2026-10-10",
			wantErr: true,
		},
		{
			name:    "invalid Sunday date",
			input:   "2026-10-11",
			wantErr: true,
		},
		{
			name:    "another valid Monday",
			input:   "2026-10-12",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodInitialDateIsMonday(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodInitialDateIsMonday() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != domain.ErrPeriodInitialDateMustBeMonday {
				t.Errorf("ValidatePeriodInitialDateIsMonday() error = %v, want %v", err, domain.ErrPeriodInitialDateMustBeMonday)
			}
		})
	}
}

func TestValidatePeriodInitialDateIsFuture(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid future date",
			input:   "2026-10-05",
			wantErr: false,
		},
		{
			name:    "invalid past date",
			input:   "2020-10-05",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodInitialDateIsFuture(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodInitialDateIsFuture() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != domain.ErrPeriodInitialDateMustBeFuture {
				t.Errorf("ValidatePeriodInitialDateIsFuture() error = %v, want %v", err, domain.ErrPeriodInitialDateMustBeFuture)
			}
		})
	}
}

func TestValidatePeriodWeeksCount(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{
			name:    "valid 8 weeks",
			input:   8,
			wantErr: false,
		},
		{
			name:    "valid 16 weeks",
			input:   16,
			wantErr: false,
		},
		{
			name:    "invalid 4 weeks",
			input:   4,
			wantErr: true,
		},
		{
			name:    "invalid 10 weeks",
			input:   10,
			wantErr: true,
		},
		{
			name:    "invalid 32 weeks",
			input:   32,
			wantErr: true,
		},
		{
			name:    "invalid 0 weeks",
			input:   0,
			wantErr: true,
		},
		{
			name:    "invalid negative weeks",
			input:   -8,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePeriodWeeksCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePeriodWeeksCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != domain.ErrPeriodWeeksCountInvalid {
				t.Errorf("ValidatePeriodWeeksCount() error = %v, want %v", err, domain.ErrPeriodWeeksCountInvalid)
			}
		})
	}
}

func TestCalculatePeriodFinalDate(t *testing.T) {
	tests := []struct {
		name           string
		initialDate    string
		weeksCount     int
		expectedResult string
		wantErr        bool
	}{
		{
			name:           "calculate final date for 8 weeks",
			initialDate:    "2026-10-05",
			weeksCount:     8,
			expectedResult: "2026-11-29",
			wantErr:        false,
		},
		{
			name:           "calculate final date for 16 weeks",
			initialDate:    "2026-10-05",
			weeksCount:     16,
			expectedResult: "2027-01-24",
			wantErr:        false,
		},
		{
			name:           "calculate final date starting from different date",
			initialDate:    "2026-10-12",
			weeksCount:     8,
			expectedResult: "2026-12-06",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := domain.CalculatePeriodFinalDate(tt.initialDate, tt.weeksCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePeriodFinalDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result != tt.expectedResult {
				t.Errorf("CalculatePeriodFinalDate() result = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestCalculatePeriodInscriptionFinalDate(t *testing.T) {
	tests := []struct {
		name           string
		initialDate    string
		expectedResult string
		wantErr        bool
	}{
		{
			name:           "calculate inscription date one day before",
			initialDate:    "2026-10-05",
			expectedResult: "2026-10-04",
			wantErr:        false,
		},
		{
			name:           "calculate inscription date for another date",
			initialDate:    "2026-10-12",
			expectedResult: "2026-10-11",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := domain.CalculatePeriodInscriptionFinalDate(tt.initialDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePeriodInscriptionFinalDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result != tt.expectedResult {
				t.Errorf("CalculatePeriodInscriptionFinalDate() result = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
