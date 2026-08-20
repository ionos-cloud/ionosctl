package nodepool

import (
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"
)

func TestParseTaints(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	effPtr := func(e ionoscloud.TaintEffect) *ionoscloud.TaintEffect { return &e }

	t.Run("key value and effect", func(t *testing.T) {
		got, err := parseTaints([]string{"dedicated=gpu:NoSchedule"})
		assert.NoError(t, err)
		assert.Equal(t, []ionoscloud.KubernetesNodePoolTaint{
			{Key: strPtr("dedicated"), Value: strPtr("gpu"), Effect: effPtr(ionoscloud.NO_SCHEDULE)},
		}, got)
	})

	t.Run("key and effect without value", func(t *testing.T) {
		got, err := parseTaints([]string{"dedicated:NoExecute"})
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, "dedicated", *got[0].Key)
		assert.Nil(t, got[0].Value)
		assert.Equal(t, ionoscloud.NO_EXECUTE, *got[0].Effect)
	})

	t.Run("multiple taints", func(t *testing.T) {
		got, err := parseTaints([]string{"a=1:NoSchedule", "b:PreferNoSchedule"})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, ionoscloud.PREFER_NO_SCHEDULE, *got[1].Effect)
	})

	t.Run("empty input", func(t *testing.T) {
		got, err := parseTaints([]string{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("missing effect", func(t *testing.T) {
		_, err := parseTaints([]string{"dedicated=gpu"})
		assert.Error(t, err)
	})

	t.Run("empty effect", func(t *testing.T) {
		_, err := parseTaints([]string{"dedicated="})
		assert.Error(t, err)
	})

	t.Run("invalid effect", func(t *testing.T) {
		_, err := parseTaints([]string{"dedicated=gpu:Nope"})
		assert.Error(t, err)
	})

	t.Run("empty key", func(t *testing.T) {
		_, err := parseTaints([]string{"=gpu:NoSchedule"})
		assert.Error(t, err)
	})
}
