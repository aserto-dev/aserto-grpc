package grpcutil

import (
	"context"

	"github.com/aserto-dev/header"
	"github.com/rs/zerolog"
)

// CompleteLogger returns a logger that contains the
func CompleteLogger(ctx context.Context, log *zerolog.Logger) *zerolog.Logger {
	values := header.KnownContextValueStrings(ctx)
	completeLogger := log.With().Fields(values).Logger()
	return &completeLogger
}
