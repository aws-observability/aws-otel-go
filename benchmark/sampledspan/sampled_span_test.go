package sampledspan

import (
	"testing"
)

func BenchmarkStartAndEndSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startAndEndSampledSpan()
	}
}

func BenchmarkStartAndEndNestedSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startAndEndNestedSampledSpan()
	}
}

func BenchmarkGetCurrentSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCurrentSampledSpan()
	}
}

func BenchmarkAddAttributesToSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addAttributesToSampledSpan()
	}
}
