package mq

import "testing"

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s == nil {
		t.Fatal("nil server")
	}
}
