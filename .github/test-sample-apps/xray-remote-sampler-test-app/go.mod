module github.com/aws-observability/aws-otel-go/.github/test-sample-apps/xray-remote-sampler-test-app

go 1.23.0

toolchain go1.23.2

replace github.com/aws-observability/aws-otel-go/samplers/aws/xray => ../../../samplers/aws/xray

require (
	github.com/aws-observability/aws-otel-go/samplers/aws/xray v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/contrib/propagators/aws v1.23.0
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/sdk v1.35.0
	go.opentelemetry.io/otel/trace v1.35.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
