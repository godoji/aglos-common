package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestOBVCross(t *testing.T) {
	ritmic.RunShortTestSet(Evaluate, [][]float64{{}}, Params)
}
