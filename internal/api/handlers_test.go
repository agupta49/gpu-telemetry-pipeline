package api
import "testing"
func TestGetGPUs(t *testing.T) { if len(GetGPUs()) != 0 { t.Fatal() } }
