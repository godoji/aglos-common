package main

import (
	"github.com/godoji/algocore/pkg/ritmic"
	"testing"
)

func TestTripleTop(t *testing.T) {
	ritmic.RunTestShort(Evaluate, [][]float64{{7}}, Params)
}
