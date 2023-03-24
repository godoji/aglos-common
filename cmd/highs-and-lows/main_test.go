package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestHighsAndLows(t *testing.T) {
	ritmic.RunTestShort(Evaluate, [][]float64{{7}}, Params)
}

func BenchmarkHighsAndLows(b *testing.B) {
	for i := 0; i < 1000; i++ {
		ritmic.RunTestShort(Evaluate, [][]float64{{7}}, Params)
	}
}
