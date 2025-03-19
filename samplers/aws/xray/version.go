// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package xray // import "github.com/aws-observability/aws-otel-go/samplers/aws/xray"

// Version is the current release version of the AWS XRay remote sampler.
func Version() string {
	return "0.24.0"
	// This string is updated by the pre_release.sh script during release
}

// SemVersion is the semantic version to be supplied to tracer/meter creation.
//
// Deprecated: Use [Version] instead.
func SemVersion() string {
	return Version()
}
