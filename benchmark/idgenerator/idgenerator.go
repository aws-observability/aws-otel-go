package idgenerator

import (
	"context"
	
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/trace"
)

func generateNewTraceID() trace.TraceID {
	idg := xray.NewIDGenerator()
	traceID, _ := idg.NewIDs(context.Background())
	return traceID
}
