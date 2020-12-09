module github.com/aws-observability/aws-otel-go/benchmark

go 1.15

replace (
	go.opentelemetry.io/contrib/propagators => /Users/wilbeguo/go/src/github.com/open-o11y/opentelemetry-go-contrib/propagators
	go.opentelemetry.io/otel/sdk => /Users/wilbeguo/go/src/github.com/Aneurysm9/opentelemetry-go/sdk
)

require (
	go.opentelemetry.io/contrib/propagators v0.14.0
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/otlp v0.14.0
	go.opentelemetry.io/otel/exporters/stdout v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
)
