package v1

import (
	"context"
	"net/http"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"go.uber.org/zap"
)

type RefsService interface {
	GetTradePoints(ctx context.Context) ([]md.TradePoint, error)
	GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error)
	GetBoards(ctx context.Context) ([]quik.Board, error)
	GetBoardByID(ctx context.Context, id uint8) (quik.Board, error)
}

type Handler struct {
	commonHandler handler.Adapter
	refsSvc       RefsService
	logger        *zap.Logger
}

func NewHandler(commonHandler handler.Adapter, refsSvc RefsService, log *zap.Logger) *Handler {
	return &Handler{
		logger:        log,
		refsSvc:       refsSvc,
		commonHandler: commonHandler,
	}
}

func (h *Handler) Adapt(fn handler.HandlerFunc) http.HandlerFunc {
	return h.commonHandler.Adapt(fn)
}
