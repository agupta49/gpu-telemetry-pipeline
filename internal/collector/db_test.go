package collector

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("localhost", "postgres", "telemetry")
	assert.Equal(t, "localhost", cfg.PostgresHost)
	assert.Equal(t, 5432, cfg.PostgresPort)
	assert.Equal(t, "postgres", cfg.PostgresUser)
	assert.Equal(t, "telemetry", cfg.PostgresDB)
	assert.Equal(t, "mq:50051", cfg.MQAddress)
}

func TestDSN(t *testing.T) {
	cfg := NewConfig("localhost", "postgres", "telemetry")
	dsn := cfg.DSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "dbname=telemetry")
	assert.Contains(t, dsn, "sslmode=disable")
}

func TestValidate(t *testing.T) {
	cfg := NewConfig("localhost", "postgres", "telemetry")
	assert.NoError(t, cfg.Validate())
	
	cfg2 := &Config{PostgresHost: ""}
	assert.Error(t, cfg2.Validate())
	
	cfg3 := &Config{PostgresHost: "localhost", PostgresDB: ""}
	assert.Error(t, cfg3.Validate())
}
