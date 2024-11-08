package storage

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/codes"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Store(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "store")
	defer span.End()

	span.AddEvent("store called")

	storageItemCounter.Add(ctx, 1)

	slog.Info("item stored")
	span.SetStatus(codes.Ok, "store finished")
	return nil
}
