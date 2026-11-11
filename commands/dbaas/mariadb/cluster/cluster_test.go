package cluster

import (
	"testing"
	"time"
)

// TestParseRecoveryTime covers the friendlier --recovery-time parsing: keywords,
// several date/date-time layouts, naive values treated as UTC, and errors.
func TestParseRecoveryTime(t *testing.T) {
	t.Run("keywords and empty default to ~now", func(t *testing.T) {
		for _, in := range []string{"", "now", "NOW", "latest", " Latest "} {
			got, err := parseRecoveryTime(in)
			if err != nil {
				t.Fatalf("parseRecoveryTime(%q) unexpected err: %v", in, err)
			}
			if d := time.Since(got); d < 0 || d > time.Minute {
				t.Errorf("parseRecoveryTime(%q) = %v, want ~now", in, got)
			}
			if got.Location() != time.UTC {
				t.Errorf("parseRecoveryTime(%q) not UTC: %v", in, got.Location())
			}
		}
	})

	cases := map[string]string{
		"2025-01-02T15:04:05Z": "2025-01-02T15:04:05Z",
		"2025-01-02T15:04:05":  "2025-01-02T15:04:05Z", // naive -> UTC
		"2025-01-02 15:04:05":  "2025-01-02T15:04:05Z",
		"2025-01-02 15:04":     "2025-01-02T15:04:00Z",
		"2025-01-02":           "2025-01-02T00:00:00Z",
	}
	for in, wantRFC := range cases {
		t.Run("parses "+in, func(t *testing.T) {
			got, err := parseRecoveryTime(in)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Format(time.RFC3339) != wantRFC {
				t.Errorf("parseRecoveryTime(%q) = %s, want %s", in, got.Format(time.RFC3339), wantRFC)
			}
		})
	}

	t.Run("offset timestamp normalized to UTC", func(t *testing.T) {
		got, err := parseRecoveryTime("2025-01-02T15:04:05+02:00")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if want := "2025-01-02T13:04:05Z"; got.Format(time.RFC3339) != want {
			t.Errorf("got %s, want %s", got.Format(time.RFC3339), want)
		}
	})

	t.Run("garbage errors", func(t *testing.T) {
		if _, err := parseRecoveryTime("not-a-time"); err == nil {
			t.Error("expected error for unparseable input")
		}
	})
}

// TestToInt32GB guards the RAM/storage flag conversion: unit suffixes normalize
// to GB and a bad value errors instead of silently producing 0.
func TestToInt32GB(t *testing.T) {
	cases := map[string]int32{
		"4":    4,
		"4GB":  4,
		"16GB": 16,
		"1TB":  1024,
	}
	for in, want := range cases {
		got, err := toInt32GB(in, "ram")
		if err != nil {
			t.Fatalf("toInt32GB(%q) unexpected err: %v", in, err)
		}
		if got != want {
			t.Errorf("toInt32GB(%q) = %d, want %d", in, got, want)
		}
	}

	if _, err := toInt32GB("not-a-size", "ram"); err == nil {
		t.Error("expected error for unparseable size")
	}
}

// TestValidateBackupRetentionDays enforces the API's accepted 1-365 range.
func TestValidateBackupRetentionDays(t *testing.T) {
	valid := []int32{1, 30, 365}
	for _, d := range valid {
		if err := validateBackupRetentionDays(d); err != nil {
			t.Errorf("validateBackupRetentionDays(%d) unexpected err: %v", d, err)
		}
	}
	invalid := []int32{0, -1, 366}
	for _, d := range invalid {
		if err := validateBackupRetentionDays(d); err == nil {
			t.Errorf("validateBackupRetentionDays(%d) expected err, got nil", d)
		}
	}
}
