// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var tracer = otel.Tracer("sample-app")

func main() {
	initProvider()

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("my-server"))

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []attribute.KeyValue{
		attribute.String("labelA", "chocolate"),
		attribute.String("labelB", "raspberry"),
		attribute.String("labelC", "vanilla"),
	}

	r.HandleFunc("/aws-sdk-call", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()

		xrayTraceID := getXrayTraceID(trace.SpanFromContext(ctx))
		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()
		_, _ = w.Write(payload)

	}))

	r.HandleFunc("/outgoing-http-call", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		ctx := r.Context()

		xrayTraceID, _ := func(ctx context.Context) (string, error) {

			req, _ := http.NewRequestWithContext(ctx, "GET", "https://aws.amazon.com", nil)

			res, err := client.Do(req)
			if err != nil {
				handleErr(err, "HTTP call to aws.amazon.com failed")
			}

			_, _ = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			return getXrayTraceID(trace.SpanFromContext(ctx)), err

		}(ctx)

		ctx, span := tracer.Start(
			ctx,
			"CollectorExporter-Example",
			trace.WithAttributes(commonLabels...))
		defer span.End()

		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()

		_, _ = w.Write(payload)

	}))

	http.Handle("/", r)

	// Start server
	address := os.Getenv("LISTEN_ADDRESS")
	if len(address) > 0 {
		_ = http.ListenAndServe(fmt.Sprintf("%s", address), nil)
	} else {
		// Default port 8080
		_ = http.ListenAndServe("localhost:8080", nil)
	}
}

func initProvider() {
	ctx := context.Background()

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://localhost:4317" // setting default endpoint for exporter
	}

	// Create and start new OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	handleErr(err, "failed to create new OTLP trace exporter")

	idg := xray.NewIDGenerator()

	service := os.Getenv("GO_GORILLA_SERVICE_NAME")
	if service == "" {
		service = "go-gorilla"
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String("test-service"),
	)
	handleErr(err, "failed to create resource")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
}

func getXrayTraceID(span trace.Span) string {
	xrayTraceID := span.SpanContext().TraceID().String()
	result := fmt.Sprintf("1-%s-%s", xrayTraceID[0:8], xrayTraceID[8:])
	return result
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
