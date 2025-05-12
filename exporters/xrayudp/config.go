// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package xrayudp

import "os"

type config struct {
	endpoint     string
	signalPrefix string
}

// Option sets configuration on the OTLP UDP Client.
type Option interface {
	apply(*config) *config
}

type optionFunc func(*config) *config

func (f optionFunc) apply(cfg *config) *config {
	return f(cfg)
}

// WithEndpoint sets custom daemon endpoint.
// If this option is not provided, a default endpoint will be used.
func WithEndpoint(endpoint string) Option {
	return optionFunc(func(cfg *config) *config {
		cfg.endpoint = endpoint
		return cfg
	})
}

func WithSignalPrefix(signalPrefix string) Option {
	return optionFunc(func(cfg *config) *config {
		cfg.signalPrefix = signalPrefix
		return cfg
	})
}

func newConfig(opts ...Option) *config {
	endpoint := DefaultEndpoint

	if isLambdaEnvironment() {
		// If in an AWS Lambda Environment, `AWS_XRAY_DAEMON_ADDRESS` will be defined
		endpoint = os.Getenv("AWS_XRAY_DAEMON_ADDRESS")
	}

	cfg := &config{
		endpoint:     endpoint,
		signalPrefix: DefaultFormatOtelTracesBinaryPrefix,
	}

	for _, option := range opts {
		option.apply(cfg)
	}

	return cfg
}

func isLambdaEnvironment() bool {
	// Detect if running in AWS Lambda environment
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
