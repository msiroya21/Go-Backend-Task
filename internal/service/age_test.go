package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {

	tests := []struct {
		birthDateStr string
		refDateStr   string
		expectedAge  int
	}{
		{"2000-01-01", "2025-10-27", 25},
		{"2000-12-31", "2025-10-27", 24},
		{"1989-02-28", "2024-02-29", 35},
		{"1990-05-10", "2025-05-10", 35},
	}

	for _, tc := range tests {
		birthDate, _ := time.Parse("2006-01-02", tc.birthDateStr)
		refDate, _ := time.Parse("2006-01-02", tc.refDateStr)
		age := CalculateAge(birthDate, refDate)
		if age != tc.expectedAge {
			t.Errorf("CalculateAge(%s, %s) = %d; want %d", tc.birthDateStr, tc.refDateStr, age, tc.expectedAge)
		}
	}
}
