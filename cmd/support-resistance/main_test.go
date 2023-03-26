package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestSupportResistance(t *testing.T) {
	ritmic.RunTestShort(Evaluate, [][]float64{{200}}, Params)
}

func BenchmarkSupportResistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ritmic.RunTestShort(Evaluate, [][]float64{{200}}, Params)
	}
}
