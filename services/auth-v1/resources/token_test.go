package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTokenResourceVar = "test-token-resource"
	testContractNo       = int32(2)
)

func TestNewTokenService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_tokens_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.List(testContractNo)
		assert.Error(t, err)
	})
	t.Run("get_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.Get(testTokenResourceVar, testContractNo)
		assert.Error(t, err)
	})
	t.Run("create_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.Create(testContractNo)
		assert.Error(t, err)
	})
	t.Run("deleteById_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.DeleteByID(testTokenResourceVar, testContractNo)
		assert.Error(t, err)
	})
	t.Run("deleteByCriteria_token_error", func(t *testing.T) {
		svc := getTestClient(t)
		tokenSvc := NewTokenService(svc.Get(), ctx)
		_, _, err := tokenSvc.DeleteByCriteria(testTokenResourceVar, testContractNo)
		assert.Error(t, err)
	})
}
