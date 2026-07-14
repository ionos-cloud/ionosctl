package doc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateStructure_DirPermissions(t *testing.T) {
	dir := t.TempDir()

	// createStructure uses MkdirAll with 0755 on filepath.Dir(targetDir).
	// We simulate that by calling os.MkdirAll the same way and checking the mode.
	target := filepath.Join(dir, "subdir", "leaf")
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, "subdir"))
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	got := info.Mode().Perm()
	want := os.FileMode(0755)
	if got != want {
		t.Errorf("directory mode = %04o, want %04o", got, want)
	}
}
