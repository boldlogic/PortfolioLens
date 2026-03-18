package workers

import (
	"context"
	"sync/atomic"

	"github.com/boldlogic/PortfolioLens/pkg/utils"
)

func withRollForwardDateCache(lastUpToDate *atomic.Int64, fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		today := utils.DateToYYYYMMDD(utils.Today())
		if last := lastUpToDate.Load(); last != 0 && today <= last {
			return nil
		}
		if err := fn(ctx); err != nil {
			return err
		}
		lastUpToDate.Store(today)
		return nil
	}
}
