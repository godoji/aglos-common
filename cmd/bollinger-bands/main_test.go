package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestBollingerBands(t *testing.T) {
	ritmic.RunShortTestSet(Evaluate, [][]float64{{}}, Params)
}
