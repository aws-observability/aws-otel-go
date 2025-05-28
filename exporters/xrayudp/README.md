# AWS X-Ray UDP Exporter for Usage with OpenTelemetry on Lambda

## Installation

Install with the following command:

```
go get github.com/aws-observability/aws-otel-go/exporters/xrayudp
```

## Usage

```go
	// ...

	udpExporter, _ := xrayudp.NewSpanExporter(ctx)

	tp := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(udpExporter)))

	// ...
```
