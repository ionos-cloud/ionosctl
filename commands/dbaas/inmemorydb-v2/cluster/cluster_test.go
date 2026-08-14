package cluster

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
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

// TestBuildHashedPassword guards the client-side password hashing: the In-Memory
// DB API only accepts a SHA-256 HashedPassword, so a plaintext value must be
// hashed, and a value that is already a SHA-256 hash must be forwarded verbatim
// (not double-hashed).
func TestBuildHashedPassword(t *testing.T) {
	t.Run("plaintext is SHA-256 hashed", func(t *testing.T) {
		sum := sha256.Sum256([]byte("hunter2"))
		want := hex.EncodeToString(sum[:])

		got := buildHashedPassword("hunter2")
		if got.Algorithm != "SHA-256" {
			t.Errorf("Algorithm = %q, want SHA-256", got.Algorithm)
		}
		if got.Hash != want {
			t.Errorf("Hash = %q, want %q", got.Hash, want)
		}
	})

	t.Run("existing SHA-256 hash passes through unchanged", func(t *testing.T) {
		sum := sha256.Sum256([]byte("hunter2"))
		hashed := hex.EncodeToString(sum[:]) // 64 hex chars

		got := buildHashedPassword(hashed)
		if got.Algorithm != "SHA-256" {
			t.Errorf("Algorithm = %q, want SHA-256", got.Algorithm)
		}
		if got.Hash != hashed {
			t.Errorf("Hash = %q, want passthrough %q", got.Hash, hashed)
		}
	})

	t.Run("uppercase hash is normalized to lowercase", func(t *testing.T) {
		sum := sha256.Sum256([]byte("hunter2"))
		lower := hex.EncodeToString(sum[:])
		upper := strings.ToUpper(lower)

		got := buildHashedPassword(upper)
		if got.Hash != lower {
			t.Errorf("Hash = %q, want lowercased %q", got.Hash, lower)
		}
	})

	t.Run("non-hex 64-char string is still hashed", func(t *testing.T) {
		// 64 chars but contains non-hex 'z' — must NOT be treated as a hash.
		in := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
		sum := sha256.Sum256([]byte(in))
		want := hex.EncodeToString(sum[:])

		if got := buildHashedPassword(in); got.Hash != want {
			t.Errorf("Hash = %q, want hashed %q", got.Hash, want)
		}
	})
}
