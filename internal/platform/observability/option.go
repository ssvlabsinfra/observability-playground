package observability

type (
	Exporter string

	Option func(*ObservabilityConfig)
)

const (
	Prometheus Exporter = "Prometheus"
	Stdout     Exporter = "Stdout"
	GRPC       Exporter = "gRPC"
)

func WithMetrics(exporters []Exporter) Option {
	return func(cfg *ObservabilityConfig) {
		cfg.metricsEnabled = true
		cfg.metricExporters = exporters
	}
}

func WithTraces() Option {
	return func(cfg *ObservabilityConfig) {
		cfg.tracesEnabled = true
	}
}

func WithLogger() Option {
	return func(cfg *ObservabilityConfig) {
		cfg.loggerEnabled = true
	}
}
