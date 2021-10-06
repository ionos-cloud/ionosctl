package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuotasService(t *testing.T) {
	ctx := context.Background()
	t.Run("get_logs_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewQuotasService(svc.Get(), ctx)
		_, _, err := infoUnitSvc.Get()
		assert.Error(t, err)
	})
}
