package telemetry

import "time"

type TelemetryPoint struct {
	Timestamp time.Time         `json:"timestamp"`
	MetricName string           `json:"metric_name"`
	GPUID     string            `json:"gpu_id"`
	Hostname  string            `json:"hostname"`
	ModelName string            `json:"model_name"`
	Device    string            `json:"device"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels,omitempty"`
}
