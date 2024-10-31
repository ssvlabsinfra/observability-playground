package work

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const name = "github.com/ssvlabsinfra/p2p-observability/internal/work"

var (
	tracer  = otel.Tracer(name)
	meter   = otel.Meter(name, metric.WithInstrumentationVersion("0.0.1"))
	counter metric.Int64Counter
)

func init() {
	var err error
	counter, err = meter.Int64Counter("p2p-observability", metric.WithDescription("p2p observability counter description"))
	if err != nil {
		panic(err)
	}
}

type Work struct {
}

func New() *Work {
	return &Work{}
}

func (Work) Do(ctx context.Context, uuid uuid.UUID) {
	ctx, span := tracer.Start(ctx, "work", trace.WithAttributes(attribute.String("id", uuid.String())))
	defer span.End()

	counter.Add(ctx, 1, metric.WithAttributes(attribute.String("attr", "custom_attribute")))
	span.AddEvent("work called", trace.WithAttributes(attribute.String("id", uuid.String())))
	slog.
		With("id", uuid.String()).
		Info("work finished")

	span.SetStatus(codes.Ok, "work finished")
}
