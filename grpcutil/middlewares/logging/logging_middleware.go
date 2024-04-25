package logging

import (
	"context"

	"github.com/aserto-dev/aserto-grpc/grpcutil"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type LoggingMiddleware struct {
	logger *zerolog.Logger
}

type TracingHook struct{}

type TracingFields struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()

	serviceMethod, ok := grpc.Method(ctx)
	if ok {
		e.Str("service", serviceMethod)
	}

	tracingFields := ctx.Value(TracingFields{})
	e.Fields(tracingFields)

}

func NewLoggingMiddleware(logger *zerolog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

var _ grpcutil.Middleware = &LoggingMiddleware{}

func (m *LoggingMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		apiLogger := m.logger.With().
			Logger().Hook(TracingHook{})
		apiLoggerCtx := apiLogger.WithContext(ctx)

		result, err := handler(apiLoggerCtx, req)

		return result, err
	}
}

func (m *LoggingMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		apiLogger := m.logger.With().
			Logger().Hook(TracingHook{})
		apiLoggerCtx := apiLogger.WithContext(ctx)

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = apiLoggerCtx

		return handler(srv, wrapped)
	}
}
