package nodepool

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func taint(key, value string, effect ionoscloud.TaintEffect) ionoscloud.KubernetesNodePoolTaint {
	t := ionoscloud.NewKubernetesNodePoolTaint(key, effect)
	if value != "" {
		t.SetValue(value)
	}
	return *t
}

func TestParseTaints(t *testing.T) {
	tests := []struct {
		name    string
		in      []string
		want    []ionoscloud.KubernetesNodePoolTaint
		wantErr string
	}{
		{name: "empty clears", in: []string{}, want: []ionoscloud.KubernetesNodePoolTaint{}},
		{name: "key value effect", in: []string{"dedicated=gpu:NoSchedule"}, want: []ionoscloud.KubernetesNodePoolTaint{taint("dedicated", "gpu", ionoscloud.NO_SCHEDULE)}},
		{name: "key only", in: []string{"spot:PreferNoSchedule"}, want: []ionoscloud.KubernetesNodePoolTaint{taint("spot", "", ionoscloud.PREFER_NO_SCHEDULE)}},
		{name: "prefixed key", in: []string{"example.com/team=a:NoExecute"}, want: []ionoscloud.KubernetesNodePoolTaint{taint("example.com/team", "a", ionoscloud.NO_EXECUTE)}},
		{name: "effect case-insensitive", in: []string{"a=b:noschedule"}, want: []ionoscloud.KubernetesNodePoolTaint{taint("a", "b", ionoscloud.NO_SCHEDULE)}},
		{name: "multiple", in: []string{"a:NoSchedule", " b=c:NoExecute "}, want: []ionoscloud.KubernetesNodePoolTaint{
			taint("a", "", ionoscloud.NO_SCHEDULE), taint("b", "c", ionoscloud.NO_EXECUTE),
		}},
		{name: "missing effect", in: []string{"a=b"}, wantErr: "expected KEY[=VALUE]:EFFECT"},
		{name: "empty effect", in: []string{"a=b:"}, wantErr: "expected KEY[=VALUE]:EFFECT"},
		{name: "bad effect", in: []string{"a=b:Never"}, wantErr: "effect must be one of NoSchedule, NoExecute, PreferNoSchedule"},
		{name: "empty key", in: []string{"=b:NoSchedule"}, wantErr: "key must not be empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTaints(tt.in)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTaintCompletions(t *testing.T) {
	assert.Nil(t, taintCompletions("dedicated=gpu"))
	assert.Equal(t, []string{"a:NoSchedule", "a:NoExecute", "a:PreferNoSchedule"}, taintCompletions("a:"))
	assert.Equal(t, []string{"a:NoSchedule,b=c:NoSchedule", "a:NoSchedule,b=c:NoExecute", "a:NoSchedule,b=c:PreferNoSchedule"},
		taintCompletions("a:NoSchedule,b=c:No"))
}

func TestTaintsColumn(t *testing.T) {
	var col func(map[string]any) any
	for _, c := range allK8sNodePoolCols {
		if c.Name == "Taints" {
			col = c.Format
		}
	}
	assert.NotNil(t, col)
	assert.Nil(t, col(map[string]any{"properties": map[string]any{}}))
	assert.Equal(t, "dedicated=gpu:NoSchedule,spot:PreferNoSchedule", col(map[string]any{"properties": map[string]any{
		"taints": []any{
			map[string]any{"key": "dedicated", "value": "gpu", "effect": "NoSchedule"},
			map[string]any{"key": "spot", "effect": "PreferNoSchedule"},
		},
	}}))
}

func TestGetNewK8sNodePoolTaints(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), "2048")
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageSize), "10")

		np, err := getNewK8sNodePool(cfg)
		assert.NoError(t, err)
		assert.Nil(t, np.Properties.Taints, "taints must be omitted when --taints is not set")

		viper.Set(core.GetFlagName(cfg.NS, constants.FlagTaints), []string{"dedicated=gpu:NoSchedule"})
		np, err = getNewK8sNodePool(cfg)
		assert.NoError(t, err)
		assert.Equal(t, []ionoscloud.KubernetesNodePoolTaint{taint("dedicated", "gpu", ionoscloud.NO_SCHEDULE)}, *np.Properties.Taints)

		viper.Set(core.GetFlagName(cfg.NS, constants.FlagTaints), []string{"dedicated=gpu"})
		_, err = getNewK8sNodePool(cfg)
		assert.Error(t, err)
	})
}

func TestGetNewK8sNodePoolUpdatedTaints(t *testing.T) {
	existing := []ionoscloud.KubernetesNodePoolTaint{taint("old", "v", ionoscloud.NO_EXECUTE)}
	oldNodePool := &resources.K8sNodePool{KubernetesNodePool: ionoscloud.KubernetesNodePool{
		Properties: &ionoscloud.KubernetesNodePoolProperties{Taints: &existing},
	}}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		fn := core.GetFlagName(cfg.NS, constants.FlagTaints)

		// unset: existing taints are carried over
		np, err := getNewK8sNodePoolUpdated(oldNodePool, cfg)
		assert.NoError(t, err)
		assert.Equal(t, existing, *np.Properties.Taints)

		// set: replaced
		viper.Set(fn, []string{"new:NoSchedule"})
		np, err = getNewK8sNodePoolUpdated(oldNodePool, cfg)
		assert.NoError(t, err)
		assert.Equal(t, []ionoscloud.KubernetesNodePoolTaint{taint("new", "", ionoscloud.NO_SCHEDULE)}, *np.Properties.Taints)

		// empty: cleared, sent as [] rather than omitted
		viper.Set(fn, []string{})
		np, err = getNewK8sNodePoolUpdated(oldNodePool, cfg)
		assert.NoError(t, err)
		body, err := np.Properties.MarshalJSON()
		assert.NoError(t, err)
		assert.Contains(t, string(body), `"taints":[]`)

		viper.Set(fn, []string{"bad"})
		_, err = getNewK8sNodePoolUpdated(oldNodePool, cfg)
		assert.Error(t, err)
	})
}
