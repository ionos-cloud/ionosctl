package randutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// CryptoRandN returns a cryptographically random integer in [0, n).
func CryptoRandN(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("randutil: n must be positive, got %d", n)
	}
	val, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(val.Int64()), nil
}
