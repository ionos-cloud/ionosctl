package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContractService(t *testing.T) {
	ctx := context.Background()
	t.Run("get_contract_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewContractService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get()
		assert.Error(t, err)
	})
}
