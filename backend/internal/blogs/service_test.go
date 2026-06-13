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

func TestExtractFirstImage(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected *string
	}{
		{
			name:     "No image",
			content:  "This is a blog post.",
			expected: nil,
		},
		{
			name:     "Markdown image",
			content:  "Here is an image ![alt text](https://example.com/image.png) and some text.",
			expected: ptr("https://example.com/image.png"),
		},
		{
			name:     "HTML image",
			content:  `Here is an image <img src="https://example.com/image2.png" alt="test">`,
			expected: ptr("https://example.com/image2.png"),
		},
		{
			name:     "HTML single quote image",
			content:  `Here is an image <img src='https://example.com/image3.png'>`,
			expected: ptr("https://example.com/image3.png"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractFirstImage(tc.content)
			if (result == nil && tc.expected != nil) || (result != nil && tc.expected == nil) || (result != nil && tc.expected != nil && *result != *tc.expected) {
				var r, e string
				if result != nil { r = *result } else { r = "nil" }
				if tc.expected != nil { e = *tc.expected } else { e = "nil" }
				t.Errorf("Expected %s, got %s", e, r)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}
