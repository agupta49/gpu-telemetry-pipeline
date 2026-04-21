package telemetry

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestTelemetryPoint(t *testing.T) {
	tp := TelemetryPoint{
		Timestamp: time.Now(),
		MetricName: "DCGM_FI_DEV_GPU_UTIL",
		GPUID: "GPU-123",
		Hostname: "host1",
		ModelName: "H100",
		Device: "0",
		Value: 95.5,
		Labels: map[string]string{"uuid": "GPU-123"},
	}
	assert.Equal(t, "GPU-123", tp.GPUID)
	assert.Equal(t, 95.5, tp.Value)
	assert.Equal(t, "host1", tp.Hostname)
	assert.Equal(t, "H100", tp.ModelName)
}
