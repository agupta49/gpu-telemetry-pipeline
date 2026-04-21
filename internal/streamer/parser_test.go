package streamer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseCSV(t *testing.T) {
	rec := []string{"2025-07-18T20:42:34Z", "DCGM_FI_DEV_GPU_UTIL", "1", "nvidia1",
		"GPU-bc7a12ab-4998-fdc5-0785-2678a929a142", "NVIDIA H100 80GB HBM3",
		"mtv5-dgx1-hgpu-031", "", "", "", "100"}
	tp := ParseCSV(rec)
	assert.Equal(t, "GPU-bc7a12ab-4998-fdc5-0785-2678a929a142", tp.GPUID)
	assert.Equal(t, 100.0, tp.Value)
	assert.Equal(t, "mtv5-dgx1-hgpu-031", tp.Hostname)
	assert.Equal(t, "NVIDIA H100 80GB HBM3", tp.ModelName)
	assert.Equal(t, "nvidia1", tp.Device)
	assert.Equal(t, "DCGM_FI_DEV_GPU_UTIL", tp.MetricName)
	assert.Equal(t, "GPU-bc7a12ab-4998-fdc5-0785-2678a929a142", tp.Labels["uuid"])
}

func TestParseCSVInvalidValue(t *testing.T) {
	rec := []string{"2025-07-18T20:42:34Z", "DCGM_FI_DEV_GPU_UTIL", "1", "nvidia1",
		"GPU-123", "H100", "host1", "", "", "", "not-a-number"}
	tp := ParseCSV(rec)
	assert.Equal(t, 0.0, tp.Value)
	assert.Equal(t, "GPU-123", tp.GPUID)
}

func TestValidateRecord(t *testing.T) {
	valid := []string{"1","2","3","4","5","6","7","8","9","10","11"}
	assert.True(t, ValidateRecord(valid))
	
	short := []string{"1","2","3"}
	assert.False(t, ValidateRecord(short))
	
	noGPU := []string{"1","2","3","4","","6","7","8","9","10","11"}
	assert.False(t, ValidateRecord(noGPU))
}
