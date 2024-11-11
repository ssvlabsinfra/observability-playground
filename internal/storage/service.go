package storage

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/codes"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Store(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%s.store", observabilityComponentNamespace))
	defer span.End()

	span.AddEvent("store called")

	storageItemCounter.Add(ctx, 1)

	slog.Info("item stored")
	span.SetStatus(codes.Ok, "store finished")
	return nil
}
