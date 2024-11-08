package storage

import (
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	observabilityComponentName = "github.com/ssvlabsinfra/p2p-observability/internal/storage"
)

var (
	tracer = otel.Tracer(observabilityComponentName)
	meter  = otel.Meter(observabilityComponentName)

	storageItemCounter metric.Int64Counter
)

func init() {
	var err error
	storageItemCounter, err = meter.Int64Counter(
		"store.item.total",
		metric.WithUnit("{item}"),
		metric.WithDescription("number of items stored"))
	if err != nil {
		slog.
			With("err", err).
			With("component", observabilityComponentName).
			With("metric", "store.item.total").
			Error("error instantiating metric")
	}
}
