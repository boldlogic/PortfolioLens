package v1

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	Mux     *chi.Mux
	Handler *Handler
	logger  *zap.Logger
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	router := chi.NewRouter()
	router.Post("/tasks", handler.Adapt(handler.CreateTask))
	return &Router{
		Mux:     router,
		Handler: handler,
		logger:  logger,
	}
}
