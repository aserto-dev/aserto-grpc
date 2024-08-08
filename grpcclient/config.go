package grpcclient

import (
	"github.com/aserto-dev/aserto-grpc/grpcutil/middlewares/request"
	client "github.com/aserto-dev/go-aserto"
	"github.com/pkg/errors"
)

type Config struct {
	Address          string            `json:"address"`
	CACertPath       string            `json:"ca_cert_path"`
	ClientCertPath   string            `json:"client_cert_path"`
	ClientKeyPath    string            `json:"client_key_path"`
	APIKey           string            `json:"api_key"`
	Insecure         bool              `json:"insecure"`
	TimeoutInSeconds int               `json:"timeout_in_seconds"`
	Token            string            `json:"token"`
	Headers          map[string]string `json:"headers"`
}

func (cfg *Config) ToClientOptions(dop DialOptionsProvider) ([]client.ConnectionOption, error) {
	middleware := request.NewRequestIDMiddleware()
	options := []client.ConnectionOption{
		client.WithChainUnaryInterceptor(middleware.UnaryClient()),
		client.WithChainStreamInterceptor(middleware.StreamClient()),
		client.WithInsecure(cfg.Insecure),
	}

	if cfg.APIKey != "" && cfg.Token != "" {
		return nil, errors.New("both api_key and token are set")
	}
	if cfg.Token != "" {
		options = append(options, client.WithTokenAuth(cfg.Token))
	}

	if cfg.APIKey != "" {
		options = append(options, client.WithAPIKeyAuth(cfg.APIKey))
	}

	if cfg.Address != "" {
		options = append(options, client.WithAddr(cfg.Address))
	}

	if cfg.CACertPath != "" {
		options = append(options, client.WithCACertPath(cfg.CACertPath))
	}

	opts, err := dop(cfg)
	if err != nil {
		return nil, err
	}

	options = append(options, client.WithDialOptions(opts...))
	return options, nil
}
