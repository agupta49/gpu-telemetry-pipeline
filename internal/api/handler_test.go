package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewGPU(t *testing.T) {
	g := NewGPU("GPU-123", "H100")
	assert.Equal(t, "GPU-123", g.ID)
	assert.Equal(t, "H100", g.Model)
}

func TestGPUString(t *testing.T) {
	g := NewGPU("GPU-123", "H100")
	s := g.String()
	assert.Contains(t, s, "GPU-123")
	assert.Contains(t, s, "H100")
}

func TestHealthStatus(t *testing.T) {
	assert.Equal(t, "ok", HealthStatus())
}

func TestFilterGPUs(t *testing.T) {
	gpus := []GPU{
		{ID: "1", Model: "H100"},
		{ID: "2", Model: "A100"},
		{ID: "3", Model: "H100"},
	}
	filtered := FilterGPUs(gpus, "H100")
	assert.Len(t, filtered, 2)
	assert.Equal(t, "1", filtered[0].ID)
	assert.Equal(t, "3", filtered[1].ID)
	
	empty := FilterGPUs(gpus, "V100")
	assert.Len(t, empty, 0)
}
