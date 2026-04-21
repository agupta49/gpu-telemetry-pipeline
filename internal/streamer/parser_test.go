package streamer
import "testing"
func TestParse(t *testing.T) { if Parse("a") != 1 { t.Fatal("fail") } }
