package service

import (
	"context"

	db "go-api-simple-template/internal/postgresql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	cb, err := s.q.CreateTeaBlend(ctx, db.CreateTeaBlendParams{
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

func (s *TeaBlendService) GetAll(ctx context.Context) ([]db.TeaBlend, error) {
	return s.q.ListTeaBlends(ctx)
}

func (s *TeaBlendService) GetByID(ctx context.Context, id uuid.UUID) (*db.TeaBlend, error) {
	cb, err := s.q.GetTeaBlend(ctx, id)
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

func (s *TeaBlendService) Update(ctx context.Context, id uuid.UUID, name string, description string) (*db.TeaBlend, error) {
	cb, err := s.q.UpdateTeaBlend(ctx, db.UpdateTeaBlendParams{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

func (s *TeaBlendService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteTeaBlend(ctx, id)
}
