package server

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type MessageResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (s *Server) RegisterRoutes() {
	s.Router.Get("/debug/raw-proto", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Protocol: %s", r.Proto)))
	})

	s.Router.Get("/debug/free-memory", FreeMemoryHandler)

	// Tea Blends
	huma.Register(s.API, huma.Operation{
		OperationID: "get-tea-blends",
		Method:      http.MethodGet,
		Path:        "/tea-blends",
		Summary:     "Get all tea blends",
		Tags:        []string{"Tea Blends"},
	}, s.teaBlendHandler.GetAll)

	huma.Register(s.API, huma.Operation{
		OperationID: "create-tea-blend",
		Method:      http.MethodPost,
		Path:        "/tea-blends",
		Summary:     "Create a tea blend",
		Tags:        []string{"Tea Blends"},
		Errors:      []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, s.teaBlendHandler.Create)

	huma.Register(s.API, huma.Operation{
		OperationID: "get-tea-blend",
		Method:      http.MethodGet,
		Path:        "/tea-blends/{id}",
		Summary:     "Get a tea blend by ID",
		Tags:        []string{"Tea Blends"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, s.teaBlendHandler.GetByID)

	huma.Register(s.API, huma.Operation{
		OperationID: "update-tea-blend",
		Method:      http.MethodPut,
		Path:        "/tea-blends/{id}",
		Summary:     "Update a tea blend by ID",
		Tags:        []string{"Tea Blends"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, s.teaBlendHandler.Update)

	huma.Register(s.API, huma.Operation{
		OperationID: "delete-tea-blend",
		Method:      http.MethodDelete,
		Path:        "/tea-blends/{id}",
		Summary:     "Delete a tea blend by ID",
		Tags:        []string{"Tea Blends"},
		Errors:      []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, s.teaBlendHandler.Delete)
}
