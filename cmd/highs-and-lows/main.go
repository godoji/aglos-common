package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
	candles "github.com/northberg/candlestick"
)

var Params = []string{"historySize"}

type LocalStore struct {
	Initialized bool
	History     *env.FiLoStack
}

func Evaluate(chart env.MarketSupplier, res *algo.ResultHandler, mem *env.Memory, param env.Parameters) {

	// A way of loading memory from disk
	var store *LocalStore
	if tmp := mem.Read(); tmp == nil {
		store = new(LocalStore)
	} else {
		store = tmp.(*LocalStore)
	}
	defer mem.Store(store)

	histSize := param.GetInt("historySize")
	if !store.Initialized {
		store.History = env.NewFiLoStack(histSize)
		store.Initialized = true
	}

	// Add candle to stack
	store.History.Push(chart.Interval(candles.Interval1d).Candle())

	// Skip if history is too small
	if !store.History.IsFull() {
		return
	}

	candidate := store.History.At(histSize / 2).(*candles.Candle)

	isHigh := true
	for _, candle := range store.History.ToSlice() {
		c := candle.(*candles.Candle)
		if c.Time == candidate.Time {
			continue
		}
		if candidate.High <= c.High {
			if c.High == candidate.High && c.Low < candidate.Low {
				continue
			}
			isHigh = false
			break
		}
	}
	if isHigh {
		res.NewEvent("high").SetPrice(candidate.High).SetTime(candidate.Time).SetColor("green").SetIcon("top")
	}

	isLow := true
	for _, candle := range store.History.ToSlice() {
		c := candle.(*candles.Candle)
		if c.Time == candidate.Time {
			continue
		}
		if candidate.Low >= c.Low {
			if c.Low == candidate.Low && c.High > candidate.High {
				continue
			}
			isLow = false
			break
		}
	}
	if isLow {
		res.NewEvent("low").SetPrice(candidate.Low).SetTime(candidate.Time).SetColor("red").SetIcon("bottom")
	}
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
