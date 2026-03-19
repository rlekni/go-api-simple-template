package handler

import (
	"context"
	"log/slog"
	"time"

	"go-api-simple-template/internal/service"

	"github.com/google/uuid"
)

type LocationResponse struct {
	ID        uuid.UUID `json:"id" doc:"Unique identifier of the location"`
	Name      string    `json:"name" doc:"Name of the location"`
	Quantity  int       `json:"quantity" doc:"Quantity at this location"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Timestamp when the location was last updated"`
}

type TeaBlendResponse struct {
	ID          uuid.UUID         `json:"id" doc:"Unique identifier of the tea blend"`
	Name        string            `json:"name" doc:"Name of the tea blend"`
	Description string            `json:"description" doc:"Description of the tea blend"`
	Location    *LocationResponse `json:"location,omitempty" doc:"Location information"`
	CreatedAt   time.Time         `json:"createdAt" doc:"Timestamp when the tea blend was created"`
	UpdatedAt   time.Time         `json:"updatedAt" doc:"Timestamp when the tea blend was last updated"`
}

type CreateTeaBlendRequest struct {
	Body struct {
		Name         string `json:"name" doc:"Name of the tea blend" minLength:"1" maxLength:"500" example:"Earl Grey"`
		Description  string `json:"description" doc:"Description of the tea blend" maxLength:"500" example:"A classic black tea with bergamot"`
		LocationName string `json:"locationName" doc:"Initial location name" minLength:"1" maxLength:"500" example:"Warehouse A"`
		Quantity     int    `json:"quantity" doc:"Initial quantity" minimum:"0" example:"100"`
	}
}

type CreateTeaBlendResponse struct {
	Body TeaBlendResponse
}

type GetTeaBlendsResponse struct {
	Body []TeaBlendResponse
}

type GetTeaBlendRequest struct {
	ID string `path:"id" doc:"Tea blend ID"`
}

type GetTeaBlendResponse struct {
	Body TeaBlendResponse
}

type UpdateTeaBlendRequest struct {
	ID   string `path:"id" doc:"Tea blend ID"`
	Body struct {
		Name        string `json:"name" doc:"Name of the tea blend" minLength:"1" maxLength:"500" example:"Earl Grey"`
		Description string `json:"description" doc:"Description of the tea blend" maxLength:"500" example:"A classic black tea with bergamot"`
		Quantity    int    `json:"quantity" doc:"Update quantity at the location" minimum:"0" example:"150"`
	}
}

type UpdateTeaBlendResponse struct {
	Body TeaBlendResponse
}

type DeleteTeaBlendRequest struct {
	ID string `path:"id" doc:"Tea blend ID"`
}

type TeaBlendHandler struct {
	service *service.TeaBlendService
}

func NewTeaBlendHandler(service *service.TeaBlendService) *TeaBlendHandler {
	return &TeaBlendHandler{
		service: service,
	}
}

func (h *TeaBlendHandler) Create(ctx context.Context, input *CreateTeaBlendRequest) (*CreateTeaBlendResponse, error) {
	res, err := h.service.Create(ctx, input.Body.Name, input.Body.Description, input.Body.LocationName, input.Body.Quantity)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create tea blend", "error", err)
		return nil, err
	}

	resp := &CreateTeaBlendResponse{}
	resp.Body = TeaBlendResponse{
		ID:          res.TeaBlend.ID,
		Name:        res.TeaBlend.Name,
		Description: res.TeaBlend.Description.String,
		Location: &LocationResponse{
			ID:        res.Location.ID,
			Name:      res.Location.Name,
			Quantity:  int(res.Location.Quantity),
			UpdatedAt: res.Location.UpdatedAt.Time,
		},
		CreatedAt: res.TeaBlend.CreatedAt.Time,
		UpdatedAt: res.TeaBlend.UpdatedAt.Time,
	}

	return resp, nil
}

func (h *TeaBlendHandler) GetAll(ctx context.Context, input *struct{}) (*GetTeaBlendsResponse, error) {
	blends, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]TeaBlendResponse, 0, len(blends))
	for _, b := range blends {
		items = append(items, TeaBlendResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description.String,
			CreatedAt:   b.CreatedAt.Time,
			UpdatedAt:   b.UpdatedAt.Time,
		})
	}

	resp := &GetTeaBlendsResponse{}
	resp.Body = items
	return resp, nil
}

func (h *TeaBlendHandler) GetByID(ctx context.Context, input *GetTeaBlendRequest) (*GetTeaBlendResponse, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	cb, err := h.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &GetTeaBlendResponse{}
	resp.Body = TeaBlendResponse{
		ID:          cb.ID,
		Name:        cb.Name,
		Description: cb.Description.String,
		CreatedAt:   cb.CreatedAt.Time,
		UpdatedAt:   cb.UpdatedAt.Time,
	}
	return resp, nil
}

func (h *TeaBlendHandler) Update(ctx context.Context, input *UpdateTeaBlendRequest) (*UpdateTeaBlendResponse, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	res, err := h.service.Update(ctx, id, input.Body.Name, input.Body.Description, input.Body.Quantity)
	if err != nil {
		return nil, err
	}

	resp := &UpdateTeaBlendResponse{}
	resp.Body = TeaBlendResponse{
		ID:          res.TeaBlend.ID,
		Name:        res.TeaBlend.Name,
		Description: res.TeaBlend.Description.String,
		Location: &LocationResponse{
			ID:        res.Location.ID,
			Name:      res.Location.Name,
			Quantity:  int(res.Location.Quantity),
			UpdatedAt: res.Location.UpdatedAt.Time,
		},
		CreatedAt: res.TeaBlend.CreatedAt.Time,
		UpdatedAt: res.TeaBlend.UpdatedAt.Time,
	}
	return resp, nil
}

func (h *TeaBlendHandler) Delete(ctx context.Context, input *DeleteTeaBlendRequest) (*struct{}, error) {
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, err
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
