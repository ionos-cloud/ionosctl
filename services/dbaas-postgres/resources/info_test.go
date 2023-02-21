package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInfosService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_info_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewInfosService(svc, ctx)
		_, _, err := infoUnitSvc.List()
		assert.NoError(t, err)
	})
	t.Run("get_infos_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewInfosService(svc, ctx)
		_, _, err := infoUnitSvc.Get()
		assert.NoError(t, err)
	})
}
