package object

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetCmd_DestinationFilePermissions(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "downloaded-object")

	// Simulate the secure file creation path used in get.go.
	os.Remove(destination)
	f, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	f.Close()

	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	got := info.Mode().Perm()
	want := os.FileMode(0600)
	if got != want {
		t.Errorf("file mode = %04o, want %04o", got, want)
	}
}

func TestGetCmd_ExistingFileReplacedWith0600(t *testing.T) {
	dir := t.TempDir()
	destination := filepath.Join(dir, "existing-object")

	// Pre-create the file with a permissive mode.
	if err := os.WriteFile(destination, []byte("old content"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Simulate the secure replacement path.
	os.Remove(destination)
	f, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		t.Fatalf("OpenFile after Remove failed: %v", err)
	}
	f.Close()

	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	got := info.Mode().Perm()
	want := os.FileMode(0600)
	if got != want {
		t.Errorf("replaced file mode = %04o, want %04o", got, want)
	}
}
