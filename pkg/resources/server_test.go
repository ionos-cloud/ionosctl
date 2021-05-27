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
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List(testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(
			testServerResourceVar,
			Server{},
		)
		assert.Error(t, err)
	})
	t.Run("update_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testServerResourceVar, testServerResourceVar, ServerProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("start_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Start(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("stop_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Stop(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("reboot_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Reboot(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("suspend_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Suspend(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("resume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Resume(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_token_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.GetToken(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_remote_console_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.GetRemoteConsoleUrl(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("attach_volume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.AttachVolume(testServerResourceVar, testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("list_volumes_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.ListVolumes(testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_volume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.GetVolume(testServerResourceVar, testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
	t.Run("detach_volume_server_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewServerService(svc.Get(), ctx)
		_, err := backupUnitSvc.DetachVolume(testServerResourceVar, testServerResourceVar, testServerResourceVar)
		assert.Error(t, err)
	})
}
