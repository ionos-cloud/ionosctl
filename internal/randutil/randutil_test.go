package randutil

import (
	"testing"
)

func TestCryptoRandN_ValidRange(t *testing.T) {
	n := 5
	for i := 0; i < 100; i++ {
		v, err := CryptoRandN(n)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v < 0 || v >= n {
			t.Fatalf("CryptoRandN(%d) = %d, want value in [0, %d)", n, v, n)
		}
	}
}

func TestCryptoRandN_ErrorOnZero(t *testing.T) {
	_, err := CryptoRandN(0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
}

func TestCryptoRandN_ErrorOnNegative(t *testing.T) {
	_, err := CryptoRandN(-1)
	if err == nil {
		t.Fatal("expected error for n=-1, got nil")
	}
}

func TestCryptoRandN_One(t *testing.T) {
	v, err := CryptoRandN(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 0 {
		t.Fatalf("CryptoRandN(1) = %d, want 0", v)
	}
}
