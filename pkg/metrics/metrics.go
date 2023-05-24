package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	FailedSpaces = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "not_ready_spaces",
			Help: "Not Ready Spaces",
		}, []string{"name"},
	)
	ReadySpaces = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ready_spaces",
			Help: "Ready Spaces",
		}, []string{"name"},
	)
	InProgressSpaces = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "in_progress_spaces",
			Help: "In Progress Spaces",
		}, []string{"name"},
	)
)

func init() {
	metrics.Registry.MustRegister(FailedSpaces, ReadySpaces, InProgressSpaces)
}
