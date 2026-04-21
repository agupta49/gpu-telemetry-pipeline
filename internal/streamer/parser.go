package streamer

import (
	"strconv"
	"time"
	"github.com/agupta49/gpu-telemetry-pipeline/internal/telemetry"
)

func ParseCSV(rec []string) telemetry.TelemetryPoint {
	val, _ := strconv.ParseFloat(rec[10], 64)
	return telemetry.TelemetryPoint{
		Timestamp: time.Now(),
		MetricName: rec[1],
		GPUID:      rec[4],
		Hostname:   rec[6],
		ModelName:  rec[5],
		Device:     rec[3],
		Value:      val,
		Labels:     map[string]string{"uuid": rec[4], "modelName": rec[5], "Hostname": rec[6]},
	}
}
