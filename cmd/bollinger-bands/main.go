package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/ritmic"
	"github.com/godoji/algocore/pkg/simulated"
	candles "github.com/northberg/candlestick"
)

var Params = []string{}

type CrossState = int

const (
	_ CrossState = iota
	StateUpTrend
	StateDownTrend
)

type LocalStore struct {
	State CrossState
}

func Evaluate(chart simulated.MarketSupplier, res *algo.ResultHandler, mem *simulated.Memory, param simulated.Parameters) {

	var store *LocalStore
	if tmp := mem.Read(); tmp == nil {
		store = new(LocalStore)
	} else {
		store = tmp.(*LocalStore)
	}
	defer mem.Store(store)

	bb := chart.Interval(candles.Interval1d).Indicator("bb", 20, 2)

	var nextState CrossState
	if chart.Interval(candles.Interval1d).Candle().High >= bb.Series("upper") {
		nextState = StateDownTrend
	} else if chart.Interval(candles.Interval1d).Candle().Low <= bb.Series("lower") {
		nextState = StateUpTrend
	} else {
		nextState = store.State
	}

	switch store.State {
	case StateUpTrend:
		if nextState == StateDownTrend {
			res.NewEvent("downtrend").SetColor("red").SetIcon("down")
		}
	case StateDownTrend:
		if nextState == StateUpTrend {
			res.NewEvent("uptrend").SetColor("green").SetIcon("up")
		}
	}

	store.State = nextState
}

// Run a server to use this algorithm in headless mode
func main() {
	ritmic.Serve(Evaluate, Params)
}
