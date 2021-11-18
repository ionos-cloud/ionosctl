package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testServerResourceVar = "test-server-resource"
)

func TestNewServerService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_servers_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.List(testServerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_servers_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.List(testServerResourceVar, testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.Get(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.Create(
			testServerResourceVar,
			Server{},
		)
		assert.Error(t, err)
	})
	t.Run("update_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.Update(testServerResourceVar, testServerResourceVar, ServerProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, err := serverSvc.Delete(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("start_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, err := serverSvc.Start(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("stop_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, err := serverSvc.Stop(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("reboot_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, err := serverSvc.Reboot(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("attach_volume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.AttachVolume(testServerResourceVar, testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("list_volumes_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.ListVolumes(testServerResourceVar, testServerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_volumes_server_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.ListVolumes(testServerResourceVar, testServerResourceVar, testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_volume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, _, err := serverSvc.GetVolume(testServerResourceVar, testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("detach_volume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		serverSvc := NewServerService(svc.Get(), ctx)
		_, err := serverSvc.DetachVolume(testServerResourceVar, testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
}
