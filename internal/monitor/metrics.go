package monitor

import "github.com/prometheus/client_golang/prometheus"

var checksTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "pulsewatch_checks_total",
		Help: "Total number of completed monitor checks by result.",
	},
	[]string{"status"},
)

func init() {
	prometheus.MustRegister(checksTotal)
}

func recordCheck(result CheckResult) {
	checksTotal.WithLabelValues(result.Status).Inc()
}
