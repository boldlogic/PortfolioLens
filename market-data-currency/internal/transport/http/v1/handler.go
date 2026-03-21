package v1

import (
	"context"
	"net/http"

	"github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type CurrencyService interface {
	GetCurrencies(ctx context.Context) ([]models.Currency, error)
	GetCurrency(ctx context.Context, charCode string) (models.Currency, string, error)
}

type Handler struct {
	commonHandler handler.Adapter
	currencySvc   CurrencyService
	logger        *zap.Logger
}

func NewHandler(commonHandler handler.Adapter, currencySvc CurrencyService, log *zap.Logger) *Handler {
	return &Handler{
		logger:        log,
		currencySvc:   currencySvc,
		commonHandler: commonHandler,
	}
}

func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}
