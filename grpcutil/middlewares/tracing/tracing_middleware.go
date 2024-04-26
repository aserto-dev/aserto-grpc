package tracing

import (
	"context"
	"time"

	"github.com/aserto-dev/aserto-grpc/grpcutil"
	"github.com/aserto-dev/header"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type TracingMiddleware struct {
	logger *zerolog.Logger
}

type TracingHook struct{}

type TracingFields struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()

	e.Fields(header.KnownContextValueStrings(ctx))

	serviceMethod, ok := grpc.Method(ctx)
	if ok {
		e.Str("service", serviceMethod)
	}

	tracingFields := ctx.Value(TracingFields{})
	e.Fields(tracingFields)
}

func NewTracingMiddleware(logger *zerolog.Logger) *TracingMiddleware {
	return &TracingMiddleware{
		logger: logger,
	}
}

var _ grpcutil.Middleware = &TracingMiddleware{}

func (m *TracingMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		apiLogger := m.logger.Hook(TracingHook{})
		apiLoggerCtx := apiLogger.WithContext(ctx)

		apiLogger.Trace().Interface("request", req).Msg("grpc call start")

		start := time.Now()
		result, err := handler(apiLoggerCtx, req)
		apiLogger.Trace().Dur("duration-ms", time.Since(start)).Msg("grpc call end")

		return result, err
	}
}

func (m *TracingMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		apiLogger := m.logger.Hook(TracingHook{})
		apiLoggerCtx := apiLogger.WithContext(ctx)

		apiLogger.Trace().Msg("grpc stream call")

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = apiLoggerCtx

		return handler(srv, wrapped)
	}
}
