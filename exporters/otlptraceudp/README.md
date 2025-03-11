# AWS OTLP UDP Exporter for Lambda

## Installation

Install with the following command:

```
go get github.com/aws-observability/aws-otel-go/exporters/otlptraceudp
```

## Usage

```go
    // ...

    daemonAddress := os.Getenv("AWS_XRAY_DAEMON_ADDRESS")
    udpExporter, _ := otlptraceudp.New(ctx, otlptraceudp.WithEndpoint(daemonAddress))

    tp := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(udpExporter)))

    // ...
```