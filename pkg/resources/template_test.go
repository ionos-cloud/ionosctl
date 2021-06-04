package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTemplateResourceVar = "test-template-resource"
)

func TestNewTemplateService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_templates_error", func(t *testing.T) {
		svc := getTestClient(t)
		templateSvc := NewTemplateService(svc.Get(), ctx)
		_, _, err := templateSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_template_error", func(t *testing.T) {
		svc := getTestClient(t)
		templateSvc := NewTemplateService(svc.Get(), ctx)
		_, _, err := templateSvc.Get(testTemplateResourceVar)
		assert.Error(t, err)
	})
}
