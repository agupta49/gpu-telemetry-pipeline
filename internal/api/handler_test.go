package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	assert.Equal(t, "ok", HealthStatus())
}

func HealthStatus() string {
	return "ok"
}
