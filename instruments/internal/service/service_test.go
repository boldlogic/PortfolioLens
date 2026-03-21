package service

import (
	"context"
	"testing"

	"github.com/boldlogic/PortfolioLens/instruments/internal/models"
	md "github.com/boldlogic/PortfolioLens/pkg/models"
	"github.com/boldlogic/PortfolioLens/pkg/models/quik"
	"go.uber.org/zap"
)

type instrRepoStub struct {
	selInstr   models.Instrument
	selBoard   models.InstrumentBoard
	selClass   string
	getIDErr   error
	insID      int
	setCalls   []int
	mergeBoard models.InstrumentBoard
}

func (s *instrRepoStub) SelectInstrumentFromNewCurrentQuote(ctx context.Context) (models.Instrument, models.InstrumentBoard, string, error) {
	return s.selInstr, s.selBoard, s.selClass, nil
}

func (s *instrRepoStub) InsInstrument(ctx context.Context, i models.Instrument) (int, error) {
	return s.insID, nil
}

func (s *instrRepoStub) SetInstrument(ctx context.Context, id int, ic string) error {
	s.setCalls = append(s.setCalls, id)
	return nil
}

func (s *instrRepoStub) GetInstrumentId(ctx context.Context, ticker string, tradePointId uint8) (int, error) {
	return 0, s.getIDErr
}

func (s *instrRepoStub) MergeInstrumentBoard(ctx context.Context, ib models.InstrumentBoard) error {
	s.mergeBoard = ib
	return nil
}

type refsSyncNoop struct{}

func (refsSyncNoop) SyncInstrumentTypesFromQuotes(context.Context) error   { return nil }
func (refsSyncNoop) SyncInstrumentSubTypesFromQuotes(context.Context) error { return nil }
func (refsSyncNoop) SyncBoardsFromQuotes(context.Context) error             { return nil }
func (refsSyncNoop) TagBoardsTradePointId(context.Context) error           { return nil }

type refsQueryStub struct{}

func (refsQueryStub) GetTradePoints(ctx context.Context) ([]md.TradePoint, error) {
	return nil, nil
}

func (refsQueryStub) GetTradePointByID(ctx context.Context, id uint8) (md.TradePoint, error) {
	return md.TradePoint{}, nil
}

func (refsQueryStub) GetBoardsWithTradePoint(ctx context.Context) ([]quik.Board, error) {
	return nil, nil
}

func (refsQueryStub) GetBoardByIDWithTradePoint(ctx context.Context, id uint8) (quik.Board, error) {
	return quik.Board{}, nil
}

func TestSaveInstrument_новый_ErrNotFound_цепочка_id(t *testing.T) {
	instr := &instrRepoStub{
		selInstr: models.Instrument{Ticker: "LKOH", TradePointId: 1},
		selBoard: models.InstrumentBoard{BoardId: 7},
		selClass: "TQBR",
		getIDErr: md.ErrNotFound,
		insID:    100,
	}
	svc := NewService(instr, refsSyncNoop{}, refsQueryStub{}, zap.NewNop())
	if err := svc.SaveInstrument(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(instr.setCalls) != 1 || instr.setCalls[0] != 100 {
		t.Fatalf("SetInstrument: ожидали один вызов с id=100, получили %v", instr.setCalls)
	}
	if instr.mergeBoard.InstrumentId != 100 || instr.mergeBoard.BoardId != 7 {
		t.Fatalf("MergeInstrumentBoard: %+v", instr.mergeBoard)
	}
}
