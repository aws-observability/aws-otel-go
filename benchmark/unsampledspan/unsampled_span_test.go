package unsampledspan

import (
	"testing"
)

func BenchmarkStartAndEndUnSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startAndEndUnSampledSpan()
	}
}

func BenchmarkStartAndEndNestedUnSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startAndEndNestedUnSampledSpan()
	}
}

func BenchmarkGetCurrentUnSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCurrentUnSampledSpan()
	}
}

func BenchmarkAddAttributesToUnSampledSpan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addAttributesToUnSampledSpan()
	}
}
