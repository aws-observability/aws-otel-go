package unsampledspan

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("sample-app")

func main() {
	initTracer()
}

func startAndEndUnSampledSpan() {

	var span trace.Span
	_, span = tracer.Start(
		context.Background(),
		"Example Trace",
	)

	defer span.End()
}

func startAndEndNestedUnSampledSpan() {

	var span trace.Span
	ctx, span := tracer.Start(context.Background(), "Parent operation...")
	defer span.End()

	_, span = tracer.Start(ctx, "Sub operation...")
	defer span.End()
}

func getCurrentUnSampledSpan() trace.Span {

	var span trace.Span
	ctx, span := tracer.Start(
		context.Background(),
		"Example Trace",
	)
	defer span.End()

	return trace.SpanFromContext(ctx)
}

func addAttributesToUnSampledSpan() {

	var span trace.Span
	_, span = tracer.Start(
		context.Background(),
		"Example Trace",
	)
	defer span.End()

	span.SetAttributes(label.Key("example label 1").String("value 1"))
	span.SetAttributes(label.Key("example label 2").String("value 2"))
}

func initTracer() {

	idg := xray.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
}
