package main

import (
	"github.com/godoji/algocore/pkg/algo"
	"github.com/godoji/algocore/pkg/env"
	"github.com/godoji/algocore/pkg/ritmic"
	"github.com/northberg/candlestick"
	"math"
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

	if !store.Initialized {
		store.History = env.NewFiLoStack(5)
		store.Initialized = true
	}

	highsAndLows := chart.Algorithm("highs-and-lows", param.Get("historySize"))
	if highsAndLows.HasEvents() {
		lastEvent := highsAndLows.LastEvents()[0]
		store.History.Push(lastEvent)
	}

	if !store.History.IsFull() {
		return
	}

	t1 := store.History.At(0).(*algo.Event)
	b1 := store.History.At(1).(*algo.Event)
	t2 := store.History.At(2).(*algo.Event)
	b2 := store.History.At(3).(*algo.Event)
	t3 := store.History.At(4).(*algo.Event)

	maxTop := math.Max(t1.Price, t2.Price)

	if !(t1.Label == "high" && b1.Label == "low" && t2.Label == "high" && b2.Label == "low" && t3.Label == "high") {
		return
	}

	epsilon := 1.02

	if t1.Price > t3.Price*epsilon || t1.Price < t3.Price/epsilon {
		return
	}

	if t2.Price < t1.Price*epsilon || t2.Price < t3.Price*epsilon {
		return
	}

	if b1.Price > b2.Price*epsilon || b1.Price < b2.Price/epsilon {
		return
	}

	ch := chart.Interval(candlestick.Interval1d)
	leftIndex := ch.ToIndex(t1.Time)
	rightIndex := ch.ToIndex(t3.Time)

	for i := rightIndex; i <= leftIndex; i++ {
		c := ch.FromLast(int(i))
		if c.Time == t1.Time || c.Time == t2.Time || c.Time == b1.Time || c.Time == t3.Time || c.Time == b2.Time {
			continue
		}
		if c.Missing {
			continue
		}
		if c.High > maxTop {
			return
		}
		if c.Low < b1.Price {
			return
		}
	}

	var leftCandle *candlestick.Candle
	for i := 1; i < 200; i++ {
		j := int(leftIndex) + i
		c := ch.FromLast(j)
		if c.Missing {
			continue
		}
		if c.High >= maxTop {
			break
		}
		if c.Low <= b1.Price*1.01 {
			leftCandle = c
			break
		}
	}
	if leftCandle == nil {
		return
	}

	var rightCandle *candlestick.Candle
	for i := 1; i < 200; i++ {
		j := int(rightIndex) - i
		if j < 0 {
			break
		}
		c := ch.FromLast(j)
		if c.Missing {
			continue
		}
		if c.High >= maxTop {
			break
		}
		if c.Low <= b2.Price*1.01 {
			rightCandle = c
			break
		}
	}
	if rightCandle == nil {
		return
	}

	event := res.NewEvent("head-and-shoulders")
	event.AddSegment(&algo.SegmentAnnotation{
		TimeBegin:  leftCandle.Time,
		TimeEnd:    rightCandle.Time,
		PriceBegin: b1.Price,
		PriceEnd:   maxTop,
		Style:      "region",
		Color:      "rgb(100,100,100)",
	})
	event.AddEvent(t1, "t1")
	event.AddEvent(b1, "b1")
	event.AddEvent(t2, "t2")
	event.AddEvent(b2, "b2")
	event.AddEvent(t3, "t3")
	store.History.Clear()
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
