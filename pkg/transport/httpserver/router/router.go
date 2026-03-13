package router

import (
	"fmt"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"github.com/go-chi/chi/v5"
)

type Handlers interface {
	Adapt(fn handler.HandlerFunc) http.HandlerFunc
}

type Router struct {
	Mux      *chi.Mux
	Handlers Handlers
}

func NewRouter(handler Handlers) *Router {
	r := chi.NewRouter()
	return &Router{
		Mux:      r,
		Handlers: handler,
	}
}

func (r *Router) AddGetMethod(pattern string, fn handler.HandlerFunc) {
	r.Mux.Get(fmt.Sprintf("/%s", pattern), r.Handlers.Adapt(fn))
}

func (r *Router) AddPostMethod(pattern string, fn handler.HandlerFunc) {
	r.Mux.Post(fmt.Sprintf("/%s", pattern), r.Handlers.Adapt(fn))
}
