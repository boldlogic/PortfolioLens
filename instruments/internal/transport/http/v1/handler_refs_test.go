package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"github.com/boldlogic/PortfolioLens/pkg/transport/httpserver/handler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type refsSvcNoop struct{}

func (refsSvcNoop) GetTradePoints(ctx context.Context) ([]md.TradePoint, error) {
	return nil, nil
}

func (refsSvcNoop) GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error) {
	return md.TradePoint{}, nil
}

func (refsSvcNoop) GetBoards(ctx context.Context) ([]quik.Board, error) {
	return nil, nil
}

func (refsSvcNoop) GetBoardByID(ctx context.Context, id uint8) (quik.Board, error) {
	return quik.Board{}, nil
}

func reqWithChiID(method, path, id string) *http.Request {
	r := httptest.NewRequest(method, path, nil)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", id)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rc))
}

func TestGetBoard_некорректный_id_ErrValidation(t *testing.T) {
	h := NewHandler(handler.NewHandler(), refsSvcNoop{}, zap.NewNop())
	_, _, err := h.GetBoard(reqWithChiID("GET", "/boards/x", "нечисло"))
	if !errors.Is(err, md.ErrValidation) {
		t.Fatalf("ожидали ErrValidation: %v", err)
	}
}
