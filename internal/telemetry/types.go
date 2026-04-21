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

func NewTelemetryPoint(metricName, gpuID, hostname string, value float64) TelemetryPoint {
	return TelemetryPoint{
		Timestamp: time.Now(),
		MetricName: metricName,
		GPUID: gpuID,
		Hostname: hostname,
		Value: value,
		Labels: make(map[string]string),
	}
}

func (tp *TelemetryPoint) AddLabel(key, value string) {
	if tp.Labels == nil {
		tp.Labels = make(map[string]string)
	}
	tp.Labels[key] = value
}

func (tp *TelemetryPoint) IsHighUtilization() bool {
	return tp.MetricName == "DCGM_FI_DEV_GPU_UTIL" && tp.Value > 90.0
}
