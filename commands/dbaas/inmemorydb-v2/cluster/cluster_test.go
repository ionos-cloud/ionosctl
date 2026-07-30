package cluster

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

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
