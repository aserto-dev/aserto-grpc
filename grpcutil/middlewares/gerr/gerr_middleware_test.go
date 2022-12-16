package gerr

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/aserto-dev/aserto-grpc/grpcutil/middlewares/test"
	aerr "github.com/aserto-dev/errors"
	"github.com/aserto-dev/logger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	testError = aerr.NewAsertoError("X00000", codes.Internal, http.StatusInternalServerError, "a test error")
)

func TestUnaryServerWithWrappedError(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler("output", errors.Wrap(testError, "unimportant error"))

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)
	_, err := NewErrorMiddleware().Unary()(ctx, "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)
	assert.Contains(err.Error(), "a test error")
}

func TestUnaryServerWithFields(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler(
		"output",
		errors.Wrap(testError.Str("my-field", "deadbeef"), "another error"),
	)

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	logOutput := buf.String()

	assert.Contains(logOutput, "deadbeef")
}

func TestUnaryServerWithDoubleCerr(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler(
		"output",
		aerr.ErrUnknown.Err(testError.Str("my-field", "deadbeef").Msg("old message")).Msg("new message"),
	)

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	logOutput := buf.String()

	assert.Contains(logOutput, "new message")
	assert.Contains(logOutput, "deadbeef")
}

func TestSimpleInnerError(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler("output", aerr.ErrUnknown.Err(errors.New("deadbeef")).Msg("failed to setup initial tag"))

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	logOutput := buf.String()

	assert.Contains(logOutput, "deadbeef")
}

func TestDirectResult(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler(
		"output",
		aerr.ErrUnknown.Err(testError).Msg("failed to setup initial tag"),
	)

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	s := status.Convert(err)

	errDetailsFound := false
	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			errDetailsFound = true
			assert.Contains(t.Metadata, "msg")
			assert.Contains(t.Metadata["msg"], "failed to setup")
		}
	}

	assert.True(errDetailsFound)
	assert.Contains(s.Message(), "an unknown error has occurred")
	assert.Contains(err.Error(), "failed to setup initial tag")
}
