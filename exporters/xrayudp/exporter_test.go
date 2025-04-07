// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package xrayudp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOtlpUdpExporter(t *testing.T) {
	t.Run("Create Exporter", func(t *testing.T) {
		_, err := NewSpanExporter(context.TODO(), WithEndpoint("1.2.3.4:9876"), WithSignalPrefix("E3"))
		assert.NoError(t, err)
	})
}
