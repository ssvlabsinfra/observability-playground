package work

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer(name)

type (
	metrics interface {
		Increment(ctx context.Context)
	}
	Work struct {
		metrics metrics
	}
)

func New(metrics metrics) *Work {
	return &Work{
		metrics: metrics,
	}
}

func (w *Work) Do(ctx context.Context, uuid uuid.UUID) {
	ctx, span := tracer.Start(ctx, "work", trace.WithAttributes(attribute.String("id", uuid.String())))
	defer span.End()

	w.metrics.Increment(ctx)
	span.AddEvent("work called", trace.WithAttributes(attribute.String("id", uuid.String())))
	slog.
		With("id", uuid.String()).
		Info("work finished")

	span.AddEvent("calling subwork")
	if err := w.subWork(ctx, uuid); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	span.SetStatus(codes.Ok, "work finished")
}

func (Work) subWork(ctx context.Context, uuid uuid.UUID) error {
	_, span := tracer.Start(ctx, "sub_work", trace.WithAttributes(attribute.String("id", uuid.String())))
	defer span.End()

	err := errors.New("work error")
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return err
}
