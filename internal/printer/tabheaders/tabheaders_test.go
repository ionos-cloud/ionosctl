package tabheaders

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureStderr redirects os.Stderr to a pipe, runs f, and returns what was written.
func captureStderr(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	orig := os.Stderr
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = orig

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("io.ReadAll: %v", err)
	}
	r.Close()

	return string(out)
}

var (
	allCols     = []string{"Name", "Age", "Email"}
	defaultCols = []string{"Name", "Age"}
)

func TestGetHeaders_AllValid(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeaders(allCols, defaultCols, []string{"Name", "Email"})
		want := []string{"Name", "Email"}
		if len(got) != len(want) {
			t.Errorf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	if stderr != "" {
		t.Errorf("expected no stderr output, got: %q", stderr)
	}
}

func TestGetHeaders_MixedValidInvalid(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeaders(allCols, defaultCols, []string{"Name", "Bogus"})
		want := []string{"Name"}
		if len(got) != len(want) {
			t.Errorf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	if !strings.Contains(stderr, "Warning:") {
		t.Errorf("expected warning on stderr, got: %q", stderr)
	}
	if !strings.Contains(stderr, "Bogus") {
		t.Errorf("expected 'Bogus' in warning, got: %q", stderr)
	}
}

func TestGetHeaders_AllInvalid(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeaders(allCols, defaultCols, []string{"Typo1", "Typo2"})
		// Falls back to defaultCols
		want := defaultCols
		if len(got) != len(want) {
			t.Errorf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	if !strings.Contains(stderr, "Warning:") {
		t.Errorf("expected warning on stderr, got: %q", stderr)
	}
	if !strings.Contains(stderr, "Typo1") || !strings.Contains(stderr, "Typo2") {
		t.Errorf("expected both unknown columns in warning, got: %q", stderr)
	}
}

func TestGetHeaders_ColsAll(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeaders(allCols, defaultCols, []string{"all"})
		if len(got) != len(allCols) {
			t.Errorf("got %v, want %v", got, allCols)
		}
		for i := range allCols {
			if got[i] != allCols[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], allCols[i])
			}
		}
	})

	if stderr != "" {
		t.Errorf("expected no stderr for 'all', got: %q", stderr)
	}
}

func TestGetHeaders_NilCustomColumns(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeaders(allCols, defaultCols, nil)
		if len(got) != len(defaultCols) {
			t.Errorf("got %v, want %v", got, defaultCols)
		}
		for i := range defaultCols {
			if got[i] != defaultCols[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], defaultCols[i])
			}
		}
	})

	if stderr != "" {
		t.Errorf("expected no stderr for nil customColumns, got: %q", stderr)
	}
}

func TestGetHeaders_CaseInsensitiveMatch(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeaders(allCols, defaultCols, []string{"name", "EMAIL"})
		want := []string{"Name", "Email"}
		if len(got) != len(want) {
			t.Errorf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	if stderr != "" {
		t.Errorf("expected no warning for case-insensitive match, got: %q", stderr)
	}
}

func TestGetHeadersAllDefault(t *testing.T) {
	stderr := captureStderr(t, func() {
		got := GetHeadersAllDefault(allCols, []string{"Name"})
		want := []string{"Name"}
		if len(got) != len(want) {
			t.Errorf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	if stderr != "" {
		t.Errorf("expected no stderr, got: %q", stderr)
	}
}
