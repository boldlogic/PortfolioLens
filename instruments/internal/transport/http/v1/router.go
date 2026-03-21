package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	mux    *chi.Mux
	logger *zap.Logger
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	r := chi.NewRouter()
	r.Get("/tradepoints", handler.Adapt(handler.GetTradePoints))
	r.Get("/tradepoints/{id}", handler.Adapt(handler.GetTradePoint))
	r.Get("/boards", handler.Adapt(handler.GetBoards))
	r.Get("/boards/{id}", handler.Adapt(handler.GetBoard))
	return &Router{
		mux:    r,
		logger: logger,
	}
}
