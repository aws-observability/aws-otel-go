package idgenerator

import (
	"testing"
)

func BenchmarkIDGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateNewTraceID()
	}
}
