package sampledspan

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/aws/xray/xrayidgenerator"
	awspropagator "go.opentelemetry.io/contrib/propagators/awsxray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
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

	span.SetAttributes(label.Key("example label 1").String("value 1"))
	span.SetAttributes(label.Key("example label 2").String("value 2"))
}

func initTracer() {

	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	idg := xrayidgenerator.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(awspropagator.Xray{})
}
