package grpcclient

import (
	"context"
	"crypto/tls"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type DialOptionsProvider func(*Config) ([]grpc.DialOption, error)

func NewDialOptionsProvider() DialOptionsProvider {
	return func(cfg *Config) ([]grpc.DialOption, error) {
		if (cfg.ClientCertPath != "") != (cfg.ClientKeyPath != "") {
			return nil, errors.New("both client cert and key must be specified, or both must be empty")
		}

		dialOpts := []grpc.DialOption{}

		if cfg.ClientCertPath != "" {
			certificate, err := tls.LoadX509KeyPair(cfg.ClientCertPath, cfg.ClientKeyPath)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to load client GRPC certs")
			}

			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{certificate},
				MinVersion:   tls.VersionTLS12,
			}

			dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		}

		var pairs []string
		for k, v := range cfg.Headers {
			pairs = append(pairs, k, v)
		}

		if pairs != nil {
			dialOpts = append(dialOpts,
				grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
					ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
					return invoker(ctx, method, req, reply, cc, opts...)
				}),
				grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
					ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
					return streamer(ctx, desc, cc, method, opts...)
				}),
			)
		}

		return dialOpts, nil
	}
}
