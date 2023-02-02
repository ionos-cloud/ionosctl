package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testGroupResourceVar = "test-group-resource"
)

func TestNewGroupService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_groups_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_groups_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.Get(testGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.Create(Group{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.Update(testGroupResourceVar, Group{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, err := groupSvc.Delete(testGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("listresources_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.ListResources(testGroupResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("adduser_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.AddUser(testGroupResourceVar, User{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("listusers_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.ListUsers(testGroupResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("removeuser_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, err := groupSvc.RemoveUser(testGroupResourceVar, testGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("addshare_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.AddShare(testGroupResourceVar, testGroupResourceVar, GroupShare{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("updateshare_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.UpdateShare(testGroupResourceVar, testGroupResourceVar, GroupShare{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("getshare_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.GetShare(testGroupResourceVar, testGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("listshares_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, _, err := groupSvc.ListShares(testGroupResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("removeshare_group_error", func(t *testing.T) {
		svc := getTestClient(t)
		groupSvc := NewGroupService(svc, ctx)
		_, err := groupSvc.RemoveShare(testGroupResourceVar, testGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
