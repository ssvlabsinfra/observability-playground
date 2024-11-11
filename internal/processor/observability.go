package processor

import (
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	observabilityComponentName      = "github.com/ssvlabsinfra/observability-playground/internal/processor"
	observabilityComponentNamespace = "observability_playground.processor"
)

var (
	tracer = otel.Tracer(observabilityComponentName)
	meter  = otel.Meter(observabilityComponentName)

	itemCounter metric.Int64Counter
)

func init() {
	var err error
	itemCounterMetricName := fmt.Sprintf("%s.items", observabilityComponentNamespace)
	itemCounter, err = meter.Int64Counter(
		itemCounterMetricName,
		metric.WithUnit("{item}"),
		metric.WithDescription("number of items processed"))
	if err != nil {
		slog.
			With("err", err).
			With("component", observabilityComponentName).
			With("metric", itemCounterMetricName).
			Error("error instantiating metric")
	}
}
