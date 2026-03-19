package service

import (
	"context"
	"log/slog"

	db "go-api-simple-template/internal/postgresql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("tea-blends-service")

type TeaBlendService struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewTeaBlendService(pool *pgxpool.Pool) *TeaBlendService {
	return &TeaBlendService{
		pool: pool,
		q:    db.New(pool),
	}
}

func (s *TeaBlendService) Create(ctx context.Context, name string, description string) (*db.TeaBlend, error) {
	ctx, span := tracer.Start(ctx, "TeaBlendService.Create")
	defer span.End()

	slog.InfoContext(ctx, "creating tea blend", "name", name)

	cb, err := s.q.CreateTeaBlend(ctx, db.CreateTeaBlendParams{
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create tea blend", "error", err)
		return nil, err
	}
	return &cb, nil
}

func (s *TeaBlendService) GetAll(ctx context.Context) ([]db.TeaBlend, error) {
	ctx, span := tracer.Start(ctx, "TeaBlendService.GetAll")
	defer span.End()

	slog.DebugContext(ctx, "fetching all tea blends")

	return s.q.ListTeaBlends(ctx)
}

func (s *TeaBlendService) GetByID(ctx context.Context, id uuid.UUID) (*db.TeaBlend, error) {
	ctx, span := tracer.Start(ctx, "TeaBlendService.GetByID")
	defer span.End()

	slog.DebugContext(ctx, "fetching tea blend by ID", "id", id)

	cb, err := s.q.GetTeaBlend(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch tea blend", "id", id, "error", err)
		return nil, err
	}
	return &cb, nil
}

func (s *TeaBlendService) Update(ctx context.Context, id uuid.UUID, name string, description string) (*db.TeaBlend, error) {
	ctx, span := tracer.Start(ctx, "TeaBlendService.Update")
	defer span.End()

	slog.InfoContext(ctx, "updating tea blend", "id", id, "name", name)

	cb, err := s.q.UpdateTeaBlend(ctx, db.UpdateTeaBlendParams{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to update tea blend", "id", id, "error", err)
		return nil, err
	}
	return &cb, nil
}

func (s *TeaBlendService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "TeaBlendService.Delete")
	defer span.End()

	slog.InfoContext(ctx, "deleting tea blend", "id", id)

	err := s.q.DeleteTeaBlend(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete tea blend", "id", id, "error", err)
		return err
	}
	return nil
}
