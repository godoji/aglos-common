package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
	"math/rand"
)

var Params = []string{"frequency"}

func Evaluate(chart env.MarketSupplier, res *algo.ResultHandler, mem *env.Memory, param env.Parameters) {
	if rand.Float64() > param.Get("frequency") {
		return
	}
	res.NewEvent("random").SetColor("orange").SetIcon("cross")
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
