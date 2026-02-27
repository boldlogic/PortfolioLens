package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	httputils "github.com/boldlogic/PortfolioLens/pkg/http_utils"
	"github.com/boldlogic/PortfolioLens/quik-portfolio/internal/apperrors"
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
	GetML(ctx context.Context, date time.Time) ([]models.MoneyLimit, error)
	GetSL(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	GetSLOtc(ctx context.Context, date time.Time) ([]models.SecurityLimit, error)
	SaveSL(ctx context.Context, sec models.SecurityLimit) error
	SaveSLOtc(ctx context.Context, sec models.SecurityLimit) error
	GetLimits(ctx context.Context, date time.Time) ([]models.Limit, error)
	GetPortfolio(ctx context.Context) ([]models.PortfolioItem, error)
	SaveFirm(ctx context.Context, code string, name string) (models.Firm, error)

	GetTradePoints(ctx context.Context) ([]models.TradePoint, error)
	GetTradePointByID(ctx context.Context, id uint8) (models.TradePoint, error)
	GetBoards(ctx context.Context) ([]models.Board, error)
	GetBoardByID(ctx context.Context, id uint8) (models.Board, error)
}

type HandlerFunc func(r *http.Request) (any, string, error)

func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, detail, err := h(r)
		if err != nil {
			var resp httputils.HTTPErr
			switch {
			case errors.Is(err, apperrors.ErrValidation):
				resp = httputils.BadRequest(detail)
			case errors.Is(err, apperrors.ErrBusinessValidation):
				resp = httputils.UnprocessableEntity(detail)
			case errors.Is(apperrors.ErrNotFound, err) || errors.Is(models.ErrLimitsNotFound, err) || errors.Is(models.ErrPortfolioNotFound, err):
				resp = httputils.NotFound(detail)

			case errors.Is(apperrors.ErrConflict, err):
				resp = httputils.Conflict(err.Error())
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

		if r.Method == "POST" {
			httputils.WriteResp(w, http.StatusCreated, data)
		} else {
			httputils.WriteResp(w, http.StatusOK, data)
		}

	}
}
