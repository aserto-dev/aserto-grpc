module github.com/aserto-dev/aserto-grpc

go 1.22

// replace github.com/aserto-dev/errors => ../errors

require (
	github.com/aserto-dev/errors v0.0.11
	github.com/aserto-dev/go-aserto v0.32.0
	github.com/aserto-dev/header v0.0.8
	github.com/aserto-dev/logger v0.0.5
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.19.1
	github.com/rs/zerolog v1.33.0
	github.com/stretchr/testify v1.9.0
	go.opencensus.io v0.24.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240812133136-8ffd90a71988
	google.golang.org/grpc v1.65.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240812133136-8ffd90a71988 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
