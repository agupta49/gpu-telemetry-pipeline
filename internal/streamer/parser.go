package streamer

import (
	"strconv"
	"github.com/agupta49/gpu-telemetry-pipeline/internal/telemetry"
)

func ParseCSV(rec []string) telemetry.TelemetryPoint {
	val, _ := strconv.ParseFloat(rec[10], 64)
	tp := telemetry.NewTelemetryPoint(rec[1], rec[4], rec[6], val)
	tp.ModelName = rec[5]
	tp.Device = rec[3]
	tp.AddLabel("uuid", rec[4])
	tp.AddLabel("modelName", rec[5])
	tp.AddLabel("Hostname", rec[6])
	return tp
}

func ValidateRecord(rec []string) bool {
	if len(rec) < 11 {
		return false
	}
	if rec[4] == "" {
		return false
	}
	return true
}
