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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	obsvsS3 "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go/service/s3/otels3"
	mocks "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go/service/s3/otels3/mocks"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/propagators/aws/xray/xrayidgenerator"
	awspropagator "go.opentelemetry.io/contrib/propagators/awsxray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("mux-server")
var meter = otel.Meter("test-meter")

func main() {

	tracerProvider := initProvider()

	valuerecorder := metric.Must(meter).
		NewFloat64Counter(
			"an_important_metric",
			metric.WithDescription("Measures the cumulative epicness of the app"),
		)

	_, err := obsvsS3.NewInstrumentedS3Client(
		&mocks.MockS3Client{},
		obsvsS3.WithTracerProvider(tracerProvider),
		obsvsS3.WithSpanCorrelation(true),
	)
	handleErr(err, "failed to create new S3 Client")

	tracer := tracerProvider.Tracer("http-tracer")
	outerSpanCtx, span := tracer.Start(
		context.Background(),
		"http_request_served",
	)
	defer span.End()

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("my-server"))

	r.HandleFunc("/aws-sdk-call", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		// _, _ = client.PutObjectWithContext(outerSpanCtx, &s3.PutObjectInput{
		// 	Bucket: aws.String("test-bucket"),
		// 	Key:    aws.String("010101"),
		// 	Body:   bytes.NewReader([]byte("foo")),
		// })

		// time.Sleep(time.Second * 15)

		xrayTraceID := getXrayTraceID(span)
		json := simplejson.New()
		json.Set("traceId", xrayTraceID)
		payload, _ := json.MarshalJSON()

		w.Write(payload)

	}))

	r.HandleFunc("/outgoing-http-call", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		// response, err := http.Get("https://aws.amazon.com/")
		http.Get("https://aws.amazon.com/")
		// if err != nil || response.StatusCode != http.StatusOK {
		// 	handleErr(err, "HTTP call to aws.amazon.com failed")
		// }

		valuerecorder.Add(outerSpanCtx, 1.0)

		xrayTraceID := getXrayTraceID(span)
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

func initProvider() *sdktrace.TracerProvider {

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
	// bsp := sdktrace.NewBatchSpanProcessor(exporter)
	ssp := sdktrace.NewSimpleSpanProcessor(exporter)
	idg := xrayidgenerator.NewIDGenerator()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithSpanProcessor(ssp),
		sdktrace.WithIDGenerator(idg),
	)

	pusher := push.New(
		basic.New(
			simple.NewWithExactDistribution(),
			exporter,
		),
		exporter,
		push.WithPeriod(2*time.Second),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(awspropagator.Xray{})
	otel.SetMeterProvider(pusher.MeterProvider())
	pusher.Start()

	return tp
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
