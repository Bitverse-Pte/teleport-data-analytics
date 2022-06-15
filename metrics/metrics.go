package metrics

import (
	"github.com/go-kit/kit/metrics"
	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricManager struct {
	Counter metrics.Counter
	Gauge   metrics.Gauge
}

func NewMetricManager() *MetricManager {
	var counter metrics.Counter
	var gauge metrics.Gauge
	counter = metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem: "teleport_bridge_backend",
		Name:      "counter",
		Help:      "system status",
	}, []string{"chain_name", "option"})
	gauge = metricsprometheus.NewGaugeFrom(prometheus.GaugeOpts{
		Subsystem: "teleport_bridge_backend",
		Name:      "gauge",
		Help:      "system status",
	}, []string{"chain_name", "option"})
	return &MetricManager{
		Counter: counter,
		Gauge:   gauge,
	}
}
