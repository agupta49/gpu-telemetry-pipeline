package collector

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCollectorConfig(t *testing.T) {
	cfg := Config{
		PostgresHost: "localhost",
		PostgresDB: "test",
	}
	assert.Equal(t, "localhost", cfg.PostgresHost)
	assert.Equal(t, "test", cfg.PostgresDB)
}

type Config struct {
	PostgresHost string
	PostgresDB   string
}
