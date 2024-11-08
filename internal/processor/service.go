package processor

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	storage interface {
		Store(context.Context) error
	}
	Processor struct {
		storage storage
	}
)

func New(storage storage) *Processor {
	return &Processor{
		storage: storage,
	}
}

func (p *Processor) Process(ctx context.Context, uuid uuid.UUID) {
	ctx, span := tracer.Start(ctx, "process",
		trace.WithAttributes(attribute.String("id", uuid.String())))
	defer span.End()
	span.AddEvent("process called")

	span.AddEvent("calling storage")
	if err := p.storage.Store(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	processorItemCounter.Add(ctx, 1)

	span.SetStatus(codes.Ok, "processor finished")
}
