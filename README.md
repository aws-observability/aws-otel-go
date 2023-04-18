# AWS Distro for OpenTelemetry Go Language

## Introduction

This repo hosts documentation and sample apps for the ADOT Go library which provides the AWS service integrations for traces and metrics for the [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go) library. The library can be configured to support trace applications with the AWS X-Ray service.

Please note all source code for the Go library is upstream on the OpenTelemetry project in the [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go) library repo. All features of the OpenTelemetry library are available along with its components beinbg configured to create traces which can be viewed in the AWS X-Ray console and to allow propagation of those contexts across multiple downstream AWS services.

Once traces have been generated, they can be sent to a tracing service, like AWS X-Ray, to visualize and understand exactly what happened during the traced calls. For more information about the AWS X-Ray service, see the [AWS X-Ray Developer Guide](https://docs.aws.amazon.com/xray/latest/devguide/aws-xray.html). 

To send traces to AWS X-Ray, you can use the [AWS Distro for OpenTelemetry (ADOT) Collector](https://github.com/aws-observability/aws-otel-collector). OpenTelemetry Go exports traces from the application to the ADOT Collector. The ADOT Collector is configured with [AWS credentials for the CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html), an AWS region, and which trace attributes to index so that it can send the traces to the AWS X-Ray console. Read more about the [AWS X-Ray Tracing Exporter for OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/awsxrayexporter).

ADOT is in preview for Go metrics.

## Getting Started 

Check out the getting started [documentation](https://aws-otel.github.io/docs/getting-started/go-sdk)

## Sample Application (DEPRECATED)

**WARNING:** This sample app is deprecated and is no longer maintained.  Please use the [new standardized sample apps](https://github.com/aws-observability/aws-otel-community/tree/master/sample-apps).

See the [example sample application](https://github.com/aws-observability/aws-otel-go/blob/main/sampleapp/main.go) for setup steps.

The [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go) SDK provides entry points for configuration through its API. This is can be used to configure the [id_generator](https://github.com/open-telemetry/opentelemetry-go/blob/970755bd087801619575b7473806356818e24e15/sdk/trace/id_generator.go) needed to support the X-Ray trace ID format. In addition, it also allows the use of a custom propagator, passed into the tracer provider, to generate and propagate the AWS X-Ray trace header. 

## Useful Links

* For more information on OpenTelemetry, visit their [website](https://opentelemetry.io/)
* [OpenTelemetry Go core Repo](https://github.com/open-telemetry/opentelemetry-go)
* [OpenTelemetry Go Contrib Repo](https://github.com/open-telemetry/opentelemetry-go-contrib)
* [AWS Distro for OpenTelemetry](https://aws-otel.github.io/)

## Security issue notifications
If you discover a potential security issue in this project we ask that you notify AWS/Amazon Security via our [vulnerability reporting page](http://aws.amazon.com/security/vulnerability-reporting/). Please do **not** create a public github issue.

## License

This project is licensed under the Apache-2.0 License.
