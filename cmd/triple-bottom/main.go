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

	b1 := store.History.At(0).(*algo.Event)
	t1 := store.History.At(1).(*algo.Event)
	b2 := store.History.At(2).(*algo.Event)
	t2 := store.History.At(3).(*algo.Event)
	b3 := store.History.At(4).(*algo.Event)

	supportPrice := math.Min(b1.Price, b2.Price)

	if !(b1.Label == "low" && t1.Label == "high" && b2.Label == "low" && t2.Label == "high" && b3.Label == "low") {
		return
	}

	if b1.Price > b2.Price*1.01 || b1.Price < b2.Price/1.01 || b3.Price > b2.Price*1.01 || b3.Price < b2.Price/1.01 {
		return
	}

	if t1.Price > t2.Price*1.01 || t1.Price < t2.Price/1.01 {
		return
	}

	ch := chart.Interval(candlestick.Interval1d)
	leftIndex := ch.ToIndex(b1.Time)
	rightIndex := ch.ToIndex(b3.Time)

	for i := rightIndex; i <= leftIndex; i++ {
		c := ch.FromLast(int(i))
		if c.Time == b1.Time || c.Time == b2.Time || c.Time == t1.Time || c.Time == b3.Time || c.Time == t2.Time {
			continue
		}
		if c.Missing {
			continue
		}
		if c.High > t1.Price {
			return
		}
		if c.Low < supportPrice {
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
		if c.Low <= supportPrice {
			break
		}
		if c.High >= t1.Price/1.01 {
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
		if c.Low <= supportPrice {
			break
		}
		if c.High >= t2.Price/1.01 {
			rightCandle = c
			break
		}
	}
	if rightCandle == nil {
		return
	}

	event := res.NewEvent("triple-bottom")
	event.AddSegment(&algo.SegmentAnnotation{
		TimeBegin:  leftCandle.Time,
		TimeEnd:    rightCandle.Time,
		PriceBegin: t1.Price,
		PriceEnd:   supportPrice,
		Style:      "region",
		Color:      "rgb(100,100,100)",
	})
	event.AddEvent(b1, "b1")
	event.AddEvent(t1, "t1")
	event.AddEvent(b2, "b2")
	event.AddEvent(t2, "t2")
	event.AddEvent(b3, "b3")
	store.History.Clear()
}

func main() {
	ritmic.Serve(Evaluate, Params)
}
