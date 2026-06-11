package blogs

import (
	"strings"
	"testing"
)

func TestCalculateEstimatedReadTime(t *testing.T) {
	s := &service{}

	tests := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Empty content",
			content:  "",
			expected: 0,
		},
		{
			name:     "Short content (< 200 words)",
			content:  "Ini adalah konten blog singkat.",
			expected: 1, // Minimal 1 menit jika ada konten
		},
		{
			name:     "Exactly 200 words",
			content:  strings.Repeat("kata ", 200),
			expected: 1,
		},
		{
			name:     "Exactly 400 words",
			content:  strings.Repeat("kata ", 400),
			expected: 2,
		},
		{
			name:     "450 words",
			content:  strings.Repeat("kata ", 450),
			expected: 2, // 450 / 200 = 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.CalculateEstimatedReadTime(tt.content)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
