package utils

import (
	"testing"
	"time"
)

func TestIsTomorrow(t *testing.T) {
	tests := []struct {
		name       string
		dateString string
		expected   bool
	}{
		{
			name:       "Tomorrow",
			dateString: "07:30 " + time.Now().Add(24*time.Hour).Format("02/01/2006"),
			expected:   true,
		},
		{
			name:       "Today",
			dateString: time.Now().Format("15:04 02/01/2006"),
			expected:   false,
		},
		{
			name:       "Yesterday",
			dateString: "07:30 " + time.Now().Add(-24*time.Hour).Format("02/01/2006"),
			expected:   false,
		},
		{
			name:       "DifferentTime",
			dateString: time.Now().Add(1 * time.Hour).Format("15:04 02/01/2006"),
			expected:   false,
		},
		{
			name:       "DifferentDate",
			dateString: "07:30 07/01/2023",
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsTomorrow(test.dateString)
			if actual != test.expected {
				t.Errorf("For input %q, expected %v, but got %v", test.dateString, test.expected, actual)
			}
		})
	}
}
