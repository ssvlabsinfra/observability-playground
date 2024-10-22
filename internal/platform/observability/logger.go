package observability

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

const loggerName = "p2p-observability"

func init() {
	logger := otelslog.NewLogger(loggerName)
	slog.SetDefault(logger)
}
