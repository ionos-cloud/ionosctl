package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testTokenResourceVar = "test-token-resource"

func TestNewTokenService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_tokens_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.Get(testTokenResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.Create()
		assert.Error(t, err)
	})
	t.Run("deleteById_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.DeleteByID(testTokenResourceVar)
		assert.Error(t, err)
	})
	t.Run("deleteByCriteria_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.DeleteByCriteria(testTokenResourceVar)
		assert.Error(t, err)
	})
}
