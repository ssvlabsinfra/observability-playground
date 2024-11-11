package observability

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	loggerName = "observability-playground"
)

var (
	once   sync.Once
	config ObservabilityConfig
)

func Initialize(ctx context.Context, appName, appVersion string, options ...Option) (shutdown func(context.Context) error, err error) {
	var (
		initErr       error
		shutdownFuncs []func(context.Context) error
	)

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	once.Do(func() {
		for _, option := range options {
			option(&config)
		}

		resources, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
			semconv.ServiceVersion(appVersion),
		))
		if err != nil {
			initErr = errors.Join(errors.New("failed to instantiate observability resources"), err)
			return
		}

		if config.metricsEnabled {
			var readers []metric.Reader
			for _, exporter := range config.metricExporters {
				switch exporter {
				case Prometheus:
					promExporter, err := prometheus.New()
					if err != nil {
						initErr = errors.Join(errors.New("failed to instantiate metric Prometheus exporter"), err)
						return
					}
					readers = append(readers, promExporter)
				case GRPC:
					gRPCExporter, err := otlpmetricgrpc.New(
						ctx,
						otlpmetricgrpc.WithInsecure(),
						otlpmetricgrpc.WithEndpoint("otel-collector:4317"))
					if err != nil {
						initErr = errors.Join(errors.New("failed to instantiate metric gRPC exporter"), err)
						return
					}
					reader := metric.NewPeriodicReader(gRPCExporter, metric.WithInterval(time.Second*10))
					readers = append(readers, reader)
				case Stdout:
					stdoutExporter, err := stdoutmetric.New()
					if err != nil {
						initErr = errors.Join(errors.New("failed to instantiate metric stdout exporter"), err)
						return
					}
					reader := metric.NewPeriodicReader(stdoutExporter, metric.WithInterval(time.Second*5))
					readers = append(readers, reader)
				}
			}

			var options []metric.Option
			options = append(options, metric.WithResource(resources))
			for _, reader := range readers {
				options = append(options, metric.WithReader(reader))
			}

			meterProvider := metric.NewMeterProvider(options...)

			shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
			otel.SetMeterProvider(meterProvider)
		}

		if config.tracesEnabled {
			traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
			if err != nil {
				initErr = errors.Join(errors.New("failed to instantiate traces stdout exporter"), err)
				return
			}

			traceProvider := trace.NewTracerProvider(
				trace.WithResource(resources),
				trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second*5)),
			)
			shutdownFuncs = append(shutdownFuncs, traceExporter.Shutdown)
			otel.SetTracerProvider(traceProvider)
		}

		if config.loggerEnabled {
			logExporter, err := stdoutlog.New()
			if err != nil {
				initErr = errors.Join(errors.New("failed to instantiate logger stdout exporter"), err)
				return
			}

			loggerProvider := log.NewLoggerProvider(
				log.WithResource(resources),
				log.WithProcessor(log.NewBatchProcessor(logExporter)),
			)
			shutdownFuncs = append(shutdownFuncs, logExporter.Shutdown)
			global.SetLoggerProvider(loggerProvider)
			logger := otelslog.NewLogger(loggerName)
			slog.SetDefault(logger)
		}
	})

	return shutdown, initErr
}
