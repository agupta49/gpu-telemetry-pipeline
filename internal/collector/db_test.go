package collector

import "testing"

func TestNewConfig(t *testing.T) {
	c := NewConfig("localhost", "5432", "postgres")
	if c.Host != "localhost" || c.Port != "5432" || c.User != "postgres" {
		t.Fatal("NewConfig failed")
	}
}

func TestDSN(t *testing.T) {
	c := NewConfig("localhost", "5432", "postgres")
	c.Password = "pass"
	c.DBName = "telemetry"
	if c.DSN() == "" {
		t.Fatal("DSN empty")
	}
}

func TestValidate(t *testing.T) {
	c := NewConfig("", "5432", "postgres")
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for empty host")
	}
	c = NewConfig("localhost", "5432", "postgres")
	if err := c.Validate(); err != nil {
		t.Fatal("unexpected error")
	}
}
