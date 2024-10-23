package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ssvlabsinfra/p2p-observability/internal/platform/lifecycle"
	"github.com/ssvlabsinfra/p2p-observability/internal/platform/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	appName    = "p2p-observability"
	appVersion = "0.1.0"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	otelShutdown, err := observability.SetupOTelSDK(ctx, appName, appVersion)
	if err != nil {
		panic(err.Error())
	}

	meter := otel.Meter("github.com/ssvlabsinfra/p2p-observability/main", metric.WithInstrumentationVersion("0.0.1"))
	counter, err := meter.Int64Counter("p2p-observability", metric.WithDescription("p2p observability counter description"))
	if err != nil {
		panic(err.Error())
	}
	slog.Info("OTeL SDK configured. Listening for application shutdown")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				counter.Add(ctx, 1, metric.WithAttributes(attribute.String("attr", "custom_attribute")))
			}
		}
	}()

	host := http.NewServeMux()
	host.Handle("/metrics", promhttp.Handler())

	go func() {
		err := http.ListenAndServe("0.0.0.0:8080", host)
		if err != nil {
			cancel()
			panic(err.Error())
		}
	}()

	lifecycle.ListenForApplicationShutDown(ctx, func() {
		if err = otelShutdown(ctx); err != nil {
			panic(err.Error())
		}
		cancel()
	}, make(chan os.Signal))
}
