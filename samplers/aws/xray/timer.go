// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Modifications Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.

// 2025-05-07: Begin Amazon modification.
package xray // import "github.com/aws-observability/aws-otel-go/samplers/aws/xray"
// End of Amazon modification.

import (
	"time"
)

// ticker is the same as time.Ticker except that it has jitters.
// A Ticker must be created with newTicker.
type ticker struct {
	tick     *time.Ticker
	duration time.Duration
	jitter   time.Duration
}

// newTicker creates a new Ticker that will send the current time on its channel with the passed jitter.
func newTicker(duration, jitter time.Duration) *ticker {
	t := time.NewTicker(duration - time.Duration(newGlobalRand().Int63n(int64(jitter))))

	jitteredTicker := ticker{
		tick:     t,
		duration: duration,
		jitter:   jitter,
	}

	return &jitteredTicker
}

// c returns a channel that receives when the ticker fires.
func (j *ticker) c() <-chan time.Time {
	return j.tick.C
}
