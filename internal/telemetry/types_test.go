package telemetry

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNewTelemetryPoint(t *testing.T) {
	tp := NewTelemetryPoint("DCGM_FI_DEV_GPU_UTIL", "GPU-123", "host1", 95.5)
	assert.Equal(t, "GPU-123", tp.GPUID)
	assert.Equal(t, 95.5, tp.Value)
	assert.Equal(t, "host1", tp.Hostname)
	assert.WithinDuration(t, time.Now(), tp.Timestamp, time.Second)
}

func TestAddLabel(t *testing.T) {
	tp := NewTelemetryPoint("test", "GPU-1", "host", 50.0)
	tp.AddLabel("rack", "A1")
	assert.Equal(t, "A1", tp.Labels["rack"])
	
	tp2 := TelemetryPoint{}
	tp2.AddLabel("test", "val")
	assert.Equal(t, "val", tp2.Labels["test"])
}

func TestIsHighUtilization(t *testing.T) {
	tp1 := NewTelemetryPoint("DCGM_FI_DEV_GPU_UTIL", "GPU-1", "host", 95.0)
	assert.True(t, tp1.IsHighUtilization())
	
	tp2 := NewTelemetryPoint("DCGM_FI_DEV_GPU_UTIL", "GPU-1", "host", 80.0)
	assert.False(t, tp2.IsHighUtilization())
	
	tp3 := NewTelemetryPoint("DCGM_FI_DEV_MEM_COPY_UTIL", "GPU-1", "host", 95.0)
	assert.False(t, tp3.IsHighUtilization())
}
