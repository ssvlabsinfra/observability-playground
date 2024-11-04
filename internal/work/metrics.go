package work

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "github.com/ssvlabsinfra/p2p-observability/internal/work"

var (
	meter = otel.Meter(name, metric.WithInstrumentationVersion("0.0.1"))
)

type Metrics struct {
	counter metric.Int64Counter
}

func NewMetrics() *Metrics {
	counter, err := meter.Int64Counter("p2p-observability",
		metric.WithDescription("p2p observability counter description"))
	if err != nil {
		panic(err)
	}

	return &Metrics{
		counter: counter,
	}
}

func (m *Metrics) Increment(ctx context.Context) {
	m.counter.Add(ctx, 1, metric.WithAttributes(attribute.String("attr", "custom_attribute")))
}
