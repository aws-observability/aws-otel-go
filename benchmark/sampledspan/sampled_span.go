package sampledspan

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("sample-app")

func main() {
	initTracer()
}

func startAndEndSampledSpan() {

	var span trace.Span
	_, span = tracer.Start(
		context.Background(),
		"Example Trace",
	)

	defer span.End()
}

func startAndEndNestedSampledSpan() {

	var span trace.Span
	ctx, span := tracer.Start(context.Background(), "Parent operation...")
	defer span.End()

	_, span = tracer.Start(ctx, "Sub operation...")
	defer span.End()
}

func getCurrentSampledSpan() trace.Span {

	var span trace.Span
	ctx, span := tracer.Start(
		context.Background(),
		"Example Trace",
	)
	defer span.End()

	return trace.SpanFromContext(ctx)
}

func addAttributesToSampledSpan() {

	var span trace.Span
	_, span = tracer.Start(
		context.Background(),
		"Example Trace",
	)
	defer span.End()

	span.SetAttributes(attribute.Key("example attribute 1").String("value 1"))
	span.SetAttributes(attribute.Key("example attribute 2").String("value 2"))
}

func initTracer() {

	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
}
