package v1

import (
	"context"
	"errors"
	"net/http"

	httputils "github.com/boldlogic/PortfolioLens/pkg/http_utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/models"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(svc Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: svc,
		logger:  logger,
	}
}

type Service interface {
	GetML(ctx context.Context) ([]models.MoneyLimit, error)
	GetSL(ctx context.Context) ([]models.SecurityLimit, error)
	GetLimits(ctx context.Context) ([]models.Limit, error)
	GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error)
}

type HandlerFunc func(r *http.Request) (any, error)

func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h(r)
		if err != nil {
			var resp httputils.HTTPErr
			switch {
			case errors.Is(models.ErrSLNotFound, err) || errors.Is(models.ErrMLNotFound, err) || errors.Is(models.ErrLimitsNotFound, err) || errors.Is(models.ErrPortfolioNotFound, err):
				resp = httputils.NotFound(err.Error())
			default:
				resp = httputils.Internal(err.Error())
			}
			httputils.WriteResp(w, resp.Status, resp)

			return
		}
		if data == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		httputils.WriteResp(w, http.StatusOK, data)

	}
}
