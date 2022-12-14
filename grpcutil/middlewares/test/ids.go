package test

import (
	"context"
	"testing"

	"github.com/aserto-dev/header"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func RequestIDContext(t *testing.T) context.Context {
	assert := require.New(t)
	id, err := uuid.NewUUID()
	assert.NoError(err)
	return context.WithValue(context.Background(), header.HeaderAsertoRequestID, id)
}
