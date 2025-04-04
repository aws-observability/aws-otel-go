// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws-observability/aws-otel-go/exporters/xrayudp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func webServer() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("healthcheck"))
		if err != nil {
			log.Println(err)
		}
	}))

	http.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("testTracer")

		ctx := context.Background()

		_, span := tracer.Start(ctx, "testSpan", oteltrace.WithSpanKind(oteltrace.SpanKindServer))
		traceId := span.SpanContext().TraceID().String()
		span.End()

		xrayFormatTraceId := "1-" + traceId[0:8] + "-" + traceId[8:]
		log.Printf("X-Ray Trace ID is: %s", xrayFormatTraceId)

		_, err := w.Write([]byte(xrayFormatTraceId))
		if err != nil {
			log.Println(err)
		}
	}))

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "localhost:8080"
	}
	log.Printf("App is listening on %s !", listenAddress)
	_ = http.ListenAndServe(listenAddress, nil)
}

func start_otel() error {
	ctx := context.Background()

	idg := xray.NewIDGenerator()

	myudpexporter, _ := xrayudp.NewSpanExporter(ctx)

	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(myudpexporter)),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return nil
}

func main() {
	log.Println("Starting Golang OTel Sample App...")

	err := start_otel()
	if err != nil {
		log.Fatalf("Failed to start OTel: %v", err)
		return
	}

	webServer()
}
