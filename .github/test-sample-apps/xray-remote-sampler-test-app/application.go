package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	sampler "github.com/aws-observability/aws-otel-go/samplers/aws/xray"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func getSampledSpanCount(name string, totalSpans string, attributes []attribute.KeyValue) (int, error) {
	tracer := otel.Tracer(name)

	var sampleCount = 0
	totalSamples, err := strconv.Atoi(totalSpans)
	if err != nil {
		return 0, err
	}

	ctx := context.Background()

	for i := 0; i < totalSamples; i++ {
		_, span := tracer.Start(ctx, name, oteltrace.WithSpanKind(oteltrace.SpanKindServer), oteltrace.WithAttributes(attributes...))

		if span.SpanContext().IsSampled() {
			sampleCount++
		}

		span.End()
	}

	return sampleCount, nil
}

func webServer() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("healthcheck"))
		if err != nil {
			log.Println(err)
		}
	}))

	http.Handle("/getSampled", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAttribute := r.Header.Get("User")
		required := r.Header.Get("Required")
		serviceName := r.Header.Get("Service_name")
		totalSpans := r.Header.Get("Totalspans")

		var attributes = []attribute.KeyValue{
			attribute.KeyValue{Key: "http.method", Value: attribute.StringValue(r.Method)},
			attribute.KeyValue{Key: "http.url", Value: attribute.StringValue("http://localhost:8080/getSampled")},
			attribute.KeyValue{Key: "user", Value: attribute.StringValue(userAttribute)},
			attribute.KeyValue{Key: "http.route", Value: attribute.StringValue("/getSampled")},
			attribute.KeyValue{Key: "required", Value: attribute.StringValue(required)},
			attribute.KeyValue{Key: "http.target", Value: attribute.StringValue("/getSampled")},
		}

		totalSampled, err := getSampledSpanCount(serviceName, totalSpans, attributes)
		if err != nil {
			log.Println(err)
		}

		_, err = w.Write([]byte(strconv.Itoa(totalSampled)))
		if err != nil {
			log.Println(err)
		}
	}))

	http.Handle("/importantEndpoint", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAttribute := r.Header.Get("User")
		required := r.Header.Get("Required")
		serviceName := r.Header.Get("Service_name")
		totalSpans := r.Header.Get("Totalspans")

		var attributes = []attribute.KeyValue{
			attribute.KeyValue{Key: "http.method", Value: attribute.StringValue("GET")},
			attribute.KeyValue{Key: "http.url", Value: attribute.StringValue("http://localhost:8080/importantEndpoint")},
			attribute.KeyValue{Key: "user", Value: attribute.StringValue(userAttribute)},
			attribute.KeyValue{Key: "http.route", Value: attribute.StringValue("/importantEndpoint")},
			attribute.KeyValue{Key: "required", Value: attribute.StringValue(required)},
			attribute.KeyValue{Key: "http.target", Value: attribute.StringValue("/importantEndpoint")},
		}

		totalSampled, err := getSampledSpanCount(serviceName, totalSpans, attributes)
		if err != nil {
			log.Println(err)
		}

		_, err = w.Write([]byte(strconv.Itoa(totalSampled)))
		if err != nil {
			log.Println(err)
		}
	}))

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "localhost:8080"
	}
	log.Println("App is listening on %s !", listenAddress)
	_ = http.ListenAndServe(listenAddress, nil)
}

func start_xray() error {
	ctx := context.Background()

	exporterEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if exporterEndpoint == "" {
		exporterEndpoint = "localhost:4317"
	}

	idg := xray.NewIDGenerator()

	samplerEndpoint := os.Getenv("XRAY_ENDPOINT")
	if samplerEndpoint == "" {
		samplerEndpoint = "http://localhost:2000"
	}
	endpointUrl, err := url.Parse(samplerEndpoint)

	res, err := sampler.NewRemoteSampler(ctx, "adot-integ-test", "", sampler.WithEndpoint(*endpointUrl), sampler.WithSamplingRulesPollingInterval(10*time.Second))
	if err != nil {
		log.Fatalf("Failed to create new XRay Remote Sampler: %v", err)
		return err
	}

	// attach remote sampler to tracer provider
	tp := trace.NewTracerProvider(
		trace.WithSampler(res),
		trace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return nil
}

func main() {
	log.Println("Starting Golang OTel Sample App...")

	err := start_xray()
	if err != nil {
		log.Fatalf("Failed to start XRay: %v", err)
		return
	}

	webServer()
}
