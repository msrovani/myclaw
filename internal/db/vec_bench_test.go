package db

import (
	"math/rand"
	"testing"
)

func BenchmarkCosineDistance(b *testing.B) {
	const dim = 1536
	vecA := make([]float32, dim)
	vecB := make([]float32, dim)

	for i := 0; i < dim; i++ {
		vecA[i] = rand.Float32()
		vecB[i] = rand.Float32()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CosineDistance(vecA, vecB)
	}
}
