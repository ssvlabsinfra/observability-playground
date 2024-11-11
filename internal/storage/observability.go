package storage

import (
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	observabilityComponentName      = "github.com/ssvlabsinfra/observability-playground/internal/storage"
	observabilityComponentNamespace = "observability_playground.storage"
)

var (
	tracer = otel.Tracer(observabilityComponentName)
	meter  = otel.Meter(observabilityComponentName)

	storageItemCounter metric.Int64Counter
)

func init() {
	var err error
	itemCounterMetricName := fmt.Sprintf("%s.items", observabilityComponentNamespace)
	storageItemCounter, err = meter.Int64Counter(
		itemCounterMetricName,
		metric.WithUnit("{item}"),
		metric.WithDescription("number of items stored"))
	if err != nil {
		slog.
			With("err", err).
			With("component", observabilityComponentName).
			With("metric", itemCounterMetricName).
			Error("error instantiating metric")
	}
}
