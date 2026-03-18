package v1

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	Mux    *chi.Mux
	logger *zap.Logger
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	r := chi.NewRouter()
	r.Get("/tradepoints", handler.Adapt(handler.GetTradePoints))
	r.Get("/tradepoints/{id}", handler.Adapt(handler.GetTradePoint))
	r.Get("/boards", handler.Adapt(handler.GetBoards))
	r.Get("/boards/{id}", handler.Adapt(handler.GetBoard))
	return &Router{
		Mux:    r,
		logger: logger,
	}
}
