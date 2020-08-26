package mutux

import (
	"testing"
)

func BenchmarkWork1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		work2()
	}
}
