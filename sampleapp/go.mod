module github.com/aws-observability/aws-otel-go/sampleapp

go 1.15

replace (
	go.opentelemetry.io/contrib/propagators => /Users/wilbeguo/go/src/github.com/open-o11y/opentelemetry-go-contrib/propagators
	go.opentelemetry.io/otel/sdk => /Users/wilbeguo/go/src/github.com/Aneurysm9/opentelemetry-go/sdk
)

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/contrib/propagators v0.14.0
	go.opentelemetry.io/otel v0.14.0
	go.opentelemetry.io/otel/exporters/otlp v0.14.0
	go.opentelemetry.io/otel/sdk v0.14.0
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201119123407-9b1e624d6bc4 // indirect
	google.golang.org/grpc v1.33.2
)
