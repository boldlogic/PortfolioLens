package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type CommonHandler interface {
	Adapt(fn handler.HandlerFunc) http.HandlerFunc
}

type Handler struct {
	commonHandler CommonHandler
	service       Service
	logger        *zap.Logger
}

func NewHandler(commonHandler CommonHandler, svc Service, logger *zap.Logger) *Handler {
	return &Handler{
		commonHandler: commonHandler,
		service:       svc,
		logger:        logger,
	}
}
func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}

type Service interface {
	GetML(ctx context.Context, date time.Time) ([]models.MoneyLimit, error)
	GetSL(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	GetSLOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	SaveSL(ctx context.Context, sec models.SecurityLimit) error
	SaveSLOtc(ctx context.Context, sec models.SecurityLimit) error
	GetLimits(ctx context.Context, date time.Time) ([]models.Limit, error)
	GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error)
	SaveFirm(ctx context.Context, code string, name string) (quik.Firm, error)
}
