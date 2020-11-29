package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/propagators/aws/xray/xrayidgenerator"
	awspropagator "go.opentelemetry.io/contrib/propagators/awsxray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func initProvider() func() {
	ctx := context.Background()

	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:30080"),
		otlp.WithGRPCDialOption(grpc.WithBlock()),
	)
	handleErr(err, "failed to create exporter")

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	handleErr(err, "failed to create resource")

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	idg := xrayidgenerator.NewIDGenerator()

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithIDGenerator(idg),
	)

	pusher := push.New(
		basic.New(
			simple.NewWithExactDistribution(),
			exp,
		),
		exp,
		push.WithPeriod(2*time.Second),
	)

	otel.SetTextMapPropagator(awspropagator.Xray{})
	// otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(pusher.MeterProvider())
	pusher.Start()

	return func() {
		handleErr(tracerProvider.Shutdown(ctx), "failed to shutdown provider")
		handleErr(exp.Shutdown(ctx), "failed to stop exporter")
		pusher.Stop() // pushes any last exports to the receiver
	}
}

func main() {
	log.Printf("Waiting for connection...")

	shutdown := initProvider()
	defer shutdown()

	tracer := otel.Tracer("test-tracer")
	meter := otel.Meter("test-meter")

	commonLabels := []label.KeyValue{
		label.String("labelA", "chocolate"),
		label.String("labelB", "raspberry"),
		label.String("labelC", "vanilla"),
	}

	// Recorder metric example
	valuerecorder := metric.Must(meter).
		NewFloat64Counter(
			"an_important_metric",
			metric.WithDescription("Measures the cumulative epicness of the app"),
		).Bind(commonLabels...)
	defer valuerecorder.Unbind()

	// work begins
	ctx, span := tracer.Start(
		context.Background(),
		"Example Trace",
		trace.WithAttributes(commonLabels...))
	defer span.End()

	// Create new router to handle API endpoints
	router := mux.NewRouter()

	// HTTP GET: /aws-sdk-call endpoint
	router.HandleFunc("/aws-sdk-call", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		xrayTraceID := getXrayTraceID(span)
		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()
		w.Write(payload)
	}).Methods(http.MethodGet)

	// HTTP GET: /outgoing-http-call endpoint
	router.HandleFunc("/outgoing-http-call", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response, err := http.Get("https://aws.amazon.com/")
		if err != nil || response.StatusCode != http.StatusOK {
			fmt.Println("HTTP call failed:", err)
			return
		}
		valuerecorder.Add(ctx, 1.0)
		xrayTraceID := getXrayTraceID(span)
		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()
		w.Write(payload)
	}).Methods(http.MethodGet)

	// Start server
	address := os.Getenv("LISTEN_ADDRESS")
	if len(address) > 0 {
		http.ListenAndServe(fmt.Sprintf(":%s", address), router)
	} else {
		// Default port 8000
		http.ListenAndServe(":8000", router)
	}
}

func getXrayTraceID(span trace.Span) string {
	xrayTraceID := span.SpanContext().TraceID.String()
	result := fmt.Sprintf("1-%s-%s", xrayTraceID[0:8], xrayTraceID[8:])
	return result
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
