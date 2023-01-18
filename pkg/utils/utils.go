package utils

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// ReadPublicKey from a specific path
func ReadPublicKey(path string) (key string, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(pubKey)[:]), nil
}

// StringSlicesEqual returns true if 2 slices of
// type string are equal.
func StringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Map applies a function to a slice, and returns the modified slice
func Map[V comparable, R any](s []V, f func(int, V) R) []R {
	sm := make([]R, len(s))
	for i, v := range s {
		sm[i] = f(i, v)
	}
	return sm
}
