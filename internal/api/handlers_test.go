package api

import "testing"

func TestTelemetryPoint(t *testing.T) {
	p := TelemetryPoint{GPUID: "0"}
	if p.GPUID != "0" {
		t.Fatal("bad gpu id")
	}
}
