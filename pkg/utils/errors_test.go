package utils

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequiredFlagErr(t *testing.T) {
	requiredFlag := NewRequiredFlagErr("dummy")
	assert.True(t, requiredFlag.Error() == "dummy required flag is not set")
}

func TestVarErrAction(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		ErrAction()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestVarErrAction")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
