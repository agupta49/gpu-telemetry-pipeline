package collector
import "testing"
func TestNewConfig(t *testing.T) { _ = NewConfig() }
func TestDSN(t *testing.T) { c := NewConfig(); c.DSN = "x"; if c.DSNStr() != "x" { t.Fatal() } }
func TestValidate(t *testing.T) { if NewConfig().Validate() != nil { t.Fatal() } }
