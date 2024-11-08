package processor

import (
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	observabilityComponentName = "github.com/ssvlabsinfra/p2p-observability/internal/processor"
)

var (
	tracer = otel.Tracer(observabilityComponentName)
	meter  = otel.Meter(observabilityComponentName)

	processorItemCounter metric.Int64Counter
)

func init() {
	var err error
	processorItemCounter, err = meter.Int64Counter(
		"processor.item.total",
		metric.WithUnit("{item}"),
		metric.WithDescription("number of items processed"))
	if err != nil {
		slog.
			With("err", err).
			With("component", observabilityComponentName).
			With("metric", "processor.item.total").
			Error("error instantiating metric")
	}
}
