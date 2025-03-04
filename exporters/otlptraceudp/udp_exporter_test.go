// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package otlptraceudp

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testingPrefix = "T1"

type MockConn struct {
	writtenBytes       []byte
	closeCalled        bool
	shouldErrorOnWrite bool
}

func (m *MockConn) Write(byteArray []byte) (int, error) {
	m.writtenBytes = byteArray

	if m.shouldErrorOnWrite {
		return 0, fmt.Errorf("simulate_error")
	}
	return 0, nil
}

func (m *MockConn) Close() error {
	m.closeCalled = true
	return nil
}

func TestUdpExporter(t *testing.T) {
	endpoint := "localhost:3000"
	host := "localhost"
	port := 3000

	t.Run("should parse the endpoint correctly", func(t *testing.T) {
		exporter, err := NewUdpExporter(endpoint)
		assert.Nil(t, err)
		assert.Equal(t, host, exporter.host)
		assert.Equal(t, port, exporter.port)
	})

	t.Run("should send UDP data correctly", func(t *testing.T) {
		mockConn := &MockConn{writtenBytes: nil}

		exporter := &UdpExporter{
			conn: mockConn,
		}

		data := []byte{1, 2, 3}
		prefix := testingPrefix
		err := exporter.SendData(data, prefix)

		assert.NoError(t, err)

		base64EncodedString := base64.StdEncoding.EncodeToString(data)
		message := fmt.Sprintf("%s%s%s", ProtocolHeader, prefix, base64EncodedString)

		assert.Equal(t, mockConn.writtenBytes, []byte(message))
	})

	t.Run("should handle errors when sending UDP data", func(t *testing.T) {
		mockConn := &MockConn{closeCalled: false, shouldErrorOnWrite: true}

		exporter := &UdpExporter{
			conn: mockConn,
		}

		data := []byte{1, 2, 3}
		prefix := testingPrefix
		err := exporter.SendData(data, prefix)

		assert.Error(t, err)
		assert.Equal(t, "simulate_error", err.Error())
	})

	t.Run("should close the socket on shutdown", func(t *testing.T) {
		mockConn := &MockConn{closeCalled: false}

		exporter := &UdpExporter{
			conn: mockConn,
		}

		err := exporter.Shutdown()

		assert.NoError(t, err)
		assert.True(t, mockConn.closeCalled)
	})
}
