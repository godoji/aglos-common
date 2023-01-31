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

	ema10 := chart.Interval(candles.Interval1d).Indicator("ema", 10).Value()
	ema50 := chart.Interval(candles.Interval1d).Indicator("ema", 50).Value()

	// check current trend state based on ema
	var nextState CrossState
	if ema50 > ema10 {
		nextState = StateDownTrend
	} else {
		nextState = StateUpTrend
	}

	// handle any trend flips
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

	// store next state as current state
	store.State = nextState

}

// Run a server to use this algorithm in headless mode
func main() {
	ritmic.Serve(Evaluate, Params)
}
