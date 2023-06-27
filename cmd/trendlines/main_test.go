package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestTrendLines(t *testing.T) {
	ritmic.RunTestShort(Evaluate, [][]float64{{14}}, Params)
}
