// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

// 2025-05-07: Begin Amazon modification.
package xray // import "github.com/aws-observability/aws-otel-go/samplers/aws/xray"
// End of Amazon modification.

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"
)

func newSeed() int64 {
	var seed int64
	if err := binary.Read(crand.Reader, binary.BigEndian, &seed); err != nil {
		// fallback to timestamp
		seed = time.Now().UnixNano()
	}
	return seed
}

var seed = newSeed()

func newGlobalRand() *rand.Rand {
	src := rand.NewSource(seed)
	if src64, ok := src.(rand.Source64); ok {
		return rand.New(src64) //nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand) is ignored as this is not security-sensitive.
	}
	return rand.New(src) //nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand) is ignored as this is not security-sensitive.
}
