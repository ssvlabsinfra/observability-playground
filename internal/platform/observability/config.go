package observability

type ObservabilityConfig struct {
	metricsEnabled, tracesEnabled, loggerEnabled bool
	metricExporters                              []Exporter
}
