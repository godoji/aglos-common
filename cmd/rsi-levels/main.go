package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
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

func Evaluate(chart env.MarketSupplier, res *algo.ResultHandler, mem *env.Memory, param env.Parameters) {

	var store *LocalStore
	if tmp := mem.Read(); tmp == nil {
		store = new(LocalStore)
	} else {
		store = tmp.(*LocalStore)
	}
	defer mem.Store(store)

	rsi := chart.Interval(candles.Interval1d).Indicator("rsi", 14).Value()

	var nextState CrossState
	if rsi >= 70 {
		nextState = StateDownTrend
	} else if rsi <= 30 {
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
