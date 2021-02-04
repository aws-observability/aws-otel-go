package propagator

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"google.golang.org/grpc"
	"net/http"
	"testing"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func BenchmarkPropagatorExtract(b *testing.B) {
	propagator := xray.Propagator{}

	ctx := context.Background()
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	req.Header.Set("Root", "1-8a3c60f7-d188f8fa79d48a391a778fa6")
	req.Header.Set("Parent", "53995c3f42cd8ad8")
	req.Header.Set("Sampled", "1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = propagator.Extract(ctx, req.Header)
	}
}

func BenchmarkPropagatorInject(b *testing.B) {
	propagator := xray.Propagator{}
	initTracer()
	tracer := otel.Tracer("test")

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	ctx, _ := tracer.Start(context.Background(), "Parent operation...")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		propagator.Inject(ctx, req.Header)
	}
}

func initTracer() {
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint("0.0.0.0:55680"),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)

	// Create new OTLP Exporter
	exporter, _ := otlp.NewExporter(context.Background(), driver)

	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
}
