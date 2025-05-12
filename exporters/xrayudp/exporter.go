// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package xrayudp

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
)

// NewSpanExporter constructs a new X-Ray UDP Span Exporter and starts it.
func NewSpanExporter(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	client, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return otlptrace.New(ctx, client)
}
