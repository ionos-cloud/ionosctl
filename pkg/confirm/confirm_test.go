package confirm_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
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
			actual := confirm.FAsk(in, "Test", tt.overrides...)
			if actual != tt.expected {
				t.Errorf("FAsk() = %v; want %v", actual, tt.expected)
			}
		})
	}

}

type customConfirmer struct {
}

func (c customConfirmer) Ask(_ io.Reader, _ string, _ ...bool) bool {
	return true
}

func TestFAskCustom(t *testing.T) {
	t.Run("Test with custom strategy", func(t *testing.T) {

		confirm.SetStrategy(customConfirmer{})
		defer confirm.SetStrategy(nil) // Reset to default after test

		in := bytes.NewBufferString("no\n")
		if !confirm.FAsk(in, "Test") {
			t.Errorf("FAsk() with custom strategy should always return true")
		}
	})
}
