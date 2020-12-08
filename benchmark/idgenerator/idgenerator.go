package idgenerator

import (
	"go.opentelemetry.io/contrib/propagators/aws/xray/xrayidgenerator"
	"go.opentelemetry.io/otel/trace"
)

func generateNewTraceID() trace.TraceID {
	idg := xrayidgenerator.NewIDGenerator()
	traceID := idg.NewTraceID()
	return traceID
}
