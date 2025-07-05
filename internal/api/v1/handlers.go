package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ApiHandlers struct {
	ExampleUseCase ExampleUseCase
}

func (h *ApiHandlers) Routes(r chi.Router) {
	// TODO: Add routes
	r.Get("/health", h.Health)
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/example", func(r chi.Router) {
			r.Post("/", h.CreateExample)
			r.Get("/{id}", h.GetExampleByID)
		})
	})
}

func (h *ApiHandlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type ErrorResponseBody struct {
	Error string `json:"error"`
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	render.Status(r, code)
	render.JSON(w, r, ErrorResponseBody{
		Error: err.Error(),
	})
}

func unknownErrorResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.PlainText(w, r, http.StatusText(http.StatusInternalServerError))
}
