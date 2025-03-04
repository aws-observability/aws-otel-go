// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package otlptraceudp

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
)

// New constructs a new Exporter and starts it.
func New(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	client, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return otlptrace.New(ctx, client)
}

// NewUnstarted constructs a new Exporter and does not start it.
func NewUnstarted(opts ...Option) (*otlptrace.Exporter, error) {
	client, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return otlptrace.NewUnstarted(client), nil
}
