// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package xrayudp

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

const (
	DefaultFormatOtelTracesBinaryPrefix = "T1S"
)

type OtlpTraceClient interface {
	Shutdown() error
	SendData(data []byte, signalFormatPrefix string) error
}

type client struct {
	udpExporter  OtlpTraceClient
	signalPrefix string
	endpoint     string
}

// NewClient creates a new OTLP UDP trace client.
func NewClient(opts ...Option) (otlptrace.Client, error) {
	config := newConfig(opts...)

	udpExporter, err := NewUdpExporter(config.endpoint)
	if err != nil {
		return nil, err
	}

	return &client{
		udpExporter:  udpExporter,
		signalPrefix: config.signalPrefix,
		endpoint:     config.endpoint,
	}, nil
}

// Start does nothing in an OTLP UDP client.
func (d *client) Start(ctx context.Context) error {
	// nothing to do
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// Stop shuts down the UDP Exporter.
func (d *client) Stop(ctx context.Context) error {
	return d.udpExporter.Shutdown()
}

// UploadTraces sends a batch of spans to the UDP Exporter.
func (d *client) UploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error {
	pbRequest := &coltracepb.ExportTraceServiceRequest{
		ResourceSpans: protoSpans,
	}
	rawRequest, err := proto.Marshal(pbRequest)
	if err != nil {
		return err
	}

	if err := d.udpExporter.SendData(rawRequest, d.signalPrefix); err != nil {
		otel.Handle(fmt.Errorf("error exporting spans: %w", err))
		return err
	}

	return nil
}
