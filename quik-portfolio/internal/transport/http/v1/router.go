package v1

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Router struct {
	Mux    *chi.Mux
	logger *zap.Logger
	//config  *config.Config
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	router := chi.NewRouter()
	router.Route("/quik", func(r chi.Router) {
		r.Route("/limits", func(r chi.Router) {
			r.Get("/", handler.Adapt(handler.GetLimits))
			r.Get("/money", handler.Adapt(handler.GetMoneyLimits))
			r.Get("/securities", handler.Adapt(handler.GetSecurityLimits))
			r.Post("/securities", handler.Adapt(handler.AddSecurityLimit))
			r.Get("/securities/otc", handler.Adapt(handler.GetSecurityLimitsOtc))
			r.Post("/securities/otc", handler.Adapt(handler.AddSecurityLimitOtc))
		})
		r.Get("/portfolio", handler.Adapt(handler.GetPortfolio))
		r.Post("/firms", handler.Adapt(handler.AddFirm))
	})
	return &Router{
		Mux:    router,
		logger: logger,
	}
}
