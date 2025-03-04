// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package otlptraceudp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOtlpUdpExporter(t *testing.T) {
	t.Run("Create Exporter", func(t *testing.T) {
		_, err := New(context.TODO(), WithEndpoint("1.2.3.4:9876"), WithSignalPrefix("E3"))
		assert.NoError(t, err)
	})

	t.Run("Create Unstarted Exporter", func(t *testing.T) {
		_, err := NewUnstarted(WithEndpoint("1.2.3.4:9876"), WithSignalPrefix("E3"))
		assert.NoError(t, err)
	})
}
