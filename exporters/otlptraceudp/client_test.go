// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package otlptraceudp

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"
)

type MockUdpExporter struct {
	sentData          []byte
	shouldErrorOnSend bool
}

func (m *MockUdpExporter) Shutdown() error {
	return nil
}
func (m *MockUdpExporter) SendData(data []byte, signalFormatPrefix string) error {
	if m.shouldErrorOnSend {
		return fmt.Errorf("SendDataError")
	}
	m.sentData = data
	return nil
}

func TestOtlpUdpClient(t *testing.T) {
	t.Run("Client uploads traces successfully", func(t *testing.T) {
		mockUdpExporter := &MockUdpExporter{}

		exporter := &client{
			udpExporter: mockUdpExporter,
		}

		spans := []*tracepb.ResourceSpans{&tracepb.ResourceSpans{}, &tracepb.ResourceSpans{}}
		err := exporter.UploadTraces(context.TODO(), spans)

		assert.NoError(t, err)

		pbRequest := &coltracepb.ExportTraceServiceRequest{
			ResourceSpans: spans,
		}
		rawRequest, err := proto.Marshal(pbRequest)
		assert.NoError(t, err)

		assert.Equal(t, mockUdpExporter.sentData, rawRequest)
	})

	t.Run("Return error when internal udpExporter.SendData(...) returns error", func(t *testing.T) {
		mockUdpExporter := &MockUdpExporter{shouldErrorOnSend: true}

		exporter := &client{
			udpExporter: mockUdpExporter,
		}

		spans := []*tracepb.ResourceSpans{&tracepb.ResourceSpans{}, &tracepb.ResourceSpans{}}
		err := exporter.UploadTraces(context.TODO(), spans)

		assert.Error(t, err)
		assert.Equal(t, "SendDataError", err.Error())
	})

	t.Run("Stop the OTLP Udp Client successfully", func(t *testing.T) {
		mockUdpExporter := &MockUdpExporter{}
		exporter := &client{
			udpExporter: mockUdpExporter,
		}
		err := exporter.Stop(context.TODO())

		assert.NoError(t, err)
	})
}
