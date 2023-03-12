package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestSupportResistance(t *testing.T) {
	ritmic.RunTestShort(Evaluate, [][]float64{{200}}, Params)
}
