package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationColsHaveEnabledFeatures(t *testing.T) {
	var cpuHas bool
	for _, c := range allCpuCols {
		if c.Name == "EnabledFeatures" {
			cpuHas = true
		}
	}
	assert.True(t, cpuHas, "allCpuCols must expose an EnabledFeatures column (CoCo CPU-family discovery)")

	var locHas bool
	for _, c := range allLocationCols {
		if c.Name == "CpuEnabledFeatures" {
			locHas = true
		}
	}
	assert.True(t, locHas, "allLocationCols must expose a CpuEnabledFeatures column")
}
