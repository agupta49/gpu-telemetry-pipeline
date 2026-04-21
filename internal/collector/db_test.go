package collector

import "testing"

func TestNewRepoEmptyDSN(t *testing.T) {
	_, err := NewRepo("")
	if err == nil {
		t.Fatal("expected error for empty dsn")
	}
}
