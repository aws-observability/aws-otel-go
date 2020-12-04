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
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray/xrayidgenerator"
	awspropagator "go.opentelemetry.io/contrib/propagators/awsxray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("sample-app")

func main() {

	initTracer()

	// _, err := obsvsS3.NewInstrumentedS3Client(
	// 	&mocks.MockS3Client{},
	// 	obsvsS3.WithTracerProvider(tracerProvider),
	// 	obsvsS3.WithSpanCorrelation(true),
	// )
	// handleErr(err, "failed to create new S3 Client")

	// tracer := tracerProvider.Tracer("http-tracer")
	// ctx, span := tracer.Start(
	// 	context.Background(),
	// 	"http_request_served",
	// )
	// defer span.End()

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("my-server"))

	r.HandleFunc("/aws-sdk-call", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()

		// _, _ = client.PutObjectWithContext(outerSpanCtx, &s3.PutObjectInput{
		// 	Bucket: aws.String("test-bucket"),
		// 	Key:    aws.String("010101"),
		// 	Body:   bytes.NewReader([]byte("foo")),
		// })

		// time.Sleep(time.Second * 15)

		xrayTraceID := getXrayTraceID(trace.SpanFromContext(ctx))
		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()

		w.Write(payload)

	}))

	r.HandleFunc("/outgoing-http-call", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		ctx := r.Context()

		xrayTraceID, _ := func(ctx context.Context) (string, error) {
			// ctx, span := tracer.Start(ctx, "HTTP GET Request")
			// defer span.End()

			// ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
			req, _ := http.NewRequestWithContext(ctx, "GET", "https://aws.amazon.com/", nil)

			res, err := client.Do(req)
			if err != nil {
				handleErr(err, "HTTP call to aws.amazon.com failed")
			}

			ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			return getXrayTraceID(trace.SpanFromContext(ctx)), err

		}(ctx)

		time.Sleep(10 * time.Second)

		// xrayTraceID := getXrayTraceID(span)
		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()

		w.Write(payload)

	}))

	http.Handle("/", r)

	// Start server
	address := os.Getenv("LISTEN_ADDRESS")
	if len(address) > 0 {
		http.ListenAndServe(fmt.Sprintf(":%s", address), nil)
	} else {
		// Default port 8000
		http.ListenAndServe(":8080", nil)
	}
}

func initTracer() {

	// Create new OTLP Exporter
	exporter, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:30080"),
		// otlp.WithGRPCDialOption(grpc.WithBlock()),
	)
	handleErr(err, "failed to create new OTLP exporter")

	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	idg := xrayidgenerator.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(awspropagator.Xray{})
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
