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
	r.Get("/currencies", handler.Adapt(handler.GetCurrencies))
	r.Get("/currencies/{code}", handler.Adapt(handler.GetCurrency))
	return &Router{
		Mux:    r,
		logger: logger,
	}
}
