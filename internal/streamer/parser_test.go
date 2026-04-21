package streamer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseCSV(t *testing.T) {
	rec := []string{"2025-07-18T20:42:34Z", "DCGM_FI_DEV_GPU_UTIL", "1", "nvidia1",
		"GPU-bc7a12ab-4998-fdc5-0785-2678a929a142", "NVIDIA H100 80GB HBM3",
		"mtv5-dgx1-hgpu-031", "", "", "", "100"}
	tp := parseCSV(rec)
	assert.Equal(t, "GPU-bc7a12ab-4998-fdc5-0785-2678a929a142", tp.GPUID)
	assert.Equal(t, 100.0, tp.Value)
	assert.Equal(t, "mtv5-dgx1-hgpu-031", tp.Hostname)
}
