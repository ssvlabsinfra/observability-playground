package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ssvlabsinfra/p2p-observability/internal/platform/lifecycle"
	"github.com/ssvlabsinfra/p2p-observability/internal/platform/observability"
	"github.com/ssvlabsinfra/p2p-observability/internal/processor"
	"github.com/ssvlabsinfra/p2p-observability/internal/storage"
)

const (
	appName    = "p2p-observability"
	appVersion = "0.1.0"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	shutdownFunc, err := observability.Initialize(ctx, appName, appVersion,
		observability.WithMetrics([]observability.Exporter{observability.Prometheus, observability.GRPC, observability.Stdout}),
		observability.WithTraces(),
		observability.WithLogger())
	if err != nil {
		panic(err.Error())
	}

	slog.Info("OTeL SDK configured. Listening for application shutdown")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	if err != nil {
		panic(err.Error())
	}
	storage := storage.New()

	processor := processor.New(storage)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				processor.Process(ctx, uuid.New())
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
		if err = shutdownFunc(ctx); err != nil {
			panic(err.Error())
		}
		cancel()
	}, make(chan os.Signal))
}
