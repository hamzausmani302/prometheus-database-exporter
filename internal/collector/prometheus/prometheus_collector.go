/*
contains actual implementaion of a collector
for prometheus
*/
package promcollector

import (
	schemacollector "github.com/hamzausmani302/prometheus-database-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type PrometheusGoCollector struct {
    Logger  *logrus.Logger
    Collector schemacollector.ICollector[schemacollector.MMetricResultType]
}

// Describe sends metric descriptions to Prometheus
func (p PrometheusGoCollector) Describe(ch chan<- *prometheus.Desc) {
    metrics, err := p.Collector.GetCollectedMetrics()
    if err != nil {
        p.Logger.Warn(err)
    }
    for _, m := range metrics {
		labels := []string{}
		for _, label := range m.Labels{
			labels = append(labels, label.Name)
		}
        // assume CollectorMetric has a Desc() method or fields for name/help/labels
        ch <- prometheus.NewDesc(
            m.Name,        // metric name
            m.Help,        // help text
            labels,      // variable labels
            nil,           // constant labels
        )
    }
}

// Collect sends metric values to Prometheus
func (p PrometheusGoCollector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := p.Collector.GetCollectedMetrics()
    if err != nil {
        p.Logger.Warn(err)
    }
    descLabels:= []string{}
    for _, m := range metrics {
        labels := []string{}
		for _, label := range m.Labels{
			labels = append(labels, label.Value)
            descLabels = append(descLabels, label.Name)
        }
		desc := prometheus.NewDesc(m.Name, m.Help, descLabels, nil)
        // add support for counter  nad everything else via a factory function 
        // by supplying the channel and let the factory handle putting data in channel
        metric, err := prometheus.NewConstMetric(
            desc,
            prometheus.GaugeValue,  // prometheus.GaugeValue / CounterValue etc.
            float64(m.Value),
            labels...,
        )
        if err != nil {
            p.Logger.Errorf("error creating metric %s: %v", m.Name, err)
            continue
        }

        ch <- metric
    }
}
