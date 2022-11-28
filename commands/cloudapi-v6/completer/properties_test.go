package completer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatacenterCPUFamiliesDefault(t *testing.T) {
	cpuFamilies := DatacenterCPUFamilies(context.Background(), nil, "")
	assert.Equal(t, []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cpuFamilies)
}
