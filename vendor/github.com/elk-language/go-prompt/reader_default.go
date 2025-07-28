//go:build !unix && !windows
// +build !unix,!windows

package prompt

import (
	"os"
)

// DefaultReader is a Reader implementation for environments other than Unix and Windows.
type DefaultReader struct{}

// Open should be called before starting input
func (t DefaultReader) Open() error {
	return nil
}

// Close should be called after stopping input
func (t DefaultReader) Close() error {
	return nil
}

// Read returns byte array.
func (t DefaultReader) Read(buff []byte) (int, error) {
	return os.Stdin.Read(buff)
}

// GetWinSize returns WinSize object to represent width and height of terminal.
func (t DefaultReader) GetWinSize() *WinSize {
	return &WinSize{
		Row: DefRowCount,
		Col: DefColCount,
	}
}

var _ Reader = DefaultReader{}

// NewStdinReader returns Reader object to read from stdin.
func NewStdinReader() DefaultReader {
	return DefaultReader{}
}
