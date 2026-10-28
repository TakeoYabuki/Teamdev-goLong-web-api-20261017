package handler

import (
	"context"
	"errors"
	"net"
)

func isUpstreamTimeout(err error) bool {
	var netErr net.Error
	return errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout())
}
