package confirm

import (
	"bytes"
	"testing"
)

func TestFAsk(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		overrides []bool
		expected  bool
	}{
		{
			name:      "Test empty input",
			input:     "\n",
			overrides: []bool{},
			expected:  false,
		},
		{
			name:      "Test yes input",
			input:     "yes\n",
			overrides: []bool{},
			expected:  true,
		},
		{
			name:      "Test no input",
			input:     "no\n",
			overrides: []bool{},
			expected:  false,
		},
		{
			name:      "Test uppercase yes input",
			input:     "YES\n",
			overrides: []bool{},
			expected:  true,
		},
		{
			name:      "Test uppercase no input",
			input:     "NO\n",
			overrides: []bool{},
			expected:  false,
		},
		{
			name:      "Test override true",
			input:     "no\n",
			overrides: []bool{true},
			expected:  true,
		},
		{
			name:      "Test multiple overrides",
			input:     "\n",
			overrides: []bool{false, true, false},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := bytes.NewBufferString(tt.input)
			actual := FAsk(in, "Test", tt.overrides...)
			if actual != tt.expected {
				t.Errorf("FAsk() = %v; want %v", actual, tt.expected)
			}
		})
	}
}
