package schemas

const (
	MetricStatusOk    = "ok"
	MetricStatusWarn  = "warn"
	MetricStatusError = "error"
)

type MonitoringMetric struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
