package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testResourceVar = "test-resource"
)

func TestNewUserService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_users_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_users_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.Get(testResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.Create(UserPost{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.Update(testResourceVar, UserPut{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, err := userSvc.Delete(testResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("listresources_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.ListResources()
		assert.Error(t, err)
	})
	t.Run("getresourcebytype_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.GetResourcesByType(testResourceVar)
		assert.Error(t, err)
	})
	t.Run("getresourcebytypeandId_user_error", func(t *testing.T) {
		svc := getTestClient(t)
		userSvc := NewUserService(svc, ctx)
		_, _, err := userSvc.GetResourceByTypeAndId(testResourceVar, testResourceVar)
		assert.Error(t, err)
	})
}
