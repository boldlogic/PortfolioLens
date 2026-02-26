package v1

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Router struct {
	Mux     *chi.Mux
	Handler *Handler
	logger  *zap.Logger
	//config  *config.Config
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	router := chi.NewRouter()
	router.Route("/quik", func(r chi.Router) {
		r.Route("/limits", func(r chi.Router) {
			r.Get("/", Adapt(handler.GetLimits))
			r.Get("/money", Adapt(handler.GetMoneyLimits))
			r.Get("/securities", Adapt(handler.GetSecurityLimits))
			r.Post("/securities", Adapt(handler.AddSecurityLimit))
			r.Get("/securities/otc", Adapt(handler.GetSecurityLimitsOtc))
			r.Post("/securities/otc", Adapt(handler.AddSecurityLimitOtc))
		})
		r.Get("/portfolio", Adapt(handler.GetPortfolio))
		r.Post("/firms", Adapt(handler.AddFirm))
		r.Get("/tradepoints", Adapt(handler.GetTradePoints))
		r.Get("/boards", Adapt(handler.GetBoards))
		r.Get("/boards/{id}", Adapt(handler.GetBoard))
	})
	return &Router{
		Mux:     router,
		Handler: handler,
		logger:  logger,
	}
}
