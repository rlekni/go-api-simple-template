package service

import (
	"context"
	"log/slog"

	"go-api-simple-template/internal/database"
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

type TeaBlendWithLocation struct {
	TeaBlend db.TeaBlend
	Location db.Location
}

func (s *TeaBlendService) Create(ctx context.Context, name string, description string, locationName string, quantity int) (*TeaBlendWithLocation, error) {
	ctx, span := tracer.Start(ctx, "TeaBlendService.Create")
	defer span.End()

	slog.InfoContext(ctx, "creating tea blend with location", "name", name, "location", locationName)

	var result TeaBlendWithLocation

	err := database.RunInTx(ctx, s.pool, func(q *db.Queries) error {
		cb, err := q.CreateTeaBlend(ctx, db.CreateTeaBlendParams{
			Name:        name,
			Description: pgtype.Text{String: description, Valid: description != ""},
		})
		if err != nil {
			return err
		}

		loc, err := q.CreateLocation(ctx, db.CreateLocationParams{
			TeaBlendID: cb.ID,
			Name:       locationName,
			Quantity:   int32(quantity),
		})
		if err != nil {
			return err
		}

		result.TeaBlend = cb
		result.Location = loc
		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "failed to create tea blend and location", "error", err)
		return nil, err
	}

	return &result, nil
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

func (s *TeaBlendService) Update(ctx context.Context, id uuid.UUID, name string, description string, quantity int) (*TeaBlendWithLocation, error) {
	ctx, span := tracer.Start(ctx, "TeaBlendService.Update")
	defer span.End()

	slog.InfoContext(ctx, "updating tea blend and quantity", "id", id, "name", name, "quantity", quantity)

	var result TeaBlendWithLocation

	err := database.RunInTx(ctx, s.pool, func(q *db.Queries) error {
		cb, err := q.UpdateTeaBlend(ctx, db.UpdateTeaBlendParams{
			ID:          id,
			Name:        name,
			Description: pgtype.Text{String: description, Valid: description != ""},
		})
		if err != nil {
			return err
		}

		loc, err := q.UpdateLocationQuantity(ctx, db.UpdateLocationQuantityParams{
			TeaBlendID: id,
			Quantity:   int32(quantity),
		})
		if err != nil {
			// If location doesn't exist for some reason, we might want to create it
			// but for this example, we'll just return the error
			return err
		}

		result.TeaBlend = cb
		result.Location = loc
		return nil
	})

	if err != nil {
		slog.ErrorContext(ctx, "failed to update tea blend and quantity", "id", id, "error", err)
		return nil, err
	}

	return &result, nil
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
