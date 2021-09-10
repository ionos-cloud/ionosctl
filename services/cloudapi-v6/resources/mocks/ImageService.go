// Code generated by MockGen. DO NOT EDIT.
// Source: image.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
)

// MockImagesService is a mock of ImagesService interface.
type MockImagesService struct {
	ctrl     *gomock.Controller
	recorder *MockImagesServiceMockRecorder
}

// MockImagesServiceMockRecorder is the mock recorder for MockImagesService.
type MockImagesServiceMockRecorder struct {
	mock *MockImagesService
}

// NewMockImagesService creates a new mock instance.
func NewMockImagesService(ctrl *gomock.Controller) *MockImagesService {
	mock := &MockImagesService{ctrl: ctrl}
	mock.recorder = &MockImagesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImagesService) EXPECT() *MockImagesServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockImagesService) Delete(imageId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", imageId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockImagesServiceMockRecorder) Delete(imageId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockImagesService)(nil).Delete), imageId)
}

// Get mocks base method.
func (m *MockImagesService) Get(imageId string) (*resources.Image, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", imageId)
	ret0, _ := ret[0].(*resources.Image)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockImagesServiceMockRecorder) Get(imageId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockImagesService)(nil).Get), imageId)
}

// List mocks base method.
func (m *MockImagesService) List() (resources.Images, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].(resources.Images)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockImagesServiceMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockImagesService)(nil).List))
}

// Update mocks base method.
func (m *MockImagesService) Update(imageId string, imgProp resources.ImageProperties) (*resources.Image, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", imageId, imgProp)
	ret0, _ := ret[0].(*resources.Image)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockImagesServiceMockRecorder) Update(imageId, imgProp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockImagesService)(nil).Update), imageId, imgProp)
}
