package shutdown

import (
	"context"
	"errors"
)

func IsExceeded(err error) bool {
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return false
}
