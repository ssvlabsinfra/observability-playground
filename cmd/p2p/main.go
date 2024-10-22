package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ssvlabsinfra/p2p-observability/internal/platform/lifecycle"
	"github.com/ssvlabsinfra/p2p-observability/internal/platform/observability"
)

func main() {
	ctx := context.Background()
	otelShutdown, err := observability.SetupOTelSDK(ctx)
	if err != nil {
		panic(err.Error())
	}

	slog.Info("OTeL SDK configured. Listening for application shutdown")

	lifecycle.ListenForApplicationShutDown(ctx, func() { _ = otelShutdown(ctx) }, make(chan os.Signal))
}
