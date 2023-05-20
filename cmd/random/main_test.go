package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestRandom(t *testing.T) {
	ritmic.RunTestShort(Evaluate, [][]float64{{0.005}}, Params)
}
