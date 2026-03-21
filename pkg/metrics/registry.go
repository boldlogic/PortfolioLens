package metrics

import "github.com/prometheus/client_golang/prometheus"

type Registry interface {
	prometheus.Registerer
	prometheus.Gatherer
}

func New() *prometheus.Registry {
	return prometheus.NewRegistry()
}
